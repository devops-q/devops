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

	if c.Request.Method == http.MethodPost {
		err = service.HandleRegister(c)
		if err == "" { // If there is no returned error string.
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
		"Username": "",
		"Email":    "",
		"Endpoint": "/register",
		"Flashes":  utils.RetrieveFlashes(c),
	})

}
