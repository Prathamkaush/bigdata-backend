package controllers

import (
	"bigdata-api/internal/repository"
	"bigdata-api/internal/utils"
	"context"

	"github.com/gofiber/fiber/v2"
)

func UpdateUserCreditsController(c *fiber.Ctx) error {
	id := utils.ToInt(c.Params("id"))
	credits := utils.ToInt(c.FormValue("credits"))

	if credits < 0 {
		return c.Status(400).JSON(fiber.Map{"error": "credits cannot be negative"})
	}

	err := repository.UpdateUserCredits(context.Background(), id, credits)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"credits": credits,
	})
}
