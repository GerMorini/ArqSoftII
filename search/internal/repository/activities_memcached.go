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

// FlushAll clears all entries from memcached
func (r MemcachedActivitiesRepository) FlushAll() error {
	return r.client.FlushAll()
}
