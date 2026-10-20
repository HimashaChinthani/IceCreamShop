package middleware

import (
    "strings"

    "github.com/gofiber/fiber/v2"
)

func RequireRoles(allowedRoles ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        role, ok := c.Locals("role").(string)

        if !ok || role == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(
                fiber.Map{
                    "error": "Authentication required",
                },
            )
        }

        for _, allowed := range allowedRoles {
            if strings.EqualFold(role, allowed) {
                return c.Next()
            }
        }

        return c.Status(fiber.StatusForbidden).JSON(
            fiber.Map{
                "error": "You do not have permission to perform this action",
            },
        )
    }
}