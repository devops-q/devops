package handlers

import (
	"fmt"
	"itu-minitwit/internal/models"
	"itu-minitwit/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

		if username == "" {
			err = "You have to eneter a username"
		} else if email == "" || !strings.Contains(email, "@") {
			err = "You have to enter a valid email address"
		} else if password == "" {
			err = "You have to enter a password"
		} else if password2 != password {
			err = "The two passwords do not match"
		} else if utils.UserExists(c, username) {
			err = "The username is already taken"
		} else {
			hashed, error := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if error != nil {
				fmt.Println("Error hashing password: ", error)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
				return
			}

			utils.CreateUser(c, username, email, string(hashed))

			c.Set("flash", "You were successfully registered and can log in now")
			c.Redirect(http.StatusFound, "/login")
			return
		}
	}

	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":    "Sign Up",
		"body":     "register",
		"Error":    err,
		"Username": "",
		"Email":    "",
	})
}
