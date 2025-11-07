package services_test

import (
	"errors"
	"testing"

	// Importa los paquetes del proyecto, ajusta la ruta si es necesario
	"users/internal/dao"
	"users/internal/dto"
	"users/internal/repository"
	"users/internal/services"

	"golang.org/x/crypto/bcrypt"
)

// ====================================================================
// Mock del Repositorio de Usuarios
// ====================================================================

// MockUsersRepository es una implementación de prueba de repository.UsersRepository
type MockUsersRepository struct {
	// Campos para configurar el comportamiento del mock en cada test
	GetUserByIDFn       func(id int) (dao.User, error)
	GetUserByEmailFn    func(email string) (dao.User, error)
	GetUserByUsernameFn func(username string) (dao.User, error)
	CreateFn            func(user dao.User) (dao.User, error)
	DeleteFn            func(id int) error
}

// Asegurar que MockUsersRepository implementa la interfaz
var _ repository.UsersRepository = &MockUsersRepository{}

func (m *MockUsersRepository) GetUserByID(id int) (dao.User, error) {
	return m.GetUserByIDFn(id)
}
func (m *MockUsersRepository) GetUserByEmail(email string) (dao.User, error) {
	return m.GetUserByEmailFn(email)
}
func (m *MockUsersRepository) GetUserByUsername(username string) (dao.User, error) {
	return m.GetUserByUsernameFn(username)
}
func (m *MockUsersRepository) Create(user dao.User) (dao.User, error) {
	return m.CreateFn(user)
}
func (m *MockUsersRepository) Update(id int, user dao.User) (dao.User, error) {
	// Mock básico
	return dao.User{}, nil
}
func (m *MockUsersRepository) Delete(id int) error {
	return m.DeleteFn(id)
}
func (m *MockUsersRepository) GetAll() ([]dao.User, error) {
	// Mock básico
	return nil, nil
}

// ====================================================================
// Funciones de utilidad para los tests
// ====================================================================

// hashPassword es una copia simplificada de la lógica de hash en el servicio
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ====================================================================
// Implementación de Tests
// ====================================================================

// Test_GetByID realiza pruebas para el caso de uso GetByID
func Test_GetByID(t *testing.T) {
	// Caso de prueba: Usuario de éxito
	mockUser := dao.User{
		Id: 1, Nombre: "Juan", Apellido: "Perez", Username: "jperez", Email: "j@p.com", IsAdmin: false,
	}

	// Caso de prueba: Usuario no encontrado
	var notFoundUser dao.User // Valor cero para dao.User
	notFoundErr := errors.New("user not found")

	tests := []struct {
		name     string
		userID   int
		mockRepo *MockUsersRepository
		wantDTO  dto.UserDTO
		wantErr  error
	}{
		{
			name:   "Success: User Found",
			userID: 1,
			mockRepo: &MockUsersRepository{
				GetUserByIDFn: func(id int) (dao.User, error) {
					if id == 1 {
						return mockUser, nil
					}
					return dao.User{}, errors.New("unexpected call")
				},
			},
			wantDTO: dto.UserDTO{
				Id: 1, Nombre: "Juan", Apellido: "Perez", Username: "jperez", Email: "j@p.com", IsAdmin: false,
			},
			wantErr: nil,
		},
		{
			name:   "Error: User Not Found",
			userID: 99,
			mockRepo: &MockUsersRepository{
				GetUserByIDFn: func(id int) (dao.User, error) {
					return notFoundUser, notFoundErr
				},
			},
			wantDTO: dto.UserDTO{},
			wantErr: notFoundErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup del servicio con el mock
			s := services.NewUsersService(tt.mockRepo, "testSecret")

			got, err := s.GetByID(tt.userID)

			// 1. Verificación de error
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 2. Verificación de DTO
			if got != tt.wantDTO {
				t.Errorf("GetByID() got = %+v, want %+v", got, tt.wantDTO)
			}
		})
	}
}

// Test_Create realiza pruebas para el caso de uso Create (registro de usuario)
func Test_Create(t *testing.T) {
	// Password que se usará para el DTO de entrada
	rawPassword := "SecureP@ss123"
	// Hash simulado que el servicio generaría
	hashedPassword, _ := hashPassword(rawPassword)

	// DAO de usuario que el repositorio devolvería después de crearlo
	createdUserDAO := dao.User{
		Id: 5, Nombre: "Ana", Apellido: "Gomez", Username: "agomez", Email: "a@g.com", Password: hashedPassword,
	}

	tests := []struct {
		name     string
		inputDTO dto.UserMinDTO
		mockRepo *MockUsersRepository
		wantDTO  dto.UserMinDTO
		wantErr  error
	}{
		{
			name: "Success: User Created",
			inputDTO: dto.UserMinDTO{
				Nombre: "Ana", Apellido: "Gomez", Username: "agomez", Email: "a@g.com", Password: rawPassword,
			},
			mockRepo: &MockUsersRepository{
				// Simula que no existe un usuario con ese email/username
				GetUserByEmailFn: func(email string) (dao.User, error) { return dao.User{}, errors.New("not found") },
				// Simula la creación exitosa
				CreateFn: func(user dao.User) (dao.User, error) { return createdUserDAO, nil },
			},
			wantDTO: dto.UserMinDTO{ // El servicio devuelve un UserMinDTO sin la contraseña
				Nombre: "Ana", Apellido: "Gomez", Username: "agomez", Email: "a@g.com",
			},
			wantErr: nil,
		},
		{
			name: "Error: Duplicate Email",
			inputDTO: dto.UserMinDTO{
				Nombre: "Ana", Apellido: "Gomez", Username: "agomez", Email: "a@g.com", Password: rawPassword,
			},
			mockRepo: &MockUsersRepository{
				// Simula que ya existe un usuario con ese email
				GetUserByEmailFn: func(email string) (dao.User, error) { return createdUserDAO, nil }, // Devuelve el usuario existente
				CreateFn:         func(user dao.User) (dao.User, error) { return dao.User{}, errors.New("should not be called") },
			},
			wantDTO: dto.UserMinDTO{},
			wantErr: errors.New("El email ya se encuentra registrado"), // Error asumido de negocio
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := services.NewUsersService(tt.mockRepo, "testSecret")
			got, err := s.Create(tt.inputDTO)

			// 1. Verificación de error
			if (err != nil) && (tt.wantErr != nil) && (err.Error() != tt.wantErr.Error()) {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 2. Verificación de DTO (solo en caso de éxito)
			if tt.wantErr == nil && got != tt.wantDTO {
				t.Errorf("Create() got = %+v, want %+v", got, tt.wantDTO)
			}
		})
	}
}

// Test_Login realiza pruebas para el caso de uso Login
func Test_Login(t *testing.T) {
	// Contraseña de prueba
	rawPassword := "password123"
	// Contraseña hasheada (simulando lo que estaría en la BD/DAO)
	hashedPassword, _ := hashPassword(rawPassword)

	// Usuario que el repositorio devolvería
	mockUser := dao.User{
		Id: 10, Username: "testuser", Email: "test@example.com", Password: hashedPassword, IsAdmin: false,
	}

	// Un hash para una contraseña incorrecta
	incorrectHashedPassword, _ := hashPassword("wrongpassword")
	mockUserWithWrongPassword := dao.User{
		Id: 11, Username: "baduser", Email: "bad@example.com", Password: incorrectHashedPassword, IsAdmin: false,
	}

	// Errores de negocio definidos en el servicio
	ErrIncorrectCredentials := services.ErrIncorrectCredentials // Importado desde el paquete services
	ErrLoginFormat := services.ErrLoginFormat                   // Importado desde el paquete services

	tests := []struct {
		name     string
		inputDTO dto.UserLoginDTO
		mockRepo *MockUsersRepository
		wantErr  error
	}{
		{
			name: "Success: Login with Username",
			inputDTO: dto.UserLoginDTO{
				Username: mockUser.Username,
				Password: rawPassword,
			},
			mockRepo: &MockUsersRepository{
				// Simula el éxito al obtener el usuario
				GetUserByUsernameFn: func(username string) (dao.User, error) { return mockUser, nil },
				GetUserByEmailFn:    func(email string) (dao.User, error) { return dao.User{}, errors.New("should not be called") },
			},
			wantErr: nil, // El servicio devuelve un token string en caso de éxito
		},
		{
			name: "Error: Login with both Username and Email",
			inputDTO: dto.UserLoginDTO{
				Username: mockUser.Username,
				Email:    mockUser.Email,
				Password: rawPassword,
			},
			mockRepo: &MockUsersRepository{
				GetUserByUsernameFn: func(username string) (dao.User, error) { return dao.User{}, errors.New("should not be called") },
				GetUserByEmailFn:    func(email string) (dao.User, error) { return dao.User{}, errors.New("should not be called") },
			},
			wantErr: ErrLoginFormat,
		},
		{
			name: "Error: User not found",
			inputDTO: dto.UserLoginDTO{
				Username: "nonexistent",
				Password: rawPassword,
			},
			mockRepo: &MockUsersRepository{
				GetUserByUsernameFn: func(username string) (dao.User, error) { return dao.User{}, errors.New("not found") },
			},
			wantErr: ErrIncorrectCredentials,
		},
		{
			name: "Error: Incorrect Password",
			inputDTO: dto.UserLoginDTO{
				Username: mockUser.Username,
				Password: "a_different_password",
			},
			mockRepo: &MockUsersRepository{
				// Simula que el usuario fue encontrado
				GetUserByUsernameFn: func(username string) (dao.User, error) { return mockUser, nil },
			},
			wantErr: ErrIncorrectCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := services.NewUsersService(tt.mockRepo, "testSecret")

			got, err := s.Login(tt.inputDTO)

			// 1. Verificación de error
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 2. Verificación de resultado (token generado)
			if tt.wantErr == nil {
				if got == "" {
					t.Errorf("Login() got empty token, want valid token")
				}

				// CORRECCIÓN: Usar la instancia 's' del servicio para llamar a GetClaimsFromToken.
				// GetClaimsFromToken está definido en la interfaz de servicio.
				claims, err := s.GetClaimsFromToken(got)
				if err != nil {
					t.Errorf("Login() generated invalid token: %v", err)
				}
				// El ID del token se serializa a float64 por el paquete JWT.
				if claims["user_id"] != float64(mockUser.Id) {
					t.Errorf("Login() token claims incorrect. got ID: %v, want ID: %v", claims["user_id"], mockUser.Id)
				}
			} else if got != "" {
				t.Errorf("Login() got token '%s', but want empty string on error", got)
			}
		})
	}
}
