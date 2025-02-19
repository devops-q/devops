package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"itu-minitwit/internal/api/json_models"
	"itu-minitwit/internal/service"
	"net/http"
	"strconv"
)

func GetUserFollowers(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)

	limitParam := c.DefaultQuery("no", "100")
	limit, err := strconv.Atoi(limitParam)
	username := c.Param("username")

	if err != nil {
		c.JSON(http.StatusBadRequest, json_models.ErrorResponse{
			Code:         http.StatusInternalServerError,
			ErrorMessage: "Invalid number of messages provided in param\"no\"",
		})
		return
	}

	userId, err := service.GetUserIdByUsername(db, username)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusInternalServerError, json_models.ErrorResponse{
			Code:         http.StatusInternalServerError,
			ErrorMessage: "Error fetching user",
		})
		return
	}

	followers, err := service.GetUserFollows(db, uint(userId), limit)

	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, json_models.ErrorResponse{
			Code:         http.StatusInternalServerError,
			ErrorMessage: "Error fetching followers",
		})
		return
	}

	var followsUsernames = make([]string, 0)
	for _, follower := range followers {
		followsUsernames = append(followsUsernames, follower.Username)
	}

	c.JSON(http.StatusOK, json_models.GetFollowsResponse{
		Follows: followsUsernames,
	})
}
