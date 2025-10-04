package config

import (
	"os"
)

type Config struct {
	Port  string
	MySQL MySQLConfig
}

type MySQLConfig struct {
	DB_USER   string
	DB_PASS   string
	DB_HOST   string
	DB_PORT   string
	DB_SCHEMA string
}

func Load() Config {
	return Config{
		Port: getEnv("PORT", "8080"),
		MySQL: MySQLConfig{
			DB_USER:   getEnv("DB_USER", "root"),
			DB_PASS:   getEnv("DB_PASS", "root"),
			DB_HOST:   getEnv("DB_HOST", "mysql"),
			DB_PORT:   getEnv("DB_PORT", "3306"),
			DB_SCHEMA: getEnv("DB_SCHEMA", "users"),
		},
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
