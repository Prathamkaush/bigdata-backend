package middlewares

import (
	"bigdata-api/internal/repository"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		startTime := time.Now()
		reqBody := string(c.Body())

		// continue request
		err := c.Next()

		status := c.Response().StatusCode()
		durationMs := int(time.Since(startTime).Milliseconds())

		// Extract user_id set by Auth middleware
		userID := 0
		if uid := c.Locals("user_id"); uid != nil {
			userID = uid.(int)
		}

		// SAVE LOG
		logErr := repository.LogAPIRequest(
			context.Background(),
			userID,
			c.Path(),
			c.Method(),
			status,
			durationMs,
			reqBody,
		)

		if logErr != nil {
			fmt.Println("🔥 ERROR SAVING LOG:", logErr.Error())
		}

		return err
	}
}
