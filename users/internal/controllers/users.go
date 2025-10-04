package controllers

import (
	"net/http"
	"strings"
	"users/internal/dto"
	"users/internal/services"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type UsersController struct {
	service services.UsersService
}

func NewUsersController(usersService services.UsersService) *UsersController {
	return &UsersController{
		service: usersService,
	}
}

func (c *UsersController) LoginByUsername(ctx *gin.Context) {
	var loginJSON dto.UserLoginDTO
	if err := ctx.BindJSON(&loginJSON); err != nil {
		log.Errorf("error al parsear body al loggear usuaro\nerror: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos con formato incorrecto"})
		return
	}

	token, err := c.service.LoginByUsername(loginJSON.Username, loginJSON.Password)
	if err != nil {
		if err == services.ErrIncorrectCredentials || strings.Contains(err.Error(), "record not found") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		} else {
			log.Debugf("error al loggear el usuario\nerror: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ocurrio un error en el servidor"})
			return
		}
		return
	}

	log.Debugf("usuario loggeado\nusername: %s\npassword: %s\n", loginJSON.Username, loginJSON.Password)
	ctx.JSON(http.StatusCreated, gin.H{
		"access_token": token,
		"token_type":   "bearer",
		"expires_in":   1800, // en segundos
	})
}

func (c *UsersController) LoginByEmail(ctx *gin.Context) {
	var loginJSON dto.UserLoginDTO
	if err := ctx.BindJSON(&loginJSON); err != nil {
		log.Errorf("error al parsear body al loggear usuaro\nerror: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos con formato incorrecto"})
		return
	}

	token, err := c.service.LoginByEmail(loginJSON.Email, loginJSON.Password)
	if err != nil {
		if err == services.ErrIncorrectCredentials || strings.Contains(err.Error(), "record not found") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		} else {
			log.Debugf("error al loggear el usuario\nerror: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ocurrio un error en el servidor"})
			return
		}
		return
	}

	log.Debugf("usuario loggeado\nemail: %s\npassword: %s\n", loginJSON.Email, loginJSON.Password)
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
		log.Debug("LoginDTO:", datos)
		return
	}

	user, err := c.service.Create(datos)
	if err != nil {
		log.Errorf("Error al registrar un usuario: %s\nDTO: %v", err.Error(), datos)

		errString := strings.ToLower(err.Error())
		if strings.Contains(errString, "error 1062") {
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "El usuario ya está registrado"})
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error al registrarse"})
		}

		return
	}

	// una vez registrado el usuario le generamos un token
	token, err := c.service.LoginByUsername(user.Username, user.Password)
	if err != nil {
		log.Debugf("error al loggear al usuario despues de registrarse\nerror: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ocurrio un error en el servidor"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"access_token": token,
		"token_type":   "bearer",
		"expires_in":   1800, // en segundos
	})
}
