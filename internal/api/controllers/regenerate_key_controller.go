package controllers

import (
	"bigdata-api/internal/repository"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
)

func generateNewAPIKey() string {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RegenerateAPIKeyController(c *fiber.Ctx) error {

	userID := c.Locals("user_id").(int)

	// generate new key
	newKey := generateNewAPIKey()

	hash := sha256.Sum256([]byte(newKey))
	hashedKey := hex.EncodeToString(hash[:])

	err := repository.UpdateAPIKey(context.Background(), userID, hashedKey)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to update key"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"api_key": newKey,
	})
}

func GetAdminAPIKey(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	hash, err := repository.FetchAPIKey(context.Background(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch key"})
	}

	masked := hash[:6] + "********" + hash[len(hash)-4:]

	return c.JSON(fiber.Map{
		"api_key": masked,
		"note":    "Full key cannot be retrieved. Regenerate to get a new raw API key.",
	})
}
