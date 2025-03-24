package main

import (
	"flag"
	"fmt"
	"itu-minitwit/config"
	"itu-minitwit/internal/service"
	"itu-minitwit/pkg/database"
	"itu-minitwit/pkg/logger"
	"os"
)

func main() {
	log := logger.Init()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	database.InitDb(cfg)

	username := flag.String("username", "", "API user username")
	password := flag.String("password", "", "API user password")
	flag.Parse()

	if *username == "" || *password == "" {
		fmt.Println("Usage: api_user -username=<username> -password=<password>")
		os.Exit(1)
	}

	db := database.DB
	success, err := service.CreateApiUser(db, *username, *password)
	if err != nil {
		log.Fatalf("Failed to create API user: %v", err)
	}

	if success {
		fmt.Printf("API user %s created successfully\n", *username)
	}
}
