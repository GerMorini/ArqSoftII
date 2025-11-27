package services

import (
	"activities/internal/dto"
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	ErrValidation = errors.New("validation error")
)

type ActivitiesRepository interface {
	List(ctx context.Context) ([]dto.Activity, error)
	GetMany(ctx context.Context, ids []string) ([]dto.Activity, error)
	Create(ctx context.Context, activity dto.ActivityAdministration) (dto.ActivityAdministration, error)
	GetByID(ctx context.Context, id string) (dto.ActivityAdministration, error)
	Update(ctx context.Context, id string, activity dto.ActivityAdministration) (dto.ActivityAdministration, error)
	Delete(ctx context.Context, id string) error
	Inscribir(ctx context.Context, id string, userID string) (string, error)
	Desinscribir(ctx context.Context, id string, userID string) (string, error)
	GetInscripcionesByUserID(ctx context.Context, userID string) ([]string, error)
	GetActivitiesByUserID(ctx context.Context, userID string) (dto.Activities, error)
	ListAllForAdmin(ctx context.Context) ([]dto.ActivityAdministration, error)
}

type ActivitiesService interface {
	List(ctx context.Context) ([]dto.Activity, error)
	GetMany(ctx context.Context, ids []string) ([]dto.Activity, error)
	Create(ctx context.Context, activity dto.ActivityAdministration) (dto.ActivityAdministration, error)
	GetByID(ctx context.Context, id string) (dto.ActivityAdministration, error)
	Update(ctx context.Context, id string, activity dto.ActivityAdministration) (dto.ActivityAdministration, error)
	Delete(ctx context.Context, id string) error
	Inscribir(ctx context.Context, id string, userID string) (string, error)
	Desinscribir(ctx context.Context, id string, userID string) (string, error)
	GetInscripcionesByUserID(ctx context.Context, userID string) ([]string, error)
	GetStatistics(ctx context.Context) (dto.ActivityStatistics, error)
}

type RabbitMQPublisher interface {
	Publish(ctx context.Context, action string, id string) error
}

type ActivitiesServiceImpl struct {
	repository      ActivitiesRepository
	rabbitPublisher RabbitMQPublisher
}

func NewActivitiesService(repo ActivitiesRepository, rabbit RabbitMQPublisher) *ActivitiesServiceImpl {
	return &ActivitiesServiceImpl{
		repository:      repo,
		rabbitPublisher: rabbit,
	}
}

func (s *ActivitiesServiceImpl) validateActivity(a dto.ActivityAdministration) error {
	if strings.TrimSpace(a.Nombre) == "" {
		return errors.New("titulo is required and cannot be empty")
	}
	if strings.TrimSpace(a.Profesor) == "" {
		return errors.New("instructor is required and cannot be empty")
	}
	if strings.TrimSpace(a.HoraInicio) == "" || strings.TrimSpace(a.HoraFin) == "" {
		return errors.New("hora_inicio and hora_fin are required and cannot be empty")
	}
	if a.CapacidadMax == 0 {
		return errors.New("cupo is required and cannot be empty")
	}
	if a.CapacidadMax < 0 {
		return errors.New("cupo cannot be negative")
	}
	if strings.TrimSpace(a.DiaSemana) == "" {
		return errors.New("dia is required and cannot be empty")
	}
	validDays := map[string]bool{
		"Lunes":     true,
		"Martes":    true,
		"Miércoles": true,
		"Jueves":    true,
		"Viernes":   true,
		"Sábado":    true,
		"Domingo":   true,
	}
	if !validDays[a.DiaSemana] {
		return errors.New("dia must be a valid day of the week (e.g., Lunes, Martes, etc.)")
	}

	return nil
}

// List obtiene todas las actividades
func (s *ActivitiesServiceImpl) List(ctx context.Context) ([]dto.Activity, error) {
	return s.repository.List(ctx)
}

// GetMany obtiene multiples actividades por IDs (ignora IDs no encontrados)
func (s *ActivitiesServiceImpl) GetMany(ctx context.Context, ids []string) ([]dto.Activity, error) {
	return s.repository.GetMany(ctx, ids)
}

// Create valida y crea una nueva actividad
func (s *ActivitiesServiceImpl) Create(ctx context.Context, activity dto.ActivityAdministration) (dto.ActivityAdministration, error) {
	if err := s.validateActivity(activity); err != nil {
		return dto.ActivityAdministration{}, fmt.Errorf("%w: %v", ErrValidation, err)
	}

	created, err := s.repository.Create(ctx, activity)
	if err != nil {
		return dto.ActivityAdministration{}, fmt.Errorf("error creating activity in repository: %w", err)
	}

	if err := s.rabbitPublisher.Publish(ctx, "create", created.ID); err != nil {
		log.Errorf("Failed to publish create event for activity %s: %v", created.ID, err)

		// Rollback: delete the created activity from MongoDB
		if deleteErr := s.repository.Delete(ctx, created.ID); deleteErr != nil {
			log.Errorf("CRITICAL: Failed to rollback activity %s after RabbitMQ publish failure: %v", created.ID, deleteErr)
			return dto.ActivityAdministration{}, fmt.Errorf("failed to publish event and rollback failed: publish error: %w, rollback error: %v", err, deleteErr)
		}

		log.Warnf("Successfully rolled back activity %s after RabbitMQ publish failure", created.ID)
		return dto.ActivityAdministration{}, fmt.Errorf("failed to publish create event, activity rolled back: %w", err)
	}

	log.Infof("Activity %s created and event published successfully", created.ID)
	return created, nil
}

// GetByID obtiene una actividad por ID
func (s *ActivitiesServiceImpl) GetByID(ctx context.Context, id string) (dto.ActivityAdministration, error) {
	act, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return dto.ActivityAdministration{}, fmt.Errorf("error getting activity from repository: %w", err)
	}
	return act, nil
}

// Update actualiza una actividad existente
func (s *ActivitiesServiceImpl) Update(ctx context.Context, id string, activity dto.ActivityAdministration) (dto.ActivityAdministration, error) {
	currentActivity, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return dto.ActivityAdministration{}, fmt.Errorf("activity does not exist: %w", err)
	}

	if err := s.validateActivity(activity); err != nil {
		return dto.ActivityAdministration{}, fmt.Errorf("%w: %v", ErrValidation, err)
	}

	// Validar que la nueva capacidad no sea menor a la cantidad de inscritos
	if activity.CapacidadMax > 0 {
		var inscriptosCount int = len(currentActivity.UsersInscribed)
		if activity.CapacidadMax < inscriptosCount {
			return dto.ActivityAdministration{}, fmt.Errorf("%w: cupo cannot be less than the number of inscribed users (%d)", ErrValidation, inscriptosCount)
		}
	}

	// If admin provided an explicit users list, validate it against capacity
	if activity.UsersInscribed != nil {
		// determine capacity to compare: prefer updated capacity if provided, otherwise current
		capToCheck := 0
		if activity.CapacidadMax > 0 {
			capToCheck = activity.CapacidadMax
		} else {
			capToCheck = currentActivity.CapacidadMax
		}
		if capToCheck > 0 {
			if len(activity.UsersInscribed) > capToCheck {
				return dto.ActivityAdministration{}, fmt.Errorf("%w: number of inscritos (%d) cannot exceed capacity (%d)", ErrValidation, len(activity.UsersInscribed), capToCheck)
			}
		}
	}

	updated, err := s.repository.Update(ctx, id, activity)
	if err != nil {
		return dto.ActivityAdministration{}, err
	}

	if err := s.rabbitPublisher.Publish(ctx, "update", updated.ID); err != nil {
		log.Errorf("Failed to publish update event for activity %s: %v", id, err)

		// Rollback: restore the original activity in MongoDB
		if _, restoreErr := s.repository.Update(ctx, id, currentActivity); restoreErr != nil {
			log.Errorf("CRITICAL: Failed to rollback activity %s after RabbitMQ publish failure: %v", id, restoreErr)
			return dto.ActivityAdministration{}, fmt.Errorf("failed to publish event and rollback failed: publish error: %w, rollback error: %v", err, restoreErr)
		}

		log.Warnf("Successfully rolled back activity %s after RabbitMQ publish failure", id)
		return dto.ActivityAdministration{}, fmt.Errorf("failed to publish update event, activity rolled back: %w", err)
	}

	log.Infof("Activity %s updated and event published successfully", id)
	return updated, nil
}

// Delete elimina una actividad por ID
func (s *ActivitiesServiceImpl) Delete(ctx context.Context, id string) error {
	activityToDelete, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("activity does not exist: %w", err)
	}

	if err := s.repository.Delete(ctx, id); err != nil {
		return err
	}

	if err := s.rabbitPublisher.Publish(ctx, "delete", activityToDelete.ID); err != nil {
		log.Errorf("Failed to publish delete event for activity %s: %v", id, err)

		// Rollback: restore the deleted activity in MongoDB
		if _, restoreErr := s.repository.Create(ctx, activityToDelete); restoreErr != nil {
			log.Errorf("CRITICAL: Failed to rollback activity %s after RabbitMQ publish failure: %v", id, restoreErr)
			return fmt.Errorf("failed to publish event and rollback failed: publish error: %w, rollback error: %v", err, restoreErr)
		}

		log.Warnf("Successfully rolled back activity %s after RabbitMQ publish failure", id)
		return fmt.Errorf("failed to publish delete event, activity rolled back: %w", err)
	}

	log.Infof("Activity %s deleted and event published successfully", id)
	return nil
}

// Inscribir registra al usuario en la actividad
func (s *ActivitiesServiceImpl) Inscribir(ctx context.Context, id string, userID string) (string, error) {
	return s.repository.Inscribir(ctx, id, userID)
}

// Desinscribir quita al usuario de la actividad
func (s *ActivitiesServiceImpl) Desinscribir(ctx context.Context, id string, userID string) (string, error) {
	return s.repository.Desinscribir(ctx, id, userID)
}

// GetInscripcionesByUserID obtiene las actividades inscritas por un usuario
func (s *ActivitiesServiceImpl) GetInscripcionesByUserID(ctx context.Context, userID string) ([]string, error) {
	return s.repository.GetInscripcionesByUserID(ctx, userID)
}

func (s *ActivitiesServiceImpl) GetActivitiesByUserID(ctx context.Context, userID string) (dto.Activities, error) {
	return s.repository.GetActivitiesByUserID(ctx, userID)
}

// GetStatistics calcula estadísticas de actividades usando concurrencia
func (s *ActivitiesServiceImpl) GetStatistics(ctx context.Context) (dto.ActivityStatistics, error) {
	activities, err := s.repository.ListAllForAdmin(ctx)
	if err != nil {
		log.WithError(err).Error("Error fetching activities for statistics")
		return dto.ActivityStatistics{}, err
	}

	totalEnrollmentsChan := make(chan int)
	capacityUtilizationChan := make(chan float64)
	dayDistributionChan := make(chan []dto.DayDistribution)
	mostPopularChan := make(chan *dto.Activity)
	fullActivitiesChan := make(chan int)

	var wg sync.WaitGroup

	// Goroutine 1: Calcular total de inscripciones y capacidad total
	wg.Add(1)
	go func() {
		defer wg.Done()
		totalEnrollments := 0
		totalCapacity := 0

		for _, act := range activities {
			totalEnrollments += len(act.UsersInscribed)
			totalCapacity += act.CapacidadMax
		}

		// Calcular utilización de capacidad
		var utilization float64
		if totalCapacity > 0 {
			utilization = (float64(totalEnrollments) / float64(totalCapacity)) * 100
		}

		totalEnrollmentsChan <- totalEnrollments
		capacityUtilizationChan <- utilization
	}()

	// Goroutine 2: Calcular distribución por día de la semana
	wg.Add(1)
	go func() {
		defer wg.Done()
		dayCount := make(map[string]int)

		for _, act := range activities {
			dayCount[act.DiaSemana]++
		}

		distribution := []dto.DayDistribution{}
		for day, count := range dayCount {
			distribution = append(distribution, dto.DayDistribution{
				Dia:   day,
				Count: count,
			})
		}

		dayDistributionChan <- distribution
	}()

	// Goroutine 3: Encontrar actividad más popular
	wg.Add(1)
	go func() {
		defer wg.Done()
		var mostPopular *dto.Activity
		maxEnrollments := -1

		for _, act := range activities {
			if len(act.UsersInscribed) > maxEnrollments {
				maxEnrollments = len(act.UsersInscribed)
				activity := act.Activity // Copy the embedded Activity struct
				mostPopular = &activity
			}
		}

		mostPopularChan <- mostPopular
	}()

	// Goroutine 4: Contar actividades llenas vs disponibles
	wg.Add(1)
	go func() {
		defer wg.Done()
		fullCount := 0

		for _, act := range activities {
			if len(act.UsersInscribed) >= act.CapacidadMax {
				fullCount++
			}
		}

		fullActivitiesChan <- fullCount
	}()

	// Goroutine para cerrar channels después de que todas las goroutines terminen
	go func() {
		wg.Wait()
		close(totalEnrollmentsChan)
		close(capacityUtilizationChan)
		close(dayDistributionChan)
		close(mostPopularChan)
		close(fullActivitiesChan)
	}()

	totalEnrollments := <-totalEnrollmentsChan
	capacityUtilization := <-capacityUtilizationChan
	dayDistribution := <-dayDistributionChan
	mostPopular := <-mostPopularChan
	fullActivities := <-fullActivitiesChan

	// Calcular tasa promedio de inscripción
	var avgEnrollmentRate float64
	if len(activities) > 0 {
		avgEnrollmentRate = float64(totalEnrollments) / float64(len(activities))
	}

	// Calculate total capacity
	var totalCapacity int
	for _, act := range activities {
		totalCapacity += act.CapacidadMax
	}

	stats := dto.ActivityStatistics{
		TotalActivities:       len(activities),
		TotalEnrollments:      totalEnrollments,
		AverageEnrollmentRate: avgEnrollmentRate,
		TotalCapacity:         totalCapacity,
		CapacityUtilization:   capacityUtilization,
		ActivitiesByDay:       dayDistribution,
		MostPopularActivity:   mostPopular,
		FullActivitiesCount:   fullActivities,
		AvailableActivities:   len(activities) - fullActivities,
	}

	log.WithFields(log.Fields{
		"total_activities":  stats.TotalActivities,
		"total_enrollments": stats.TotalEnrollments,
		"full_activities":   stats.FullActivitiesCount,
	}).Info("Statistics calculated successfully")

	return stats, nil
}
