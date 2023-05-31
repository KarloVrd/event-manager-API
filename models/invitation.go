package models

type Invitation struct {
	ID        uint    `json:"id" gorm:"primary_key,auto_increment"`
	Meeting   Meeting `json:"meeting,omitempty" gorm:"foreignkey:MeetingID" `
	MeetingID uint    `json:"meeting_id" required:"true"`
	Invitee   User    `json:"invitee,omitempty" gorm:"foreignkey:InviteeID"`
	InviteeID uint    `json:"invitee_id" required:"true"`
	InvStatus string  `json:"inv_status" required:"true"`
}

type InvitationPostRequest struct {
	MeetingID uint `json:"meeting_id" required:"true"`
	InviteeID uint `json:"invitee_id" required:"true"`
}

type InvitationPutRequest struct {
	InvStatus string `json:"inv_status" required:"true"`
}

type InvitationResponse struct {
	ID        uint   `json:"id"`
	InviteeID uint   `json:"invitee_id"`
	MeetingID uint   `json:"meeting_id"`
	InvStatus string `json:"inv_status"`
}

func GetInvitationByID(id uint) (Invitation, error) {
	var invitation Invitation
	err := DB.Where("id = ?", id).First(&invitation).Error
	return invitation, err
}

func GetInvitationsByUserID(id uint) ([]Invitation, error) {
	var invitations []Invitation
	err := DB.Model(invitations).Where("invitee_id = ?", id).Find(&invitations).Error
	return invitations, err
}
