package main

import (
	"fmt"
	"itu-minitwit/config"
	"itu-minitwit/internal/api"
	"itu-minitwit/internal/models"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	database.InitDb(cfg)

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	api.SetupRoutes(r, cfg)

	log.Printf("Server starting on port %d", cfg.Port)
	r.Run(fmt.Sprintf("localhost:%d", cfg.Port))
}
