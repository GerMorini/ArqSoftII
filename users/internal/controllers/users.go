package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"users/internal/dto"
	"users/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UsersController struct {
	service services.UsersService
}

func NewUsersController(usersService services.UsersService) *UsersController {
	return &UsersController{
		service: usersService,
	}
}

func (c *UsersController) Login(ctx *gin.Context) {
	var loginDTO dto.UserLoginDTO

	if err := ctx.BindJSON(&loginDTO); err != nil {
		ctx.Error(fmt.Errorf("error al parsear body al loggear usuaro: %v", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos con formato incorrecto"})
		return
	}

	token, err := c.service.Login(loginDTO)
	if err != nil {
		if err == services.ErrIncorrectCredentials || strings.Contains(err.Error(), "record not found") {
			ctx.Error(fmt.Errorf("contraseña incorrecta para el usuario: username=%s email=%s", loginDTO.Username, loginDTO.Email))
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		} else if err == services.ErrLoginFormat {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			ctx.Error(fmt.Errorf("error al realizar login de usuario: %s", err.Error()))
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ocurrio un error en el servidor"})
			return
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"access_token": token,
		"token_type":   "bearer",
		"expires_in":   1800, // en segundos
	})
}

func (c *UsersController) Create(ctx *gin.Context) {
	var datos dto.UserMinDTO

	if err := ctx.BindJSON(&datos); err != nil {
		ctx.Error(fmt.Errorf("no se pudo procesar la petición del usuario: %s", err.Error()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos con formato incorrecto"})
		return
	}

	_, err := c.service.Create(datos)
	if err != nil {
		ctx.Error(fmt.Errorf("error al registrar un usuario: %s", err.Error()))

		errString := strings.ToLower(err.Error())
		if strings.Contains(errString, "error 1062") {
			ctx.JSON(http.StatusConflict, gin.H{"error": "El usuario ya está registrado"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrarse"})
		}

		return
	}

	// una vez registrado el usuario le generamos un token
	token, err := c.service.Login(dto.UserLoginDTO{
		Username: datos.Username,
		Password: datos.Password,
	})
	if err != nil {
		ctx.Error(fmt.Errorf("error al loggear al usuario despues de registrarse: %v", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ocurrio un error en el servidor"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"access_token": token,
		"token_type":   "bearer",
		"expires_in":   1800, // en segundos
	})
}

func (c *UsersController) GetByID(ctx *gin.Context) {
	id_str := ctx.Param("id")

	id, err := strconv.Atoi(id_str)
	if err != nil {
		ctx.Error(fmt.Errorf("no se pudo obtener el ID del parámetro de la consulta: %s", err.Error()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID con formato incorrecto. Debe ser un número"})
		return
	}

	userData, err := c.service.GetByID(id)
	if err != nil {
		ctx.Error(fmt.Errorf("error al buscar un usuario por su ID: %v", err))
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al buscar usuario"})
		return
	}

	ctx.JSON(http.StatusOK, userData)
}

func (c *UsersController) GetAll(ctx *gin.Context) {
	usuarios, err := c.service.GetAll()
	if err != nil {
		ctx.Error(fmt.Errorf("error al obtener todos los usuarios: %v", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener usuarios"})
		return
	}

	ctx.JSON(http.StatusOK, usuarios)
}

func (c *UsersController) Update(ctx *gin.Context) {
	id_str := ctx.Param("id")

	id, err := strconv.Atoi(id_str)
	if err != nil {
		ctx.Error(fmt.Errorf("no se pudo obtener el ID del parámetro de la consulta: %s", err.Error()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID con formato incorrecto. Debe ser un número"})
		return
	}

	var updateDTO dto.UserUpdateDTO
	if err := ctx.BindJSON(&updateDTO); err != nil {
		ctx.Error(fmt.Errorf("error al parsear body al actualizar usuario: %v", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos con formato incorrecto"})
		return
	}

	userData, err := c.service.Update(id, updateDTO)
	if err != nil {
		ctx.Error(fmt.Errorf("error al actualizar usuario con ID %d: %v", id, err))
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar usuario"})
		return
	}

	ctx.JSON(http.StatusOK, userData)
}

func (c *UsersController) Delete(ctx *gin.Context) {
	id_str := ctx.Param("id")

	id, err := strconv.Atoi(id_str)
	if err != nil {
		ctx.Error(fmt.Errorf("no se pudo obtener el ID del parámetro de la consulta: %s", err.Error()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID con formato incorrecto. Debe ser un número"})
		return
	}

	err = c.service.Delete(id)
	if err != nil {
		ctx.Error(fmt.Errorf("error al eliminar usuario con ID %d: %v", id, err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error al eliminar usuario"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "usuario eliminado exitosamente"})
}

func (c *UsersController) IsAdmin(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token != "" {
		token = strings.TrimPrefix(token, "Bearer ")
	} else {
		ctx.Error(errors.New("usuario sin autorizacion: no se especifico header 'Authorization'"))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "se necesita especificar el campo 'Authorization'"})
		return
	}

	bool, err := c.service.IsAdmin(token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "se necesita autenticación válida"})
		return
	}

	if !bool {
		ctx.Error(errors.New("validacion negativa"))
		ctx.Status(http.StatusForbidden)
		return
	}

	ctx.Status(http.StatusOK)
}
