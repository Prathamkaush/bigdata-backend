package repository

import (
	"bigdata-api/internal/database"
	"bigdata-api/internal/models"
	"context"
)

type User struct {
	ID         int
	Username   string
	ApiKeyHash string
	Role       string
	Credits    int
	Status     string
}

// Fetch user by API Hash
func GetUserByAPIHash(ctx context.Context, hash string) (*User, error) {
	row := database.Postgres.QueryRow(ctx,
		`SELECT id, username, api_key_hash, role, credits
         FROM users 
         WHERE api_key_hash = $1`,
		hash,
	)

	var u User
	err := row.Scan(&u.ID, &u.Username, &u.ApiKeyHash, &u.Role, &u.Credits)
	if err != nil {
		return nil, err
	}

	u.Status = "active"
	return &u, nil
}

// Create new user
func CreateUser(ctx context.Context, username, apiKeyHash string) error {
	_, err := database.Postgres.Exec(ctx,
		`INSERT INTO users (username, api_key_hash) 
         VALUES ($1, $2)`,
		username, apiKeyHash,
	)
	return err
}

// Get all users
func GetAllUsers(ctx context.Context) ([]models.User, error) {

	query := `
        SELECT id, username, role, credits, created_at
        FROM users
        ORDER BY id ASC
    `

	rows, err := database.Postgres.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var u models.User

		err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.Role,
			&u.Credits,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

// Fetch user by username
func GetUserByName(ctx context.Context, username string) (*User, error) {
	row := database.Postgres.QueryRow(ctx,
		`SELECT id, username, api_key_hash, role, credits
         FROM users WHERE username = $1`,
		username,
	)

	var u User
	err := row.Scan(&u.ID, &u.Username, &u.ApiKeyHash, &u.Role, &u.Credits)
	if err != nil {
		return nil, err
	}

	u.Status = "active"
	return &u, nil
}

func UpdateAPIKey(ctx context.Context, userID int, hash string) error {
	_, err := database.Postgres.Exec(ctx,
		`UPDATE users SET api_key_hash = $1 WHERE id = $2`,
		hash, userID,
	)
	return err
}

func FetchAPIKey(ctx context.Context, userID int) (string, error) {
	row := database.Postgres.QueryRow(ctx,
		`SELECT api_key_hash FROM users WHERE id = $1`,
		userID,
	)

	var hash string
	err := row.Scan(&hash)
	if err != nil {
		return "", err
	}

	return hash, nil
}
