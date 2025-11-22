package middlewares

import (
	"bigdata-api/internal/repository"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		start := time.Now()

		// Save request body safely
		reqBody := string(c.Body())

		// Continue request
		err := c.Next()

		// After handler
		status := c.Response().StatusCode()
		duration := time.Since(start).Milliseconds()

		userID := 0
		if uid := c.Locals("user_id"); uid != nil {
			userID = uid.(int)
		}

		// ✅ NEW FIXED CALL — matches your current repository function
		repository.LogAPIRequest(
			context.Background(),
			userID,
			c.Path(),
			c.Method(),
			status,
			int(duration),
			reqBody,
		)

		return err
	}
}
