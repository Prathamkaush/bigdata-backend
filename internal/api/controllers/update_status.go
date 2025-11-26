package controllers

import (
	"bigdata-api/internal/repository"
	"bigdata-api/internal/utils"
	"context"

	"github.com/gofiber/fiber/v2"
)

func UpdateUserStatusController(c *fiber.Ctx) error {
	id := utils.ToInt(c.Params("id"))
	status := c.FormValue("status") // active | disabled

	if status == "" {
		return c.Status(400).JSON(fiber.Map{"error": "status is required"})
	}
	if status != "active" && status != "disabled" {
		return c.Status(400).JSON(fiber.Map{"error": "invalid status"})
	}

	err := repository.UpdateUserStatus(context.Background(), id, status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"status":  status,
	})
}
