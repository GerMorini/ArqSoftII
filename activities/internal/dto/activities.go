package dto

import "time"

type Activity struct {
	ID           string `json:"id_actividad"`
	Nombre       string `json:"titulo"`
	Descripcion  string `json:"descripcion"`
	Profesor     string `json:"instructor"`
	DiaSemana    string `json:"dia"`
	HoraInicio   string `json:"hora_inicio"`
	HoraFin      string `json:"hora_fin"`
	CapacidadMax string `json:"cupo"`
	FotoUrl      string `json:"foto_url"`
}

type Activities []Activity

type ActivityAdministration struct {
	Activity
	UsersInscribed []string // Array de User IDs
	FechaCreacion  time.Time
}

type ActivitiesAdministrations []ActivityAdministration
