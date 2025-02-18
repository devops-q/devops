package handlers

import (
	"gorm.io/gorm"
	"itu-minitwit/internal/models"
	"itu-minitwit/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UnfollowHandler(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	username := c.Param("username")

	value, userLoggedIn := c.Get("user")
	if !userLoggedIn {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	currentUser := value.(*models.User)

	var userToUnfollow models.User
	if err := db.Where("username = ?", username).First(&userToUnfollow).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err := db.Model(currentUser).Association("Following").Delete(&userToUnfollow); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	utils.SetFlashes(c, "You are no longer following "+username)

	c.Redirect(http.StatusTemporaryRedirect, "/"+username)
	return
}
