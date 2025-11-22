package middlewares

import (
	"bigdata-api/internal/database"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	RateLimitRequests  = 10 // 10 requests
	RateLimitWindowSec = 60 // per 60 seconds
)

func RateLimitMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		apiKey := c.Get("x-api-key")
		if apiKey == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "API key missing",
			})
		}

		redis := database.Redis

		bucketKey := fmt.Sprintf("ratelimit:%s", apiKey)

		// Increment request count
		reqCount, err := redis.Incr(context.Background(), bucketKey).Result()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "redis increment failed",
			})
		}

		// Set expiration if first request in window
		if reqCount == 1 {
			redis.Expire(context.Background(), bucketKey, time.Duration(RateLimitWindowSec)*time.Second)
		}

		// Check limit
		if reqCount > RateLimitRequests {
			ttl, _ := redis.TTL(context.Background(), bucketKey).Result()

			return c.Status(429).JSON(fiber.Map{
				"error":           "rate limit exceeded",
				"retry_after_sec": int(ttl.Seconds()),
			})
		}

		return c.Next()
	}
}
