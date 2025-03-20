package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        int
	Environment string
	PerPage     int

	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBSSLMode  string

	InitialApiUser     string
	InitialApiPassword string
}

func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	config := &Config{}

	config.Port = getEnvAsInt("PORT", 8080)
	config.Environment = getEnv("ENVIRONMENT", "development")
	config.PerPage = getEnvAsInt("PER_PAGE", 30)

	config.DBHost = mustGetEnv("DB_HOST")
	config.DBUser = mustGetEnv("DB_USER")
	config.DBPassword = mustGetEnv("DB_PASSWORD")
	config.DBName = mustGetEnv("DB_NAME")
	config.DBPort = mustGetEnv("DB_PORT")
	config.DBSSLMode = getEnv("DB_SSL_MODE", "disable")

	config.InitialApiUser = getEnv("INITIAL_API_USER", "")
	config.InitialApiPassword = getEnv("INITIAL_API_PASSWORD", "")

	return config, nil
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func mustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		panic("Environment variable " + key + " not set")
	}

	return value
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}
