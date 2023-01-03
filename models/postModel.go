package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title  string
	Desc   string
	UserID uint
	User   User
}

type PostApi struct {
	ID     uint
	Title  string
	Desc   string
	UserID uint
}