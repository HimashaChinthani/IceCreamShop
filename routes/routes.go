package routes

import (
	"github.com/gofiber/fiber/v2"

	"icecream-shop-api/handlers"
	"icecream-shop-api/middleware"
)

func RegisterInventoryRoutes(app *fiber.App, secret []byte) {
	inventory := app.Group(
		"/api/inventory",
		middleware.AuthRequired(secret),
	)

	inventory.Get(
		"/",
		middleware.RequireRoles("Owner", "Manager", "Cashier"),
		handlers.GetInventory,
	)

	inventory.Get(
		"/:id",
		middleware.RequireRoles("Owner", "Manager", "Cashier"),
		handlers.GetInventoryByID,
	)

	inventory.Post(
		"/",
		middleware.RequireRoles("Owner", "Manager"),
		handlers.CreateInventory,
	)

	inventory.Put(
		"/:id",
		middleware.RequireRoles("Owner", "Manager"),
		handlers.UpdateInventory,
	)

	inventory.Delete(
		"/:id",
		middleware.RequireRoles("Owner"),
		handlers.DeleteInventory,
	)
}
