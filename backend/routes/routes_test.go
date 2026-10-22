package routes

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"icecream-shop-api/database"
	"icecream-shop-api/middleware"
	"icecream-shop-api/models"

	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

const testSecret = "01234567890123456789012345678901"

func routeTestApp(t *testing.T) *fiber.App {
	t.Helper()

	dsn := "file:" + strings.ReplaceAll(t.Name(), "/", "_") + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Inventory{}, &models.Customer{}); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}
	database.DB = db

	app := fiber.New()
	RegisterInventoryRoutes(app, []byte(testSecret))
	RegisterCustomerRoutes(app, []byte(testSecret))
	return app
}

func routeToken(t *testing.T, role string) string {
	t.Helper()
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, middleware.TokenClaims{
		UserID: 1,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "icecream-shop-api",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		},
	})
	signed, err := token.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("sign test token: %v", err)
	}
	return signed
}

func authorizedRequest(t *testing.T, method, path, role, body string) int {
	t.Helper()
	app := routeTestApp(t)
	req := httptest.NewRequest(method, path, nil)
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+routeToken(t, role))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	return resp.StatusCode
}

func TestProtectedRoutesRequireAuthentication(t *testing.T) {
	app := routeTestApp(t)
	for _, path := range []string{"/api/inventory/", "/api/customers/"} {
		resp, err := app.Test(httptest.NewRequest("GET", path, nil))
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("%s: expected 401, got %d", path, resp.StatusCode)
		}
	}
}

func TestInventoryRolePermissions(t *testing.T) {
	tests := []struct {
		name, method, path, role string
		want                     int
	}{
		{"manager cannot delete", "DELETE", "/api/inventory/1", "Manager", fiber.StatusForbidden},
		{"cashier cannot create", "POST", "/api/inventory/", "Cashier", fiber.StatusForbidden},
		{"cashier cannot delete", "DELETE", "/api/inventory/1", "Cashier", fiber.StatusForbidden},
		{"cashier can update", "PUT", "/api/inventory/1", "Cashier", fiber.StatusNotFound},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := `{"name":"Chocolate","category":"Chocolate","sku":"CHO-001","quantity":10,"price":3.5}`
			if got := authorizedRequest(t, test.method, test.path, test.role, body); got != test.want {
				t.Fatalf("expected %d, got %d", test.want, got)
			}
		})
	}
}

func TestCustomerRolePermissions(t *testing.T) {
	tests := []struct {
		name, method, path, role string
		want                     int
	}{
		{"manager cannot delete", "DELETE", "/api/customers/1", "Manager", fiber.StatusForbidden},
		{"cashier cannot delete", "DELETE", "/api/customers/1", "Cashier", fiber.StatusForbidden},
		{"cashier can create", "POST", "/api/customers/", "Cashier", fiber.StatusCreated},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := `{"name":"Alice","email":"alice@example.com"}`
			if got := authorizedRequest(t, test.method, test.path, test.role, body); got != test.want {
				t.Fatalf("expected %d, got %d", test.want, got)
			}
		})
	}
}
