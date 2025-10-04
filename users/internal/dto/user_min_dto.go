package dto

type UserMinDTO struct {
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UsersMinDTO []UserMinDTO
