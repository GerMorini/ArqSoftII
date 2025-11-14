package services

import (
	"activities/internal/dto"
	"context"
	"errors"
	"testing"
)

// Mock implementations
type mockRepo struct {
	listFunc                     func(ctx context.Context) ([]dto.Activity, error)
	getManyFunc                  func(ctx context.Context, ids []string) ([]dto.Activity, error)
	createFunc                   func(ctx context.Context, activity dto.ActivityAdministration) (dto.ActivityAdministration, error)
	getByIDFunc                  func(ctx context.Context, id string) (dto.ActivityAdministration, error)
	updateFunc                   func(ctx context.Context, id string, activity dto.ActivityAdministration) (dto.ActivityAdministration, error)
	deleteFunc                   func(ctx context.Context, id string) error
	inscribirFunc                func(ctx context.Context, id string, userID string) (string, error)
	desinscribirFunc             func(ctx context.Context, id string, userID string) (string, error)
	getInscripcionesByUserIDFunc func(ctx context.Context, userID string) ([]string, error)
	listAllForAdminFunc          func(ctx context.Context) ([]dto.ActivityAdministration, error)
}

func (m *mockRepo) List(ctx context.Context) ([]dto.Activity, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx)
	}
	return nil, nil
}

func (m *mockRepo) GetMany(ctx context.Context, ids []string) ([]dto.Activity, error) {
	if m.getManyFunc != nil {
		return m.getManyFunc(ctx, ids)
	}
	return nil, nil
}

func (m *mockRepo) Create(ctx context.Context, activity dto.ActivityAdministration) (dto.ActivityAdministration, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, activity)
	}
	return dto.ActivityAdministration{}, nil
}

func (m *mockRepo) GetByID(ctx context.Context, id string) (dto.ActivityAdministration, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return dto.ActivityAdministration{}, nil
}

func (m *mockRepo) Update(ctx context.Context, id string, activity dto.ActivityAdministration) (dto.ActivityAdministration, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, activity)
	}
	return dto.ActivityAdministration{}, nil
}

func (m *mockRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func (m *mockRepo) Inscribir(ctx context.Context, id string, userID string) (string, error) {
	if m.inscribirFunc != nil {
		return m.inscribirFunc(ctx, id, userID)
	}
	return "", nil
}

func (m *mockRepo) Desinscribir(ctx context.Context, id string, userID string) (string, error) {
	if m.desinscribirFunc != nil {
		return m.desinscribirFunc(ctx, id, userID)
	}
	return "", nil
}

func (m *mockRepo) GetInscripcionesByUserID(ctx context.Context, userID string) ([]string, error) {
	if m.getInscripcionesByUserIDFunc != nil {
		return m.getInscripcionesByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *mockRepo) GetActivitiesByUserID(ctx context.Context, userID string) (dto.Activities, error) {
	return nil, nil
}

func (m *mockRepo) ListAllForAdmin(ctx context.Context) ([]dto.ActivityAdministration, error) {
	if m.listAllForAdminFunc != nil {
		return m.listAllForAdminFunc(ctx)
	}
	return nil, nil
}

type mockRabbit struct {
	publishFunc func(ctx context.Context, action string, id string, nombre string, descripcion string, dia string) error
}

func (m *mockRabbit) Publish(ctx context.Context, action string, id string, nombre string, descripcion string, dia string) error {
	if m.publishFunc != nil {
		return m.publishFunc(ctx, action, id, nombre, descripcion, dia)
	}
	return nil
}

// TestList tests the List method
func TestList(t *testing.T) {
	ctx := context.Background()

	// Happy path
	t.Run("success", func(t *testing.T) {
		mockRepo := &mockRepo{
			listFunc: func(ctx context.Context) ([]dto.Activity, error) {
				return []dto.Activity{
					{ID: "1", Nombre: "Yoga"},
					{ID: "2", Nombre: "Pilates"},
				}, nil
			},
		}
		mockRabbit := &mockRabbit{}
		service := NewActivitiesService(mockRepo, mockRabbit)

		result, err := service.List(ctx)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 activities, got %d", len(result))
		}
	})

	// Error from repository
	t.Run("repository error", func(t *testing.T) {
		mockRepo := &mockRepo{
			listFunc: func(ctx context.Context) ([]dto.Activity, error) {
				return nil, errors.New("db error")
			},
		}
		mockRabbit := &mockRabbit{}
		service := NewActivitiesService(mockRepo, mockRabbit)

		_, err := service.List(ctx)

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

// TestCreate tests the Create method
func TestCreate(t *testing.T) {
	ctx := context.Background()

	validActivity := dto.ActivityAdministration{
		Activity: dto.Activity{
			Nombre:       "Yoga",
			Profesor:     "Juan Perez",
			HoraInicio:   "10:00",
			HoraFin:      "11:00",
			CapacidadMax: 20,
			DiaSemana:    "Lunes",
		},
	}

	// Happy path
	t.Run("success", func(t *testing.T) {
		mockRepo := &mockRepo{
			createFunc: func(ctx context.Context, activity dto.ActivityAdministration) (dto.ActivityAdministration, error) {
				activity.ID = "123"
				return activity, nil
			},
		}
		mockRabbit := &mockRabbit{
			publishFunc: func(ctx context.Context, action, id, nombre, descripcion, dia string) error {
				return nil
			},
		}
		service := NewActivitiesService(mockRepo, mockRabbit)

		result, err := service.Create(ctx, validActivity)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result.ID != "123" {
			t.Errorf("expected ID 123, got %s", result.ID)
		}
	})

	// Validation error
	t.Run("validation error - empty nombre", func(t *testing.T) {
		invalidActivity := validActivity
		invalidActivity.Nombre = ""

		mockRepo := &mockRepo{}
		mockRabbit := &mockRabbit{}
		service := NewActivitiesService(mockRepo, mockRabbit)

		_, err := service.Create(ctx, invalidActivity)

		if err == nil {
			t.Error("expected validation error, got nil")
		}
	})

	// Repository error
	t.Run("repository error", func(t *testing.T) {
		mockRepo := &mockRepo{
			createFunc: func(ctx context.Context, activity dto.ActivityAdministration) (dto.ActivityAdministration, error) {
				return dto.ActivityAdministration{}, errors.New("db error")
			},
		}
		mockRabbit := &mockRabbit{}
		service := NewActivitiesService(mockRepo, mockRabbit)

		_, err := service.Create(ctx, validActivity)

		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	// RabbitMQ error with rollback
	t.Run("rabbitmq error with rollback", func(t *testing.T) {
		mockRepo := &mockRepo{
			createFunc: func(ctx context.Context, activity dto.ActivityAdministration) (dto.ActivityAdministration, error) {
				activity.ID = "123"
				return activity, nil
			},
			deleteFunc: func(ctx context.Context, id string) error {
				return nil
			},
		}
		mockRabbit := &mockRabbit{
			publishFunc: func(ctx context.Context, action, id, nombre, descripcion, dia string) error {
				return errors.New("rabbitmq error")
			},
		}
		service := NewActivitiesService(mockRepo, mockRabbit)

		_, err := service.Create(ctx, validActivity)

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

// TestUpdate tests the Update method
func TestUpdate(t *testing.T) {
	ctx := context.Background()

	existingActivity := dto.ActivityAdministration{
		Activity: dto.Activity{
			ID:           "1",
			Nombre:       "Yoga",
			Profesor:     "Juan Perez",
			HoraInicio:   "10:00",
			HoraFin:      "11:00",
			CapacidadMax: 20,
			DiaSemana:    "Lunes",
		},
		UsersInscribed: []int{1, 2, 3},
	}

	validUpdate := dto.ActivityAdministration{
		Activity: dto.Activity{
			Nombre:       "Yoga Avanzado",
			Profesor:     "Juan Perez",
			HoraInicio:   "10:00",
			HoraFin:      "11:00",
			CapacidadMax: 25,
			DiaSemana:    "Martes",
		},
	}

	// Happy path
	t.Run("success", func(t *testing.T) {
		mockRepo := &mockRepo{
			getByIDFunc: func(ctx context.Context, id string) (dto.ActivityAdministration, error) {
				return existingActivity, nil
			},
			updateFunc: func(ctx context.Context, id string, activity dto.ActivityAdministration) (dto.ActivityAdministration, error) {
				activity.ID = id
				return activity, nil
			},
		}
		mockRabbit := &mockRabbit{
			publishFunc: func(ctx context.Context, action, id, nombre, descripcion, dia string) error {
				return nil
			},
		}
		service := NewActivitiesService(mockRepo, mockRabbit)

		result, err := service.Update(ctx, "1", validUpdate)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result.Nombre != "Yoga Avanzado" {
			t.Errorf("expected name 'Yoga Avanzado', got %s", result.Nombre)
		}
	})

	// Not found
	t.Run("not found", func(t *testing.T) {
		mockRepo := &mockRepo{
			getByIDFunc: func(ctx context.Context, id string) (dto.ActivityAdministration, error) {
				return dto.ActivityAdministration{}, errors.New("not found")
			},
		}
		mockRabbit := &mockRabbit{}
		service := NewActivitiesService(mockRepo, mockRabbit)

		_, err := service.Update(ctx, "999", validUpdate)

		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	// Capacity less than inscribed users
	t.Run("capacity less than inscribed users", func(t *testing.T) {
		invalidUpdate := validUpdate
		invalidUpdate.CapacidadMax = 2

		mockRepo := &mockRepo{
			getByIDFunc: func(ctx context.Context, id string) (dto.ActivityAdministration, error) {
				return existingActivity, nil
			},
		}
		mockRabbit := &mockRabbit{}
		service := NewActivitiesService(mockRepo, mockRabbit)

		_, err := service.Update(ctx, "1", invalidUpdate)

		if err == nil {
			t.Error("expected error for capacity less than inscribed users, got nil")
		}
	})

	// Validation error
	t.Run("validation error", func(t *testing.T) {
		invalidUpdate := validUpdate
		invalidUpdate.Nombre = ""

		mockRepo := &mockRepo{
			getByIDFunc: func(ctx context.Context, id string) (dto.ActivityAdministration, error) {
				return existingActivity, nil
			},
		}
		mockRabbit := &mockRabbit{}
		service := NewActivitiesService(mockRepo, mockRabbit)

		_, err := service.Update(ctx, "1", invalidUpdate)

		if err == nil {
			t.Error("expected validation error, got nil")
		}
	})

	// Repository error
	t.Run("repository error", func(t *testing.T) {
		mockRepo := &mockRepo{
			getByIDFunc: func(ctx context.Context, id string) (dto.ActivityAdministration, error) {
				return existingActivity, nil
			},
			updateFunc: func(ctx context.Context, id string, activity dto.ActivityAdministration) (dto.ActivityAdministration, error) {
				return dto.ActivityAdministration{}, errors.New("db error")
			},
		}
		mockRabbit := &mockRabbit{}
		service := NewActivitiesService(mockRepo, mockRabbit)

		_, err := service.Update(ctx, "1", validUpdate)

		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	// RabbitMQ error with rollback
	t.Run("rabbitmq error with rollback", func(t *testing.T) {
		mockRepo := &mockRepo{
			getByIDFunc: func(ctx context.Context, id string) (dto.ActivityAdministration, error) {
				return existingActivity, nil
			},
			updateFunc: func(ctx context.Context, id string, activity dto.ActivityAdministration) (dto.ActivityAdministration, error) {
				activity.ID = id
				return activity, nil
			},
		}
		mockRabbit := &mockRabbit{
			publishFunc: func(ctx context.Context, action, id, nombre, descripcion, dia string) error {
				return errors.New("rabbitmq error")
			},
		}
		service := NewActivitiesService(mockRepo, mockRabbit)

		_, err := service.Update(ctx, "1", validUpdate)

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

// TestDelete tests the Delete method
func TestDelete(t *testing.T) {
	ctx := context.Background()

	existingActivity := dto.ActivityAdministration{
		Activity: dto.Activity{
			ID:           "1",
			Nombre:       "Yoga",
			Profesor:     "Juan Perez",
			HoraInicio:   "10:00",
			HoraFin:      "11:00",
			CapacidadMax: 20,
			DiaSemana:    "Lunes",
		},
	}

	// Happy path
	t.Run("success", func(t *testing.T) {
		mockRepo := &mockRepo{
			getByIDFunc: func(ctx context.Context, id string) (dto.ActivityAdministration, error) {
				return existingActivity, nil
			},
			deleteFunc: func(ctx context.Context, id string) error {
				return nil
			},
		}
		mockRabbit := &mockRabbit{
			publishFunc: func(ctx context.Context, action, id, nombre, descripcion, dia string) error {
				return nil
			},
		}
		service := NewActivitiesService(mockRepo, mockRabbit)

		err := service.Delete(ctx, "1")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	// Activity not found
	t.Run("not found", func(t *testing.T) {
		mockRepo := &mockRepo{
			getByIDFunc: func(ctx context.Context, id string) (dto.ActivityAdministration, error) {
				return dto.ActivityAdministration{}, errors.New("not found")
			},
		}
		mockRabbit := &mockRabbit{}
		service := NewActivitiesService(mockRepo, mockRabbit)

		err := service.Delete(ctx, "999")

		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	// Repository delete error
	t.Run("repository error", func(t *testing.T) {
		mockRepo := &mockRepo{
			getByIDFunc: func(ctx context.Context, id string) (dto.ActivityAdministration, error) {
				return existingActivity, nil
			},
			deleteFunc: func(ctx context.Context, id string) error {
				return errors.New("db error")
			},
		}
		mockRabbit := &mockRabbit{}
		service := NewActivitiesService(mockRepo, mockRabbit)

		err := service.Delete(ctx, "1")

		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	// RabbitMQ error with rollback
	t.Run("rabbitmq error with rollback", func(t *testing.T) {
		mockRepo := &mockRepo{
			getByIDFunc: func(ctx context.Context, id string) (dto.ActivityAdministration, error) {
				return existingActivity, nil
			},
			deleteFunc: func(ctx context.Context, id string) error {
				return nil
			},
			createFunc: func(ctx context.Context, activity dto.ActivityAdministration) (dto.ActivityAdministration, error) {
				return activity, nil
			},
		}
		mockRabbit := &mockRabbit{
			publishFunc: func(ctx context.Context, action, id, nombre, descripcion, dia string) error {
				return errors.New("rabbitmq error")
			},
		}
		service := NewActivitiesService(mockRepo, mockRabbit)

		err := service.Delete(ctx, "1")

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

// TestInscribir tests the Inscribir method
func TestInscribir(t *testing.T) {
	ctx := context.Background()

	// Happy path
	t.Run("success", func(t *testing.T) {
		mockRepo := &mockRepo{
			inscribirFunc: func(ctx context.Context, id, userID string) (string, error) {
				return "inscribed", nil
			},
		}
		mockRabbit := &mockRabbit{}
		service := NewActivitiesService(mockRepo, mockRabbit)

		result, err := service.Inscribir(ctx, "1", "100")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result != "inscribed" {
			t.Errorf("expected 'inscribed', got %s", result)
		}
	})

	// Repository error
	t.Run("repository error", func(t *testing.T) {
		mockRepo := &mockRepo{
			inscribirFunc: func(ctx context.Context, id, userID string) (string, error) {
				return "", errors.New("activity full")
			},
		}
		mockRabbit := &mockRabbit{}
		service := NewActivitiesService(mockRepo, mockRabbit)

		_, err := service.Inscribir(ctx, "1", "100")

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

// TestGetInscripcionesByUserID tests the GetInscripcionesByUserID method
func TestGetInscripcionesByUserID(t *testing.T) {
	ctx := context.Background()

	// Happy path
	t.Run("success", func(t *testing.T) {
		mockRepo := &mockRepo{
			getInscripcionesByUserIDFunc: func(ctx context.Context, userID string) ([]string, error) {
				return []string{"1", "2", "3"}, nil
			},
		}
		mockRabbit := &mockRabbit{}
		service := NewActivitiesService(mockRepo, mockRabbit)

		result, err := service.GetInscripcionesByUserID(ctx, "100")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(result) != 3 {
			t.Errorf("expected 3 inscriptions, got %d", len(result))
		}
	})

	// Repository error
	t.Run("repository error", func(t *testing.T) {
		mockRepo := &mockRepo{
			getInscripcionesByUserIDFunc: func(ctx context.Context, userID string) ([]string, error) {
				return nil, errors.New("db error")
			},
		}
		mockRabbit := &mockRabbit{}
		service := NewActivitiesService(mockRepo, mockRabbit)

		_, err := service.GetInscripcionesByUserID(ctx, "100")

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
