package database

import (
	"context"
	"crypto/tls"
	"log"
	"time"

	"bigdata-api/internal/config"

	ch "github.com/ClickHouse/clickhouse-go/v2"
)

var ClickHouse ch.Conn

func ConnectClickHouse(cfg *config.Config) {
	conn, err := ch.Open(&ch.Options{
		Addr: []string{cfg.ClickHouseHost}, // host:port only e.g. zs62xxxxx:8443

		Auth: ch.Auth{
			Database: "default",
			Username: cfg.ClickHouseUser,
			Password: cfg.ClickHousePassword,
		},

		Protocol: ch.HTTP,       // ✔️ REQUIRED for ClickHouse Cloud
		TLS:      &tls.Config{}, // ✔️ HTTPS enabled

		Settings: ch.Settings{
			"max_execution_time": 60,
		},

		DialTimeout:  10 * time.Second,
		MaxIdleConns: 5,
		MaxOpenConns: 10,
	})

	if err != nil {
		log.Fatalf("❌ ClickHouse connection failed: %v", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		log.Fatalf("❌ ClickHouse ping failed: %v", err)
	}

	ClickHouse = conn
	log.Println("🟢 Connected to ClickHouse Cloud (HTTPS)")
}
