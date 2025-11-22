package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"bigdata-api/internal/config"
	"bigdata-api/internal/database"
	"bigdata-api/internal/ingestion"
)

// XML structure that can capture ANY fields dynamically
type XMLRecord struct {
	Fields []XMLField `xml:",any"`
}

type XMLField struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

type XMLRecords struct {
	Records []XMLRecord `xml:"record"`
}

// Convert XMLRecord → map[string]string
func (r XMLRecord) ToMap() map[string]string {
	m := map[string]string{}
	for _, f := range r.Fields {
		m[f.XMLName.Local] = f.Value
	}
	return m
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("❌ Usage: go run scripts/ingest_xml.go <file.xml>")
	}

	filePath := os.Args[1]
	fmt.Println("📑 Reading XML:", filePath)

	cfg := config.LoadConfig()
	database.ConnectClickHouse(&cfg)

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("❌ Failed to read XML:", err)
	}

	var records XMLRecords

	// Try parsing multiple <record> elements
	err = xml.Unmarshal(data, &records)
	if err != nil || len(records.Records) == 0 {
		// Try parsing single <record>
		var single XMLRecord
		if err2 := xml.Unmarshal(data, &single); err2 != nil {
			log.Fatal("❌ Invalid XML:", err2)
		}
		records.Records = append(records.Records, single)
	}

	fmt.Println("📦 Total XML records:", len(records.Records))

	batch := [][]interface{}{}
	inserted := 0
	batchSize := 5000

	for _, xr := range records.Records {
		raw := xr.ToMap()

		rec := ingestion.NormalizeRecord(raw)
		rec.DedupeKey = ingestion.GenerateDedupeKey(&rec)

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
			if err := ingestion.BatchInsert(batch); err != nil {
				log.Fatal("❌ Batch insert failed:", err)
			}
			inserted += len(batch)
			batch = [][]interface{}{}
		}
	}

	if len(batch) > 0 {
		if err := ingestion.BatchInsert(batch); err != nil {
			log.Fatal("❌ Final batch failed:", err)
		}
		inserted += len(batch)
	}

	fmt.Printf("✅ Successfully inserted %d XML records into ClickHouse\n", inserted)
}
