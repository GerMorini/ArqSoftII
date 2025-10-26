package dto

type UserDTO struct {
	Id       int    `json:"id_usuario"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
	Email    string `json:"email"`
}

type UsersDTO []UserDTO

type UserUpdateDTO struct {
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	IsAdmin  bool   `json:"is_admin"`
}
