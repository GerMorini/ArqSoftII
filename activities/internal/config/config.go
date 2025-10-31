package config

import (
	"os"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Port      string
	Mongo     MongoConfig
	RabbitMQ  RabbitMQConfig
	JwtSecret string
}

type MongoConfig struct {
	URI string
	DB  string
}

type RabbitMQConfig struct {
	Host      string
	Port      string
	User      string
	Pass      string
	QueueName string
}

func Load() Config {
	log.SetOutput(os.Stderr)
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "02/01/2006-15:04:05.000",
	})

	var secret = getEnv("JWT_SECRET", "")
	if secret == "" {
		log.Fatalf("no se pudo iniciar la aplicación, se debe especificar la variable de entorno JWT_SECRET")
	}

	cfg := Config{
		Port: getEnv("PORT_ACTIVIDADES_API", "8080"),
		Mongo: MongoConfig{
			URI: getEnv("MONGO_URI", "mongodb://mongo_activities_api:27017"),
			DB:  getEnv("MONGO_DB", "demo"),
		},
		RabbitMQ: RabbitMQConfig{
			Host:      getEnv("RABBITMQ_HOST", "rabbit-search-api"),
			Port:      getEnv("RABBITMQ_PORT", "5672"),
			User:      getEnv("RABBITMQ_USER", "admin"),
			Pass:      getEnv("RABBITMQ_PASS", "admin"),
			QueueName: getEnv("RABBITMQ_QUEUE_NAME", "items"),
		},
		// Solr indexing is handled by the search service; activities service
		// does not need Solr configuration anymore.
		JwtSecret: secret,
	}

	log.Infoln("=== variables de entorno ===")
	log.Infoln("PORT:", cfg.Port)
	log.Infoln("MONGO_URI:", cfg.Mongo.URI)
	log.Infoln("MONGO_DB:", cfg.Mongo.DB)
	log.Infoln("RABBITMQ_HOST:", cfg.RabbitMQ.Host)
	log.Infoln("RABBITMQ_PORT:", cfg.RabbitMQ.Port)
	log.Infoln("RABBITMQ_QUEUE:", cfg.RabbitMQ.QueueName)
	log.Infoln("JWT_SECRET:", cfg.JwtSecret)
	log.Infoln("==================================")
	return cfg
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
