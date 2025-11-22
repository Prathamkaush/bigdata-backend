package main

import (
	"context"
	"log"
	"time"

	"bigdata-api/internal/config"
	"bigdata-api/internal/database"
)

func main() {
	// load config (reuse your config loader)
	cfg := config.LoadConfig()

	// connect (ensure database.ConnectClickHouse called)
	database.ConnectClickHouse(cfg)
	db := database.ClickHouse
	ctx := context.Background()

	// wait a bit to ensure connection ready
	time.Sleep(500 * time.Millisecond)

	ddl := `
CREATE TABLE IF NOT EXISTS default.master_records
(
    dedupe_key String,
    source String,
    source_id String,
    customer_id String DEFAULT '',
    first_name String,
    last_name String,
    email String,
    phone String,
    address String,
    city String,
    state String,
    postal_code String,
    country String,
    date_of_birth Date DEFAULT '1970-01-01',
    gender String,
    record_ts DateTime,
    ingest_ts DateTime,
    raw_payload String,
    data_quality_score Float32 DEFAULT 0.0,
    tenant_id String DEFAULT '',
    version UInt64 DEFAULT toUInt64(now())
)
ENGINE = ReplacingMergeTree(version)
PARTITION BY toYYYYMM(ingest_ts)
ORDER BY (dedupe_key, ingest_ts)
SETTINGS index_granularity = 8192;
`

	if err := db.Exec(ctx, ddl); err != nil {
		log.Fatalf("❌ failed to create master_records: %v", err)
	}
	log.Println("✅ master_records table created/verified")

	// Create ingest_daily_summary table + materialized view
	ddl2 := `
CREATE TABLE IF NOT EXISTS default.ingest_daily_summary
(
    day Date,
    source String,
    records UInt64
)
ENGINE = SummingMergeTree()
PARTITION BY toYYYYMM(day)
ORDER BY (day, source);
`
	if err := db.Exec(ctx, ddl2); err != nil {
		log.Fatalf("❌ failed to create ingest_daily_summary: %v", err)
	}
	log.Println("✅ ingest_daily_summary created")

	ddl3 := `
CREATE MATERIALIZED VIEW IF NOT EXISTS default.mv_ingest_daily_summary
TO default.ingest_daily_summary
AS
SELECT
    toDate(ingest_ts) AS day,
    source,
    count() AS records
FROM default.master_records
GROUP BY day, source;
`
	if err := db.Exec(ctx, ddl3); err != nil {
		log.Fatalf("❌ failed to create mv_ingest_daily_summary: %v", err)
	}
	log.Println("✅ mv_ingest_daily_summary created")

	// Optional projection creation (may require ALTER permission)
	// Note: Projections are supported in modern ClickHouse; skip if not supported
	proj := `
ALTER TABLE default.master_records
ADD PROJECTION IF NOT EXISTS pr_country_city
(
    SELECT
        dedupe_key,
        source,
        email,
        phone,
        city,
        country,
        ingest_ts
    ORDER BY (country, city, ingest_ts)
);
`
	_ = db.Exec(ctx, proj) // ignore error if not supported in older CH versions
	log.Println("ℹ️ attempted to add projection pr_country_city (may be skipped on older CH versions)")
}
