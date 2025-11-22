package controllers

import "github.com/gofiber/fiber/v2"

func GetCredits(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"credits": 0,
	})
}
