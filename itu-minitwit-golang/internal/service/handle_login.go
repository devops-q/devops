package service

import (
	"itu-minitwit/internal/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleLogin(c *gin.Context, username string, password string) (string, error) {

	user, err := utils.FindUserWithName(c, username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "Invalid username", err
		} else {
			return "Database error", err
		}
	} else if hashedPS, error := utils.CheckPassword(user.PwHash, password); !hashedPS {
		return "Invalid password", error
	} else {
		// Login successful, set session
		session := sessions.Default(c)
		session.Set("user_id", user.ID)


		return "", nil
	}
}
