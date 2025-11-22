package repository

import (
	"bigdata-api/internal/database"
	"context"
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
