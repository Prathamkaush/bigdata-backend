package database

import (
	"bigdata-api/internal/config"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Postgres *pgxpool.Pool

func ConnectPostgres(cfg *config.Config) {
	pool, err := pgxpool.New(context.Background(), cfg.PostgresURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect Postgres: %v", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("❌ Postgres ping failed: %v", err)
	}

	Postgres = pool
	log.Println("🟢 Connected to Postgres")
}
