package setup

import (
	"itu-minitwit/config"
	"itu-minitwit/internal/api"
	"itu-minitwit/internal/api/middlewares"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
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
	return r
}
