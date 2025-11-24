package repository

import (
	"bigdata-api/internal/database"
	"context"
	"time"
)

// IncrementDailyUsage increments today's requests and credits_used for a user.
func IncrementDailyUsage(ctx context.Context, userID int, credits int) error {
	_, err := database.Postgres.Exec(ctx, `
INSERT INTO daily_usage (user_id, date, requests, credits_used)
VALUES ($1, CURRENT_DATE, 1, $2)
ON CONFLICT (user_id, date)
DO UPDATE SET
requests = daily_usage.requests + 1,
credits_used = daily_usage.credits_used + EXCLUDED.credits_used
`, userID, credits)

	return err
}

func GetDailyUsageHistory(ctx context.Context, userID int) ([]map[string]interface{}, error) {
	rows, err := database.Postgres.Query(ctx, `
        SELECT date, requests, credits_used
        FROM daily_usage
        WHERE user_id = $1
        ORDER BY date DESC
        LIMIT 30
    `, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []map[string]interface{}

	for rows.Next() {
		var date time.Time
		var req, credits int

		rows.Scan(&date, &req, &credits)

		list = append(list, map[string]interface{}{
			"date":         date.Format("2006-01-02"),
			"requests":     req,
			"credits_used": credits,
		})
	}

	return list, nil
}
