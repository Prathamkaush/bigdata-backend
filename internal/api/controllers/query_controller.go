package controllers

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"bigdata-api/internal/metrics"
	"bigdata-api/internal/models"
	"bigdata-api/internal/repository"
	"bigdata-api/internal/services"
	"bigdata-api/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func QueryController(c *fiber.Ctx) error {
	start := time.Now()

	var req models.QueryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// default rules
	if req.Limit <= 0 {
		req.Limit = 100
	}
	if req.Limit > 1000 {
		req.Limit = 1000
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// ---------------------------------------------
	// 1️⃣ BUILD CACHE KEY (hash of request body)
	// ---------------------------------------------
	bodyBytes, _ := json.Marshal(req)
	key := "cache_query:" + fmt.Sprintf("%x", sha256.Sum256(bodyBytes))

	// ---------------------------------------------
	// 2️⃣ CACHE CHECK (Redis)
	// ---------------------------------------------
	if cached, err := utils.CacheGet(key); err == nil && cached != "" {
		metrics.CacheHits++

		var resp fiber.Map
		json.Unmarshal([]byte(cached), &resp)

		// Also set header so CreditsMiddleware works
		if meta, ok := resp["metadata"].(map[string]interface{}); ok {
			if returned, ok := meta["returned"].(float64); ok {
				c.Set("X-Records-Returned", strconv.Itoa(int(returned)))
			}
		}

		return c.JSON(resp)
	}

	metrics.CacheMisses++

	// ---------------------------------------------
	// 3️⃣ BUILD CLICKHOUSE QUERIES
	// ---------------------------------------------
	selectCols := []string{
		"dedupe_key", "source", "source_id", "customer_id",
		"first_name", "last_name", "email", "phone",
		"city", "state", "country", "ingest_ts",
	}

	sqlSelect, args := services.BuildSelectQuery(
		selectCols,
		req.Filters,
		req.Range,
		req.Fuzzy,
		req.Sort,
		req.Limit,
		req.Offset,
	)

	sqlCount, countArgs := services.BuildCountQuery(
		req.Filters,
		req.Range,
		req.Fuzzy,
	)

	// ---------------------------------------------
	// 4️⃣ RUN CLICKHOUSE QUERY
	// ---------------------------------------------
	ctx := context.Background()
	rows, err := services.SearchRecords(ctx, sqlSelect, args)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Query failed"})
	}

	total, err := services.CountRecords(ctx, sqlCount, countArgs)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Count failed"})
	}

	// ---------------------------------------------
	// 5️⃣ BUILD RESPONSE
	// ---------------------------------------------
	resp := fiber.Map{
		"metadata": fiber.Map{
			"total_records": total,
			"returned":      len(rows),
			"credits_used":  1 + len(rows)/100,
			"response_time": time.Since(start).Milliseconds(),
		},
		"data": rows,
	}

	// Set header so CreditsMiddleware can deduct credits
	c.Set("X-Records-Returned", strconv.Itoa(len(rows)))

	// ---------------------------------------------
	// 6️⃣ SAVE TO CACHE
	// ---------------------------------------------
	utils.CacheSet(key, resp, 60*time.Second)

	// 7️⃣ UPDATE DAILY USAGE (async)
	userID := c.Locals("user_id")

	var uid int
	switch v := userID.(type) {
	case int:
		uid = v
	case int64:
		uid = int(v)
	}

	if uid > 0 {
		creditsUsed := 1 + len(rows)/100
		go repository.IncrementDailyUsage(context.Background(), uid, creditsUsed)
	}

	return c.JSON(resp)
}
