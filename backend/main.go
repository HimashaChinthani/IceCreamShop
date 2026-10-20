package main

import (
	"log"
	"os"
	"strings"

	"icecream-shop-api/database"
	"icecream-shop-api/handlers"
	authMiddleware "icecream-shop-api/middleware"
	"icecream-shop-api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load local environment variables.
	// In production, configure environment variables in the host.
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found; using system environment variables")
	}

	secret := strings.TrimSpace(os.Getenv("JWT_SECRET"))

	if len([]byte(secret)) < 32 {
		log.Fatal("JWT_SECRET must contain at least 32 bytes")
	}

	// Connect to SQLite.
	db := database.Connect()

	// Create the Fiber application.
	app := fiber.New()

	// Public endpoints.
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

	// Authentication endpoints.
	auth := handlers.NewAuthHandler(db)
	app.Post("/api/auth/login", auth.Login)
	routes.RegisterInventoryRoutes(app, []byte(secret))

	// Protected endpoints.
	protected := app.Group(
		"/api",
		authMiddleware.AuthRequired([]byte(secret)),
	)

	protected.Get("/auth/me", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Authentication successful",
			"user_id": c.Locals("userID"),
			"role":    c.Locals("role"),
		})
	})

	log.Fatal(app.Listen(":4000"))
}
