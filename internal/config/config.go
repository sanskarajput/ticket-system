package config

import (
	"os"
)

type Config struct {
	Port      string
	DBPath    string
	JWTSecret string
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func Load() *Config {
	
	return &Config{
		Port:      getEnv("PORT", "8080"),
		DBPath:    getEnv("DB_PATH", "./tickets.db"),
		JWTSecret: getEnv("JWT_SECRET", "super-secret-key-change-in-prod"),
	}
}

