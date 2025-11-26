package controllers

import (
	"bigdata-api/internal/repository"
	"bigdata-api/internal/utils"
	"context"

	"github.com/gofiber/fiber/v2"
)

func ChangeUserRoleController(c *fiber.Ctx) error {
	id := utils.ToInt(c.Params("id"))
	newRole := c.FormValue("role")

	if newRole == "" {
		return c.Status(400).JSON(fiber.Map{"error": "role is required"})
	}

	err := repository.ChangeUserRole(context.Background(), id, newRole)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "role updated",
		"user_id": id,
		"role":    newRole,
	})
}
