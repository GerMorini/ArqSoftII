package main

import (
	"activities/internal/clients"
	"activities/internal/config"
	"activities/internal/repository"
	"context"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	// Load configuration from environment variables
	cfg := config.Load()

	// JWT_SECRET is not required for reindexing, but config.Load() requires it
	// Set a dummy value if not set to avoid panic
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "reindex-dummy-secret")
	}

	// Initialize context
	ctx := context.Background()

	// Initialize MongoDB repository to fetch activities
	log.Info("Connecting to MongoDB...")
	activitiesRepo := repository.NewMongoActivitiesRepository(ctx, cfg.Mongo.URI, cfg.Mongo.DB, "activities")
	if activitiesRepo == nil {
		log.Fatal("Failed to initialize MongoDB repository")
	}
	log.Info("MongoDB connection established")

	// Initialize RabbitMQ client to publish messages
	log.Info("Connecting to RabbitMQ...")
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
	log.Info("RabbitMQ connection established")

	// Fetch all activities from MongoDB
	log.Info("Fetching all activities from MongoDB...")
	activities, err := activitiesRepo.List(ctx)
	if err != nil {
		log.Fatalf("Failed to fetch activities: %v", err)
	}
	log.Infof("Found %d activities to reindex", len(activities))

	// Publish each activity to RabbitMQ with "create" action
	successCount := 0
	errorCount := 0

	for i, activity := range activities {
		err := rabbitClient.Publish(
			ctx,
			"create",
			activity.ID,
			activity.Nombre,
			activity.Descripcion,
			activity.DiaSemana,
		)
		if err != nil {
			log.Errorf("Failed to publish activity %s (%s): %v", activity.ID, activity.Nombre, err)
			errorCount++
		} else {
			successCount++
		}

		// Log progress every 10 activities
		if (i+1)%10 == 0 || (i+1) == len(activities) {
			log.Infof("Progress: %d/%d activities processed (success: %d, errors: %d)",
				i+1, len(activities), successCount, errorCount)
		}
	}

	// Final summary
	log.Info("=== Reindex Summary ===")
	log.Infof("Total activities: %d", len(activities))
	log.Infof("Successfully published: %d", successCount)
	log.Infof("Errors: %d", errorCount)
	log.Info("=======================")

	if errorCount > 0 {
		log.Warn("Reindexing completed with errors")
		os.Exit(1)
	}

	log.Info("Reindexing completed successfully!")
}
