package dto

type UserLoginDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UsersLoginDTO []UserLoginDTO
