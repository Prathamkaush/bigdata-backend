package repository

import (
	"bigdata-api/internal/database"
	"bigdata-api/internal/models"
	"bigdata-api/internal/utils"
	"context"
)

// ------------------------------
// AUTH USER STRUCT (internal)
// ------------------------------
type User struct {
	ID         int
	Username   string
	ApiKeyHash string
	Role       string
	Credits    int
	Status     string
}

// ------------------------------
// 1. Fetch user by API hash
// ------------------------------
func GetUserByAPIHash(ctx context.Context, hash string) (*User, error) {
	row := database.Postgres.QueryRow(ctx,
		`SELECT id, username, api_key_hash, role, credits
         FROM users WHERE api_key_hash = $1`,
		hash,
	)

	var u User
	if err := row.Scan(&u.ID, &u.Username, &u.ApiKeyHash, &u.Role, &u.Credits); err != nil {
		return nil, err
	}

	u.Status = "active"
	return &u, nil
}

// ------------------------------
// 2. Create user
// ------------------------------
func CreateUser(ctx context.Context, username string) (*models.User, string, error) {

	rawKey := utils.GenerateApiKey()
	hashed := utils.HashString(rawKey)

	row := database.Postgres.QueryRow(ctx, `
        INSERT INTO users (username, api_key_hash, credits, role)
        VALUES ($1, $2, 100, 'user')
        RETURNING id, username, credits, role, created_at
    `, username, hashed)

	var u models.User
	if err := row.Scan(&u.ID, &u.Username, &u.Credits, &u.Role, &u.CreatedAt); err != nil {
		return nil, "", err
	}

	u.ApiKey = hashed
	return &u, rawKey, nil
}

// ------------------------------
// 3. Get all users
// ------------------------------
func GetAllUsers(ctx context.Context) ([]models.User, error) {
	rows, err := database.Postgres.Query(ctx, `
        SELECT id, username, credits, created_at, api_key_hash, role
        FROM users ORDER BY id DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.User

	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Credits, &u.CreatedAt, &u.ApiKey, &u.Role); err == nil {
			list = append(list, u)
		}
	}

	return list, nil
}

// ------------------------------
// 4. Fetch user by name
// ------------------------------
func GetUserByName(ctx context.Context, username string) (*User, error) {
	row := database.Postgres.QueryRow(ctx,
		`SELECT id, username, api_key_hash, role, credits
         FROM users WHERE username = $1`,
		username,
	)

	var u User
	if err := row.Scan(&u.ID, &u.Username, &u.ApiKeyHash, &u.Role, &u.Credits); err != nil {
		return nil, err
	}

	u.Status = "active"
	return &u, nil
}

// ------------------------------
// 5. Update API Key
// ------------------------------
func UpdateAPIKey(ctx context.Context, userID int, hash string) error {
	_, err := database.Postgres.Exec(ctx,
		`UPDATE users SET api_key_hash = $1 WHERE id = $2`,
		hash, userID,
	)
	return err
}

// ------------------------------
// 6. Fetch stored hashed key
// ------------------------------
func FetchAPIKey(ctx context.Context, id int) (string, error) {
	row := database.Postgres.QueryRow(ctx,
		`SELECT api_key_hash FROM users WHERE id = $1`,
		id,
	)

	var key string
	if err := row.Scan(&key); err != nil {
		return "", err
	}

	return key, nil
}

// ------------------------------
// 7. Get user full details (Admin panel)
// ------------------------------
func GetUserDetails(ctx context.Context, id int) (map[string]interface{}, error) {
	row := database.Postgres.QueryRow(ctx, `
        SELECT id, username, credits, role, created_at, api_key_hash 
        FROM users WHERE id = $1
    `, id)

	var u models.User
	if err := row.Scan(&u.ID, &u.Username, &u.Credits, &u.Role, &u.CreatedAt, &u.ApiKey); err != nil {
		return nil, err
	}

	today, _ := GetDailyUsage(ctx, id)
	total, _ := TotalRequestsByUser(ctx, id)
	used, _ := TotalCreditsUsedByUser(ctx, id)

	return map[string]interface{}{
		"id":             u.ID,
		"username":       u.Username,
		"credits":        u.Credits,
		"status":         u.Role,
		"created_at":     u.CreatedAt,
		"api_key":        u.ApiKey,
		"today_requests": today.Requests,
		"total_requests": total,
		"credits_used":   used,
	}, nil
}

// ------------------------------
// 8. Fetch user by ID
// ------------------------------
func GetUserByID(ctx context.Context, id int) (*User, error) {
	row := database.Postgres.QueryRow(ctx,
		`SELECT id, username, api_key_hash, role, credits
         FROM users WHERE id = $1`,
		id,
	)

	var u User
	if err := row.Scan(&u.ID, &u.Username, &u.ApiKeyHash, &u.Role, &u.Credits); err != nil {
		return nil, err
	}

	u.Status = "active"
	return &u, nil
}

// ------------------------------
// 9. Delete user by ID
// ------------------------------
func DeleteUser(ctx context.Context, id int) error {
	_, err := database.Postgres.Exec(ctx,
		`DELETE FROM users WHERE id = $1`, id)
	return err
}
