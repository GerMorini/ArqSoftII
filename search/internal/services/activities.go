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

type ActiviesServiceImpl struct {
	cache     ActivitiesRepository
	search    ActivitiesRepository
	publisher ActivitiesPublisher
	consumer  ActivitiesConsumer
}

func NewActivitysService(cache ActivitiesRepository, search ActivitiesRepository, publisher ActivitiesPublisher, consumer ActivitiesConsumer) ActiviesServiceImpl {
	return ActiviesServiceImpl{
		cache:     cache,
		search:    search,
		publisher: publisher,
		consumer:  consumer,
	}
}

func (s *ActiviesServiceImpl) List(ctx context.Context, filters dto.SearchFilters) (dto.PaginatedResponse, error) {

	return s.search.List(ctx, filters)
}

func (s *ActiviesServiceImpl) Create(ctx context.Context, activity dto.Activity) (dto.Activity, error) {
	if err := s.publisher.Publish(ctx, "create", created.ID); err != nil {
		return dto.Activity{}, fmt.Errorf("error publishing activity creation: %w", err)
	}

	_, err = s.cache.Create(ctx, created)
	if err != nil {
		return dto.Activity{}, fmt.Errorf("error creating activity in cache: %w", err)
	}

	return created, nil
}

func (s *ActiviesServiceImpl) GetByID(ctx context.Context, id string) (dto.Activity, error) {
	activity, err := s.cache.GetByID(ctx, id)
	if err != nil {
		activity, err := s.repository.GetByID(ctx, id)
		if err != nil {
			return dto.Activity{}, fmt.Errorf("error getting activity from repository: %w", err)
		}

		_, err = s.cache.Create(ctx, activity)
		if err != nil {
			return dto.Activity{}, fmt.Errorf("error creating activity in cache: %w", err)
		}

		return activity, nil
	}
	return activity, nil
}

func (s *ActiviesServiceImpl) Update(ctx context.Context, id string, activity dto.Activity) (dto.Activity, error) {

	if err := s.validateActivity(activity); err != nil {
		return dto.Activity{}, fmt.Errorf("invalid activity: %w", err)
	}

	updated, err := s.repository.Update(ctx, id, activity)
	if err != nil {
		return dto.Activity{}, fmt.Errorf("error updating activity in repository: %w", err)
	}

	if err := s.publisher.Publish(ctx, "update", updated.ID); err != nil {
		return dto.Activity{}, fmt.Errorf("error publishing activity update: %w", err)
	}

	if _, err := s.cache.Update(ctx, id, updated); err != nil {
		return dto.Activity{}, fmt.Errorf("error updating activity in cache: %w", err)
	}

	return updated, nil
}

func (s *ActiviesServiceImpl) Delete(ctx context.Context, id string) error {

	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting activity from repository: %w", err)
	}

	if err := s.publisher.Publish(ctx, "delete", id); err != nil {
		return fmt.Errorf("error publishing activity deletion: %w", err)
	}

	if err := s.cache.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting activity from cache: %w", err)
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

		activity, err := s.repository.GetByID(ctx, message.ActivityID)
		if err != nil {
			slog.Error("❌ Error getting activity for indexing",
				slog.String("activity_id", message.ActivityID),
				slog.String("error", err.Error()))
			return fmt.Errorf("error getting activity for indexing: %w", err)
		}

		if _, err := s.search.Create(ctx, activity); err != nil {
			slog.Error("❌ Error indexing activity in search",
				slog.String("activity_id", message.ActivityID),
				slog.String("error", err.Error()))
		}

		slog.Info("🔍 Activity indexed in search engine", slog.String("activity_id", message.ActivityID))
	case "update":
		slog.Info("✏️ Activity updated", slog.String("activity_id", message.ActivityID))

		activity, err := s.repository.GetByID(ctx, message.ActivityID)
		if err != nil {
			slog.Error("❌ Error getting activity for reindexing",
				slog.String("activity_id", message.ActivityID),
				slog.String("error", err.Error()))
			return fmt.Errorf("error getting activity for indexing: %w", err)
		}

		_, err = s.search.Update(ctx, message.ActivityID, activity)
		if err != nil {
			slog.Error("❌ Error reindexing activity in search",
				slog.String("activity_id", message.ActivityID),
				slog.String("error", err.Error()))
		}

		slog.Info("🔍 Activity reindexed in search engine", slog.String("activity_id", message.ActivityID))
	case "delete":
		slog.Info("🗑️ Activity deleted", slog.String("activity_id", message.ActivityID))
		err := s.search.Delete(ctx, message.ActivityID)

		if err != nil {
			slog.Error("❌ Error deleting activity in search",
				slog.String("activity_id", message.ActivityID),
				slog.String("error", err.Error()))
			return fmt.Errorf("error deleting activity in search: %w", err)
		}
	default:
		slog.Info("⚠️ Unknown action", slog.String("action", message.Action))
	}

	return nil
}
