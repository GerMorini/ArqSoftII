package dto

type Activity struct {
	ID                 string `json:"id_actividad"`
	Titulo             string `json:"titulo"`
	Descripcion        string `json:"descripcion"`
	Profesor           string `json:"instructor"`
	DiaSemana          string `json:"dia"`
	HoraInicio         string `json:"hora_inicio"`
	HoraFin            string `json:"hora_fin"`
	CapacidadMax       int    `json:"cupo"`
	LugaresDisponibles int    `json:"lugares_disponibles"`
	FotoUrl            string `json:"foto_url"`
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
