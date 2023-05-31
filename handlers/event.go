package handlers

import (
	"net/http"

	"github.com/Karlovrd/event-manager-API/models"

	"github.com/gin-gonic/gin"
)

func CreateEvent(c *gin.Context) {
	var event models.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// test if event title already exists
	if models.EventTitleExists(event.Title) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event with that title already exists"})
		return
	}

	models.DB.Create(&event)

	c.JSON(http.StatusOK, gin.H{"message": "Event created successfully", "event": event})
}
