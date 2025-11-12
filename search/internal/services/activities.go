package services

import (
	"context"
	"fmt"
	"log/slog"
	"search/internal/dto"

	log "github.com/sirupsen/logrus"
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

type ActivitiesCacheRepository interface {
	List(ctx context.Context, filters dto.SearchFilters) (dto.PaginatedResponse, error)
	SetPaginatedResult(filters dto.SearchFilters, result dto.PaginatedResponse) error
	FlushAll() error
}

type ActivitiesConsumer interface {
	Consume(ctx context.Context, handler func(ctx context.Context, message ActivityEvent) error) error
}

type ActiviesServiceImpl struct {
	localCache ActivitiesCacheRepository
	memCached  ActivitiesCacheRepository
	search     ActivitiesRepository
	consumer   ActivitiesConsumer
}

func NewActivitiesService(localCache ActivitiesCacheRepository, cache ActivitiesCacheRepository, search ActivitiesRepository, consumer ActivitiesConsumer) ActiviesServiceImpl {
	return ActiviesServiceImpl{
		localCache: localCache,
		memCached:  cache,
		search:     search,
		consumer:   consumer,
	}
}

func (s *ActiviesServiceImpl) List(ctx context.Context, filters dto.SearchFilters) (dto.PaginatedResponse, error) {
	var localCacheMiss, memcacheMiss bool = false, false

	result, err := s.localCache.List(ctx, filters)
	if err == nil {
		log.Infof("cache hit en localcache: %v", filters)
		return result, nil
	}
	log.Warnf("no se encontro actividad en cache local")
	localCacheMiss = true

	result, err = s.memCached.List(ctx, filters)
	if err == nil {
		log.Infof("cache hit en memcached: %v", filters)
		return result, nil
	}
	log.Warnf("no se encontro actividad en memcached")
	memcacheMiss = true

	result, err = s.search.List(ctx, filters)
	if err == nil {
		log.Infof("actividad buscada exitosamente en solr")

		// Cache the entire paginated response using the filters as the key
		if localCacheMiss && result.Total != 0 {
			if err := s.localCache.SetPaginatedResult(filters, result); err != nil {
				log.Errorf("error cacheando resultado en cache local: %s", err.Error())
			} else {
				log.Infof("resultado cacheado exitosamente en cache local")
			}
		}

		if memcacheMiss && result.Total != 0 {
			if err := s.memCached.SetPaginatedResult(filters, result); err != nil {
				log.Errorf("error cacheando resultado en memcached: %s", err.Error())
			} else {
				log.Infof("resultado cacheado exitosamente en memcached")
			}
		}

		return result, err
	}

	log.Infof("no se encontro actividad en solr")
	return dto.PaginatedResponse{}, err
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

		// Invalidate all cache to ensure consistency
		if err := s.localCache.FlushAll(); err != nil {
			slog.Warn("⚠️ Error flushing local cache",
				slog.String("error", err.Error()))
		}

		if err := s.memCached.FlushAll(); err != nil {
			slog.Warn("⚠️ Error flushing memcached",
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

		// Invalidate all cache to ensure consistency
		if err := s.localCache.FlushAll(); err != nil {
			slog.Warn("⚠️ Error flushing local cache",
				slog.String("error", err.Error()))
		}

		if err := s.memCached.FlushAll(); err != nil {
			slog.Warn("⚠️ Error flushing memcached",
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

		// Invalidate all cache to ensure consistency
		if err := s.localCache.FlushAll(); err != nil {
			slog.Warn("⚠️ Error flushing local cache",
				slog.String("error", err.Error()))
		}

		if err := s.memCached.FlushAll(); err != nil {
			slog.Warn("⚠️ Error flushing memcached",
				slog.String("error", err.Error()))
		}

		slog.Info("🗑️ Activity deleted from search engine", slog.String("activity_id", message.ID))
	default:
		slog.Info("⚠️ Unknown action", slog.String("action", message.Action))
	}

	return nil
}
