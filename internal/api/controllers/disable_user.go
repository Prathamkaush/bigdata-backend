package controllers

import (
	"bigdata-api/internal/repository"
	"bigdata-api/internal/utils"
	"context"

	"github.com/gofiber/fiber/v2"
)

func DisableUserController(c *fiber.Ctx) error {
	id := utils.ToInt(c.Params("id"))

	err := repository.UpdateUserStatus(context.Background(), id, "disabled")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "user disabled",
	})
}
