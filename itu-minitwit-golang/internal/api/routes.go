package api

import (
	"github.com/gin-gonic/gin"
	"itu-minitwit/config"
	"itu-minitwit/internal/api/handlers"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("web/templates/*")

	r.GET("/ping", handlers.PingHandler)
	r.GET("/hello/:name", handlers.HelloHandler)

	// Add more routes here
}
