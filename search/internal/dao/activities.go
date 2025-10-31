package dao

import (
	"search/internal/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ActivityDAO struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Nombre      string             `bson:"nombre"`
	Descripcion string             `bson:"descripcion"`
	DiaSemana   string             `bson:"dia_semana"`
}

func (dao ActivityDAO) ToDomain() dto.Activity {
	return dto.Activity{
		ID:          dao.ID.Hex(),
		Titulo:      dao.Nombre,
		Descripcion: dao.Descripcion,
		DiaSemana:   dao.DiaSemana,
	}
}

func FromDomainDAO(a dto.Activity) ActivityDAO {
	return ActivityDAO{
		Nombre:      a.Titulo,
		Descripcion: a.Descripcion,
		DiaSemana:   a.DiaSemana,
	}
}
