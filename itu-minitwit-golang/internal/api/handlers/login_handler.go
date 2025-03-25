package handlers

import (
	"itu-minitwit/internal/service"
	"itu-minitwit/internal/utils"
	"itu-minitwit/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	log := logger.Init()
	// Check if the user is already logged in
	user := utils.GetUserFomContext(c)
	if user != nil {
		log.Info("[LoginHandler] User already logged, redirecting to /")
		c.Redirect(http.StatusFound, "/")
		return
	}

	if c.Request.Method == http.MethodPost {
		username := c.PostForm("username")
		password := c.PostForm("password")

		// Find user in DB

		msg, err := service.HandleLogin(c, username, password)

		if err != nil {
			log.Error("[Login Method] Error signing user in: %v", err)
		}

		if msg != "" {
			utils.SetFlashes(c, msg)
		} else {
			utils.SetFlashes(c, "You were logged in")
			c.Redirect(http.StatusFound, "/")
			log.Info("[LoginHandler] Successfully logged in")
			return
		}

	}

	flashes := utils.GetFlashes(c)

	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":   "Sign In",
		"body":    "login",
		"Flashes": flashes,
	})
}
