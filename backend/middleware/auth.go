package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func AuthRequired(secret []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") ||
			strings.TrimSpace(parts[1]) == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid Authorization header",
			})
		}

		tokenString := strings.TrimSpace(parts[1])

		claims := &TokenClaims{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, fiber.NewError(
						fiber.StatusUnauthorized,
						"Invalid token signing method",
					)
				}

				return secret, nil
			},
			jwt.WithValidMethods([]string{"HS256"}),
			jwt.WithIssuer("icecream-shop-api"),
		)

		if err != nil || token == nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		if claims.UserID == 0 ||
			claims.Role == "" ||
			claims.ExpiresAt == nil ||
			claims.IssuedAt == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		// Make the authenticated user's identity available to handlers.
		c.Locals("userID", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}
