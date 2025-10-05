package main

import (
	"log"
	"net/http"
	"time"
	"users/internal/config"
	"users/internal/controllers"
	"users/internal/middleware"
	"users/internal/repository"
	"users/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	usersMySQLRepo := repository.NewMySQLUsersRepository(cfg.MySQL)
	userService := services.NewUsersService(usersMySQLRepo, cfg.JwtSecret)
	userController := controllers.NewUsersController(&userService)

	router := gin.Default()
	router.Use(middleware.CORSMiddleware)

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.POST("/register", userController.Create)
	router.POST("/login", userController.Login)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("🚀 API listening on port %s", cfg.Port)
	log.Printf("📊 Health check: http://localhost:%s/healthz", cfg.Port)
	log.Printf("📚 users API: http://localhost:%s/users", cfg.Port)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
