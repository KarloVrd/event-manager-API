package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql" // import mysql driver
)

type Event struct {
	ID          uint   `json:"id" gorm:"primary_key,auto_increment"`
	Title       string `json:"title" required:"true" gorm:"unique"`
	Description string `json:"description" required:"false"`
	Location    string `json:"location" required:"false"`
}

// EventIdExists func
// Returns true if event exists, false otherwise
func EventIdExists(id uint) bool {
	var event Event
	if err := DB.Where("id = ?", id).First(&event).Error; err != nil {
		return false
	}
	return true
}

func EventTitleExists(title string) bool {
	var event Event
	if err := DB.Where("title = ?", title).First(&event).Error; err != nil {
		return false
	}
	return true
}

func GetEventByID(id uint) (Event, error) {
	var event Event
	if err := DB.Where("id = ?", id).First(&event).Error; err != nil {
		return event, err
	}
	return event, nil
}
