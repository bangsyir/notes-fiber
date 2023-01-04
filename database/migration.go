package database

import "github.com/bangsyir/notes/models"

func DbMigration() {
	DB.AutoMigrate(&models.Post{}, &models.User{})
}
