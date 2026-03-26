package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	DB_SSLMODE  string
	HTTP_PORT   string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("env. не найден")
	}
	return &Config{
		DB_HOST:     getEnv("DB_HOST", "localhost"),
		DB_PORT:     getEnv("DB_PORT", "5432"),
		DB_USER:     getEnv("DB_USER", "postgres"),
		DB_PASSWORD: getEnv("DB_PASSWORD", "1234"),
		DB_NAME:     getEnv("DB_NAME", "soar"),
		DB_SSLMODE:  getEnv("DB_SSLMODE", "disable"),
		HTTP_PORT:   getEnv("HTTP_PORT", "8080"),
	}

}
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}