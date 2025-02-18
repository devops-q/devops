package api

import (
	"itu-minitwit/config"
	"itu-minitwit/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("web/templates/*")

	r.GET("/ping", handlers.PingHandler)
	r.GET("/hello/:name", handlers.HelloHandler)
	r.GET("/register", handlers.RegisterHandler)	
	r.POST("/register", handlers.RegisterHandler)
	r.POST("/login", handlers.LoginHandler)


	r.GET("/login", handlers.LoginHandler)
	r.GET("/public", handlers.PublicHandler)
	r.GET("/", handlers.TimelineHandler)

	r.GET("/logout", handlers.LogoutHandler)

	// Examples of how to use the ORM in endpoints
	r.GET("/users", handlers.GetUsersHandler)
	r.GET("/user/create/:name", handlers.CreateUserHandler)
	r.GET("/user/find/:name", handlers.FindUserWithName)
	r.GET("/users/messages", handlers.GetAllUsersWithNonFlaggedMessages)

	r.GET("/user/current", handlers.GetUserInSession)
	r.GET("/user/force-login/:id", handlers.ForceSetUserId)

	// Add more routes here
}
