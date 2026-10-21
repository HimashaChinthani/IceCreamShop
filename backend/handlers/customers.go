package handlers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"icecream-shop-api/database"
	"icecream-shop-api/models"
)

func GetCustomers(c *fiber.Ctx) error {
	var customers []models.Customer
	query := c.Query("search")
	db := database.DB

	if strings.TrimSpace(query) != "" {
		like := "%" + strings.TrimSpace(query) + "%"
		db = db.Where("name LIKE ? OR email LIKE ? OR phone LIKE ?", like, like, like)
	}

	if err := db.Order("name asc").Find(&customers).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve customers"})
	}

	return c.JSON(customers)
}

func GetCustomerByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid customer ID"})
	}

	var customer models.Customer
	if err := database.DB.First(&customer, uint(id)).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Customer not found"})
	}

	return c.JSON(customer)
}

func CreateCustomer(c *fiber.Ctx) error {
	var customer models.Customer
	if err := c.BodyParser(&customer); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	normalizeCustomer(&customer)
	if customer.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name is required"})
	}

	if err := database.DB.Create(&customer).Error; err != nil {
		return c.Status(409).JSON(fiber.Map{"error": "Could not create customer; check for duplicate email"})
	}

	return c.Status(201).JSON(customer)
}

func UpdateCustomer(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid customer ID"})
	}

	var customer models.Customer
	if err := database.DB.First(&customer, uint(id)).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Customer not found"})
	}

	var input models.Customer
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	normalizeCustomer(&input)
	if input.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name is required"})
	}

	customer.Name = input.Name
	customer.Email = input.Email
	customer.Phone = input.Phone
	customer.Address = input.Address

	if err := database.DB.Save(&customer).Error; err != nil {
		return c.Status(409).JSON(fiber.Map{"error": "Could not update customer; check for duplicate email"})
	}

	return c.JSON(customer)
}

func DeleteCustomer(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid customer ID"})
	}

	result := database.DB.Delete(&models.Customer{}, uint(id))
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete customer"})
	}
	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Customer not found"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func normalizeCustomer(customer *models.Customer) {
	customer.Name = strings.TrimSpace(customer.Name)
	customer.Email = strings.TrimSpace(strings.ToLower(customer.Email))
	customer.Phone = strings.TrimSpace(customer.Phone)
	customer.Address = strings.TrimSpace(customer.Address)
}
