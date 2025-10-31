package main

import (
	"context"
	"log"
	"net/http"
	"search/internal/clients"
	"search/internal/config"
	"search/internal/controllers"
	"search/internal/middleware"
	"search/internal/repository"
	"search/internal/services"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	activiesMemcachedRepo := repository.NewMemcachedActivitysRepository(
		cfg.Memcached.Host,
		cfg.Memcached.Port,
		time.Duration(cfg.Memcached.TTLSeconds)*time.Second,
	)

	activitysSolrRepo := repository.NewSolrActivitysRepository(
		cfg.Solr.Host,
		cfg.Solr.Port,
		cfg.Solr.Core,
	)

	activiesQueue := clients.NewRabbitMQClient(
		cfg.RabbitMQ.Username,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.QueueName,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)

	activityService := services.NewActivitysService(activiesMemcachedRepo, activitysSolrRepo, activiesQueue, activiesQueue)
	go activityService.InitConsumer(ctx)

	activityController := controllers.NewActivitiesController(&activityService)
	router := gin.Default()

	router.Use(middleware.CORSMiddleware)

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/activitys", activityController.List)
	router.POST("/activitys", activityController.CreateActivity)
	router.GET("/activitys/:id", activityController.GetActivityByID)
	router.PUT("/activitys/:id", activityController.UpdateActivity)
	router.DELETE("/activitys/:id", activityController.DeleteActivity)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("🚀 API listening on port %s", cfg.Port)
	log.Printf("📊 Health check: http://localhost:%s/healthz", cfg.Port)
	log.Printf("📚 Activitys API: http://localhost:%s/activitys", cfg.Port)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
