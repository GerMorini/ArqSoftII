package main

import (
	"context"
	"net/http"
	"search/internal/clients"
	"search/internal/config"
	"search/internal/controllers"
	"search/internal/middleware"
	"search/internal/repository"
	"search/internal/services"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	activitiesLocalCacheRepo := repository.NewActivitysLocalCacheRepository(1 * time.Hour)

	activiesMemcachedRepo := repository.NewMemcachedActivitiesRepository(
		cfg.Memcached.Host,
		cfg.Memcached.Port,
		time.Duration(cfg.Memcached.TTLSeconds)*time.Second,
	)

	activitiesSolrRepo := repository.NewSolrActivitysRepository(
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

	activityService := services.NewActivitiesService(activiesMemcachedRepo, activitiesSolrRepo, activiesQueue)
	go activityService.InitConsumer(ctx)

	activityController := controllers.NewActivitiesController(&activityService)
	router := gin.Default()

	router.Use(middleware.CORSMiddleware)

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/activitys", activityController.List)

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
