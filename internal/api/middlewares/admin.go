// internal/api/middlewares/admin.go

package middlewares

import (
	"bigdata-api/internal/repository"
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/gofiber/fiber/v2"
)

func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		apiKey := c.Get("x-api-key")
		if apiKey == "" {
			return c.Status(401).JSON(fiber.Map{"error": "missing api key"})
		}

		hash := sha256.Sum256([]byte(apiKey))
		hashedKey := hex.EncodeToString(hash[:])

		user, err := repository.GetUserByAPIHash(context.Background(), hashedKey)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid api key"})
		}

		// ⬅ Allow both "super_admin" & "admin"
		if user.Role != "super_admin" && user.Role != "admin" {
			return c.Status(403).JSON(fiber.Map{
				"error": "admin access required",
			})
		}

		// Set session info
		c.Locals("user_id", user.ID)
		c.Locals("role", user.Role)

		return c.Next()
	}
}
