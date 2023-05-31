package models

import (
	"time"
)

type Meeting struct {
	ID        uint      `json:"id" gorm:"primary_key,auto_increment"`
	InviterID uint      `json:"inviter_id" required:"true"`
	Inviter   User      `json:"inviter,omitempty" gorm:"foreignkey:InviterID"`
	EventID   uint      `json:"event_id" required:"true"`
	Event     Event     `json:"event,omitempty" gorm:"foreignkey:EventID"`
	StartTime time.Time `json:"start_time" required:"true"`
	EndTime   time.Time `json:"end_time" required:"true"`
	Status    string    `json:"status" required:"true"`
}

// struct for posting meeting request
type MeetingPostRequest struct {
	EventID   uint      `json:"event_id" required:"true"`
	InviterID uint      `json:"inviter_id" required:"true"`
	StartTime time.Time `json:"start_time" required:"true"`
	EndTime   time.Time `json:"end_time" required:"true"`
	Invitees  []uint    `json:"invitees" required:"false"`
}

// struct for getting meeting response
type MeetingResponse struct {
	ID        uint      `json:"id"`
	InviterID uint      `json:"inviter_id"`
	EventID   uint      `json:"event_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
}

func GetMeetingByID(id uint) (Meeting, error) {
	var meeting Meeting
	if err := DB.Where("id = ?", id).First(&meeting).Error; err != nil {
		return meeting, err
	}
	return meeting, nil
}

func ConvertMeetingRequest(meetingPostRequest MeetingPostRequest) Meeting {
	return Meeting{
		EventID:   meetingPostRequest.EventID,
		InviterID: meetingPostRequest.InviterID,
		StartTime: meetingPostRequest.StartTime,
		EndTime:   meetingPostRequest.EndTime,
		Status:    "pending",
	}
}
