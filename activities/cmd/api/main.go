package main

import (
	"activities/internal/clients"
	"activities/internal/config"
	"activities/internal/controllers"
	"activities/internal/middleware"
	"activities/internal/repository"
	"activities/internal/services"
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 📋 Cargar configuración desde las variables de entorno
	cfg := config.Load()

	// 🏗️ Inicializar capas de la aplicación (Dependency Injection)
	// Patrón: Repository -> Service -> Controller
	// Cada capa tiene una responsabilidad específica

	// Context
	ctx := context.Background()

	// Capa de datos: maneja operaciones DB
	activitiesMongoRepo := repository.NewMongoActivitiesRepository(ctx, cfg.Mongo.URI, cfg.Mongo.DB, "activities")

	// RabbitMQ client: para publicar eventos de actividades
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

	// Capa de lógica de negocio: validaciones, transformaciones
	activityService := services.NewActivitiesService(activitiesMongoRepo, rabbitClient)

	// Capa de controladores: maneja HTTP requests/responses
	activityController := controllers.NewActivitiesController(activityService)

	// Cache (ejercicio: ajustar TTL y agregar "índice" de claves)
	// cache := cache.NewMemcached(memAddr)

	// 🌐 Configurar router HTTP con Gin
	router := gin.Default()

	// Middleware: funciones que se ejecutan en cada request
	router.Use(middleware.CORSMiddleware)

	// 🏥 Health check endpoint
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 📚 Rutas de Activities API
	// GET /activities - listar todos los activities (✅ implementado)
	router.GET("/activities", activityController.GetActivities)

	// GET /activities/many?ids=id1,id2,id3 - obtener multiples activities por IDs (público)
	router.GET("/activities/many", activityController.GetManyActivities)

	// POST /activities - crear nuevo activity (protegido)
	router.POST("/activities", middleware.AuthMiddleware(cfg.JwtSecret), activityController.CreateActivity)

	// GET /activities/:id - obtener activity por ID (devuelve DTO admin o público según rol)
	router.GET("/activities/:id", middleware.AuthMiddleware(cfg.JwtSecret), activityController.GetActivityByID)

	// PUT /activities/:id - actualizar activity existente (protegido)
	router.PUT("/activities/:id", middleware.AuthMiddleware(cfg.JwtSecret), activityController.UpdateActivity)

	// DELETE /activities/:id - eliminar activity (protegido)
	router.DELETE("/activities/:id", middleware.AuthMiddleware(cfg.JwtSecret), activityController.DeleteActivity)

	// POST /activities/:id/inscribir - inscribir usuario (protegido)
	router.POST("/activities/:id/inscribir", middleware.AuthMiddleware(cfg.JwtSecret), activityController.Inscribir)

	// POST /activities/:id/desinscribir - desinscribir usuario (protegido)
	router.POST("/activities/:id/desinscribir", middleware.AuthMiddleware(cfg.JwtSecret), activityController.Desinscribir)

	// GET /inscriptions/:userId - obtener actividades inscritas por usuario (protegido)
	router.GET("/inscriptions/:userId", middleware.AuthMiddleware(cfg.JwtSecret), activityController.GetInscripcionesByUserID)

	// Configuración del server HTTP
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("🚀 API listening on port %s", cfg.Port)
	log.Printf("📊 Health check: http://localhost:%s/healthz", cfg.Port)
	log.Printf("📚 Activities API: http://localhost:%s/activities", cfg.Port)

	// Iniciar servidor (bloquea hasta que se pare el servidor)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
