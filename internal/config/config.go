package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string

	// ClickHouse Cloud (HTTPS)
	ClickHouseHost     string
	ClickHouseUser     string
	ClickHousePassword string

	// Neon Postgres
	PostgresURL string

	// Upstash Redis (rediss://)
	RedisURL string

	ApiRateLimit int
}

func LoadConfig() Config {
	// Load .env if exists (local dev)
	godotenv.Load()

	return Config{
		ServerPort: get("SERVER_PORT", "8080"),

		// ✔ NEW ClickHouse format
		ClickHouseHost:     get("CLICKHOUSE_HOST", "localhost:8443"),
		ClickHouseUser:     get("CLICKHOUSE_USER", "default"),
		ClickHousePassword: get("CLICKHOUSE_PASSWORD", ""),

		// Neon Postgres
		PostgresURL: get("POSTGRES_URL", ""),

		// Upstash Redis
		RedisURL: get("REDIS_URL", ""),

		ApiRateLimit: 60,
	}
}

func get(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	log.Println("[CONFIG] Using fallback for:", key)
	return fallback
}
