package repository

import (
	"bigdata-api/internal/database"
	"bigdata-api/internal/models"
	"context"
)

type Feedback struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Message   string `json:"message"`
	Rating    int    `json:"rating"`
	CreatedAt string `json:"created_at"`
}

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

	// Count total
	row := database.Postgres.QueryRow(ctx, `SELECT COUNT(*) FROM feedback`)
	var total int
	row.Scan(&total)

	return feedbacks, total, nil
}
