package handlers

import (
	"itu-minitwit/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUsersHandler(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)

	var users []models.User
	// Finds all that matches the given array and fills it
	db.Find(&users)

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func CreateUserHandler(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	name := c.Param("name")

	// All not-null fields must be set
	var newUser = &models.User{
		Username: name,
		Email:    "test@email",
		PwHash:   "Test hash",
	}

	// after this uninitialized fields will be initialized by the db
	// i.e. fields ID, CreatedAt and UpdatedAt will have a value on
	// the struct afterwards
	db.Create(newUser)

	c.JSON(http.StatusCreated, newUser)
}

func FindUserWithName(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	name := c.Param("name")

	var user *models.User

	// This will ignore zero/falsy-values, i.e. for flagged if searching for flagged=false
	db.First(&user, models.User{Username: name}) // This throws a non-panic error if nothing was found
	// equal to
	// db.Where(&models.User{Username: name}).First(&user)

	// If we want to search for zero/falsy values do like this
	// (specifying flagged is false in struct is not needed, as it is default)
	// db.Where(&models.Message{Flagged: false}, "Flagged")

	// To check if something was found, we need to look at the PK
	// as that will be populated if something was found
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, nil)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func GetAllUsersWithNonFlaggedMessages(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)

	var users []models.User
	// By default nested objects will not be inclueded
	// to include them, preload them on the query
	db.Preload("Messages", "flagged=0").Find(&users)
	c.JSON(http.StatusOK, users)
}
