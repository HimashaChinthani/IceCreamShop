
package main

import (
	"log"

	"icecream-shop-api/database"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Connect to SQLite and create the tables.
	database.Connect()

	// Create the Fiber application.
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Ice Cream Shop API is running",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
		})
	})

	log.Fatal(app.Listen(":4000"))
}
