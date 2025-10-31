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
	return dto.PaginatedResponse{}, fmt.Errorf("list is not supported in memcached")
}

func (r ActivitiesLocalCacheRepository) Create(ctx context.Context, activity dto.Activity) (dto.Activity, error) {
	r.client.Set(activity.ID, activity, r.ttl)
	return activity, nil
}

func (r ActivitiesLocalCacheRepository) Update(ctx context.Context, id string, activity dto.Activity) (dto.Activity, error) {
	r.client.Set(activity.ID, activity, r.ttl)
	return activity, nil
}

func (r ActivitiesLocalCacheRepository) Delete(ctx context.Context, id string) error {
	if !r.client.Delete(id) {
		return errors.New("error deleting activity from cache")
	}

	return nil
}
