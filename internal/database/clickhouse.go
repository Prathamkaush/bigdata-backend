package database

import (
	"context"
	"crypto/tls"
	"log"

	"github.com/ClickHouse/clickhouse-go/v2"
)

var ClickHouse clickhouse.Conn

func ConnectClickHouse(host, user, password string) {
	ctx := context.Background()

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{host}, // Example: "zs62tr9pg6.ap-south-1.aws.clickhouse.cloud:9440"

		Auth: clickhouse.Auth{
			Database: "default",
			Username: user,
			Password: password,
		},

		// TLS FIX HERE
		TLS: &tls.Config{
			InsecureSkipVerify: true,
		},

		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
	})

	if err != nil {
		log.Fatalf("❌ ClickHouse connect error: %v", err)
	}

	// Ping requires context
	if err := conn.Ping(ctx); err != nil {
		log.Fatalf("❌ ClickHouse ping error: %v", err)
	}

	log.Println("✅ ClickHouse connected successfully!")
	ClickHouse = conn
}
