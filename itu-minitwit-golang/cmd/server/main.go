package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"itu-minitwit/config"
	"itu-minitwit/internal/api"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	api.SetupRoutes(r, cfg)

	log.Printf("Server starting on port %d", cfg.Port)
	r.Run(fmt.Sprintf(":%d", cfg.Port))
}
