package repository

import (
	"bigdata-api/internal/database"
	"context"
	"time"
)

type DailyUsage struct {
	Requests    int
	CreditsUsed int
}

// GetDailyUsage returns today's usage for a particular user (or 0s if none)
func GetDailyUsage(ctx context.Context, userID int) (*DailyUsage, error) {
	row := database.Postgres.QueryRow(ctx, `
SELECT COALESCE(requests,0), COALESCE(credits_used,0)
FROM daily_usage
WHERE user_id = $1 AND date = CURRENT_DATE
`, userID)

	var du DailyUsage
	if err := row.Scan(&du.Requests, &du.CreditsUsed); err != nil {
		// return zero values instead of error (caller can handle)
		return &DailyUsage{Requests: 0, CreditsUsed: 0}, nil
	}

	return &du, nil
}

// CountUsers returns total users
func CountUsers(ctx context.Context) (int, error) {
	var count int
	err := database.Postgres.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&count)
	return count, err
}

// TotalRequests returns sum of all requests
func TotalRequests(ctx context.Context) (int, error) {
	var total int
	err := database.Postgres.QueryRow(ctx, `SELECT COALESCE(SUM(requests),0) FROM daily_usage`).Scan(&total)
	return total, err
}

// DailyHistoryRow represents a chart row
type DailyHistoryRow struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// GetDailyHistory returns last 30 days of aggregated requests
func GetDailyHistory(ctx context.Context) ([]DailyHistoryRow, error) {
	rows, err := database.Postgres.Query(ctx, `
SELECT date, requests
FROM daily_usage
ORDER BY date ASC
LIMIT 30
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []DailyHistoryRow
	for rows.Next() {
		var d time.Time
		var count int
		if err := rows.Scan(&d, &count); err == nil {
			list = append(list, DailyHistoryRow{Date: d.Format("2006-01-02"), Count: count})
		}
	}

	return list, nil
}
