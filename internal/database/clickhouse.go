package database

import (
	"crypto/tls"
	"log"
	"os"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func ConnectClickHouse() clickhouse.Conn {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{
			os.Getenv("CLICKHOUSE_HOST"), // host:8443
		},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: os.Getenv("CLICKHOUSE_USER"),
			Password: os.Getenv("CLICKHOUSE_PASSWORD"),
		},
		TLS: &tls.Config{
			InsecureSkipVerify: true, // REQUIRED for ClickHouse Cloud
		},
		Protocol: clickhouse.HTTP, // IMPORTANT!!
	})

	if err != nil {
		log.Fatal("❌ ClickHouse connect failed:", err)
	}

	// ping
	if err := conn.Ping(); err != nil {
		log.Fatal("❌ ClickHouse ping failed:", err)
	}

	log.Println("✅ ClickHouse connected!")
	return conn
}
