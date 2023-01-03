package router

import (
	"github.com/bangsyir/notes/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func AuthRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())
	api.Post("/register", handler.Register)
}
