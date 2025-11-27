package controllers

import (
	"bigdata-api/internal/repository"
	"context"

	"github.com/gofiber/fiber/v2"
)

type VerifyKeyRequest struct {
	Hash string `json:"hash"`
}

func VerifyKeyController(c *fiber.Ctx) error {
	var body VerifyKeyRequest

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	if body.Hash == "" {
		return c.Status(400).JSON(fiber.Map{"error": "missing hash"})
	}

	// Fetch user from DB using hashed key
	user, err := repository.GetUserByAPIHash(context.Background(), body.Hash)
	if err != nil || user == nil {
		return c.JSON(fiber.Map{
			"valid": false,
		})
	}

	// SUCCESS — return role too
	return c.JSON(fiber.Map{
		"valid":    true,
		"role":     user.Role, // <-- IMPORTANT
		"id":       user.ID,
		"username": user.Username,
	})
}
