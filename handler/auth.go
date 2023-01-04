package handler

import (
	"errors"
	"net/http"
	"net/mail"
	"os"
	"time"

	"github.com/bangsyir/notes/database"
	"github.com/bangsyir/notes/helper"
	"github.com/bangsyir/notes/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type User struct {
	ID       uint
	Name     string
	Username string
	Email    string
	Password string
}

func GetUserByEmail(e string) (*models.User, error) {
	db := database.DB
	var user models.User
	if err := db.Where(&models.User{Email: e}).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func GetUserByUsername(u string) (*models.User, error) {
	db := database.DB
	var user models.User
	if err := db.Where(&models.User{Username: u}).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func Register(c *fiber.Ctx) error {
	// get data from request body
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
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
	type LoginInput struct {
		Identiy  string `json:"identity"`
		Password string `json:"password"`
	}

	input := new(LoginInput)
	var ud User

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
	}

	identity := input.Identiy
	password := input.Password

	user, email, err := new(models.User), new(models.User), *new(error)
	if valid(identity) {
		email, err = GetUserByEmail(identity)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Error on email", "data": err})
		}
	} else {
		user, err = GetUserByUsername(identity)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Error on username", "data": err})
		}
	}

	if email == nil && user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "User not found", "data": err})
	}

	if email != nil {
		ud = User{
			ID:       email.ID,
			Name:     email.Name,
			Username: email.Username,
			Password: email.Password,
		}
	}
	if user != nil {
		ud = User{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Password: user.Password,
		}
	}

	CheckPassword := helper.CheckPasswordhash(password, ud.Password)
	if !CheckPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid password", "data": nil})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to create token"})
	}

	// c.Cookie(&fiber.Cookie{Name: "authorization", Value: tokenString, SameSite: "lax", HTTPOnly: true})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": nil, "data": fiber.Map{"accessToken": tokenString}})
}
