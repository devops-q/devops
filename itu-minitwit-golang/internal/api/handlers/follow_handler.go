package handlers

import (
	"gorm.io/gorm"
	"itu-minitwit/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UnfollowHandler(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	username := c.Param("username")

	value, userLoggedIn := c.Get("user")
	if !userLoggedIn {
		c.Redirect(http.StatusUnauthorized, "/public")
		return
	}

	currentUser := value.(*models.User)

	var userToUnfollow models.User
	if err := db.Where("username = ?", username).First(&userToUnfollow).Error; err != nil {
		c.Redirect(http.StatusNotFound, "/public")
		return
	}

	if err := db.Model(currentUser).Association("Following").Delete(&userToUnfollow); err != nil {
		c.Redirect(http.StatusInternalServerError, "/public")
		return
	}

	c.SetCookie("flash", "You are no longer following "+username, 3600, "/", "", false, true)
	c.Redirect(http.StatusFound, "/"+username)
}
