package dao

import "time"

type Activity struct {
	ID                string // _id (ObjectID)
	Nombre            string // ej: "Yoga Principiantes"
	Descripcion       string
	ProfesorID        uint   // FK a Profesor (MySQL)
	DiaSemana         string // "Lunes", "Martes", etc
	HoraInicio        string // "09:00"
	HoraFin           string // "10:30"
	UsuariosInscritos []int  // Array de User IDs
	CapacidadMax      int
	Activa            bool
	FechaCreacion     time.Time
}
