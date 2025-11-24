package middlewares

import (
	"bigdata-api/internal/repository"
	"bigdata-api/internal/utils"
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Credit pricing model:
//
// Base cost = 1
// +1 credit per 50 rows
// Cap = 20 credits
// If rows <= 10 → only 1 credit
func calculateCredits(rows int) int {
	base := 1

	if rows <= 10 {
		return base
	}

	extra := rows / 50
	total := base + extra

	if total > 20 {
		return 20
	}
	return total
}

func CreditsMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		// user_id must be set by AuthMiddleware
		userID := c.Locals("user_id")
		if userID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "user not authenticated",
			})
		}

		uid, ok := userID.(int)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid user id type",
			})
		}

		// PRE-CHECK: user must have at least 1 credit to run the query
		credits, err := repository.GetCredits(context.Background(), uid)
		if err != nil {
			utils.Error("Failed to fetch user credits: " + err.Error())
			return c.Status(500).JSON(fiber.Map{"error": "cannot check credits"})
		}

		if credits < 1 {
			return c.Status(402).JSON(fiber.Map{
				"error": "insufficient credits",
			})
		}

		// Continue to actual endpoint
		err = c.Next()

		// If handler returned an error → do not deduct credits
		if err != nil {
			return err
		}

		// If cached response → no deduction
		if c.GetRespHeader("X-From-Cache") == "1" {
			return nil
		}

		// AFTER handler: determine returned rows
		headerRows := c.Get("X-Records-Returned", "0")
		returnedRows, _ := strconv.Atoi(headerRows)
		creditsToDeduct := calculateCredits(returnedRows)
		// Deduct credits
		err = repository.DeductCredits(context.Background(), uid, creditsToDeduct)
		if err != nil {
			utils.Error("Failed credit deduction: " + err.Error())
			return c.Status(402).JSON(fiber.Map{
				"error": "insufficient credits",
			})
		}

		// 📌 NEW: Increment daily usage stats
		err = repository.IncrementDailyUsage(context.Background(), uid, creditsToDeduct)
		if err != nil {
			utils.Error("Failed to update daily usage: " + err.Error())
		}

		// Log credit usage (async fire-and-forget)
		go repository.LogCreditUsage(context.Background(), uid, creditsToDeduct, c.Path())

		return nil
	}
}
