package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)

		for _, r := range allowedRoles {
			if role == r {
				return c.Next()
			}
		}

		return c.Status(403).JSON(fiber.Map{
			"error": "Permission denied",
		})
	}
}
