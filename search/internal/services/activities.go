package services

import (
	"context"
	"fmt"
	"log/slog"
	"search/internal/dto"
)

type ActivityEvent struct {
	Action      string `json:"action"`
	ID          string `json:"id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Dia         string `json:"dia"`
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
	localCache ActivitiesRepository
	cache      ActivitiesRepository
	search     ActivitiesRepository
	consumer   ActivitiesConsumer
}

func NewActivitiesService(localCache ActivitiesRepository, cache ActivitiesRepository, search ActivitiesRepository, consumer ActivitiesConsumer) ActiviesServiceImpl {
	return ActiviesServiceImpl{
		localCache: localCache,
		cache:      cache,
		search:     search,
		consumer:   consumer,
	}
}

func (s *ActiviesServiceImpl) List(ctx context.Context, filters dto.SearchFilters) (dto.PaginatedResponse, error) {
	return s.search.List(ctx, filters)
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
		slog.String("id", message.ID),
		slog.String("nombre", message.Nombre),
		slog.String("descripcion", message.Descripcion),
		slog.String("dia", message.Dia),
	)

	activity := dto.Activity{
		ID:          message.ID,
		Titulo:      message.Nombre,
		Descripcion: message.Descripcion,
		DiaSemana:   message.Dia,
	}

	switch message.Action {
	case "create":
		slog.Info("✅ Activity created", slog.String("activity_id", message.ID))

		// Index in SolR
		if _, err := s.search.Create(ctx, activity); err != nil {
			slog.Error("❌ Error indexing activity in search",
				slog.String("activity_id", message.ID),
				slog.String("error", err.Error()))
			return fmt.Errorf("error indexing activity: %w", err)
		}

		if _, err := s.cache.Create(ctx, activity); err != nil {
			slog.Warn("⚠️ Error caching activity in localcache",
				slog.String("activity_id", message.ID),
				slog.String("error", err.Error()))
		}

		if _, err := s.cache.Create(ctx, activity); err != nil {
			slog.Warn("⚠️ Error caching activity",
				slog.String("activity_id", message.ID),
				slog.String("error", err.Error()))
		}

		slog.Info("🔍 Activity indexed in search engine", slog.String("activity_id", message.ID))
	case "update":
		slog.Info("✏️ Activity updated", slog.String("activity_id", message.ID))

		// Reindex in SolR
		_, err := s.search.Update(ctx, message.ID, activity)
		if err != nil {
			slog.Error("❌ Error reindexing activity in search",
				slog.String("activity_id", message.ID),
				slog.String("error", err.Error()))
			return fmt.Errorf("error reindexing activity: %w", err)
		}

		// Update cache
		if _, err := s.localCache.Update(ctx, message.ID, activity); err != nil {
			slog.Warn("⚠️ Error updating activity in local cache",
				slog.String("activity_id", message.ID),
				slog.String("error", err.Error()))
		}

		if _, err := s.cache.Update(ctx, message.ID, activity); err != nil {
			slog.Warn("⚠️ Error updating activity in cache",
				slog.String("activity_id", message.ID),
				slog.String("error", err.Error()))
		}

		slog.Info("🔍 Activity reindexed in search engine", slog.String("activity_id", message.ID))
	case "delete":
		slog.Info("🗑️ Activity deleted", slog.String("activity_id", message.ID))

		// Delete from SolR
		err := s.search.Delete(ctx, message.ID)
		if err != nil {
			slog.Error("❌ Error deleting activity in search",
				slog.String("activity_id", message.ID),
				slog.String("error", err.Error()))
			return fmt.Errorf("error deleting activity in search: %w", err)
		}

		// Delete from cache
		if err := s.localCache.Delete(ctx, message.ID); err != nil {
			slog.Warn("⚠️ Error deleting activity from local cache",
				slog.String("activity_id", message.ID),
				slog.String("error", err.Error()))
		}

		if err := s.cache.Delete(ctx, message.ID); err != nil {
			slog.Warn("⚠️ Error deleting activity from cache",
				slog.String("activity_id", message.ID),
				slog.String("error", err.Error()))
		}

		slog.Info("🗑️ Activity deleted from search engine", slog.String("activity_id", message.ID))
	default:
		slog.Info("⚠️ Unknown action", slog.String("action", message.Action))
	}

	return nil
}
