package service

import (
	"errors"
	"itu-minitwit/internal/utils"
	"itu-minitwit/pkg/logger"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleLogin(c *gin.Context, username string, password string) (string, error) {
	log := logger.Init()
	user, err := utils.FindUserWithName(c, username)
	if err != nil {
		log.Error("[HandleLogin] FindUserWithName", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "Invalid username", err
		} else {
			return "Database error", err
		}
	} else if hashedPS, err2 := utils.CheckPassword(user.PwHash, password); !hashedPS {
		log.Error("[HandleLogin] CheckPassword error", log)
		return "Invalid password", err2
	} else {
		log.Info("[HandleLogin], successful login")
		session := sessions.Default(c)
		session.Set("user_id", user.ID)

		return "", nil
	}
}
