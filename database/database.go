package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{QueryFields: true})
	if err != nil {
		panic("failed to connect database")
	}
}
