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

	var error string

	if c.Request.Method == http.MethodPost {
		username := c.PostForm("username")
		password := c.PostForm("password")

		// Find user in DB
		error = service.HandleLogin(c, username, password)
	}

	// Render login page (moved outside the if block)
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":    "Sign In",
		"body":     "login",
		"Error":    error,
		"Username": "",
		"Email":    "",
	})
}
