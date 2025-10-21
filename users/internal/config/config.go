package config

import (
	"os"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Port      string
	MySQL     MySQLConfig
	JwtSecret string
}

type MySQLConfig struct {
	DB_USER   string
	DB_PASS   string
	DB_HOST   string
	DB_PORT   string
	DB_SCHEMA string
}

func Load() Config {
	log.SetOutput(os.Stderr)
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "02/01/2006-15:04:05:000",
	})

	var secret string

	if secret = getEnv("JWT_SECRET", ""); secret == "" {
		log.Fatalf("no se pudo iniciar la aplicación, se debe especificar la variable de entorno JWT_SECRET")
	}

	cfg := Config{
		Port: getEnv("PORT_USERS_API", "8080"),
		MySQL: MySQLConfig{
			DB_USER:   getEnv("DB_USER", "root"),
			DB_PASS:   getEnv("DB_PASS", "root"),
			DB_HOST:   getEnv("DB_HOST", "mysql"),
			DB_PORT:   getEnv("DB_PORT", "3306"),
			DB_SCHEMA: getEnv("DB_SCHEMA", "users"),
		},
		JwtSecret: secret,
	}

	log.Infoln("=== variables de entorno ===")
	log.Infoln()
	log.Infoln("PORT:", cfg.Port)
	log.Infoln("DB_USER:", cfg.MySQL.DB_USER)
	log.Infoln("DB_PASS:", cfg.MySQL.DB_PASS)
	log.Infoln("DB_HOST:", cfg.MySQL.DB_HOST)
	log.Infoln("DB_PORT:", cfg.MySQL.DB_PORT)
	log.Infoln("DB_SCHEMA:", cfg.MySQL.DB_SCHEMA)
	log.Infoln("JWT_SECRET:", cfg.JwtSecret)
	log.Infoln()
	log.Infoln("==================================")

	return cfg
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
