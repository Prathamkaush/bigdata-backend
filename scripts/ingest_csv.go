package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"bigdata-api/internal/config"
	"bigdata-api/internal/database"
	"bigdata-api/internal/ingestion"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("❌ Usage: go run scripts/ingest_csv.go <file.csv>")
	}

	filePath := os.Args[1]
	fmt.Println("📑 Reading CSV:", filePath)

	cfg := config.LoadConfig()
	database.ConnectClickHouse(&cfg)

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("❌ Failed to open CSV:", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.TrimLeadingSpace = true

	// Read header
	headers, err := reader.Read()
	if err != nil {
		log.Fatal("❌ Failed to read CSV header:", err)
	}

	batch := [][]interface{}{}
	batchSize := 5000
	inserted := 0

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		// Convert row to map[string]string
		raw := map[string]string{}
		for i, h := range headers {
			raw[h] = record[i]
		}

		// Normalize
		rec := ingestion.NormalizeRecord(raw)
		rec.DedupeKey = ingestion.GenerateDedupeKey(&rec)

		// Convert to CH row
		row := []interface{}{
			rec.DedupeKey,
			rec.Source,
			rec.SourceID,
			rec.CustomerID,
			rec.FirstName,
			rec.LastName,
			rec.Email,
			rec.Phone,
			rec.City,
			rec.State,
			rec.Country,
			rec.IngestTS,
		}

		batch = append(batch, row)

		if len(batch) >= batchSize {
			err := ingestion.BatchInsert(batch)
			if err != nil {
				log.Fatal("❌ Batch insert failed:", err)
			}
			inserted += len(batch)
			batch = [][]interface{}{}
		}
	}

	// Insert remaining rows
	if len(batch) > 0 {
		err := ingestion.BatchInsert(batch)
		if err != nil {
			log.Fatal("❌ Final batch failed:", err)
		}
		inserted += len(batch)
	}

	fmt.Printf("✅ Successfully inserted %d records into ClickHouse\n", inserted)
}
