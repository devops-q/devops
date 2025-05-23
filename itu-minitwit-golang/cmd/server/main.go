package main

import (
	"fmt"
	"html/template"
	"itu-minitwit/config"
	"itu-minitwit/internal/api"
	"itu-minitwit/internal/api/middlewares"
	"itu-minitwit/internal/utils"
	"itu-minitwit/pkg/database"
	"itu-minitwit/pkg/logger"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	log := logger.Init()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	database.InitDb(cfg)
	err = database.InitApiUserIfNotExists(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize API user: %v", err)
	}

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	store := cookie.NewStore([]byte("secret"))                       // TODO this should be an actual secret if we dont want people to read our sessions
	store.Options(sessions.Options{MaxAge: 60 * 60 * 12, Path: "/"}) // Cookie will last max 12 hours
	r.Use(sessions.Sessions("itu-minitwit-session", store))

	r.SetFuncMap(template.FuncMap{
		"GravatarURL":     utils.GravatarURL,
		"FormatDateTime":  utils.FormatDateTime,
		"ToISODateString": utils.ToISODateString,
	})

	r.Use(middlewares.PrometheusMiddleware(r).Instrument())
	r.Use(middlewares.SetConfigMiddleware(cfg))
	r.Use(middlewares.SetDbMiddleware())
	r.Use(middlewares.SetUserContext())
	r.Use(middlewares.UpdateLatestMiddleware())
	api.SetupRoutes(r, cfg)

	log.Info(fmt.Sprintf("Server starting on port %d", cfg.Port))
	err = r.Run(fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
