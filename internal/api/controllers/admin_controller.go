package controllers

import (
	"bigdata-api/internal/repository"
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var ctx = context.Background()

// =====================================================
// CREATE USER (ADMIN)
// =====================================================
type CreateUserRequest struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Credits  int    `json:"credits"`
}

func CreateUserController(c *fiber.Ctx) error {
	var body CreateUserRequest

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if body.Username == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Username is required"})
	}

	if body.Role == "" {
		body.Role = "sub_admin"
	}

	if body.Credits == 0 {
		body.Credits = 100
	}

	user, rawKey, err := repository.CreateUser(ctx, body.Username, body.Role, body.Credits)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.JSON(fiber.Map{
		"user":    user,
		"api_key": rawKey, // return RAW key
	})
}

// =====================================================
// LIST USERS
// =====================================================
func GetUsersController(c *fiber.Ctx) error {
	users, err := repository.GetAllUsers(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch users"})
	}
	return c.JSON(users)
}

// =====================================================
// USER DETAILS
// =====================================================
func GetUserDetails(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user id"})
	}

	user, err := repository.GetUserDetails(ctx, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}

	return c.JSON(user)
}

// =====================================================
// USER LOGS
// =====================================================
func GetUserLogs(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user id"})
	}

	logs, err := repository.FetchLogsByUser(ctx, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch logs"})
	}

	return c.JSON(logs)
}

// =====================================================
// USER DAILY USAGE
// =====================================================
func GetUserUsage(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user id"})
	}

	usage, err := repository.GetDailyUsageHistory(ctx, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch usage"})
	}

	return c.JSON(usage)
}

// =====================================================
// ADD CREDITS
// =====================================================
type AddCreditsRequest struct {
	UserID  int `json:"user_id"`
	Credits int `json:"credits"`
}

func AddCreditsController(c *fiber.Ctx) error {
	var body AddCreditsRequest

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}

	if body.Credits <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "credits must be positive"})
	}

	err := repository.AddCredits(ctx, body.UserID, body.Credits)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "credits added successfully",
		"user_id": body.UserID,
		"credits": body.Credits,
	})
}

// =====================================================
// GLOBAL LOGS
// =====================================================
func GetLogsController(c *fiber.Ctx) error {
	logs, err := repository.FetchLogs(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch logs"})
	}
	return c.JSON(logs)
}

// =====================================================
// FEEDBACK MANAGEMENT
// =====================================================
func AdminGetFeedback(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}

	limit := 20
	offset := (page - 1) * limit

	feedback, total, err := repository.GetFeedback(ctx, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to load feedback",
		})
	}

	totalPages := (total + limit - 1) / limit

	return c.JSON(fiber.Map{
		"data":       feedback,
		"total":      total,
		"page":       page,
		"totalPages": totalPages,
	})
}
