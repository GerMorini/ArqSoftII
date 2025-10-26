package services

import (
	"activities/internal/dto"
	"context"
	"errors"
	"fmt"
	"strings"
)

// ActivitiesRepository define las operaciones de datos para Activities
type ActivitiesRepository interface {
	List(ctx context.Context) ([]dto.Activity, error)
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
	Create(ctx context.Context, activity dto.ActivityAdministration) (dto.ActivityAdministration, error)
	GetByID(ctx context.Context, id string) (dto.ActivityAdministration, error)
	Update(ctx context.Context, id string, activity dto.ActivityAdministration) (dto.ActivityAdministration, error)
	Delete(ctx context.Context, id string) error
	Inscribir(ctx context.Context, id string, userID string) (string, error)
	Desinscribir(ctx context.Context, id string, userID string) (string, error)
	GetInscripcionesByUserID(ctx context.Context, userID string) ([]string, error)
}

// ActivitiesServiceImpl implementa ActivitiesService
type ActivitiesServiceImpl struct {
	repository ActivitiesRepository
}

// NewActivitiesService crea una nueva instancia del service
func NewActivitiesService(repo ActivitiesRepository) *ActivitiesServiceImpl {
	return &ActivitiesServiceImpl{repository: repo}
}

// List obtiene todas las actividades
func (s *ActivitiesServiceImpl) List(ctx context.Context) ([]dto.Activity, error) {
	return s.repository.List(ctx)
}

// Create valida y crea una nueva actividad
func (s *ActivitiesServiceImpl) Create(ctx context.Context, activity dto.ActivityAdministration) (dto.ActivityAdministration, error) {
	if err := s.validateActivity(activity); err != nil {
		return dto.ActivityAdministration{}, fmt.Errorf("validation error: %w", err)
	}

	created, err := s.repository.Create(ctx, activity)
	if err != nil {
		return dto.ActivityAdministration{}, fmt.Errorf("error creating activity in repository: %w", err)
	}

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
	if activity.CapacidadMax != "" {
		var newCapacity int
		if _, err := fmt.Sscanf(activity.CapacidadMax, "%d", &newCapacity); err == nil {
			inscriptosCount := len(currentActivity.UsersInscribed)
			if newCapacity < inscriptosCount {
				return dto.ActivityAdministration{}, fmt.Errorf("capacidadMax cannot be less than the number of inscribed users (%d)", inscriptosCount)
			}
		}
	}

	updated, err := s.repository.Update(ctx, id, activity)
	if err != nil {
		return dto.ActivityAdministration{}, fmt.Errorf("error updating activity in repository: %w", err)
	}

	return updated, nil
}

// Delete elimina una actividad por ID
func (s *ActivitiesServiceImpl) Delete(ctx context.Context, id string) error {
	_, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("activity does not exist: %w", err)
	}

	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting activity in repository: %w", err)
	}

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
	if strings.TrimSpace(a.CapacidadMax) == "" {
		return errors.New("capacidadMax is required and cannot be empty")
	}
	if _, err := fmt.Sscanf(a.CapacidadMax, "%d", new(int)); err != nil {
		return errors.New("capacidadMax must be a valid integer")
	}
	for i := 0; i < len(a.UsersInscribed); i++ {
		if strings.TrimSpace(a.UsersInscribed[i]) == "" {
			return fmt.Errorf("user ID at position %d cannot be empty", i)
		}
		if _, err := fmt.Sscanf(a.UsersInscribed[i], "%d", new(int)); err != nil {
			return fmt.Errorf("user ID at position %d must be a valid integer", i)
		}
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
