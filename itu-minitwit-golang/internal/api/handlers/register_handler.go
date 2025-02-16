package handlers

import (
	"fmt"
	"itu-minitwit/internal/models"
	"itu-minitwit/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterHandler(c *gin.Context) {
	// currently assumes that the user is not logged in.

	var err string

	if user, exists := c.Get("user"); exists && user != nil {
		c.Redirect(http.StatusFound, "/timeline")
	}

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

			user := &models.User{
				Username: username,
				Email:    email,
				PwHash:   string(hashed),
			}

            db := c.MustGet("DB").(*gorm.DB)

            if error := db.Create(user).Error; error != nil {
                fmt.Println("Error creating user:", error)
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
                return
            }

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