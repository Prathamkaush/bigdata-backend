package repository

import (
	"bigdata-api/internal/database"
	"context"
	"time"
)

// LogAPIRequest inserts a simple log entry (matching your DB schema)
func LogAPIRequest(ctx context.Context, userID int, endpoint string, params string, creditsUsed int) error {
	_, err := database.Postgres.Exec(ctx, `
INSERT INTO api_logs (user_id, endpoint, params, credits_used, created_at)
VALUES ($1, $2, $3, $4, NOW())
`, userID, endpoint, params, creditsUsed)

	return err
}

// FetchLogs returns recent logs in a simple shape the frontend can use
func FetchLogs(ctx context.Context) ([]map[string]interface{}, error) {
	rows, err := database.Postgres.Query(ctx, `
SELECT id, user_id, endpoint, params, credits_used, created_at
FROM api_logs
ORDER BY created_at DESC
LIMIT 200
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id, userID, creditsUsed int
		var endpoint, params string
		var createdAt time.Time

		if err := rows.Scan(&id, &userID, &endpoint, &params, &creditsUsed, &createdAt); err != nil {
			continue
		}

		result = append(result, map[string]interface{}{
			"id":           id,
			"user_id":      userID,
			"endpoint":     endpoint,
			"params":       params,
			"credits_used": creditsUsed,
			"created_at":   createdAt.Format(time.RFC3339),
		})
	}

	return result, nil
}
