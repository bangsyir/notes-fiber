package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title  string `gorm:"not null" json:"title"`
	Desc   string `gorm:"not null" json:"description"`
	UserID uint   `gorm:"not null" json:"user_id"`
	User   User   `gorm:"not null" json:"user"`
}
