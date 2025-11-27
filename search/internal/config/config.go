package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Port             string
	Memcached        MemcachedConfig
	RabbitMQ         RabbitMQConfig
	Solr             SolrConfig
	ActivitiesAPIURL string
}

type MemcachedConfig struct {
	Host       string
	Port       string
	TTLSeconds int
}

type RabbitMQConfig struct {
	Username  string
	Password  string
	QueueName string
	Host      string
	Port      string
}

type SolrConfig struct {
	Host string
	Port string
	Core string
}

var config *Config

func Load() *Config {
	if config != nil {
		return config
	}

	log.SetOutput(os.Stderr)
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "02/01/2006-15:04:05.000",
	})

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	memcachedTTL, err := strconv.Atoi(getEnv("MEMCACHED_TTL_SECONDS", "60"))
	if err != nil {
		memcachedTTL = 60
	}

	config = &Config{
		Port: getEnv("PORT", "8080"),
		Memcached: MemcachedConfig{
			Host:       getEnv("MEMCACHED_HOST", "localhost"),
			Port:       getEnv("MEMCACHED_PORT", "11211"),
			TTLSeconds: memcachedTTL,
		},
		RabbitMQ: RabbitMQConfig{
			Username:  getEnv("RABBITMQ_USER", "admin"),
			Password:  getEnv("RABBITMQ_PASS", "admin"),
			QueueName: getEnv("RABBITMQ_QUEUE_NAME", "items-news"),
			Host:      getEnv("RABBITMQ_HOST", "localhost"),
			Port:      getEnv("RABBITMQ_PORT", "5672"),
		},
		Solr: SolrConfig{
			Host: getEnv("SOLR_HOST", "localhost"),
			Port: getEnv("SOLR_PORT", "8983"),
			Core: getEnv("SOLR_CORE", "demo"),
		},
		ActivitiesAPIURL: getEnv("ACTIVITIES_API_URL", "http://activities-api:8080"),
	}

	log.Infoln("========== CONFIGURACIÓN ==========")
	log.Infoln("PORT:", config.Port)
	log.Infoln("MEMCACHED_HOST:", config.Memcached.Host)
	log.Infoln("MEMCACHED_PORT:", config.Memcached.Port)
	log.Infoln("MEMCACHED_TTL_SECONDS:", config.Memcached.TTLSeconds)
	log.Infoln("RABBITMQ_USER:", config.RabbitMQ.Username)
	log.Infoln("RABBITMQ_PASS:", config.RabbitMQ.Password)
	log.Infoln("RABBITMQ_QUEUE_NAME:", config.RabbitMQ.QueueName)
	log.Infoln("RABBITMQ_HOST:", config.RabbitMQ.Host)
	log.Infoln("RABBITMQ_PORT:", config.RabbitMQ.Port)
	log.Infoln("SOLR_HOST", config.Solr.Host)
	log.Infoln("SOLR_PORT", config.Solr.Port)
	log.Infoln("SOLR_CORE", config.Solr.Core)
	log.Infoln("ACTIVITIES_API_URL:", config.ActivitiesAPIURL)
	log.Infoln("===================================")

	return config
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
