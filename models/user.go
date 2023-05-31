package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql" // import mysql driver
)

// User struct
type User struct {
	ID             uint         `json:"id" gorm:"primary_key,auto_increment"`
	Firstname      string       `json:"firstname" required:"true"`
	Lastname       string       `json:"lastname" required:"true"`
	Email          string       `json:"email" required:"true" gorm:"unique"`
	Password       string       `json:"password" required:"true"`
	OrganizationID uint         `json:"organization_id" required:"true"`
	Organization   Organization `json:"organization" required:"false" gorm:"foreignkey:OrganizationID"`
}

// GetUserByEmail
func GetUserByEmail(email string) (User, error) {
	var user User
	if err := DB.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

// GetUserByID
func GetUserByID(id uint) (User, error) {
	var user User
	if err := DB.Where("id = ?", id).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

// UserIdExists
func UserIdExists(id uint) bool {
	var user User
	if err := DB.Where("id = ?", id).First(&user).Error; err != nil {
		return false
	}
	return true
}
