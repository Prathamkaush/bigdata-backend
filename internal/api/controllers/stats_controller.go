package controllers

import (
	"context"
	"log"

	"bigdata-api/internal/repository"

	"github.com/gofiber/fiber/v2"
)

func StatsController(c *fiber.Ctx) error {
	ctx := context.Background()

	totalUsers, _ := repository.CountUsers(ctx)
	totalRequests, _ := repository.TotalRequests(ctx)
	todayUsage, _ := repository.GetGlobalDailyUsage(ctx)
	totalCreditsUsedAll, _ := repository.TotalCreditsUsedAll(ctx)

	newUsersToday, _ := repository.NewUsersToday(ctx)
	lowCreditUsers, _ := repository.LowCreditUsers(ctx)
	newFeedbackToday, _ := repository.NewFeedbackCount(ctx)

	last30Days, _ := repository.Get30DayUsage(ctx)
	log.Println("🔥 NEW StatsController loaded!!!", last30Days)

	return c.JSON(fiber.Map{
		"total_users":    totalUsers,
		"total_requests": totalRequests,

		"today_requests":        todayUsage.Requests,
		"credits_used":          todayUsage.CreditsUsed,
		"credits_used_all_time": totalCreditsUsedAll,

		"new_users_today":    newUsersToday,
		"low_credit_users":   lowCreditUsers,
		"new_feedback_today": newFeedbackToday,

		"daily_usage": last30Days,
	})
}
