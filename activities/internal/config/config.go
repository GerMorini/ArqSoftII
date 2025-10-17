package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
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
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}
	var secret = getEnv("JWT_SECRET", "")
	if secret == "" {
		log.Fatalf("no se pudo iniciar la aplicación, se debe especificar la variable de entorno JWT_SECRET")
	}

	return Config{
		Port: getEnv("PORT_ACTIVIDADES_API", "8080"),
		Mongo: MongoConfig{
			URI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
			DB:  getEnv("MONGO_DB", "demo"),
		},

		JwtSecret: secret,
	}

}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
