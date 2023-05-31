package handlers

import (
	"net/http"

	"github.com/Karlovrd/event-manager-API/models"

	"github.com/gin-gonic/gin"
)

// CreateUserEvent func
func CreateUserEvent(c *gin.Context) {
	var userEvent models.UserEvent

	if err := c.ShouldBindJSON(&userEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// test if user exists
	if !models.UserIdExists(userEvent.UserID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// test if event exists
	if !models.EventIdExists(userEvent.EventID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event not found"})
		return
	}

	// test if user already in event
	if models.UserEventExists(userEvent.UserID, userEvent.EventID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already in event"})
		return
	}

	// create user event
	models.DB.Create(&userEvent)

	c.JSON(http.StatusOK, gin.H{"message": "User event created successfully"})
}
