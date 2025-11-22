package controllers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"bigdata-api/internal/repository"
	"bigdata-api/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ------------------------------------------------------
// 1. CREATE USER (Admin) — SECURE (hashed API key)
// ------------------------------------------------------

type CreateUserRequest struct {
	Name string `json:"name"`
}

func CreateUserController(c *fiber.Ctx) error {
	var body CreateUserRequest

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}

	// Generate raw key
	rawKey := uuid.New().String()

	// Hash before storing
	hash := sha256.Sum256([]byte(rawKey))
	hashedKey := hex.EncodeToString(hash[:])

	// Save hashed key
	err := repository.CreateUser(context.Background(), body.Name, hashedKey)
	if err != nil {
		utils.Error("CreateUser failed: " + err.Error())
		return c.Status(500).JSON(fiber.Map{"error": "failed to create user"})
	}

	return c.JSON(fiber.Map{
		"message": "user created",
		"api_key": rawKey, // return RAW key only once
	})
}

// ------------------------------------------------------
// 2. ADD CREDITS (Admin)
// ------------------------------------------------------

type AddCreditsRequest struct {
	Username string `json:"username"`
	Credits  int    `json:"credits"`
}

func AddCreditsController(c *fiber.Ctx) error {
	var body AddCreditsRequest

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}

	// find by username
	user, err := repository.GetUserByName(context.Background(), body.Username)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}

	// add credits
	err = repository.AddCredits(context.Background(), user.ID, body.Credits)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to add credits"})
	}

	return c.JSON(fiber.Map{
		"message": "credits added",
		"user_id": user.ID,
		"credits": body.Credits,
	})
}

// ------------------------------------------------------
// 3. GET LOGS (Admin)
// ------------------------------------------------------

func GetLogsController(c *fiber.Ctx) error {
	ctx := context.Background()
	logs, err := repository.FetchLogs(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch logs"})
	}
	return c.JSON(logs)
}
func GetUsersController(c *fiber.Ctx) error {

	users, err := repository.GetAllUsers(context.Background())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(users)
}
