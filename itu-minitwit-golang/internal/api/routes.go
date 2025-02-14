package api

import (
	"github.com/gin-gonic/gin"
	"itu-minitwit/config"
	"itu-minitwit/internal/api/handlers"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	r.GET("/ping", handlers.PingHandler)

	// Add more routes here
}
