package handlers

import (
	"errors"
	"fmt"
	"itu-minitwit/internal/api/json_models"
	"itu-minitwit/internal/models"
	"itu-minitwit/internal/service"
	"itu-minitwit/internal/utils"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func FollowHandler(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	username := c.Param("username")

	userLoggedIn := utils.GetUserFomContext(c)
	if userLoggedIn == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var userToFollow models.User
	if err := db.First(&userToFollow, models.User{Username: username}).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err := db.Model(&userLoggedIn).Association("Following").Append(&userToFollow); err != nil {
		_ = c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	utils.SetFlashes(c, fmt.Sprintf("You are now following \"%s\"", userToFollow.Username))
	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/%s", userToFollow.Username))
}

func UnfollowHandler(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	username := c.Param("username")

	userLoggedIn := utils.GetUserFomContext(c)
	if userLoggedIn == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var userToUnfollow models.User
	if err := db.Where("username = ?", username).First(&userToUnfollow).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err := db.Model(userLoggedIn).Association("Following").Delete(&userToUnfollow); err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	utils.SetFlashes(c, fmt.Sprintf("You are no longer following \"%s\"", username))

	c.Redirect(http.StatusTemporaryRedirect, "/"+username)
	return
}

func GetUserFollowersAPI(c *gin.Context) {
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

func FollowUnfollowHandlerAPI(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	username := c.Param("username")

	whoId, err := service.GetUserIdByUsername(db, username)

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

	var body json_models.FollowUnfollowBody

	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Println("Error binding json", err)
		_ = c.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var action string
	var whomUsername string

	if body.Follow != nil {
		action = "follow"
		whomUsername = *body.Follow
	} else if body.Unfollow != nil {
		action = "unfollow"
		whomUsername = *body.Unfollow
	} else {
		c.JSON(http.StatusBadRequest, json_models.ErrorResponse{
			Code:         http.StatusBadRequest,
			ErrorMessage: "Invalid action provided",
		})
		return
	}

	whomId, err := service.GetUserIdByUsername(db, whomUsername)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	var fllwErr error
	if action == "follow" {
		fllwErr = service.FollowUser(db, uint(whoId), uint(whomId))

	} else {
		fllwErr = service.UnfollowUser(db, uint(whoId), uint(whomId))
	}

	if fllwErr != nil && strings.Contains(fllwErr.Error(), "UNIQUE constraint failed") {
		var message string
		if action == "follow" {
			message = "You are already following this user"
		} else {
			message = "You are not following this user"
		}

		c.JSON(http.StatusBadRequest, json_models.ErrorResponse{
			Code:         http.StatusBadRequest,
			ErrorMessage: message,
		})

		return
	}

	if fllwErr != nil {
		c.Error(fllwErr)
		c.JSON(http.StatusInternalServerError, json_models.ErrorResponse{
			Code:         http.StatusInternalServerError,
			ErrorMessage: "Error following user",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
