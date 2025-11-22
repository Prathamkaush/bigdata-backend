package controllers

import (
	"bigdata-api/internal/repository"
	"context"

	"github.com/gofiber/fiber/v2"
)

func StatsController(c *fiber.Ctx) error {
	// admin routes will have AdminMiddleware which sets user_id
	// we treat admin as global stats (pass 0 to repo functions if needed)
	var userID int
	if id := c.Locals("user_id"); id != nil {
		userID = id.(int)
	}

	ctx := context.Background()

	// 1. Today's usage (use userID; admin may pass 0 and repo should handle)
	today, _ := repository.GetDailyUsage(ctx, userID)

	// 2. Total users
	totalUsers, _ := repository.CountUsers(ctx)

	// 3. Total API requests
	totalRequests, _ := repository.TotalRequests(ctx)

	// 4. Daily usage history
	dailyHistory, _ := repository.GetDailyHistory(ctx)

	return c.JSON(fiber.Map{
		"total_requests": totalRequests,
		"today_requests": today.Requests,
		"credits_used":   today.CreditsUsed,
		"total_users":    totalUsers,
		"daily_usage":    dailyHistory,
	})
}
