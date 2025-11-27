package services

import (
	"errors"
	"strings"
	"time"
	"users/internal/dao"
	"users/internal/dto"
	"users/internal/repository"

	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UsersService interface {
	Login(loginDTO dto.UserLoginDTO) (string, error)
	Create(datos dto.UserMinDTO) (dto.UserMinDTO, error)
	GetByID(id int) (dto.UserDTO, error)
	GetAll() ([]dto.UserDTO, error)
	Update(id int, updateDTO dto.UserUpdateDTO) (dto.UserDTO, error)
	Delete(id int) error
	IsAdmin(token string) (bool, error)
	GenerateToken(userdata dao.User) (string, error)
	GetClaimsFromToken(tokenString string) (jwt.MapClaims, error)
}

var (
	ErrIncorrectCredentials error = errors.New("credenciales incorrectas")
	ErrLoginFormat          error = errors.New("se debe especificar solo uno de los siguientes: username, email")
	ErrInvalidTokenClaims   error = errors.New("error al obtener los claims")
)

type UsersServiceImpl struct {
	repository repository.UsersRepository

	jwtSecret []byte
}

func NewUsersService(repository repository.UsersRepository, jwtSecret string) UsersServiceImpl {
	return UsersServiceImpl{
		repository: repository,
		jwtSecret:  []byte(jwtSecret),
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

	// Verify password with bcrypt
	if !s.verifyPassword(userdata.Password, loginDTO.Password) {
		return "", ErrIncorrectCredentials
	}

	return s.GenerateToken(userdata)
}

func (s *UsersServiceImpl) Create(datos dto.UserMinDTO) (dto.UserMinDTO, error) {
	err := validateUser(datos)
	if err != nil {
		return dto.UserMinDTO{}, err
	}

	// Hash password with bcrypt
	hashedPassword, err := s.hashPassword(datos.Password)
	if err != nil {
		log.Errorf("error hashing password: %v", err)
		return dto.UserMinDTO{}, err
	}

	var daoUser dao.User = dao.User{
		Nombre:   datos.Nombre,
		Apellido: datos.Apellido,
		Username: datos.Username,
		Email:    datos.Email,
		Password: hashedPassword,
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

	return token.SignedString(s.jwtSecret)
}

func (s *UsersServiceImpl) GetClaimsFromToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		log.Errorf("error al parsear el token: %v\n", err)
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Errorf("error al obtener los claims\ntokenString: %s\n", tokenString)
		return nil, ErrInvalidTokenClaims
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

	if user.Password == "" {
		return errors.New("se requiere especificar una contraseña")
	}

	return nil
}

// hashPassword generates a bcrypt hash with automatic salt
func (s *UsersServiceImpl) hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost, // Cost = 10 (2^10 iterations)
	)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// verifyPassword compares plaintext password with bcrypt hash
func (s *UsersServiceImpl) verifyPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(plainPassword),
	)
	return err == nil
}

func (s *UsersServiceImpl) GetAll() ([]dto.UserDTO, error) {
	usuarios, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}

	var result []dto.UserDTO
	for _, u := range usuarios {
		result = append(result, dto.UserDTO{
			Id:       u.Id,
			Nombre:   u.Nombre,
			Apellido: u.Apellido,
			Username: u.Username,
			Email:    u.Email,
			IsAdmin:  u.IsAdmin,
		})
	}

	return result, nil
}

func (s *UsersServiceImpl) Update(id int, updateDTO dto.UserUpdateDTO) (dto.UserDTO, error) {
	usuarioActual, err := s.repository.GetUserByID(id)
	if err != nil {
		return dto.UserDTO{}, err
	}

	usuarioActual.Nombre = updateDTO.Nombre
	usuarioActual.Apellido = updateDTO.Apellido
	usuarioActual.Email = updateDTO.Email
	usuarioActual.IsAdmin = updateDTO.IsAdmin

	// Si se proporciona una nueva contraseña, hashearla y actualizar
	if updateDTO.Password != "" {
		hashedPassword, err := s.hashPassword(updateDTO.Password)
		if err != nil {
			log.Errorf("error hashing password: %v", err)
			return dto.UserDTO{}, err
		}
		usuarioActual.Password = hashedPassword
	}

	usuarioActualizado, err := s.repository.Update(id, usuarioActual)
	if err != nil {
		return dto.UserDTO{}, err
	}

	return dto.UserDTO{
		Id:       usuarioActualizado.Id,
		Nombre:   usuarioActualizado.Nombre,
		Apellido: usuarioActualizado.Apellido,
		Username: usuarioActualizado.Username,
		Email:    usuarioActualizado.Email,
		IsAdmin:  usuarioActualizado.IsAdmin,
	}, nil
}

func (s *UsersServiceImpl) Delete(id int) error {
	return s.repository.Delete(id)
}

func (s *UsersServiceImpl) IsAdmin(token string) (bool, error) {
	claims, err := s.GetClaimsFromToken(token)
	if err != nil {
		log.Warnf("error al obtener claims del token para validación de admin")
		return false, err
	}

	isAdmin, ok := claims["is_admin"].(bool)
	if !ok {
		log.Warnf("error al obtener is_admin de los claims")
		return false, err
	}

	return isAdmin, nil
}
