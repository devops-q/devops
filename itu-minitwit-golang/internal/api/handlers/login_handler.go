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

		msg, error2 := service.HandleLogin(c, username, password)

		if msg != "" {
			error = error2.Error()
		} else {
			c.Set("flash", "You were successfully logged in")
			c.Redirect(http.StatusFound, "/")
			return
		}

	}

	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":    "Sign In",
		"body":     "login",
		"Error":    error,
		"Username": "",
		"Email":    "",
	})
}
