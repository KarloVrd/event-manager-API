package main

import (
	"github.com/Karlovrd/event-manager-API/models"

	"github.com/Karlovrd/event-manager-API/handlers"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/users", handlers.CreateUser)
	r.POST("/events", handlers.CreateEvent)
	r.POST("/events/:event_id/users", handlers.CreateUserEvent)

	r.POST("/meetings", handlers.CreateMeeting)
	r.GET("users/:user_id/meetings", handlers.GetAllMeetings)

	r.GET("users/:user_id/invitations", handlers.GetInvitation)
	r.GET("users/:user_id/invitations/:invitation_id", handlers.GetInvitation)
	r.POST("users/:user_id/invitations", handlers.CreateInvitation)
	r.PUT("users/:user_id/invitations/:invitation_id", handlers.UpdateInvitation)

	return r
}

func main() {
	models.Setup() // connect to database
	r := setupRouter()
	r.Run(":8080") // listen and serve on
}
