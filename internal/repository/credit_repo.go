package repository

import (
	"bigdata-api/internal/database"
	"context"
	"errors"
)

// ===============================
// GET USER CREDITS
// ===============================
func GetCredits(ctx context.Context, userID int) (int, error) {
	row := database.Postgres.QueryRow(ctx,
		`SELECT credits FROM users WHERE id = $1`,
		userID,
	)

	var credits int
	err := row.Scan(&credits)
	if err != nil {
		return 0, err
	}

	return credits, nil
}

// ===============================
// DEDUCT CREDITS
// ===============================
func DeductCredits(ctx context.Context, userID int, amount int) error {
	// Check balance
	row := database.Postgres.QueryRow(ctx,
		`SELECT credits FROM users WHERE id = $1`,
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
		`UPDATE users 
         SET credits = credits - $1 
         WHERE id = $2`,
		amount, userID,
	)

	return err
}

// ===============================
// ADD CREDITS
// ===============================
func AddCredits(ctx context.Context, userID int, credits int) error {

	// 1. Log credit addition
	_, err := database.Postgres.Exec(ctx,
		`INSERT INTO credit_logs (user_id, change_amount, reason)
         VALUES ($1, $2, $3)`,
		userID, credits, "admin_added_credits",
	)
	if err != nil {
		return err
	}

	// 2. Update balance
	_, err = database.Postgres.Exec(ctx,
		`UPDATE users
         SET credits = credits + $1
         WHERE id = $2`,
		credits, userID,
	)

	return err
}

// ===============================
// LOG CREDIT USAGE
// ===============================
func LogCreditUsage(ctx context.Context, userID int, amount int, endpoint string) {
	database.Postgres.Exec(ctx,
		`INSERT INTO credit_logs (user_id, change_amount, reason)
         VALUES ($1, $2, $3)`,
		userID, -amount, endpoint,
	)
}
