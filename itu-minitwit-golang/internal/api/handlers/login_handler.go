package handlers

import (
	"itu-minitwit/internal/utils"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		user, err := utils.FindUserWithName(c, username)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				error = "Invalid username"
			} else {
				error = "Database error"
			}
		} else if !utils.CheckPassword(user.PwHash, password) {
			error = "Invalid password"
		} else {
			// Login successful, set session
			session := sessions.Default(c)
			session.Set("user_id", user.ID)

			c.Set("flash", "You were successfully logged in")
			c.Redirect(http.StatusFound, "/")
			return
		}
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
