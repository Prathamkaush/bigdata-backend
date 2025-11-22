package controllers

import (
	"bigdata-api/internal/database"
	"bigdata-api/internal/metrics"
	"bigdata-api/internal/utils"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

var startTime = time.Now()

func MetricsController(c *fiber.Ctx) error {

	// Check DB statuses
	chErr := database.ClickHouse.Ping(context.Background())
	_, pgErr := database.Postgres.Exec(context.Background(), "SELECT 1")
	_, redisErr := database.Redis.Ping(context.Background()).Result()

	return c.SendString(
		"# HELP api_uptime_seconds API uptime in seconds\n" +
			"# TYPE api_uptime_seconds counter\n" +
			"api_uptime_seconds " + formatSeconds(time.Since(startTime)) + "\n\n" +

			"# HELP api_requests_total Total API requests\n" +
			"# TYPE api_requests_total counter\n" +
			"api_requests_total " + formatUint(metrics.ApiRequestsTotal) + "\n" +

			"# HELP api_requests_2xx Successful requests\n" +
			"# TYPE api_requests_2xx counter\n" +
			"api_requests_2xx " + formatUint(metrics.ApiRequests2xx) + "\n" +

			"# HELP api_requests_4xx Client errors\n" +
			"# TYPE api_requests_4xx counter\n" +
			"api_requests_4xx " + formatUint(metrics.ApiRequests4xx) + "\n" +

			"# HELP api_requests_5xx Server errors\n" +
			"# TYPE api_requests_5xx counter\n" +
			"api_requests_5xx " + formatUint(metrics.ApiRequests5xx) + "\n" +

			"# HELP rate_limit_blocked_total Rate limited requests\n" +
			"# TYPE rate_limit_blocked_total counter\n" +
			"rate_limit_blocked_total " + formatUint(metrics.RateLimitBlocked) + "\n\n" +

			"# HELP cache_hits Number of cache hits\n" +
			"# TYPE cache_hits counter\n" +
			"cache_hits " + formatUint(utils.CacheHits) + "\n" +

			"# HELP cache_misses Number of cache misses\n" +
			"# TYPE cache_misses counter\n" +
			"cache_misses " + formatUint(utils.CacheMisses) + "\n\n" +

			"# HELP clickhouse_status ClickHouse health (1=up)\n" +
			"# TYPE clickhouse_status gauge\n" +
			"clickhouse_status " + boolTo01(chErr == nil) + "\n" +

			"# HELP postgres_status Postgres health (1=up)\n" +
			"# TYPE postgres_status gauge\n" +
			"postgres_status " + boolTo01(pgErr == nil) + "\n" +

			"# HELP redis_status Redis health (1=up)\n" +
			"# TYPE redis_status gauge\n" +
			"redis_status " + boolTo01(redisErr == nil) + "\n",
	)
}

func formatSeconds(d time.Duration) string {
	return formatUint(uint64(d.Seconds()))
}

func formatUint(v uint64) string {
	return fmt.Sprintf("%d", v)
}

func boolTo01(ok bool) string {
	if ok {
		return "1"
	}
	return "0"
}
