package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Karlovrd/event-manager-API/models"

	"github.com/gin-gonic/gin"
)

// checks all invitations for a meeting and sets meeting status to accepted if all are accepted,
// rejected if any are rejected, pending otherwise
func CheckMeetingStatus(meetingID uint) {
	// get all invitations for meeting and check if all are accepted
	var invitations []models.Invitation
	models.DB.Model(&models.Invitation{}).Where("meeting_id = ?", meetingID).Find(&invitations)

	allAccepted := true
	anyRejected := false

	// iterate through invitations and check status
	for _, invitation := range invitations {
		if invitation.InvStatus == "rejected" {
			anyRejected = true
			allAccepted = false
		} else if invitation.InvStatus == "pending" {
			allAccepted = false
		}
	}
	// get meeting
	meeting, _ := models.GetMeetingByID(meetingID)

	// set meeting status
	if allAccepted {
		meeting.Status = "scheduled"
	} else if anyRejected {
		meeting.Status = "cancelled"
	} else {
		return
	}
	models.DB.Save(&meeting)
}

func UserIsFree(user_id uint, start_time time.Time, end_time time.Time) bool {
	// check if user has created a meeting at that time
	var otherMeetings []models.Meeting
	models.DB.Model(&models.Meeting{}).Where("inviter_id = ? AND status <> 'rejected' AND NOT (start_time > ? OR end_time < ?)",
		user_id,
		end_time,
		start_time).Find(&otherMeetings)

	// test if there are any meetings
	if len(otherMeetings) > 0 {
		return false
	}

	// check if invitee already has an accepted invitation at that time
	var invitations []models.Invitation
	models.DB.Model(&models.Invitation{}).Where("invitee_id = ? AND status = 'accepted'", user_id).
		Joins("JOIN meetings ON invitations.meeting_id = meetings.id").
		Where("NOT (meetings.start_time > ? OR meetings.end_time < ?)", end_time, start_time).
		Find(&invitations)

	// test if there are any meetings
	return len(invitations) == 0
}

func CreateMeeting(c *gin.Context) {
	var meetingRequest models.MeetingPostRequest

	if err := c.ShouldBindJSON(&meetingRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	meeting := models.ConvertMeetingRequest(meetingRequest)

	// test if inviter is in event
	if !models.UserEventExists(meeting.InviterID, meeting.EventID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Inviter not in event"})
		return
	}

	// test if inviter is free
	if !UserIsFree(meeting.InviterID, meeting.StartTime, meeting.EndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User has other meetings/accepted invitations at that time"})
		return
	}

	if meetingRequest.Invitees == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Meeting created", "meeting": meeting})
	}

	users := meetingRequest.Invitees
	// test if users are in the save event as inviter
	for _, user_id := range users {
		if !models.UserEventExists(user_id, meeting.EventID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invitee not in event"})
			return
		}
	}

	models.DB.Create(&meeting)

	// create invitations for each user
	for _, user_id := range users {
		invitation := models.Invitation{
			InviteeID: user_id,
			MeetingID: meeting.ID,
			InvStatus: "pending",
		}
		models.DB.Create(&invitation)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Meeting created, invitations sent", "meeting": meeting})
}

func GetAllMeetings(c *gin.Context) {
	var meetings []models.MeetingResponse
	var err error

	// parse user id
	u64, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
	}
	user_id := uint(u64)

	// get all meetings that are created by user
	if err = models.DB.Model(&models.Meeting{}).Where("inviter_id = ?", user_id).Find(&meetings).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No meetings found"})
		return
	}

	// get all meetings that user is invited to and accepted
	var meetingsInvitedTo []models.MeetingResponse
	models.DB.Model(&models.Invitation{}).Where("invitee_id = ? AND inv_status = 'accepted'", user_id).Joins("JOIN meetings ON invitations.meeting_id = meetings.id").Find(&meetingsInvitedTo)

	// append meetingsInvitedTo to meetings
	meetings = append(meetings, meetingsInvitedTo...)

	c.JSON(http.StatusOK, gin.H{"meetings": meetings})
}
