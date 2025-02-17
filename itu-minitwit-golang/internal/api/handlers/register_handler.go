package handlers

import (
	"itu-minitwit/internal/models"
	"itu-minitwit/internal/service"
	"itu-minitwit/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	var user *models.User = utils.GetUserFomContext(c)

	if user != nil {
		// Already logged in ! 
		c.Redirect(http.StatusFound, "/")
	}

	var err string

	if c.Request.Method == http.MethodPost {
		username := c.PostForm("username")
		email := c.PostForm("email")
		password := c.PostForm("password")
		password2 := c.PostForm("password2")

		err = service.HandleRegister(c,username,email,password,password2)
	}



	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":    "Sign Up",
		"body":     "register",
		"Error":    err,
		"Username": "",
		"Email":    "",
		"Endpoint": "/register",
	})
}
