package main

import (
	"encoding/json"
	"log"

	"github.com/bangsyir/notes/database"
	"github.com/bangsyir/notes/router"
	"github.com/gofiber/fiber/v2"
)

func init() {
	database.ConnectToDatabase()
	database.DbMigration()
}

func main() {
	// createing new fiber app
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// fiber routing
	router.SetupRoutes(app)

	// run fiber app
	log.Fatal(app.Listen(":8000"))
}
