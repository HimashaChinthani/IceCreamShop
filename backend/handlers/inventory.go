package handlers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"icecream-shop-api/database"
	"icecream-shop-api/models"
)

func GetInventory(c *fiber.Ctx) error {
	var items []models.Inventory
	query := strings.TrimSpace(c.Query("search"))
	db := database.DB

	if query != "" {
		db = db.Where("name LIKE ? OR sku LIKE ?", "%"+query+"%", "%"+query+"%")
	}

	if err := db.Order("name asc").Find(&items).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve inventory",
		})
	}

	return c.JSON(items)
}

func GetInventoryByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid inventory ID",
		})
	}

	var item models.Inventory

	if err := database.DB.First(&item, uint(id)).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Inventory item not found",
		})
	}

	return c.JSON(item)
}

func CreateInventory(c *fiber.Ctx) error {
	var item models.Inventory

	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	item.Name = strings.TrimSpace(item.Name)
	item.SKU = strings.TrimSpace(item.SKU)

	if item.Name == "" || item.SKU == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Name and SKU are required",
		})
	}

	if item.Quantity < 0 || item.Price < 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Quantity and price cannot be negative",
		})
	}

	item.ID = 0
	userID, ok := c.Locals("userID").(uint)
	if !ok || userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authentication required"})
	}
	item.CreatedBy = &userID

	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(409).JSON(fiber.Map{
			"error": "Could not create inventory item; check for duplicate SKU",
		})
	}

	return c.Status(201).JSON(item)
}

func UpdateInventory(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid inventory ID",
		})
	}

	var item models.Inventory

	if err := database.DB.First(&item, uint(id)).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Inventory item not found",
		})
	}

	var input struct {
		Name     string  `json:"name"`
		SKU      string  `json:"sku"`
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	input.Name = strings.TrimSpace(input.Name)
	input.SKU = strings.TrimSpace(input.SKU)

	if input.Name == "" || input.SKU == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Name and SKU are required",
		})
	}

	if input.Quantity < 0 || input.Price < 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Quantity and price cannot be negative",
		})
	}

	item.Name = input.Name
	item.SKU = input.SKU
	item.Quantity = input.Quantity
	item.Price = input.Price

	if err := database.DB.Save(&item).Error; err != nil {
		return c.Status(409).JSON(fiber.Map{
			"error": "Could not update inventory item; check for duplicate SKU",
		})
	}

	return c.JSON(item)
}

func DeleteInventory(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid inventory ID",
		})
	}

	result := database.DB.Delete(&models.Inventory{}, uint(id))

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete inventory item",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Inventory item not found",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
