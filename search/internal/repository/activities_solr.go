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

	// Si hay ID, retornar búsqueda exacta por ID
	if filters.ID != "" {
		return fmt.Sprintf("id:%s", filters.ID)
	}

	// Construir query con filtros disponibles
	if filters.Titulo != "" {
		parts = append(parts, fmt.Sprintf("titulo:*%s*", filters.Titulo))
	}

	if filters.Descripcion != "" {
		parts = append(parts, fmt.Sprintf("descripcion:*%s*", filters.Descripcion))
	}

	if filters.DiaSemana != "" {
		parts = append(parts, fmt.Sprintf("dia:*%s*", filters.DiaSemana))
	}

	// Si no hay ningún filtro, devolver todos los documentos
	if len(parts) == 0 {
		return "*:*"
	}

	return strings.Join(parts, " AND ")
}
