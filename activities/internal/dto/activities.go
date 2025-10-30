package dto

import "time"

type Activity struct {
	ID                 string `json:"id_actividad"`
	Nombre             string `json:"titulo"`
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

type ActivityAdministration struct {
	Activity
	UsersInscribed []int `json:"usuarios_inscritos,omitempty"` // Array de User IDs (JSON: usuarios_inscritos)
	FechaCreacion  time.Time
}

type ActivitiesAdministrations []ActivityAdministration
