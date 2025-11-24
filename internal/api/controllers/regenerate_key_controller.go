package controllers

import (
	"bigdata-api/internal/repository"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	// GET USER ID
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	// FETCH USER
	user, err := repository.GetUserByID(context.Background(), id)
	if err != nil || user == nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	// DO NOT ALLOW ADMIN KEY REGEN
	if user.Role == "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Admin API key cannot be regenerated",
		})
	}

	// CREATE NEW RAW KEY
	rawKey := uuid.New().String()

	// HASH IT
	hash := sha256.Sum256([]byte(rawKey))
	hashedKey := hex.EncodeToString(hash[:])

	// UPDATE IN DB
	err = repository.UpdateAPIKey(context.Background(), user.ID, hashedKey)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to update key",
		})
	}

	return c.JSON(fiber.Map{
		"message": "API key regenerated",
		"api_key": rawKey,
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
