package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        int
	DBPath      string
	Environment string
}

func LoadConfig(testing bool) (*Config, error) {
	// Load .env file if it exists

	if testing {
		err := godotenv.Load(".env.test")
		if err != nil {
			return nil, err
		}
	} else {
		godotenv.Load()
	}

	config := &Config{}

	config.Port = getEnvAsInt("PORT", 8080)
	config.DBPath = getEnv("DB_PATH", "database.sqlite")
	config.Environment = getEnv("ENVIRONMENT", "development")

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
