package repository

import (
	"context"
	"errors"
	"fmt"
	"search/internal/dto"
	"time"

	"github.com/karlseguin/ccache"
)

type ActivitiesLocalCacheRepository struct {
	client *ccache.Cache
	ttl    time.Duration
}

func NewActivitysLocalCacheRepository(ttl time.Duration) *ActivitiesLocalCacheRepository {
	return &ActivitiesLocalCacheRepository{
		client: ccache.New(ccache.Configure()),
		ttl:    ttl,
	}
}

func (r ActivitiesLocalCacheRepository) List(ctx context.Context, filters dto.SearchFilters) (dto.PaginatedResponse, error) {
	key := fmt.Sprintf("%s:%s:%s:%d:%d", filters.Titulo, filters.Descripcion, filters.DiaSemana, filters.Page, filters.Count)
	item := r.client.Get(key)
	if item == nil {
		return dto.PaginatedResponse{}, errors.New("cache miss")
	}
	if item.Expired() {
		return dto.PaginatedResponse{}, errors.New("cache expired")
	}
	result, ok := item.Value().(dto.PaginatedResponse)
	if !ok {
		return dto.PaginatedResponse{}, errors.New("invalid cache value type")
	}
	return result, nil
}

func (r ActivitiesLocalCacheRepository) Create(ctx context.Context, activity dto.Activity) (dto.Activity, error) {
	r.client.Set(
		fmt.Sprintf("%s:%s:%s", activity.Titulo, activity.Descripcion, activity.DiaSemana),
		activity,
		r.ttl,
	)
	return activity, nil
}

// SetPaginatedResult stores a paginated response in cache using search filters as key
func (r ActivitiesLocalCacheRepository) SetPaginatedResult(filters dto.SearchFilters, result dto.PaginatedResponse) error {
	key := fmt.Sprintf("%s:%s:%s:%d:%d", filters.Titulo, filters.Descripcion, filters.DiaSemana, filters.Page, filters.Count)
	r.client.Set(key, result, r.ttl)
	return nil
}

func (r ActivitiesLocalCacheRepository) Update(ctx context.Context, id string, activity dto.Activity) (dto.Activity, error) {
	r.client.Set(
		fmt.Sprintf("%s:%s:%s", activity.Titulo, activity.Descripcion, activity.DiaSemana),
		activity,
		r.ttl)
	return activity, nil
}

func (r ActivitiesLocalCacheRepository) Delete(ctx context.Context, id string) error {
	if !r.client.Delete(id) {
		return errors.New("error deleting activity from cache")
	}

	return nil
}
