package dto

import "time"

type Activity struct {
	ID           string // _id (ObjectID)
	Nombre       string // ej: "Yoga Principiantes"
	Descripcion  string
	Profesor     string // FK a Profesor (MySQL)
	DiaSemana    string // "Lunes", "Martes", etc
	HoraInicio   string // "09:00"
	HoraFin      string // "10:30"
	CapacidadMax string
}

type Activities []Activity

type ActivityAdministration struct {
	Activity
	UsersInscribed []string // Array de User IDs
	FechaCreacion  time.Time
}

type ActivitiesAdministrations []ActivityAdministration
