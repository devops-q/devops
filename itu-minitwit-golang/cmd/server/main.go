package main

import (
	"fmt"
	"itu-minitwit/config"
	"itu-minitwit/pkg/database"
	"itu-minitwit/setup"
	"log"
)

func main() {
	cfg, err := config.LoadConfig(false)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	database.InitDb(cfg)

	r := setup.SetupRouter(cfg)

	log.Printf("Server starting on port %d", cfg.Port)
	r.Run(fmt.Sprintf(":%d", cfg.Port))
}
