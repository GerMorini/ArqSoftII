package config

import (
	"os"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Port      string
	Mongo     MongoConfig
	JwtSecret string
}

type MongoConfig struct {
	URI string
	DB  string
}

func Load() Config {
	var secret = getEnv("JWT_SECRET", "")
	if secret == "" {
		log.Fatalf("no se pudo iniciar la aplicación, se debe especificar la variable de entorno JWT_SECRET")
	}

	cfg := Config{
		Port: getEnv("PORT_ACTIVIDADES_API", "8080"),
		Mongo: MongoConfig{
			URI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
			DB:  getEnv("MONGO_DB", "demo"),
		},

		JwtSecret: secret,
	}

	log.Infoln("=== variables de entorno ===")
	log.Infoln("PORT:", cfg.Port)
	log.Infoln("MONGO_URI:", cfg.Mongo.URI)
	log.Infoln("MONGO_DB:", cfg.Mongo.DB)
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
