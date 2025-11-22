package controllers

import (
	"bigdata-api/internal/database"
	"context"

	"github.com/gofiber/fiber/v2"
)

func VerifyKeyController(c *fiber.Ctx) error {
	var body struct {
		Hash string `json:"hash"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"valid": false,
		})
	}

	var count int
	err := database.Postgres.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM users WHERE api_key_hash=$1 AND role='admin'",
		body.Hash,
	).Scan(&count)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"valid": false,
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"valid": count > 0,
	})
}
