package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/bangsyir/notes/database"
	"github.com/bangsyir/notes/helper"
	"github.com/bangsyir/notes/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name     string
	Username string
	Email    string
	Password string
}

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

func Login(c *fiber.Ctx) error {
	var login struct {
		Email    string
		Password string
	}
	if err := c.BodyParser(&login); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	user := models.User{}
	database.DB.Where("email = ?", login.Email).First(&user)

	if user.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user not found"})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to create token"})
	}

	c.Cookie(&fiber.Cookie{Name: "authorization", Value: tokenString, SameSite: "lax", HTTPOnly: true})
	return c.Status(fiber.StatusOK).JSON("login success")
}
