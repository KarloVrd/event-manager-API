package handlers

import (
	"net/http"
	"strconv"

	"github.com/Karlovrd/event-manager-API/models"

	"github.com/gin-gonic/gin"
)

func CreateInvitation(c *gin.Context) {
	var invitationRequest models.InvitationPostRequest

	if err := c.ShouldBindJSON(&invitationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get inviter's event
	meeting, err := models.GetMeetingByID(invitationRequest.MeetingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Meeting not found"})
		return
	}

	// test if meeting is pending
	if meeting.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Meeting not pending"})
		return
	}

	// test if invitee is in meeting event
	if !models.UserEventExists(invitationRequest.InviteeID, meeting.EventID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invitee not in event"})
		return
	}

	// create invitation
	invitation := models.Invitation{
		MeetingID: meeting.ID,
		InviteeID: invitationRequest.InviteeID,
		InvStatus: "pending",
	}
	models.DB.Create(&invitation)

	// return invitation
	invitationResponse := models.InvitationResponse{
		ID:        invitation.ID,
		MeetingID: invitation.MeetingID,
		InvStatus: invitation.InvStatus,
	}

	c.JSON(http.StatusOK, gin.H{"data": invitationResponse, "message": "Invitation created"})
}

func UpdateInvitation(c *gin.Context) {
	var invitationRequest models.InvitationPutRequest

	if err := c.ShouldBindJSON(&invitationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// parse user id
	u64, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
	}
	user_id := uint(u64)

	// parse invitation id
	u64, err = strconv.ParseUint(c.Param("invitation_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invitation id"})
	}
	inv_id := uint(u64)

	// get invitation
	invitation, err := models.GetInvitationByID(inv_id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invitation not found"})
		return
	}

	// test if user is invitee
	if invitation.InviteeID != user_id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not invitee"})
		return
	}

	// status must be accepted or rejected
	if invitationRequest.InvStatus != "accepted" && invitationRequest.InvStatus != "rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invitation status must be accepted or rejected"})
		return
	}

	// test if invitation is pending
	if invitation.InvStatus != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invitation already accepted/rejected"})
		return
	}

	// test if user is free at that time
	if invitationRequest.InvStatus == "accepted" &&
		!UserIsFree(invitation.InviteeID, invitation.Meeting.StartTime, invitation.Meeting.EndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User has other meetings at that time"})
		return
	}

	// update invitation
	invitation.InvStatus = invitationRequest.InvStatus
	models.DB.Save(&invitation)

	// update meeting status if needed
	CheckMeetingStatus(invitation.MeetingID)

	invitationResponse := models.InvitationResponse{
		ID:        invitation.ID,
		MeetingID: invitation.MeetingID,
		InvStatus: invitation.InvStatus,
	}

	c.JSON(http.StatusOK, gin.H{"invitation": invitationResponse, "message": "Invitation updated"})
}

func GetInvitation(c *gin.Context) {
	// parse user
	u64, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
	}
	user_id := uint(u64)

	// test if user exists
	if _, err := models.GetUserByID(user_id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// parse invitation
	u64, err = strconv.ParseUint(c.Param("invitation_id"), 10, 64)
	inv_id := uint(u64)

	// check if inv_param exists
	if err == nil {
		// get single invitation
		invitation, err := models.GetInvitationByID(inv_id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invitation not found"})
			return
		}

		invitationResponse := models.InvitationResponse{
			ID:        invitation.ID,
			InviteeID: invitation.InviteeID,
			MeetingID: invitation.MeetingID,
			InvStatus: invitation.InvStatus,
		}

		c.JSON(http.StatusOK, gin.H{"invitation": invitationResponse, "message": "Invitation retrieved"})
		return
	}

	// get all invitations
	invitations, _ := models.GetInvitationsByUserID(user_id)

	// convert to InvitationResponse
	var invitationsResponse []models.InvitationResponse
	for _, invitation := range invitations {
		invitationResponse := models.InvitationResponse{
			ID:        invitation.ID,
			InviteeID: invitation.InviteeID,
			MeetingID: invitation.MeetingID,
			InvStatus: invitation.InvStatus,
		}
		invitationsResponse = append(invitationsResponse, invitationResponse)
	}
	c.JSON(http.StatusOK, gin.H{"invitations": invitationsResponse, "message": "Invitation retrieved"})
}
