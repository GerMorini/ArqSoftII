package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"users/internal/dto"
	"users/internal/services"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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
		log.Warnf("error al parsear body al loggear usuaro: %v\n", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Datos con formato incorrecto"})
		return
	}

	token, err := c.service.Login(loginDTO)
	if err != nil {
		if err == services.ErrIncorrectCredentials || strings.Contains(err.Error(), "record not found") {
			log.Warnf("contraseña incorrecta para el usuario:\nusername: %s\nemail: %s\npassword: %s\n", loginDTO.Username, loginDTO.Email, loginDTO.Password)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		} else if err == services.ErrLoginFormat {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			log.Errorf("error al realizar login de usuario: %s", err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Ocurrio un error en el servidor"})
			return
		}
		return
	}

	log.Infof("usuario loggeado\nusername: %s\npassword: %s\n", loginDTO.Username, loginDTO.Password)
	ctx.JSON(http.StatusCreated, gin.H{
		"access_token": token,
		"token_type":   "bearer",
		"expires_in":   1800, // en segundos
	})
}

func (c *UsersController) Create(ctx *gin.Context) {
	var datos dto.UserMinDTO

	if err := ctx.BindJSON(&datos); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Datos con formato incorrecto"})
		log.Warnf("no se pudo procesar la petición del usuario: %s\nLoginDTO: %v", err.Error(), datos)
		return
	}

	_, err := c.service.Create(datos)
	if err != nil {
		log.Warnf("error al registrar un usuario: %s\nDTO: %v", err.Error(), datos)

		errString := strings.ToLower(err.Error())
		if strings.Contains(errString, "error 1062") {
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "El usuario ya está registrado"})
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error al registrarse"})
		}

		return
	}

	// una vez registrado el usuario le generamos un token
	token, err := c.service.Login(dto.UserLoginDTO{
		Username: datos.Username,
		Password: datos.Password,
	})
	if err != nil {
		log.Errorf("error al loggear al usuario despues de registrarse: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ocurrio un error en el servidor"})
		return
	}

	log.Infof("usuario registrado exitosamente: %v", datos)
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
		log.Warnf("no se pudo obtener el ID del parámetro de la consulta: %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID con formato incorrecto. Debe ser un número"})
		return
	}

	userData, err := c.service.GetByID(id)
	if err != nil {
		log.Warnf("error al buscar un usuario por su ID: %v", err)
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error al buscar usuario"})
		return
	}

	ctx.JSON(http.StatusOK, userData)
}
