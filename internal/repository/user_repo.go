package repository

import (
	"bigdata-api/internal/database"
	"bigdata-api/internal/models"
	"bigdata-api/internal/utils"
	"context"

	"strings"
)

// VALID ROLES
var validRoles = map[string]bool{
	"super_admin": true,
	"admin":       true,
	"sub_admin":   true,
	"viewer":      true,
}

// ------------------------------
// INTERNAL USER STRUCT
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
// Fetch user by API Hash
// ------------------------------
func GetUserByAPIHash(ctx context.Context, hash string) (*User, error) {

	row := database.Postgres.QueryRow(ctx, `
        SELECT id, username, api_key_hash, role, credits, status
        FROM users WHERE api_key_hash = $1
    `, hash)

	var u User
	if err := row.Scan(&u.ID, &u.Username, &u.ApiKeyHash, &u.Role, &u.Credits, &u.Status); err != nil {
		return nil, err
	}

	return &u, nil
}

// ------------------------------
// Create Admin/Sub-Admin User
// ------------------------------
func CreateUser(ctx context.Context, username string, role string, credits int) (*models.User, string, error) {

	role = strings.ToLower(strings.TrimSpace(role))
	if !validRoles[role] {
		role = "sub_admin" // default
	}

	if credits <= 0 {
		credits = 100
	}

	rawKey := utils.GenerateApiKey()
	hashed := utils.HashString(rawKey)

	row := database.Postgres.QueryRow(ctx, `
        INSERT INTO users (username, api_key_hash, credits, role, status)
        VALUES ($1, $2, $3, $4, 'active')
        RETURNING id, username, credits, role, status, created_at
    `, username, hashed, credits, role)

	var u models.User
	if err := row.Scan(&u.ID, &u.Username, &u.Credits, &u.Role, &u.Status, &u.CreatedAt); err != nil {
		return nil, "", err
	}

	return &u, rawKey, nil
}

// ------------------------------
// Get ALL Users for Admin Panel
// ------------------------------
func GetAllUsers(ctx context.Context) ([]models.User, error) {

	rows, err := database.Postgres.Query(ctx, `
        SELECT id, username, credits, created_at, api_key_hash, role, status
        FROM users ORDER BY id DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.User

	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID, &u.Username, &u.Credits, &u.CreatedAt,
			&u.APIKey, &u.Role, &u.Status,
		); err == nil {
			list = append(list, u)
		}
	}

	return list, nil
}

// ------------------------------
// Fetch user by username
// ------------------------------
func GetUserByName(ctx context.Context, username string) (*User, error) {

	row := database.Postgres.QueryRow(ctx, `
        SELECT id, username, api_key_hash, role, credits, status
        FROM users WHERE username = $1
    `, username)

	var u User
	if err := row.Scan(&u.ID, &u.Username, &u.ApiKeyHash, &u.Role, &u.Credits, &u.Status); err != nil {
		return nil, err
	}

	return &u, nil
}

// ------------------------------
// Update API Key
// ------------------------------
func UpdateAPIKey(ctx context.Context, userID int, hash string) error {
	_, err := database.Postgres.Exec(ctx,
		`UPDATE users SET api_key_hash = $1 WHERE id = $2`,
		hash, userID,
	)
	return err
}

// ------------------------------
// Fetch stored hashed key
// ------------------------------
func FetchAPIKey(ctx context.Context, id int) (string, error) {

	row := database.Postgres.QueryRow(ctx,
		`SELECT api_key_hash FROM users WHERE id = $1`,
		id,
	)

	var key string
	return key, row.Scan(&key)
}

// ------------------------------
// Admin Panel - User Details
// ------------------------------
func GetUserDetails(ctx context.Context, id int) (map[string]interface{}, error) {

	row := database.Postgres.QueryRow(ctx, `
        SELECT id, username, credits, role, status, created_at, api_key_hash
        FROM users WHERE id = $1
    `, id)

	var u models.User
	if err := row.Scan(&u.ID, &u.Username, &u.Credits, &u.Role, &u.Status, &u.CreatedAt, &u.APIKey); err != nil {
		return nil, err
	}

	today, _ := GetDailyUsage(ctx, id)
	total, _ := TotalRequestsByUser(ctx, id)
	used, _ := TotalCreditsUsedByUser(ctx, id)

	return map[string]interface{}{
		"id":             u.ID,
		"username":       u.Username,
		"credits":        u.Credits,
		"role":           u.Role,
		"status":         u.Status,
		"created_at":     u.CreatedAt,
		"api_key":        u.APIKey,
		"today_requests": today.Requests,
		"total_requests": total,
		"credits_used":   used,
	}, nil
}

// ------------------------------
// Get by ID
// ------------------------------
func GetUserByID(ctx context.Context, id int) (*User, error) {

	row := database.Postgres.QueryRow(ctx, `
        SELECT id, username, api_key_hash, role, credits, status
        FROM users WHERE id = $1
    `, id)

	var u User
	if err := row.Scan(&u.ID, &u.Username, &u.ApiKeyHash, &u.Role, &u.Credits, &u.Status); err != nil {
		return nil, err
	}

	return &u, nil
}

// ------------------------------
// Delete user
// ------------------------------
func DeleteUser(ctx context.Context, id int) error {
	_, err := database.Postgres.Exec(ctx,
		`DELETE FROM users WHERE id = $1`, id)
	return err
}

// =========================================================
// EXTRA FEATURES (NEW)
// =========================================================

// ------------------------------
// Change user role
// ------------------------------
func ChangeUserRole(ctx context.Context, id int, role string) error {
	_, err := database.Postgres.Exec(ctx,
		`UPDATE users SET role = $1 WHERE id = $2`,
		role, id,
	)
	return err
}

// ------------------------------
// Update Credits
// ------------------------------
func UpdateCredits(ctx context.Context, id int, credits int) error {

	_, err := database.Postgres.Exec(ctx,
		`UPDATE users SET credits = $1 WHERE id = $2`,
		credits, id,
	)
	return err
}

// ------------------------------
// Disable User
// ------------------------------
func DisableUser(ctx context.Context, id int) error {
	_, err := database.Postgres.Exec(ctx,
		`UPDATE users SET status = 'disabled' WHERE id = $1`, id,
	)
	return err
}

// Update user status
func UpdateUserStatus(ctx context.Context, id int, status string) error {
	_, err := database.Postgres.Exec(ctx,
		`UPDATE users SET status = $1 WHERE id = $2`,
		status, id,
	)
	return err
}

// Update user credits
func UpdateUserCredits(ctx context.Context, id int, credits int) error {
	_, err := database.Postgres.Exec(ctx,
		`UPDATE users SET credits = $1 WHERE id = $2`,
		credits, id,
	)
	return err
}
