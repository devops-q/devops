package handlers

import (
	"fmt"
	"itu-minitwit/internal/models"
	"itu-minitwit/internal/service"
	"itu-minitwit/internal/utils"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	var user *models.User = utils.GetUserFomContext(c)

	if user != nil {
		// Already logged in !
		c.Redirect(http.StatusFound, "/")
	}

	var err string
	username := ""
	email := ""

	if c.Request.Method == http.MethodPost {
		username := c.PostForm("username")
		email := c.PostForm("email")
		password := c.PostForm("password")
		password2 := c.PostForm("password2")
		success, err := service.HandleRegister(c, username, email, password, password2)
		if success {
			sessions.Default(c).AddFlash("You were successfully registered and can log in now")
			c.Redirect(http.StatusFound, "/login")
			return

		} else {
			fmt.Println(err)
		}
	}

	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":    "Sign Up",
		"body":     "register",
		"Error":    err,
		"Username": username,
		"Email":    email,
		"Endpoint": "/register",
		"Flashes":  utils.RetrieveFlashes(c),
	})

}
