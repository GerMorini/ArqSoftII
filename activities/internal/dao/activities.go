package dao

import (
	"strconv"
	"time"

	"activities/internal/domain"

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
	HoraInicio        string             `bson:"hora_inicio"`
	HoraFin           string             `bson:"hora_fin"`
	UsuariosInscritos []int              `bson:"usuarios_inscritos"`
	CapacidadMax      int                `bson:"capacidad_max"`
	Activa            bool               `bson:"activa"`
	FechaCreacion     time.Time          `bson:"fecha_creacion"`
}

// ToDomain convierte el modelo de persistencia a domain.Activity
func (a ActivityDAO) ToDomain() domain.Activity {
	var idHex string
	if a.ID != primitive.NilObjectID {
		idHex = a.ID.Hex()
	}

	return domain.Activity{
		ID:           idHex,
		Nombre:       a.Nombre,
		Descripcion:  a.Descripcion,
		Profesor:     a.Profesor,
		DiaSemana:    a.DiaSemana,
		HoraInicio:   a.HoraInicio,
		HoraFin:      a.HoraFin,
		CapacidadMax: strconv.Itoa(a.CapacidadMax),
	}
}

// ToDomainAdmin convierte el modelo de persistencia a domain.ActivityAdministration
func (a ActivityDAO) ToDomainAdmin() domain.ActivityAdministration {
	var idHex string
	if a.ID != primitive.NilObjectID {
		idHex = a.ID.Hex()
	}
	strs := make([]string, len(a.UsuariosInscritos))
	for i, v := range a.UsuariosInscritos {
		strs[i] = strconv.Itoa(v)
	}

	return domain.ActivityAdministration{
		Activity: domain.Activity{
			ID:           idHex,
			Nombre:       a.Nombre,
			Descripcion:  a.Descripcion,
			Profesor:     a.Profesor,
			DiaSemana:    a.DiaSemana,
			HoraInicio:   a.HoraInicio,
			HoraFin:      a.HoraFin,
			CapacidadMax: strconv.Itoa(a.CapacidadMax),
		},
		UsersInscribed: strs,
		FechaCreacion:  a.FechaCreacion,
	}
}

// FromDomainDAO convierte domain.Activity a ActivityDAO (para persistencia)
func FromDomainDAO(d domain.ActivityAdministration) ActivityDAO {
	var objectID primitive.ObjectID
	if d.ID != "" {
		objectID, _ = primitive.ObjectIDFromHex(d.ID)
	}

	capacidad, err := strconv.Atoi(d.CapacidadMax)
	if err != nil {
		capacidad = 0
	}

	return ActivityDAO{
		ID:                objectID,
		Nombre:            d.Nombre,
		Descripcion:       d.Descripcion,
		Profesor:          d.Profesor,
		DiaSemana:         d.DiaSemana,
		HoraInicio:        d.HoraInicio,
		HoraFin:           d.HoraFin,
		UsuariosInscritos: nil,
		CapacidadMax:      capacidad,
		Activa:            true,
		FechaCreacion:     time.Now().UTC(),
	}
}
