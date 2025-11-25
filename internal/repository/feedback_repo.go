package repository

import (
	"bigdata-api/internal/database"
	"bigdata-api/internal/models"
	"context"
	"time"
)

func CreateFeedback(ctx context.Context, userID int, msg string, rating int) error {
	_, err := database.Postgres.Exec(ctx,
		`INSERT INTO feedback (user_id, message, rating)
         VALUES ($1, $2, $3)`,
		userID, msg, rating,
	)
	return err
}

func GetFeedback(ctx context.Context, limit, offset int) ([]models.Feedback, int, error) {
	rows, err := database.Postgres.Query(ctx, `
        SELECT id, user_id, message, rating, created_at
        FROM feedback
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var feedbacks []models.Feedback
	for rows.Next() {
		var f models.Feedback
		rows.Scan(&f.ID, &f.UserID, &f.Message, &f.Rating, &f.CreatedAt)
		feedbacks = append(feedbacks, f)
	}

	// total count
	var total int
	_ = database.Postgres.QueryRow(ctx, `SELECT COUNT(*) FROM feedback`).Scan(&total)

	return feedbacks, total, nil
}

type DailyFeedbackHistory struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

func GetFeedbackHistory(ctx context.Context) ([]DailyFeedbackHistory, error) {
	rows, err := database.Postgres.Query(ctx, `
		SELECT DATE(created_at), COUNT(*)
		FROM feedback
		GROUP BY DATE(created_at)
		ORDER BY DATE(created_at) ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []DailyFeedbackHistory

	for rows.Next() {
		var row DailyFeedbackHistory
		var d time.Time

		if err := rows.Scan(&d, &row.Count); err == nil {
			row.Date = d.Format("2006-01-02")
			history = append(history, row)
		}
	}

	return history, nil
}
