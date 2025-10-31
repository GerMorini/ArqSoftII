package dto

type Activity struct {
	ID          string `json:"id"`
	Titulo      string `json:"titulo"`
	Descripcion string `json:"descripcion"`
	DiaSemana   string `json:"dia"`
}

type Activities []Activity

type SearchFilters struct {
	ID          string `json:"id"`
	Titulo      string `json:"titulo"`
	Descripcion string `json:"descripcion"`
	DiaSemana   string `json:"dia"`
	SortBy      string `json:"sort_by"`
	Page        int    `json:"page"`
	Count       int    `json:"count"`
}

type PaginatedResponse struct {
	Page    int        `json:"page"`
	Count   int        `json:"count"`
	Total   int        `json:"total"`
	Results Activities `json:"results"`
}
