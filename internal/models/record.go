package models

type NormalizedRecord struct {
	DedupeKey  string `json:"dedupe_key"`
	Source     string `json:"source"`
	SourceID   string `json:"source_id"`
	CustomerID string `json:"customer_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	IngestTS   string `json:"ingest_ts"`
}
