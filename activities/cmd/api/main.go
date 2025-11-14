package main

import (
	"activities/internal/clients"
	"activities/internal/config"
	"activities/internal/controllers"
	"activities/internal/middleware"
	"activities/internal/repository"
	"activities/internal/services"
	"context"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	activitiesMongoRepo := repository.NewMongoActivitiesRepository(ctx, cfg.Mongo.URI, cfg.Mongo.DB, "activities")
	rabbitClient, err := clients.NewRabbitMQClient(
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Pass,
		cfg.RabbitMQ.QueueName,
	)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ client: %v", err)
	}
	defer rabbitClient.Close()

	activityService := services.NewActivitiesService(activitiesMongoRepo, rabbitClient)
	activityController := controllers.NewActivitiesController(activityService)

	router := gin.Default()
	router.Use(middleware.CORSMiddleware)

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 📚 Rutas de Activities API
	// GET /activities - listar todos los activities (✅ implementado)
	router.GET("/activities", activityController.GetActivities)

	// GET /activities/many?ids=id1,id2,id3 - obtener multiples activities por IDs (público)
	router.GET("/activities/many", activityController.GetManyActivities)

	// POST /activities - crear nuevo activity (protegido)
	router.POST("/activities", middleware.AuthMiddleware(cfg.JwtSecret, "http://users-api:8080/auth"), activityController.CreateActivity)

	// GET /activities/:id - obtener activity por ID (devuelve DTO admin o público según rol)
	router.GET("/activities/:id", middleware.AuthMiddleware(cfg.JwtSecret, "http://users-api:8080/auth"), activityController.GetActivityByID)

	// PUT /activities/:id - actualizar activity existente (protegido)
	router.PUT("/activities/:id", middleware.AuthMiddleware(cfg.JwtSecret, "http://users-api:8080/auth"), activityController.UpdateActivity)

	// DELETE /activities/:id - eliminar activity (protegido)
	router.DELETE("/activities/:id", middleware.AuthMiddleware(cfg.JwtSecret, "http://users-api:8080/auth"), activityController.DeleteActivity)

	// POST /activities/:id/inscribir - inscribir usuario (protegido)
	router.POST("/activities/:id/inscribir", middleware.AuthMiddleware(cfg.JwtSecret, "http://users-api:8080/auth"), activityController.Inscribir)

	// POST /activities/:id/desinscribir - desinscribir usuario (protegido)
	router.POST("/activities/:id/desinscribir", middleware.AuthMiddleware(cfg.JwtSecret, "http://users-api:8080/auth"), activityController.Desinscribir)

	// GET /inscriptions/:userId - obtener actividades inscritas por usuario (protegido)
	router.GET("/inscriptions/:userId", middleware.AuthMiddleware(cfg.JwtSecret, "http://users-api:8080/auth"), activityController.GetInscripcionesByUserID)

	// GET /inscriptions/data/:userId - obtener datos completos de actividades inscritas por usuario (protegido)
	router.GET("/inscriptions/data/:userId", middleware.AuthMiddleware(cfg.JwtSecret, "http://users-api:8080/auth"), activityController.GetInscribedActivities)

	// GET /activities/statistics - obtener estadísticas de actividades (protegido - solo admin)
	router.GET("/activities/statistics", middleware.AuthMiddleware(cfg.JwtSecret, "http://users-api:8080/auth"), activityController.GetStatistics)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("🚀 API listening on port %s", cfg.Port)
	log.Printf("📊 Health check: http://localhost:%s/healthz", cfg.Port)
	log.Printf("📚 Activities API: http://localhost:%s/activities", cfg.Port)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
