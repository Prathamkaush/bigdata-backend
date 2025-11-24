package controllers

import (
	"bigdata-api/internal/repository"
	"context"

	"github.com/gofiber/fiber/v2"
)

func StatsController(c *fiber.Ctx) error {
	ctx := context.Background()

	// GLOBAL STATS (admin)
	today, _ := repository.GetGlobalDailyUsage(ctx)
	totalUsers, _ := repository.CountUsers(ctx)
	totalRequests, _ := repository.TotalRequests(ctx)
	dailyHistory, _ := repository.GetDailyHistory(ctx)

	return c.JSON(fiber.Map{
		"total_requests": totalRequests,
		"today_requests": today.Requests,
		"credits_used":   today.CreditsUsed,
		"total_users":    totalUsers,
		"daily_usage":    dailyHistory,
	})
}
