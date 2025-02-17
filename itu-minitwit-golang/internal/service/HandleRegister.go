package service

import (
	"itu-minitwit/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandleRegister(c *gin.Context, username string, email string, password string, password2 string) string {
	if username == "" {
		return "You have to eneter a username"
	} else if email == "" || !strings.Contains(email, "@") {
		return "You have to enter a valid email address"
	} else if password == "" {
		return "You have to enter a password"
	} else if password2 != password {
		return "The two passwords do not match"
	} else if utils.UserExists(c, username) {
		return "The username is already taken"
	} else {
		utils.CreateUser(c, username, email, password)
		c.Set("flash", "You were successfully registered and can log in now")
		c.Redirect(http.StatusFound, "/login")
		return ""
		

	}

	

}
