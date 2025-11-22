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

		err := c.Next()

		duration := time.Since(start).Milliseconds()

		userID := 0
		if uid := c.Locals("user_id"); uid != nil {
			userID = uid.(int)
		}

		endpoint := c.Path()
		method := c.Method()
		status := c.Response().StatusCode()

		// Call new simplified logger
		go repository.LogAPIRequest(
			context.Background(),
			userID,
			endpoint,
			method,
			status,
		)

		// Debug log
		println("LOG:", userID, endpoint, method, status, duration, "ms")

		return err
	}
}
