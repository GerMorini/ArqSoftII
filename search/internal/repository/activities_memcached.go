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
	return dto.PaginatedResponse{}, fmt.Errorf("list is not supported in memcached")
}

func (r MemcachedActivitiesRepository) Create(ctx context.Context, activity dto.Activity) (dto.Activity, error) {
	bytes, err := json.Marshal(activity)
	if err != nil {
		return dto.Activity{}, fmt.Errorf("error marshalling activity to JSON: %w", err)
	}
	if err := r.client.Set(&memcache.Item{
		Key:        activity.ID,
		Value:      bytes,
		Expiration: int32(r.ttl.Seconds()),
	}); err != nil {
		return dto.Activity{}, fmt.Errorf("error setting activity in memcached: %w", err)
	}
	return activity, nil
}

func (r MemcachedActivitiesRepository) GetByID(ctx context.Context, id string) (dto.Activity, error) {
	bytes, err := r.client.Get(id)
	if err != nil {
		return dto.Activity{}, fmt.Errorf("error getting activity from memcached: %w", err)
	}
	var activity dto.Activity
	if err := json.Unmarshal(bytes.Value, &activity); err != nil {
		return dto.Activity{}, fmt.Errorf("error unmarshalling activity from JSON: %w", err)
	}
	return activity, nil
}

func (r MemcachedActivitiesRepository) Update(ctx context.Context, id string, activity dto.Activity) (dto.Activity, error) {
	bytes, err := json.Marshal(activity)
	if err != nil {
		return dto.Activity{}, fmt.Errorf("error marshalling activity to JSON: %w", err)
	}
	if err := r.client.Set(&memcache.Item{
		Key:        activity.ID,
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
