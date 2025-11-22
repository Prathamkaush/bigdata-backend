package main

import (
	"bigdata-api/internal/database"
	"context"
	"log"
)

func main() {
	db := database.Postgres
	ctx := context.Background()

	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			api_key TEXT UNIQUE NOT NULL,
			status TEXT DEFAULT 'active',
			created_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS user_credits (
			user_id INT REFERENCES users(id) ON DELETE CASCADE,
			credits INT NOT NULL DEFAULT 0,
			PRIMARY KEY (user_id)
		);`,
		`CREATE TABLE IF NOT EXISTS credit_logs (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES users(id),
			credits_used INT NOT NULL,
			endpoint TEXT,
			created_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS api_logs (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES users(id),
			endpoint TEXT,
			request_body JSONB,
			response_time_ms INT,
			created_at TIMESTAMP DEFAULT NOW()
		);`,
	}

	for _, q := range queries {
		_, err := db.Exec(ctx, q)
		if err != nil {
			log.Fatal("Migration failed:", err)
		}
	}

	log.Println("Postgres migration completed! ✅")
}
