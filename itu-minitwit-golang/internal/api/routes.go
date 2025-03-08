package api

import (
	"github.com/mcuadros/go-gin-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"itu-minitwit/config"
	"itu-minitwit/internal/api/handlers"
	"itu-minitwit/internal/service"
	"itu-minitwit/pkg/database"
	"log"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config) {

	p := ginprometheus.NewPrometheus("gin")

	p.Use(r)
	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("web/templates/*")
	r.GET("/register", handlers.RegisterHandler)
	r.POST("/register", handlers.RegisterHandler)
	r.POST("/login", handlers.LoginHandler)
	r.GET("/login", handlers.LoginHandler)
	r.GET("/", handlers.TimelineHandler)
	r.GET("/public", handlers.PublicTimelineHandler)
	r.GET("/:username", handlers.UserTimelineHandler)
	r.GET("/:username/follow", handlers.FollowHandler)
	r.GET("/:username/unfollow", handlers.UnfollowHandler)
	r.GET("/logout", handlers.LogoutHandler)
	r.GET("/metrics", p.HandlerFunc()) // This exposes the Prometheus metrics

	r.POST("/add_message", handlers.MessageHandler)

	// API endpoints

	db := database.DB
	apiUsers, err := service.GetApiUsers(db)

	if err != nil {
		log.Printf("Error getting API users: %v", err)
		panic(err)
	}

	apiV1 := r.Group("/api/v1", gin.BasicAuth(apiUsers))

	{
		apiV1.GET("/metrics", func(c *gin.Context) {
			handler := promhttp.Handler()
			handler.ServeHTTP(c.Writer, c.Request)
		}) // Tried again?
		apiV1.GET("/latest", handlers.GetLatest)
		apiV1.POST("/register", handlers.RegisterHandlerAPI)
		apiV1.GET("/msgs", handlers.MessagesHandlerAPI)
		apiV1.GET("/msgs/:username", handlers.MessagesPerUserHandlerAPI)
		apiV1.POST("/msgs/:username", handlers.MessagesCreateHandlerAPI)
		apiV1.GET("/fllws/:username", handlers.GetUserFollowersAPI)
		apiV1.POST("/fllws/:username", handlers.FollowUnfollowHandlerAPI)
	}

}
