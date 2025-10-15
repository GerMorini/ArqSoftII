package dao

import (
	"activities/internal/dto"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ActivityDAO es el modelo usado solamente para la capa de persistencia (MongoDB)
// Tiene etiquetas `bson` y usa primitive.ObjectID para el campo ID.
type ActivityDAO struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	Nombre            string             `bson:"nombre"`
	Descripcion       string             `bson:"descripcion"`
	Profesor          string             `bson:"profesor_id"`
	DiaSemana         string             `bson:"dia_semana"`
	HoraInicio        string             `bson:"hora_inicio"` // capaz cambiar a time.Time
	HoraFin           string             `bson:"hora_fin"`    // capaz cambiar a time.Time
	UsuariosInscritos []int              `bson:"usuarios_inscritos"`
	CapacidadMax      int                `bson:"capacidad_max"`
	Activa            bool               `bson:"activa"`
	FechaCreacion     time.Time          `bson:"fecha_creacion"`
}

// ToDomain convierte ActivityDAO a Activity (DTO/Domain)
func (dao ActivityDAO) ToDomain() dto.Activity {
	return dto.Activity{
		ID:           dao.ID.Hex(),
		Nombre:       dao.Nombre,
		Descripcion:  dao.Descripcion,
		Profesor:     dao.Profesor,
		DiaSemana:    dao.DiaSemana,
		HoraInicio:   dao.HoraInicio,
		HoraFin:      dao.HoraFin,
		CapacidadMax: fmt.Sprintf("%d", dao.CapacidadMax),
	}
}

func FromDomainDAO(a dto.ActivityAdministration) ActivityDAO {
	// Convertir []string de User IDs a []int
	userIDs := make([]int, len(a.UsersInscribed))
	for i, idStr := range a.UsersInscribed {
		var id int
		fmt.Sscanf(idStr, "%d", &id)
		userIDs[i] = id
	}
	capMax := 0
	fmt.Sscanf(a.CapacidadMax, "%d", &capMax)
	return ActivityDAO{
		// ID se asigna automáticamente en Create si es vacío
		Nombre:            a.Nombre,
		Descripcion:       a.Descripcion,
		Profesor:          a.Profesor,
		DiaSemana:         a.DiaSemana,
		HoraInicio:        a.HoraInicio,
		HoraFin:           a.HoraFin,
		UsuariosInscritos: userIDs,
		CapacidadMax:      capMax,
		Activa:            true, // Por defecto al crear es activa
		FechaCreacion:     time.Now().UTC(),
	}
}

func ToDomainAdministration(dao ActivityDAO) dto.ActivityAdministration {
	// Convertir []int de User IDs a []string
	userIDs := make([]string, len(dao.UsuariosInscritos))
	for i, id := range dao.UsuariosInscritos {
		userIDs[i] = fmt.Sprintf("%d", id)
	}
	return dto.ActivityAdministration{
		Activity: dto.Activity{
			ID:           dao.ID.Hex(),
			Nombre:       dao.Nombre,
			Descripcion:  dao.Descripcion,
			Profesor:     dao.Profesor,
			DiaSemana:    dao.DiaSemana,
			HoraInicio:   dao.HoraInicio,
			HoraFin:      dao.HoraFin,
			CapacidadMax: fmt.Sprintf("%d", dao.CapacidadMax),
		},
		UsersInscribed: userIDs,
		FechaCreacion:  dao.FechaCreacion,
	}
}
