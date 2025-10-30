package dao

import (
	"activities/internal/dto"
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
	FotoUrl           string             `bson:"foto_url"`
}

// ToDomain convierte ActivityDAO a Activity (DTO/Domain)
func (dao ActivityDAO) ToDomain() dto.Activity {
	lugaresDisponibles := dao.CapacidadMax - len(dao.UsuariosInscritos)
	return dto.Activity{
		ID:                 dao.ID.Hex(),
		Nombre:             dao.Nombre,
		Descripcion:        dao.Descripcion,
		Profesor:           dao.Profesor,
		DiaSemana:          dao.DiaSemana,
		HoraInicio:         dao.HoraInicio,
		HoraFin:            dao.HoraFin,
		FotoUrl:            dao.FotoUrl,
		CapacidadMax:       dao.CapacidadMax,
		LugaresDisponibles: lugaresDisponibles,
	}
}

func FromDomainDAO(a dto.ActivityAdministration) ActivityDAO {
	return ActivityDAO{
		// ID se asigna automáticamente en Create si es vacío
		Nombre:            a.Nombre,
		Descripcion:       a.Descripcion,
		Profesor:          a.Profesor,
		DiaSemana:         a.DiaSemana,
		HoraInicio:        a.HoraInicio,
		HoraFin:           a.HoraFin,
		UsuariosInscritos: a.UsersInscribed,
		CapacidadMax:      a.CapacidadMax,
		FotoUrl:           a.FotoUrl,
		Activa:            true, // Por defecto al crear es activa
		FechaCreacion:     time.Now().UTC(),
	}
}

func ToDomainAdministration(dao ActivityDAO) dto.ActivityAdministration {
	lugaresDisponibles := dao.CapacidadMax - len(dao.UsuariosInscritos)
	return dto.ActivityAdministration{
		Activity: dto.Activity{
			ID:                 dao.ID.Hex(),
			Nombre:             dao.Nombre,
			Descripcion:        dao.Descripcion,
			Profesor:           dao.Profesor,
			DiaSemana:          dao.DiaSemana,
			HoraInicio:         dao.HoraInicio,
			HoraFin:            dao.HoraFin,
			FotoUrl:            dao.FotoUrl,
			CapacidadMax:       dao.CapacidadMax,
			LugaresDisponibles: lugaresDisponibles,
		},
		UsersInscribed: dao.UsuariosInscritos,
		FechaCreacion:  dao.FechaCreacion,
	}
}
