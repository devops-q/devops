package handlers

import (
	"itu-minitwit/internal/service"
	"itu-minitwit/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	// Check if the user is already logged in
	user := utils.GetUserFomContext(c)
	if user != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	if c.Request.Method == http.MethodPost {
		username := c.PostForm("username")
		password := c.PostForm("password")

		// Find user in DB

		msg, _ := service.HandleLogin(c, username, password)

		if msg != "" {
			utils.SetFlashes(c, msg)
		} else {
			utils.SetFlashes(c, "You were logged in")
			c.Redirect(http.StatusFound, "/")
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
