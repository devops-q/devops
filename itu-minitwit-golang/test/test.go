package test

import (
	"itu-minitwit/pkg/database"
	"log"
	"os"
)

func setup() {
	var db, err = os.CreateTemp("", ".db")
	if err != nil {
		log.Fatalf("Failed to create temp database: %v", err)
	}
	// TODO define app, init app db?
	database.InitDb(cfg)
}

func register(username string, password string, password2 string, email string) {
	if password2 == "" {
		password2 = password
	}
	if email == "" {
		email = username + "@example.com"
	}
	// TODO app post regest to /register
}
