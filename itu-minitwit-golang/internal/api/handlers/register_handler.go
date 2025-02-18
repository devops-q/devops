package handlers

import (
	"fmt"
	"gorm.io/gorm"
	"itu-minitwit/internal/models"
	"itu-minitwit/internal/service"
	"itu-minitwit/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	var db = c.MustGet("DB").(*gorm.DB)
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
		success, err := service.RegisterUser(db, username, email, password, password2)
		if success {
			utils.SetFlashes(c, "You were successfully registered and can log in now")
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

func RegisterHandlerAPI(c *gin.Context) {
	var db = c.MustGet("DB").(*gorm.DB)

	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	password2 := c.PostForm("password2")

	success, err := service.RegisterUser(db, username, email, password, password2)
	if success {
		c.JSON(http.StatusNoContent, nil)
		return
	} else {
		c.JSON(http.StatusBadRequest, utils.ErrorCodeMessageResponse{Code: http.StatusBadRequest, ErrorMessage: err})
	}
}
