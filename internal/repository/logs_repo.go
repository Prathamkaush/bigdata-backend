package repository

import (
	"bigdata-api/internal/database"
	"context"
	"encoding/json"
	"time"
)

func LogAPIRequest(
	ctx context.Context,
	userID int,
	endpoint string,
	method string,
	statusCode int,
	responseMs int,
	requestBody string,
) error {

	params := map[string]interface{}{
		"method":       method,
		"status_code":  statusCode,
		"response_ms":  responseMs,
		"request_body": requestBody,
	}

	// IMPORTANT: Convert params map → JSON
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return err
	}

	_, err = database.Postgres.Exec(ctx, `
        INSERT INTO api_logs (user_id, endpoint, params, credits_used, created_at)
        VALUES ($1, $2, $3::jsonb, $4, NOW())
    `,
		userID,
		endpoint,
		jsonParams,
		1, // default credit
	)

	return err
}

// FetchLogs returns recent logs in a simple shape the frontend can use
func FetchLogs(ctx context.Context) ([]map[string]interface{}, error) {
	rows, err := database.Postgres.Query(ctx, `
        SELECT id, user_id, endpoint, params, credits_used, created_at
        FROM api_logs
        ORDER BY created_at DESC
        LIMIT 200
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []map[string]interface{}

	for rows.Next() {
		var id, userID, creditsUsed int
		var endpoint string
		var paramsJson []byte
		var createdAt time.Time

		err := rows.Scan(&id, &userID, &endpoint, &paramsJson, &creditsUsed, &createdAt)
		if err != nil {
			continue
		}

		logs = append(logs, map[string]interface{}{
			"id":           id,
			"user_id":      userID,
			"endpoint":     endpoint,
			"params":       string(paramsJson),
			"credits_used": creditsUsed,
			"created_at":   createdAt.Format(time.RFC3339),
		})
	}

	return logs, nil
}

func FetchLogsByUser(ctx context.Context, userID int) ([]map[string]interface{}, error) {
	rows, err := database.Postgres.Query(ctx, `
        SELECT id, endpoint, params, credits_used, created_at
        FROM api_logs
        WHERE user_id = $1
        ORDER BY created_at DESC
        LIMIT 50
    `, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []map[string]interface{}

	for rows.Next() {
		var id int
		var endpoint string
		var paramsJson []byte
		var credits int
		var createdAt time.Time

		err := rows.Scan(&id, &endpoint, &paramsJson, &credits, &createdAt)
		if err != nil {
			continue
		}

		logs = append(logs, map[string]interface{}{
			"id":           id,
			"endpoint":     endpoint,
			"credits_used": credits,
			"params":       string(paramsJson),
			"created_at":   createdAt.Format(time.RFC3339),
		})
	}

	return logs, nil
}
