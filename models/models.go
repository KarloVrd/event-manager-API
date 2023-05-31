package models

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "github.com/jinzhu/gorm/dialects/mysql" // import mysql driver
)

const DB_NAME = "event-manager.db"
const TEST_DB_NAME = "test.db"

var DB *gorm.DB

func Setup() {
	var err error

	DB, err = gorm.Open(sqlite.Open(DB_NAME), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Disable logger
	DB.Logger.LogMode(logger.Error)

	// Create table if not exists
	DB.AutoMigrate(&Event{}, &User{}, &UserEvent{}, &Organization{}, &Invitation{}, &Meeting{})

}

func TestSetup() {
	var err error

	// Remove old database if exists
	if err = os.Remove(TEST_DB_NAME); err == nil {
		fmt.Println("Removed old database")
	}

	// Create new database
	DB, err = gorm.Open(sqlite.Open(TEST_DB_NAME), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Disable logger
	DB.Logger.LogMode(logger.Error)

	// Create tables if not exists
	DB.AutoMigrate(&Event{}, &User{}, &UserEvent{}, &Organization{}, &Invitation{}, &Meeting{})

}
