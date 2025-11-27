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
	cfg := config.Load()

	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "reindex-dummy-secret")
	}

	ctx := context.Background()

	log.Info("Connecting to MongoDB...")
	activitiesRepo := repository.NewMongoActivitiesRepository(ctx, cfg.Mongo.URI, cfg.Mongo.DB, "activities")
	if activitiesRepo == nil {
		log.Fatal("Failed to initialize MongoDB repository")
	}
	log.Info("MongoDB connection established")

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

	log.Info("Fetching all activities from MongoDB...")
	activities, err := activitiesRepo.List(ctx)
	if err != nil {
		log.Fatalf("Failed to fetch activities: %v", err)
	}
	log.Infof("Found %d activities to reindex", len(activities))

	successCount := 0
	errorCount := 0

	for i, activity := range activities {
		err := rabbitClient.Publish(
			ctx,
			"create",
			activity.ID,
		)
		if err != nil {
			log.Errorf("Failed to publish activity %s (%s): %v", activity.ID, activity.Nombre, err)
			errorCount++
		} else {
			successCount++
		}

		if (i+1)%10 == 0 || (i+1) == len(activities) {
			log.Infof("Progress: %d/%d activities processed (success: %d, errors: %d)",
				i+1, len(activities), successCount, errorCount)
		}
	}

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
