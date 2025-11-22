package models

type NormalizedRecord struct {
	DedupeKey  string
	Source     string
	SourceID   string
	CustomerID string
	FirstName  string
	LastName   string
	Email      string
	Phone      string
	City       string
	State      string
	Country    string
	IngestTS   string
}
