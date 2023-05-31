package models

// UserEvent struct
type UserEvent struct {
	UserID  uint  `json:"user_id" required:"true" gorm:"primary_key"`
	EventID uint  `json:"event_id" required:"true" gorm:"primary_key"`
	User    User  `json:"user,omitempty" required:"false" gorm:"foreignkey:UserID"`
	Event   Event `json:"event,omitempty" required:"false" gorm:"foreignkey:EventID"`
}

// UserEventExists func
func UserEventExists(user_id uint, event_id uint) bool {
	var userEvent UserEvent
	if err := DB.Where("user_id = ? AND event_id = ?", user_id, event_id).First(&userEvent).Error; err != nil {
		return false
	}
	return true
}
