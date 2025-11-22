package controllers

import (
	"bigdata-api/internal/database"
	"time"

	"github.com/gofiber/fiber/v2"
)

func HealthCheck(c *fiber.Ctx) error {
	start := time.Now()

	// Check ClickHouse
	chErr := database.ClickHouse.Ping(c.Context())

	// Check Postgres
	pgErr := database.Postgres.Ping(c.Context())

	// Check Redis
	_, redisErr := database.Redis.Ping(c.Context()).Result()

	status := "ok"

	if chErr != nil || pgErr != nil || redisErr != nil {
		status = "degraded"
	}

	return c.JSON(fiber.Map{
		"status": status,
		"services": fiber.Map{
			"clickhouse": boolToStatus(chErr == nil),
			"postgres":   boolToStatus(pgErr == nil),
			"redis":      boolToStatus(redisErr == nil),
		},
		"uptime_ms": time.Since(start).Milliseconds(),
	})
}

func boolToStatus(ok bool) string {
	if ok {
		return "healthy"
	}
	return "unhealthy"
}
