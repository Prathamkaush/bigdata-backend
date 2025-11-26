package controllers

import (
	"bigdata-api/internal/repository"
	"context"

	"github.com/gofiber/fiber/v2"
)

func CountRecordsController(c *fiber.Ctx) error {
	ctx := context.Background()

	q := "SELECT COUNT(*) FROM normalized_records"

	count, err := repository.CountRecords(ctx, q, []interface{}{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"count": count,
	})
}
