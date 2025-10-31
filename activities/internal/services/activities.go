package services

import (
	"activities/internal/dto"
	"context"
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// ActivitiesRepository define las operaciones de datos para Activities
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
}

// ActivitiesService define la capa de servicios usada por controllers
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
}

// RabbitMQPublisher interface para publicar eventos
type RabbitMQPublisher interface {
	Publish(ctx context.Context, action string, id string, nombre string, descripcion string, dia string) error
}

// ActivitiesServiceImpl implementa ActivitiesService
type ActivitiesServiceImpl struct {
	repository      ActivitiesRepository
	rabbitPublisher RabbitMQPublisher
}

// NewActivitiesService crea una nueva instancia del service
func NewActivitiesService(repo ActivitiesRepository, rabbit RabbitMQPublisher) *ActivitiesServiceImpl {
	return &ActivitiesServiceImpl{
		repository:      repo,
		rabbitPublisher: rabbit,
	}
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
		return dto.ActivityAdministration{}, fmt.Errorf("validation error: %w", err)
	}

	// Step 1: Create in MongoDB
	created, err := s.repository.Create(ctx, activity)
	if err != nil {
		return dto.ActivityAdministration{}, fmt.Errorf("error creating activity in repository: %w", err)
	}

	// Step 2: Publish event to RabbitMQ (include activity data for search indexing)
	if err := s.rabbitPublisher.Publish(ctx, "create", created.ID, created.Nombre, created.Descripcion, created.DiaSemana); err != nil {
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
		return dto.ActivityAdministration{}, fmt.Errorf("validation error: %w", err)
	}

	// Validar que la nueva capacidad no sea menor a la cantidad de inscritos
	if activity.CapacidadMax > 0 {
		var inscriptosCount int = len(currentActivity.UsersInscribed)
		if activity.CapacidadMax < inscriptosCount {
			return dto.ActivityAdministration{}, fmt.Errorf("capacidadMax cannot be less than the number of inscribed users (%d)", inscriptosCount)
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
				return dto.ActivityAdministration{}, fmt.Errorf("number of inscritos (%d) cannot exceed capacity (%d)", len(activity.UsersInscribed), capToCheck)
			}
		}
	}

	// Step 1: Update in MongoDB
	updated, err := s.repository.Update(ctx, id, activity)
	if err != nil {
		return dto.ActivityAdministration{}, fmt.Errorf("error updating activity in repository: %w", err)
	}

	// Step 2: Publish event to RabbitMQ
	if err := s.rabbitPublisher.Publish(ctx, "update", updated.ID, updated.Nombre, updated.Descripcion, updated.DiaSemana); err != nil {
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
	// Step 0: Get the activity to be able to restore it if needed
	activityToDelete, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("activity does not exist: %w", err)
	}

	// Step 1: Delete from MongoDB
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting activity in repository: %w", err)
	}

	// Step 2: Publish event to RabbitMQ
	if err := s.rabbitPublisher.Publish(ctx, "delete", activityToDelete.ID, activityToDelete.Nombre, activityToDelete.Descripcion, activityToDelete.DiaSemana); err != nil {
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

// validateActivity aplica validaciones simples sobre la actividad
func (s *ActivitiesServiceImpl) validateActivity(a dto.ActivityAdministration) error {
	if strings.TrimSpace(a.Nombre) == "" {
		return errors.New("nombre is required and cannot be empty")
	}
	if strings.TrimSpace(a.Profesor) == "" {
		return errors.New("profesor is required and cannot be empty")
	}
	if strings.TrimSpace(a.HoraInicio) == "" || strings.TrimSpace(a.HoraFin) == "" {
		return errors.New("horaInicio and horaFin are required and cannot be empty")
	}
	if a.CapacidadMax == 0 {
		return errors.New("capacidadMax is required and cannot be empty")
	}
	if a.CapacidadMax < 0 {
		return errors.New("capacidadMax cannot be negative")
	}
	if strings.TrimSpace(a.DiaSemana) == "" {
		return errors.New("diaSemana is required and cannot be empty")
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
		return errors.New("diaSemana must be a valid day of the week (e.g., Lunes, Martes, etc.)")
	}

	// Más validaciones pueden agregarse aquí (horarios, profesor, etc.)
	return nil
}

// GetInscripcionesByUserID obtiene las actividades inscritas por un usuario
func (s *ActivitiesServiceImpl) GetInscripcionesByUserID(ctx context.Context, userID string) ([]string, error) {
	return s.repository.GetInscripcionesByUserID(ctx, userID)
}
