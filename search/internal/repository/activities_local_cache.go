package repository

import (
	"context"
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

func (r ActivitiesLocalCacheRepository) Create(ctx context.Context, Activity dto.Activity) (dto.Activity, error) {
	r.client.Set(Activity.ID, Activity, r.ttl)
	return Activity, nil
}

func (r ActivitiesLocalCacheRepository) GetByID(ctx context.Context, id string) (dto.Activity, error) {
	it := r.client.Get(id)
	if it == nil {
		return dto.Activity{}, fmt.Errorf("Activity not found in cache")
	}
	Activity, ok := it.Value().(dto.Activity)
	if !ok {
		return dto.Activity{}, fmt.Errorf("error asserting Activity type from cache")
	}
	return Activity, nil
}

func (r ActivitiesLocalCacheRepository) Update(ctx context.Context, id string, Activity dto.Activity) (dto.Activity, error) {
	// TODO implement me
	panic("implement me")
}

func (r ActivitiesLocalCacheRepository) Delete(ctx context.Context, id string) error {
	// TODO implement me
	panic("implement me")
}
