package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"icecream-shop-api/database"
	"icecream-shop-api/models"

	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupHandlerTest(t *testing.T) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.AutoMigrate(&models.Inventory{}, &models.Customer{}, &models.User{}); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}
	database.DB = db
}

func requestJSON(t *testing.T, method, path, body string, handler fiber.Handler) *http.Response {
	t.Helper()
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", uint(1))
		return c.Next()
	})
	app.Add(method, path, handler)
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	return resp
}

func TestInventorySearch(t *testing.T) {
	setupHandlerTest(t)
	database.DB.Create(&models.Inventory{Name: "Chocolate", SKU: "CHO-1", Quantity: 4, Price: 3.5})
	database.DB.Create(&models.Inventory{Name: "Vanilla", SKU: "VAN-1", Quantity: 2, Price: 3})

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", uint(1))
		return c.Next()
	})
	app.Get("/api/inventory", GetInventory)
	req := httptest.NewRequest("GET", "/api/inventory?search=choc", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	if !strings.Contains(string(body), "Chocolate") || strings.Contains(string(body), "Vanilla") {
		t.Fatalf("search returned unexpected results: %s", body)
	}
}

func TestCustomerCRUDValidation(t *testing.T) {
	setupHandlerTest(t)
	resp := requestJSON(t, "POST", "/api/customers", `{"email":"test@example.com"}`, CreateCustomer)
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected 400 for missing customer name, got %d", resp.StatusCode)
	}

	resp = requestJSON(t, "POST", "/api/customers", `{"name":"Test Customer","email":"TEST@example.com"}`, CreateCustomer)
	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected 201 for valid customer, got %d", resp.StatusCode)
	}
}

func TestLoginAuthenticatesSeededUser(t *testing.T) {
	setupHandlerTest(t)
	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	database.DB.Create(&models.User{
		FullName:     "Test Owner",
		Email:        "owner@example.com",
		PasswordHash: string(hash),
		Role:         "OWNER",
	})

	handler := &AuthHandler{DB: database.DB, JWTSecret: []byte("01234567890123456789012345678901")}
	resp := requestJSON(t, "POST", "/api/auth/login", `{"email":"OWNER@example.com","password":"correct-password"}`, handler.Login)
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200 for valid credentials, got %d", resp.StatusCode)
	}

	resp = requestJSON(t, "POST", "/api/auth/login", `{"email":"owner@example.com","password":"wrong-password"}`, handler.Login)
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected 401 for invalid credentials, got %d", resp.StatusCode)
	}
}

func TestLogoutReturnsNoContent(t *testing.T) {
	app := fiber.New()
	app.Post("/api/auth/logout", (&AuthHandler{}).Logout)
	resp, err := app.Test(httptest.NewRequest("POST", "/api/auth/logout", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}
}
