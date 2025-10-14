package services

import (
	"activities/internal/domain"
	"context"
	"errors"
	"fmt"
	"strings"
)

// ActivitiesRepository define las operaciones de datos para Activities
type ActivitiesRepository interface {
	List(ctx context.Context) ([]domain.Activity, error)
	Create(ctx context.Context, activity domain.Activity) (domain.Activity, error)
	GetByID(ctx context.Context, id string) (domain.Activity, error)
	Update(ctx context.Context, id string, activity domain.Activity) (domain.Activity, error)
	Delete(ctx context.Context, id string) error
}

// ActivitiesService define la capa de servicios usada por controllers
type ActivitiesService interface {
	List(ctx context.Context) ([]domain.Activity, error)
	Create(ctx context.Context, activity domain.Activity) (domain.Activity, error)
	GetByID(ctx context.Context, id string) (domain.Activity, error)
	Update(ctx context.Context, id string, activity domain.Activity) (domain.Activity, error)
	Delete(ctx context.Context, id string) error
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
func (s *ActivitiesServiceImpl) List(ctx context.Context) ([]domain.Activity, error) {
	return s.repository.List(ctx)
}

// Create valida y crea una nueva actividad
func (s *ActivitiesServiceImpl) Create(ctx context.Context, activity domain.Activity) (domain.Activity, error) {
	if err := s.validateActivity(activity); err != nil {
		return domain.Activity{}, fmt.Errorf("validation error: %w", err)
	}

	created, err := s.repository.Create(ctx, activity)
	if err != nil {
		return domain.Activity{}, fmt.Errorf("error creating activity in repository: %w", err)
	}

	return created, nil
}

// GetByID obtiene una actividad por ID
func (s *ActivitiesServiceImpl) GetByID(ctx context.Context, id string) (domain.Activity, error) {
	act, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return domain.Activity{}, fmt.Errorf("error getting activity from repository: %w", err)
	}
	return act, nil
}

// Update actualiza una actividad existente
func (s *ActivitiesServiceImpl) Update(ctx context.Context, id string, activity domain.Activity) (domain.Activity, error) {
	_, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return domain.Activity{}, fmt.Errorf("activity does not exist: %w", err)
	}

	if err := s.validateActivity(activity); err != nil {
		return domain.Activity{}, fmt.Errorf("validation error: %w", err)
	}

	updated, err := s.repository.Update(ctx, id, activity)
	if err != nil {
		return domain.Activity{}, fmt.Errorf("error updating activity in repository: %w", err)
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

// validateActivity aplica validaciones simples sobre la actividad
func (s *ActivitiesServiceImpl) validateActivity(a domain.Activity) error {
	if strings.TrimSpace(a.Nombre) == "" {
		return errors.New("nombre is required and cannot be empty")
	}
	// Más validaciones pueden agregarse aquí (horarios, profesor, etc.)
	return nil
}
