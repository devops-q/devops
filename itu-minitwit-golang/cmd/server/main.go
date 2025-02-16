package main

import (
	"fmt"
	"itu-minitwit/config"
	"itu-minitwit/internal/api"
	"itu-minitwit/internal/api/middlewares"
	"itu-minitwit/pkg/database"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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

	store := cookie.NewStore([]byte("secret"))            // TODO this should be an actual secret if we dont want people to read our sessions
	store.Options(sessions.Options{MaxAge: 60 * 60 * 12}) // Cookie will last max 12 hours
	r.Use(sessions.Sessions("itu-minitwit-session", store))

	r.Use(middlewares.SetDbMiddleware())
	r.Use(middlewares.SetUserContext())

	api.SetupRoutes(r, cfg)

	log.Printf("Server starting on port %d", cfg.Port)
	r.Run(fmt.Sprintf("localhost:%d", cfg.Port))
}
