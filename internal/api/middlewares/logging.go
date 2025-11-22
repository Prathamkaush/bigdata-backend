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
		path := c.Path()
		method := c.Method()

		// Read body safely
		reqBody := string(c.Body())

		// Continue
		err := c.Next()

		status := c.Response().StatusCode()
		duration := time.Since(start).Milliseconds()

		userID := 0
		if uid := c.Locals("user_id"); uid != nil {
			userID = uid.(int)
		}

		repository.LogAPIRequest(
			context.Background(),
			userID,
			path,
			method,
			status,
			int(duration),
			reqBody,
		)

		return err
	}
}
