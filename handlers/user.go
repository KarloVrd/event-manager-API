package handlers

import (
	"net/http"

	"github.com/Karlovrd/event-manager-API/models"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// test if email already exists
	if _, err := models.GetUserByEmail(user.Email); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	models.DB.Create(&user)

	// remove password from response
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": user})
}
