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

type DayDistribution struct {
	Dia   string `json:"dia"`
	Count int    `json:"count"`
}

type ActivityStatistics struct {
	TotalActivities       int                `json:"total_actividades"`
	TotalEnrollments      int                `json:"total_inscripciones"`
	AverageEnrollmentRate float64            `json:"tasa_promedio_inscripcion"`
	TotalCapacity         int                `json:"capacidad_total"`
	CapacityUtilization   float64            `json:"utilizacion_capacidad"`
	ActivitiesByDay       []DayDistribution  `json:"actividades_por_dia"`
	MostPopularActivity   *Activity          `json:"actividad_mas_popular"`
	FullActivitiesCount   int                `json:"actividades_llenas"`
	AvailableActivities   int                `json:"actividades_disponibles"`
}
