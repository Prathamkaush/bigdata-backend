package middlewares

import (
	"bigdata-api/internal/repository"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func ApiKeyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		fmt.Println("DEBUG HEADER:", c.Get("x-api-key"))

		apiKey := c.Get("x-api-key")
		if apiKey == "" {
			return c.Status(401).JSON(fiber.Map{"error": "missing api key"})
		}

		// 1️⃣ FIRST TRY RAW KEY (backward compatible)
		user, err := repository.GetUserByAPIHash(context.Background(), apiKey)
		if err == nil && user != nil && (user.Status == "active" || user.Status == "admin") {
			c.Locals("user_id", user.ID)
			return c.Next()
		}

		// 2️⃣ THEN TRY HASHED VERSION (for new users)
		hash := sha256.Sum256([]byte(apiKey))
		hashedKey := hex.EncodeToString(hash[:])

		user, err = repository.GetUserByAPIHash(context.Background(), hashedKey)
		if err == nil && user != nil && (user.Status == "active" || user.Status == "admin") {
			c.Locals("user_id", user.ID)
			return c.Next()
		}

		// ❌ If neither works → invalid key
		return c.Status(401).JSON(fiber.Map{"error": "invalid api key"})
	}
}
