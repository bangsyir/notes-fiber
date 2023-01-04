package handler

import (
	"net/http"
	"time"

	"github.com/bangsyir/notes/database"
	"github.com/bangsyir/notes/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type APIUser struct {
	ID       uint
	Name     string
	Username string
	Email    string
}

type APIpost struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Desc      string    `json:"desc"`
	CreatedAt time.Time `json:"created_at"`
	User      APIUser   `json:"user"`
}

func CreateResponsePost(post models.Post) APIpost {
	return APIpost{ID: post.ID, Title: post.Title, Desc: post.Desc, CreatedAt: post.CreatedAt, User: APIUser{
		ID:       post.User.ID,
		Name:     post.User.Name,
		Username: post.User.Username,
		Email:    post.User.Email,
	}}
}

func CreatePost(c *fiber.Ctx) error {
	post := new(models.Post)
	if err := c.BodyParser(post); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	database.DB.Create(&post)

	type Response struct {
		ID        uint
		Title     string
		Desc      string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	res := Response{ID: post.ID, Title: post.Title, Desc: post.Desc, CreatedAt: post.CreatedAt, UpdatedAt: post.UpdatedAt}

	return c.Status(http.StatusCreated).JSON(res)
}

func GetPost(c *fiber.Ctx) error {
	postId := c.Params("id")

	var post models.Post

	database.DB.Model(&models.Post{}).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "username", "email")
	}).Where("id = ?", postId).Find(&post)
	if post.ID == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": "post not found"})
	}
	response := CreateResponsePost(post)

	return c.JSON(response)
}
