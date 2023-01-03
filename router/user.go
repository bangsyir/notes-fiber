package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func UserRouter(app *fiber.App) {
	api := app.Group("/api", logger.New())

	// protected route
	api.Get("/user", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON("success")
	})
}
