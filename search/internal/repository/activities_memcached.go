package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"search/internal/dto"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcachedActivitiesRepository struct {
	ttl    time.Duration
	client *memcache.Client
}

func NewMemcachedActivitiesRepository(host string, port string, ttl time.Duration) MemcachedActivitiesRepository {
	client := memcache.New(fmt.Sprintf("%s:%s", host, port))

	return MemcachedActivitiesRepository{
		client: client,
		ttl:    ttl,
	}
}

func (r MemcachedActivitiesRepository) List(ctx context.Context, filters dto.SearchFilters) (dto.PaginatedResponse, error) {
	key := fmt.Sprintf("%s:%s:%s:%d:%d", filters.Titulo, filters.Descripcion, filters.DiaSemana, filters.Page, filters.Count)
	item, err := r.client.Get(key)
	if err != nil {
		return dto.PaginatedResponse{}, fmt.Errorf("cache miss: %w", err)
	}
	var result dto.PaginatedResponse
	if err := json.Unmarshal(item.Value, &result); err != nil {
		return dto.PaginatedResponse{}, fmt.Errorf("error unmarshalling cached data: %w", err)
	}
	return result, nil
}

func (r MemcachedActivitiesRepository) Create(ctx context.Context, activity dto.Activity) (dto.Activity, error) {
	bytes, err := json.Marshal(activity)
	if err != nil {
		return dto.Activity{}, fmt.Errorf("error marshalling activity to JSON: %w", err)
	}
	if err := r.client.Set(&memcache.Item{
		Key:        fmt.Sprintf("%s:%s:%s", activity.Titulo, activity.Descripcion, activity.DiaSemana),
		Value:      bytes,
		Expiration: int32(r.ttl.Seconds()),
	}); err != nil {
		return dto.Activity{}, fmt.Errorf("error setting activity in memcached: %w", err)
	}
	return activity, nil
}

// SetPaginatedResult stores a paginated response in cache using search filters as key
func (r MemcachedActivitiesRepository) SetPaginatedResult(filters dto.SearchFilters, result dto.PaginatedResponse) error {
	key := fmt.Sprintf("%s:%s:%s:%d:%d", filters.Titulo, filters.Descripcion, filters.DiaSemana, filters.Page, filters.Count)
	bytes, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("error marshalling paginated response to JSON: %w", err)
	}
	if err := r.client.Set(&memcache.Item{
		Key:        key,
		Value:      bytes,
		Expiration: int32(r.ttl.Seconds()),
	}); err != nil {
		return fmt.Errorf("error setting paginated response in memcached: %w", err)
	}
	return nil
}

func (r MemcachedActivitiesRepository) Update(ctx context.Context, id string, activity dto.Activity) (dto.Activity, error) {
	bytes, err := json.Marshal(activity)
	if err != nil {
		return dto.Activity{}, fmt.Errorf("error marshalling activity to JSON: %w", err)
	}
	if err := r.client.Set(&memcache.Item{
		Key:        fmt.Sprintf("%s:%s:%s", activity.Titulo, activity.Descripcion, activity.DiaSemana),
		Value:      bytes,
		Expiration: int32(r.ttl.Seconds()),
	}); err != nil {
		return dto.Activity{}, fmt.Errorf("error setting activity in memcached: %w", err)
	}
	return activity, nil
}

func (r MemcachedActivitiesRepository) Delete(ctx context.Context, id string) error {
	return r.client.Delete(id)
}
