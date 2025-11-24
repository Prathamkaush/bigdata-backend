package controllers

import (
	"bigdata-api/internal/repository"
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type FeedbackBody struct {
	Message string `json:"message"`
	Rating  int    `json:"rating"`
}

func SubmitFeedback(c *fiber.Ctx) error {
	var body FeedbackBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}

	if body.Message == "" || body.Rating < 1 || body.Rating > 5 {
		return c.Status(400).JSON(fiber.Map{"error": "message + rating (1–5) required"})
	}

	userID := c.Locals("user_id").(int)

	err := repository.CreateFeedback(context.Background(), userID, body.Message, body.Rating)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to save feedback"})
	}

	return c.JSON(fiber.Map{
		"message": "feedback submitted",
	})
}

func GetFeedback(c *fiber.Ctx) error {
	pageStr := c.Query("page", "1")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	limit := 10
	offset := (page - 1) * limit

	data, total, err := repository.GetFeedback(context.Background(), limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch feedback"})
	}

	return c.JSON(fiber.Map{
		"data":       data,
		"total":      total,
		"page":       page,
		"totalPages": (total + limit - 1) / limit,
	})
}
