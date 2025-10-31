package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"search/internal/dto"
	"strings"
)

type ActivityEvent struct {
	Action      string `json:"action"`
	Nombre      string
	Descripcion string
	Dia         string
}

type ActivitiesRepository interface {
	List(ctx context.Context, filters dto.SearchFilters) (dto.PaginatedResponse, error)
	Create(ctx context.Context, activity dto.Activity) (dto.Activity, error)
	Update(ctx context.Context, id string, activity dto.Activity) (dto.Activity, error)
	Delete(ctx context.Context, id string) error
}

type ActivitiesConsumer interface {
	Consume(ctx context.Context, handler func(ctx context.Context, message ActivityEvent) error) error
}

type ActiviesServiceImpl struct {
	cache    ActivitiesRepository
	search   ActivitiesRepository
	consumer ActivitiesConsumer
}

func NewActivitysService(cache ActivitiesRepository, search ActivitiesRepository, consumer ActivitiesConsumer) ActiviesServiceImpl {
	return ActiviesServiceImpl{
		cache:    cache,
		search:   search,
		consumer: consumer,
	}
}

func (s *ActiviesServiceImpl) List(ctx context.Context, filters dto.SearchFilters) (dto.PaginatedResponse, error) {
	return s.search.List(ctx, filters)
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
