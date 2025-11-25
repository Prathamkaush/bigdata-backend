package repository

import (
	"bigdata-api/internal/database"
	"context"
	"time"
)

func NewUsersToday(ctx context.Context) (int, error) {
	var count int
	err := database.Postgres.QueryRow(ctx,
		`SELECT COUNT(*) FROM users WHERE DATE(created_at)=CURRENT_DATE`,
	).Scan(&count)
	return count, err
}

func LowCreditUsers(ctx context.Context) (int, error) {
	var count int
	err := database.Postgres.QueryRow(ctx,
		`SELECT COUNT(*) FROM users WHERE credits < 50`,
	).Scan(&count)
	return count, err
}

func NewFeedbackCount(ctx context.Context) (int, error) {
	var count int
	err := database.Postgres.QueryRow(ctx,
		`SELECT COUNT(*) FROM feedback WHERE DATE(created_at)=CURRENT_DATE`,
	).Scan(&count)
	return count, err
}

func TotalCreditsUsedAll(ctx context.Context) (int, error) {
	var total int
	err := database.Postgres.QueryRow(ctx,
		`SELECT COALESCE(SUM(credits_used),0)
         FROM daily_usage`,
	).Scan(&total)
	return total, err
}

type FullDailyUsage struct {
	Date        string `json:"date"`
	Requests    int    `json:"requests"`
	CreditsUsed int    `json:"credits_used"`
	NewUsers    int    `json:"new_users"`
}

func Get30DayUsage(ctx context.Context) ([]FullDailyUsage, error) {
	rows, err := database.Postgres.Query(ctx, `
        SELECT date,
               SUM(requests) AS requests,
               SUM(credits_used) AS credits_used
        FROM daily_usage
        GROUP BY date
        ORDER BY date ASC
        LIMIT 30
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []FullDailyUsage

	for rows.Next() {
		var d time.Time
		var req, cred int

		if err := rows.Scan(&d, &req, &cred); err == nil {
			list = append(list, FullDailyUsage{
				Date:        d.Format("2006-01-02"),
				Requests:    req,
				CreditsUsed: cred,
				NewUsers:    0, // future proof (we add below)
			})
		}
	}

	return list, nil
}
