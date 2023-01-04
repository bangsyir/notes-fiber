package router

import (
	"github.com/bangsyir/notes/handler"
	"github.com/bangsyir/notes/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	// initiate
	api := app.Group("/api", logger.New())

	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)

	note := api.Group("/notes")
	note.Post("/", handler.CreatePost)
	note.Get("/:id", middleware.Protected(), handler.GetPost)

}
