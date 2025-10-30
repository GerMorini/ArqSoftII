package dao

type User struct {
	Id       int    `gorm:"column:id_usuario;primaryKey;autoIncrement"`
	Nombre   string `gorm:"type:varchar(30);not null"`
	Apellido string `gorm:"type:varchar(30);not null"`
	Username string `gorm:"type:varchar(30);unique;not null"`
	Email    string `gorm:"type:varchar(60);unique;not null"`
	Password string `gorm:"type:varchar(60);collation:utf8mb4_bin;not null"`
	IsAdmin  bool   `gorm:"column:is_admin;default:false;not null"`
}

type Users []User
