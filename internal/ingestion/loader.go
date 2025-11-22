package ingestion

import (
	"bigdata-api/internal/database"
	"context"
	"log"
)

// BatchInsert inserts many rows into ClickHouse efficiently
func BatchInsert(rows [][]interface{}) error {
	if len(rows) == 0 {
		return nil
	}

	ctx := context.Background()

	// Prepare batch insert
	batch, err := database.ClickHouse.PrepareBatch(ctx, `
		INSERT INTO master_records (
			dedupe_key,
			source,
			source_id,
			customer_id,
			first_name,
			last_name,
			email,
			phone,
			city,
			state,
			country,
			ingest_ts
		)
	`)
	if err != nil {
		log.Println("❌ Failed to prepare batch:", err)
		return err
	}

	// Add all rows
	for _, row := range rows {
		if err := batch.Append(row...); err != nil {
			log.Println("❌ Batch append failed:", err)
			return err
		}
	}

	// Send them to ClickHouse
	if err := batch.Send(); err != nil {
		log.Println("❌ Batch send failed:", err)
		return err
	}

	return nil
}
