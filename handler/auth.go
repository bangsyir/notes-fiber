package handler

import (
	"net/http"

	"github.com/bangsyir/notes/database"
	"github.com/bangsyir/notes/helper"
	"github.com/bangsyir/notes/models"
	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	// get data from request body
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	errors := models.ValidateStruct(*user)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	// check email is available
	userFind := database.DB.Where("email = ?", user.Email).First(&user)
	if userFind.RowsAffected > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "email already registered."})
	}
	hash, err := helper.GeneratePassword(user.Password)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Failed to hash password"})
	}
	newUser := models.User{Name: user.Name, Username: user.Username, Email: user.Email, Password: hash}
	// save data to database
	database.DB.Create(&newUser)
	// response data
	return c.JSON(newUser)
}
