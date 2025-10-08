package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"
	"users/internal/dao"
	"users/internal/dto"
	"users/internal/repository"

	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
)

type UsersService interface {
	Login(loginDTO dto.UserLoginDTO) (string, error)
	Create(datos dto.UserMinDTO) (dto.UserMinDTO, error)
	GetByID(id int) (dto.UserDTO, error)

	GenerateToken(userdata dao.User) (string, error)
	GetClaimsFromToken(tokenString string) (jwt.MapClaims, error)
}

var (
	ErrIncorrectCredentials error = errors.New("credenciales incorrectas")
	ErrLoginFormat          error = errors.New("se debe especificar solo uno de los siguientes: username, email")
)

type UsersServiceImpl struct {
	repository repository.UsersRepository

	jwtSecret string
}

func NewUsersService(repository repository.UsersRepository, jwtSecret string) UsersServiceImpl {
	return UsersServiceImpl{
		repository: repository,
		jwtSecret:  jwtSecret,
	}
}

func (s *UsersServiceImpl) Login(loginDTO dto.UserLoginDTO) (string, error) {
	var userdata dao.User
	var err error

	if loginDTO.Email == "" && loginDTO.Username != "" { // si se especifica solo username
		userdata, err = s.repository.GetUserByUsername(loginDTO.Username)
	} else if loginDTO.Email != "" && loginDTO.Username == "" { // si se especifica solo el email
		userdata, err = s.repository.GetUserByEmail(loginDTO.Email)
	} else {
		return "", ErrLoginFormat
	}

	if err != nil {
		return "", err
	}

	if calculateSHA256(loginDTO.Password) != userdata.Password {
		log.Debugf("Contraseña incorrecta para el usuario:\nusername: %s\nemail: %s\npassword: %s\n", loginDTO.Username, loginDTO.Email, loginDTO.Password)
		return "", ErrIncorrectCredentials
	}

	return s.GenerateToken(userdata)
}

func (s *UsersServiceImpl) Create(datos dto.UserMinDTO) (dto.UserMinDTO, error) {
	err := validateUser(datos)
	if err != nil {
		return dto.UserMinDTO{}, err
	}

	var daoUser dao.User = dao.User{
		Nombre:   datos.Nombre,
		Apellido: datos.Apellido,
		Username: datos.Username,
		Email:    datos.Email,
		Password: calculateSHA256(datos.Password),
	}

	_, err = s.repository.Create(daoUser)
	if err != nil {
		return dto.UserMinDTO{}, err
	}

	return datos, err
}

func (s *UsersServiceImpl) GetByID(id int) (dto.UserDTO, error) {
	var userData dao.User

	userData, err := s.repository.GetUserByID(id)
	if err != nil {
		return dto.UserDTO{}, err
	}

	return dto.UserDTO{
		Id:       id,
		Nombre:   userData.Nombre,
		Apellido: userData.Apellido,
		Username: userData.Username,
		Email:    userData.Email,
		IsAdmin:  userData.IsAdmin,
	}, nil
}

func (s *UsersServiceImpl) GenerateToken(userdata dao.User) (string, error) {
	claims := jwt.MapClaims{
		"iss":        "users-api",
		"exp":        time.Now().Add(30 * time.Minute).Unix(),
		"id_usuario": userdata.Id,
		"nombre":     userdata.Nombre,
		"apellido":   userdata.Apellido,
		"username":   userdata.Username,
		"email":      userdata.Email,
		"is_admin":   userdata.IsAdmin,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.jwtSecret))
}

func (s *UsersServiceImpl) GetClaimsFromToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		log.Errorf("error al parsear el token: %v\n", err)
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Errorf("error al obtener los claims\ntokenString: %s\n", tokenString)
		return nil, errors.New("error al obtener los claims")
	}

	return claims, nil
}

func validateUser(user dto.UserMinDTO) error {
	if strings.TrimSpace(user.Nombre) == "" {
		return errors.New("se requiere especificar un nombre")
	}

	if strings.TrimSpace(user.Apellido) == "" {
		return errors.New("se requiere especificar un apellido")
	}

	if strings.TrimSpace(user.Username) == "" {
		return errors.New("se requiere especificar un nombre de usuario")
	}

	if strings.TrimSpace(user.Email) == "" {
		return errors.New("se requiere especificar un email")
	}

	if strings.TrimSpace(user.Password) == "" {
		return errors.New("se requiere especificar una contraseña")
	}

	return nil
}

func calculateSHA256(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}
