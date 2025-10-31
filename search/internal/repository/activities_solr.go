package repository

import (
	"context"
	"fmt"
	"search/internal/clients"
	"search/internal/dto"
	"strings"
)

type SolrClient interface {
}

type SolrActivitysRepository struct {
	client *clients.SolrClient
}

func NewSolrActivitysRepository(host, port, core string) *SolrActivitysRepository {
	client := clients.NewSolrClient(host, port, core)
	return &SolrActivitysRepository{
		client: client,
	}
}

func (r *SolrActivitysRepository) List(ctx context.Context, filters dto.SearchFilters) (dto.PaginatedResponse, error) {
	query := buildQuery(filters)
	return r.client.Search(ctx, query, filters.Page, filters.Count)
}

func (r *SolrActivitysRepository) GetByID(ctx context.Context, id string) (dto.Activity, error) {
	results, err := r.List(ctx, dto.SearchFilters{ID: id})
	if err != nil {
		return dto.Activity{}, fmt.Errorf("error searching activity by ID in solr: %w", err)
	}
	if results.Total == 0 {
		return dto.Activity{}, fmt.Errorf("activity with ID %s not found", id)
	}
	return results.Results[0], nil
}

func (r *SolrActivitysRepository) Create(ctx context.Context, activity dto.Activity) (dto.Activity, error) {
	if err := r.client.Index(ctx, activity); err != nil {
		return dto.Activity{}, fmt.Errorf("error indexing activity in solr: %w", err)
	}
	return activity, nil
}

func (r *SolrActivitysRepository) Update(ctx context.Context, id string, activity dto.Activity) (dto.Activity, error) {
	// En Solr, actualizar es equivalente a re-indexar con el mismo ID
	activity.ID = id
	if err := r.client.Index(ctx, activity); err != nil {
		return dto.Activity{}, fmt.Errorf("error updating activity in solr: %w", err)
	}
	return activity, nil
}

func (r *SolrActivitysRepository) Delete(ctx context.Context, id string) error {
	if err := r.client.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting activity from solr: %w", err)
	}
	return nil
}

func buildQuery(filters dto.SearchFilters) string {
	var parts []string

	// Si no hay filtros, devolvemos todo
	if filters.ID == "" || filters.Titulo == "" && filters.Descripcion == "" && filters.DiaSemana == "" {
		return "*:*"
	}

	if filters.ID != "" {
		parts = append(parts, fmt.Sprintf("id:%s", filters.ID))
	}

	if filters.Titulo != "" {
		parts = append(parts, fmt.Sprintf("name:*%s*", filters.Titulo))
	}

	if filters.Descripcion != "" {
		parts = append(parts, fmt.Sprintf("name:*%s*", filters.Descripcion))
	}

	if filters.DiaSemana != "" {
		parts = append(parts, fmt.Sprintf("name:*%s*", filters.DiaSemana))
	}

	return strings.Join(parts, " AND ")
}
