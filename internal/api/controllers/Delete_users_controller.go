package controllers

import (
	"bigdata-api/internal/repository"
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func DeleteUserController(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user id"})
	}

	err = repository.DeleteUser(context.Background(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete user"})
	}

	return c.JSON(fiber.Map{
		"message": "user deleted",
		"user_id": id,
	})
}
