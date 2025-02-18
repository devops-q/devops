package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	Port        int
	DBPath      string
	Environment string
	PerPage     int
}

func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	config := &Config{}

	config.Port = getEnvAsInt("PORT", 8080)
	config.DBPath = getEnv("DB_PATH", "database.sqlite")
	config.Environment = getEnv("ENVIRONMENT", "development")
	config.PerPage = getEnvAsInt("PER_PAGE", 30)

	return config, nil
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}
