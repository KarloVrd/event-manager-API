package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Karlovrd/event-manager-API/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type UserResp struct {
	Message string `json:"message"`
	User    models.User
}

type EventResp struct {
	Message string `json:"message"`
	Event   models.Event
}

var r *gin.Engine

func TestMain(m *testing.M) {
	models.TestSetup() // connect to test database
	r = setupRouter()  // setup router

	m.Run()

	teardown()
}

func teardown() {
	// models.DB.Migrator().DropTable(&models.Event{})
	// models.DB.Migrator().DropTable(&models.User{})
	// models.DB.Migrator().DropTable(&models.UserEvent{})
	// models.DB.Migrator().DropTable(&models.Organization{})
	// models.DB.Migrator().DropTable(&models.Invitation{})
	// models.DB.Migrator().DropTable(&models.Meeting{})

	// close database connection
	sqlDB, _ := models.DB.DB()
	sqlDB.Close()
}

func TestCreateUser(t *testing.T) {
	// create user
	user := models.User{
		Firstname: "Test",
		Lastname:  "User",
		Email:     "test@example.com",
		Password:  "password",
	}
	jsonUser, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonUser))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	// parse response
	var userResp UserResp
	json.Unmarshal(resp.Body.Bytes(), &userResp)
	respUser := userResp.User

	// test response
	assert.Equal(t, user.Email, respUser.Email)
	assert.Equal(t, uint(1), uint(respUser.ID))
	assert.Equal(t, "", respUser.Password)

	// test database
	var dbUser models.User
	dbUser, _ = models.GetUserByEmail(user.Email)
	assert.Equal(t, user.Firstname, dbUser.Firstname)
	assert.Equal(t, user.Lastname, dbUser.Lastname)
	assert.Equal(t, user.Email, dbUser.Email)
	assert.Equal(t, user.Password, dbUser.Password)
}

func TestCreateEvent(t *testing.T) {
	// create event
	event := models.Event{
		Title:       "Test Event",
		Description: "Test Event Description",
	}
	jsonEvent, _ := json.Marshal(event)
	req, _ := http.NewRequest("POST", "/events", bytes.NewBuffer(jsonEvent))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	// parse response
	var eventResp EventResp
	json.Unmarshal(resp.Body.Bytes(), &eventResp)
	respEvent := eventResp.Event
	// test response
	assert.Equal(t, event.Title, respEvent.Title)
	assert.Equal(t, uint(1), uint(respEvent.ID))

	// test database
	var dbEvent models.Event
	dbEvent, _ = models.GetEventByID(respEvent.ID)
	assert.Equal(t, event.Title, dbEvent.Title)
	assert.Equal(t, event.Description, dbEvent.Description)
}

func TestCreateMeeting(t *testing.T) {
	user1 := models.User{
		Firstname: "Test1",
		Lastname:  "User",
		Email:     "test1@example.com",
		Password:  "password",
	}
	user2 := models.User{
		Firstname: "Test2",
		Lastname:  "User",
		Email:     "test2@example.com",
		Password:  "password",
	}
	user3 := models.User{
		Firstname: "Test3",
		Lastname:  "User",
		Email:     "test3@example.com",
		Password:  "password",
	}

	// create users
	models.DB.Create(&user1)
	models.DB.Create(&user2)
	models.DB.Create(&user3)

	// create event
	event := models.Event{
		Title:       "Test Event 2",
		Description: "Test Event Description",
	}
	models.DB.Create(&event)

	// create user events
	userEvent1 := models.UserEvent{
		UserID:  user1.ID,
		EventID: event.ID,
	}
	models.DB.Create(&userEvent1)

	userEvent2 := models.UserEvent{
		UserID:  user2.ID,
		EventID: event.ID,
	}
	models.DB.Create(&userEvent2)

	userEvent3 := models.UserEvent{
		UserID:  user3.ID,
		EventID: event.ID,
	}
	models.DB.Create(&userEvent3)

	meetingRequest := models.MeetingPostRequest{
		EventID:   event.ID,
		InviterID: user1.ID,
		StartTime: time.Now().Add(time.Hour * 2).Truncate(time.Second),
		EndTime:   time.Now().Add(time.Hour * 3).Truncate(time.Second),
		Invitees:  []uint{user2.ID, user3.ID},
	}

	jsonMeetingRequest, _ := json.Marshal(meetingRequest)
	req, _ := http.NewRequest("POST", "/meetings", bytes.NewBuffer(jsonMeetingRequest))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	type Response struct {
		Message string         `json:"message"`
		Meeting models.Meeting `json:"meeting"`
	}

	// parse response
	var meetingResp Response
	json.Unmarshal(resp.Body.Bytes(), &meetingResp)
	respMeeting := meetingResp.Meeting

	// test response
	assert.Equal(t, meetingRequest.EventID, respMeeting.EventID)
	assert.Equal(t, meetingRequest.InviterID, respMeeting.InviterID)
	assert.Equal(t, meetingRequest.StartTime, respMeeting.StartTime)
	assert.Equal(t, meetingRequest.EndTime, respMeeting.EndTime)

	// test database
	var dbMeeting models.Meeting
	dbMeeting, _ = models.GetMeetingByID(respMeeting.ID)
	assert.Equal(t, meetingRequest.EventID, dbMeeting.EventID)
	assert.Equal(t, meetingRequest.InviterID, dbMeeting.InviterID)
	assert.Equal(t, meetingRequest.StartTime, dbMeeting.StartTime.Local())
	assert.Equal(t, meetingRequest.EndTime, dbMeeting.EndTime.Local())

	// test database for invitees
	var dbInvitations []models.Invitation
	models.DB.Where("meeting_id = ?", dbMeeting.ID).Find(&dbInvitations)
	assert.Equal(t, 2, len(dbInvitations))
	assert.Equal(t, user2.ID, dbInvitations[0].InviteeID)
	assert.Equal(t, user3.ID, dbInvitations[1].InviteeID)
}

func TestUpdateInvitation(t *testing.T) {
	invitationUpdate := models.InvitationPutRequest{
		InvStatus: "accepted",
	}
	jsonInvitationUpdate, _ := json.Marshal(invitationUpdate)

	req, _ := http.NewRequest("PUT", "/users/3/invitations/1", bytes.NewBuffer(jsonInvitationUpdate))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)

	type Response struct {
		Message    string            `json:"message"`
		Invitation models.Invitation `json:"invitation"`
	}

	// parse response
	var invitationResp Response
	json.Unmarshal(resp.Body.Bytes(), &invitationResp)
	respInvitation := invitationResp.Invitation

	// test response
	assert.Equal(t, invitationUpdate.InvStatus, respInvitation.InvStatus)

	// test database
	var dbInvitation models.Invitation
	dbInvitation, _ = models.GetInvitationByID(1)
	assert.Equal(t, invitationUpdate.InvStatus, dbInvitation.InvStatus)
}

// update 2. invitation so meeting is scheduled
func TestUpdateMeeting4(t *testing.T) {
	invitationUpdate := models.InvitationPutRequest{
		InvStatus: "accepted",
	}
	jsonInvitationUpdate, _ := json.Marshal(invitationUpdate)

	req, _ := http.NewRequest("PUT", "/users/4/invitations/2", bytes.NewBuffer(jsonInvitationUpdate))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	// test database
	var dbMeeting models.Meeting
	dbMeeting, _ = models.GetMeetingByID(1)
	assert.Equal(t, "scheduled", dbMeeting.Status)
}

func TestGetInvitation(t *testing.T) {
	req, _ := http.NewRequest("GET", "/users/3/invitations/1", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	type Response struct {
		Message            string            `json:"message"`
		InvitationResponse models.Invitation `json:"invitation"`
	}

	// parse response
	var response Response
	json.Unmarshal(resp.Body.Bytes(), &response)
	invitation := response.InvitationResponse

	// test response
	assert.Equal(t, "accepted", invitation.InvStatus)
	assert.Equal(t, uint(1), invitation.ID)
	assert.Equal(t, uint(1), invitation.MeetingID)
	assert.Equal(t, uint(3), invitation.InviteeID)
}

func TestGetAllInvitations(t *testing.T) {
	req, _ := http.NewRequest("GET", "/users/3/invitations", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	type Response struct {
		Message             string                      `json:"message"`
		InvitationsResponse []models.InvitationResponse `json:"invitations"`
	}

	// parse response
	var response Response
	json.Unmarshal(resp.Body.Bytes(), &response)
	invitations := response.InvitationsResponse

	// test response
	assert.Equal(t, 1, len(invitations))
	assert.Equal(t, uint(1), invitations[0].ID)
}

func TestGetALLMeetings(t *testing.T) {
	req, _ := http.NewRequest("GET", "/users/2/meetings", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	type Response struct {
		Message         string                   `json:"message"`
		MeetingResponse []models.MeetingResponse `json:"meetings"`
	}

	// parse response
	var response Response
	json.Unmarshal(resp.Body.Bytes(), &response)
	meetings := response.MeetingResponse

	// test response
	assert.Equal(t, 1, len(meetings))
	assert.Equal(t, uint(1), meetings[0].ID)
}
