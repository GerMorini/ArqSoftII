package dao

import "time"

type User struct {
	Id            int       `gorm:"column:id_usuario;primaryKey;autoIncrement"`
	Nombre        string    `gorm:"type:varchar(30);not null"`
	Apellido      string    `gorm:"type:varchar(30);not null"`
	Username      string    `gorm:"type:varchar(30);unique;not null"`
	Password      string    `gorm:"type:char(64);collation:ascii_bin;not null"`
	IsAdmin       bool      `gorm:"column:is_admin;default:false;not null"`
	FechaRegistro time.Time `gorm:"column:fecha_registro;type:timestamp;default:CURRENT_TIMESTAMP;not null"`
}

type Users []User
