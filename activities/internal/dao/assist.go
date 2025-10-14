package dao

import "time"

type Asistencia struct {
	ID            string    // _id (ObjectID)
	ClaseID       string    // FK a Clase (MongoDB)
	UserID        uint      // FK a User (PostgreSQL)
	Fecha         time.Time // Fecha específica de la clase
	Asistio       bool
	Observaciones string // opcional
}
