package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"search/internal/dto"
	"strings"
)

type ActivitiesRepository interface {
	List(ctx context.Context, filters dto.SearchFilters) (dto.PaginatedResponse, error)
	Create(ctx context.Context, activity dto.Activity) (dto.Activity, error)
	GetByID(ctx context.Context, id string) (dto.Activity, error)
	Update(ctx context.Context, id string, activity dto.Activity) (dto.Activity, error)
	Delete(ctx context.Context, id string) error
}

// TODO: esto deberia estar en el micro de actividades, para publicar mensajes
type ActivitiesPublisher interface {
	Publish(ctx context.Context, action string, activityID string) error
}

type ActivitiesConsumer interface {
	Consume(ctx context.Context, handler func(ctx context.Context, message ActivityEvent) error) error
}

type ActivitiesAPIClient interface {
	GetActivityByID(ctx context.Context, id string) (*dto.Activity, error)
	GetActivitiesByIDs(ctx context.Context, ids []string) ([]dto.Activity, error)
}

type ActiviesServiceImpl struct {
	cache          ActivitiesRepository
	search         ActivitiesRepository
	publisher      ActivitiesPublisher
	consumer       ActivitiesConsumer
	activitiesAPI  ActivitiesAPIClient
}

func NewActivitysService(cache ActivitiesRepository, search ActivitiesRepository, publisher ActivitiesPublisher, consumer ActivitiesConsumer, activitiesAPI ActivitiesAPIClient) ActiviesServiceImpl {
	return ActiviesServiceImpl{
		cache:         cache,
		search:        search,
		publisher:     publisher,
		consumer:      consumer,
		activitiesAPI: activitiesAPI,
	}
}

func (s *ActiviesServiceImpl) List(ctx context.Context, filters dto.SearchFilters) (dto.PaginatedResponse, error) {

	return s.search.List(ctx, filters)
}

func (s *ActiviesServiceImpl) Create(ctx context.Context, activity dto.Activity) (dto.Activity, error) {
	if err := s.publisher.Publish(ctx, "create", activity.ID); err != nil {
		return dto.Activity{}, fmt.Errorf("error publishing activity creation: %w", err)
	}

	_, err := s.cache.Create(ctx, activity)
	if err != nil {
		return dto.Activity{}, fmt.Errorf("error creating activity in cache: %w", err)
	}

	return activity, nil
}

func (s *ActiviesServiceImpl) GetByID(ctx context.Context, id string) (dto.Activity, error) {
	// Try cache first
	activity, err := s.cache.GetByID(ctx, id)
	if err != nil {
		// Cache miss - fetch from Activities API
		activityPtr, err := s.activitiesAPI.GetActivityByID(ctx, id)
		if err != nil {
			return dto.Activity{}, fmt.Errorf("error getting activity from Activities API: %w", err)
		}

		// Cache the result
		_, err = s.cache.Create(ctx, *activityPtr)
		if err != nil {
			slog.Warn("error caching activity", slog.String("activity_id", id), slog.String("error", err.Error()))
		}

		return *activityPtr, nil
	}
	return activity, nil
}

func (s *ActiviesServiceImpl) Update(ctx context.Context, id string, activity dto.Activity) (dto.Activity, error) {
	// This method is for testing only - real updates come via RabbitMQ
	if err := s.validateActivity(activity); err != nil {
		return dto.Activity{}, fmt.Errorf("invalid activity: %w", err)
	}

	// Update in search index
	updated, err := s.search.Update(ctx, id, activity)
	if err != nil {
		return dto.Activity{}, fmt.Errorf("error updating activity in search: %w", err)
	}

	// Update cache
	if _, err := s.cache.Update(ctx, id, updated); err != nil {
		slog.Warn("error updating activity in cache", slog.String("activity_id", id), slog.String("error", err.Error()))
	}

	return updated, nil
}

func (s *ActiviesServiceImpl) Delete(ctx context.Context, id string) error {
	// This method is for testing only - real deletes come via RabbitMQ

	// Delete from search index
	if err := s.search.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting activity from search: %w", err)
	}

	// Delete from cache
	if err := s.cache.Delete(ctx, id); err != nil {
		slog.Warn("error deleting activity from cache", slog.String("activity_id", id), slog.String("error", err.Error()))
	}

	return nil
}

func (s *ActiviesServiceImpl) validateActivity(activity dto.Activity) error {
	if strings.TrimSpace(activity.Titulo) == "" {
		return errors.New("titulo is required and cannot be empty")
	}

	if strings.TrimSpace(activity.Descripcion) == "" {
		return errors.New("descripcion is required and cannot be empty")
	}

	if strings.TrimSpace(activity.DiaSemana) == "" {
		return errors.New("diaSemana is required and cannot be empty")
	}

	return nil
}

type ActivityEvent struct {
	Action     string `json:"action"`
	ActivityID string `json:"activity_id"`
}

func (s *ActiviesServiceImpl) InitConsumer(ctx context.Context) {
	slog.Info("🐰 Starting RabbitMQ consumer...")

	if err := s.consumer.Consume(ctx, s.handleMessage); err != nil {
		slog.Error("❌ Error in RabbitMQ consumer: %v", err)
	}
	slog.Info("🐰 RabbitMQ consumer stopped.")
}

func (s *ActiviesServiceImpl) handleMessage(ctx context.Context, message ActivityEvent) error {
	slog.Info("📨 Processing message",
		slog.String("action", message.Action),
		slog.String("activity_id", message.ActivityID),
	)

	switch message.Action {
	case "create":
		slog.Info("✅ Activity created", slog.String("activity_id", message.ActivityID))

		// Fetch activity from Activities API
		activityPtr, err := s.activitiesAPI.GetActivityByID(ctx, message.ActivityID)
		if err != nil {
			slog.Error("❌ Error getting activity from Activities API for indexing",
				slog.String("activity_id", message.ActivityID),
				slog.String("error", err.Error()))
			return fmt.Errorf("error getting activity for indexing: %w", err)
		}

		// Index in SolR
		if _, err := s.search.Create(ctx, *activityPtr); err != nil {
			slog.Error("❌ Error indexing activity in search",
				slog.String("activity_id", message.ActivityID),
				slog.String("error", err.Error()))
			return fmt.Errorf("error indexing activity: %w", err)
		}

		// Cache activity
		if _, err := s.cache.Create(ctx, *activityPtr); err != nil {
			slog.Warn("⚠️ Error caching activity",
				slog.String("activity_id", message.ActivityID),
				slog.String("error", err.Error()))
		}

		slog.Info("🔍 Activity indexed in search engine", slog.String("activity_id", message.ActivityID))
	case "update":
		slog.Info("✏️ Activity updated", slog.String("activity_id", message.ActivityID))

		// Fetch updated activity from Activities API
		activityPtr, err := s.activitiesAPI.GetActivityByID(ctx, message.ActivityID)
		if err != nil {
			slog.Error("❌ Error getting activity from Activities API for reindexing",
				slog.String("activity_id", message.ActivityID),
				slog.String("error", err.Error()))
			return fmt.Errorf("error getting activity for reindexing: %w", err)
		}

		// Reindex in SolR
		_, err = s.search.Update(ctx, message.ActivityID, *activityPtr)
		if err != nil {
			slog.Error("❌ Error reindexing activity in search",
				slog.String("activity_id", message.ActivityID),
				slog.String("error", err.Error()))
			return fmt.Errorf("error reindexing activity: %w", err)
		}

		// Update cache
		if _, err := s.cache.Update(ctx, message.ActivityID, *activityPtr); err != nil {
			slog.Warn("⚠️ Error updating activity in cache",
				slog.String("activity_id", message.ActivityID),
				slog.String("error", err.Error()))
		}

		slog.Info("🔍 Activity reindexed in search engine", slog.String("activity_id", message.ActivityID))
	case "delete":
		slog.Info("🗑️ Activity deleted", slog.String("activity_id", message.ActivityID))

		// Delete from SolR
		err := s.search.Delete(ctx, message.ActivityID)
		if err != nil {
			slog.Error("❌ Error deleting activity in search",
				slog.String("activity_id", message.ActivityID),
				slog.String("error", err.Error()))
			return fmt.Errorf("error deleting activity in search: %w", err)
		}

		// Delete from cache
		if err := s.cache.Delete(ctx, message.ActivityID); err != nil {
			slog.Warn("⚠️ Error deleting activity from cache",
				slog.String("activity_id", message.ActivityID),
				slog.String("error", err.Error()))
		}

		slog.Info("🗑️ Activity deleted from search engine", slog.String("activity_id", message.ActivityID))
	default:
		slog.Info("⚠️ Unknown action", slog.String("action", message.Action))
	}

	return nil
}
