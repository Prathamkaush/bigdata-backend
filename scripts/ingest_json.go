package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"bigdata-api/internal/config"
	"bigdata-api/internal/database"
	"bigdata-api/internal/ingestion"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("❌ Usage: go run scripts/ingest_json.go <file.json>")
	}

	filePath := os.Args[1]
	fmt.Println("📑 Reading JSON:", filePath)

	// Load config
	cfg := config.LoadConfig()

	// Connect ClickHouse
	database.ConnectClickHouse(&cfg)

	// Open file
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("❌ Failed to open JSON:", err)
	}
	defer f.Close()

	decoder := json.NewDecoder(bufio.NewReader(f))

	// Detect if it's a JSON array
	t, err := decoder.Token()
	if err != nil {
		log.Fatal("❌ Failed to read JSON:", err)
	}

	isArray := false
	if delim, ok := t.(json.Delim); ok && delim.String() == "[" {
		isArray = true
	}

	batch := [][]interface{}{}
	batchSize := 5000
	inserted := 0

	processRecord := func(raw map[string]string) {
		// Normalize
		rec := ingestion.NormalizeRecord(raw)

		// Generate dedupe key
		rec.DedupeKey = ingestion.GenerateDedupeKey(&rec)

		// Convert to ClickHouse insert row
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

	if isArray {
		// JSON Array
		var raw map[string]string
		for decoder.More() {
			err = decoder.Decode(&raw)
			if err != nil {
				log.Fatal("❌ JSON decode error:", err)
			}
			processRecord(raw)
		}
		decoder.Token() // consume closing ]
	} else {
		// JSONL mode
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Bytes()
			if len(line) == 0 {
				continue
			}

			var raw map[string]string
			if err := json.Unmarshal(line, &raw); err != nil {
				log.Fatal("❌ JSONL decode failed:", err)
			}

			processRecord(raw)
		}
	}

	// Final batch
	if len(batch) > 0 {
		err := ingestion.BatchInsert(batch)
		if err != nil {
			log.Fatal("❌ Final batch insert failed:", err)
		}
		inserted += len(batch)
	}

	fmt.Printf("✅ Successfully inserted %d JSON records into ClickHouse\n", inserted)
}
