package repository

import (
	"bigdata-api/internal/database"
	"context"
	"errors"
)

func GetCredits(ctx context.Context, userID int) (int, error) {
	row := database.Postgres.QueryRow(ctx,
		`SELECT credits FROM user_credits WHERE user_id = $1`,
		userID,
	)

	var credits int
	err := row.Scan(&credits)
	if err != nil {
		return 0, err
	}

	return credits, nil
}

func DeductCredits(ctx context.Context, userID int, amount int) error {
	// Check balance
	row := database.Postgres.QueryRow(ctx,
		`SELECT credits FROM user_credits WHERE user_id = $1`,
		userID,
	)

	var credits int
	if err := row.Scan(&credits); err != nil {
		return err
	}

	if credits < amount {
		return errors.New("insufficient credits")
	}

	// Deduct
	_, err := database.Postgres.Exec(ctx,
		`UPDATE user_credits 
		 SET credits = credits - $1 
		 WHERE user_id = $2`,
		amount, userID,
	)

	return err
}

func AddCredits(ctx context.Context, userID int, amount int) error {
	_, err := database.Postgres.Exec(ctx,
		`UPDATE users 
         SET credits = credits + $1
         WHERE id = $2`,
		amount, userID,
	)

	return err
}

func LogCreditUsage(ctx context.Context, userID int, amount int, endpoint string) {
	database.Postgres.Exec(ctx,
		`INSERT INTO credit_logs (user_id, credits_used, endpoint)
		 VALUES ($1, $2, $3)`,
		userID, amount, endpoint,
	)
}
