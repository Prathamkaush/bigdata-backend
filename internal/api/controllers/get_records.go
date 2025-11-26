package controllers

import (
	"bigdata-api/internal/services"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func GetRecords(c *fiber.Ctx) error {
	ctx := context.Background()

	search := c.Query("search", "")
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	// BASE QUERY (use ONLY ? placeholders)
	sql := `
SELECT 
    id,
    full_name,
    email,
    phone,
    source,
    dedupe_key,
    uploaded_at,
    city,
    state,
    country,
    age,
    gender,
    score,
    normalization_status
FROM normalized_records
WHERE 1 = 1
`

	args := []interface{}{}

	// SEARCH
	if search != "" {
		sql += `
AND (
    lower(full_name) LIKE lower(?)
    OR lower(email) LIKE lower(?)
    OR phone LIKE ?
)
`
		args = append(args, "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// FINAL ORDER + PAGINATION
	sql += `
ORDER BY uploaded_at DESC
LIMIT ? OFFSET ?
`

	args = append(args, limit, offset)

	// EXECUTE
	records, err := services.SearchRecords(ctx, sql, args)
	fmt.Println("🟩 FINAL SQL:\n", sql)
	fmt.Println("🟧 ARGS:", args)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
			"query": sql,
			"args":  args,
		})
	}

	return c.JSON(records)
}
