package service

import (
	"itu-minitwit/internal/utils"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleLogin(c *gin.Context, username string, password string) string {

	user, err := utils.FindUserWithName(c, username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "Invalid username"
		} else {
			return "Database error"
		}
	} else if !utils.CheckPassword(user.PwHash, password) {
		return "Invalid password"
	} else {
		// Login successful, set session
		session := sessions.Default(c)
		session.Set("user_id", user.ID)

		c.Set("flash", "You were successfully logged in")
		c.Redirect(http.StatusFound, "/")
		return ""
	}
}
