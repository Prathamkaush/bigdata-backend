package middlewares

import (
	"bigdata-api/internal/repository"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"

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

		log.Println("🔑 Received API Key:", apiKey)
		log.Println("🔐 SHA256:", hashedKey)

		user, err := repository.GetUserByAPIHash(context.Background(), hashedKey)
		if err != nil || user == nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid api key"})
		}

		// FIX: check user.Role
		if user.Role != "admin" {
			return c.Status(403).JSON(fiber.Map{"error": "admin access required"})
		}

		c.Locals("user_id", user.ID)
		return c.Next()
	}
}
