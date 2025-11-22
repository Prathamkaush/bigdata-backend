package ingestion

import (
	"strings"
	"time"

	"bigdata-api/internal/models"
)

// NormalizeString ensures safe lowercase, trimmed values
func NormalizeString(s string) string {
	return strings.TrimSpace(s)
}

// NormalizeEmail converts to lowercase + trims
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// NormalizePhone removes spaces, +, -, ()
func NormalizePhone(phone string) string {
	p := strings.TrimSpace(phone)
	p = strings.ReplaceAll(p, " ", "")
	p = strings.ReplaceAll(p, "-", "")
	p = strings.ReplaceAll(p, "+", "")
	p = strings.ReplaceAll(p, "(", "")
	p = strings.ReplaceAll(p, ")", "")
	return p
}

// NormalizeTimestamp returns ClickHouse-compatible format
func NormalizeTimestamp(ts string) string {
	if strings.TrimSpace(ts) == "" {
		return time.Now().UTC().Format("2006-01-02 15:04:05")
	}

	// Replace T and Z if user provides ISO timestamps
	ts = strings.ReplaceAll(ts, "T", " ")
	ts = strings.ReplaceAll(ts, "Z", "")

	return ts
}

// NormalizeRecord converts raw map[string]string into NormalizedRecord
func NormalizeRecord(raw map[string]string) models.NormalizedRecord {

	record := models.NormalizedRecord{
		DedupeKey:  "", // will be replaced in next step
		Source:     NormalizeString(raw["source"]),
		SourceID:   NormalizeString(raw["source_id"]),
		CustomerID: NormalizeString(raw["customer_id"]),

		FirstName: NormalizeString(
			firstNonEmpty(raw, []string{"first_name", "fname", "first", "given_name"}),
		),

		LastName: NormalizeString(
			firstNonEmpty(raw, []string{"last_name", "lname", "last", "surname"}),
		),

		Email: NormalizeEmail(
			firstNonEmpty(raw, []string{"email", "email_address", "mail"}),
		),

		Phone: NormalizePhone(
			firstNonEmpty(raw, []string{"phone", "mobile", "phone_number"}),
		),

		City: NormalizeString(
			firstNonEmpty(raw, []string{"city", "town"}),
		),

		State: NormalizeString(
			firstNonEmpty(raw, []string{"state", "province", "region"}),
		),

		Country: NormalizeString(
			firstNonEmpty(raw, []string{"country", "country_name"}),
		),

		IngestTS: NormalizeTimestamp(raw["ingest_ts"]),
	}

	// ⭐ IMPORTANT — generate dedupe key here
	record.DedupeKey = GenerateDedupeKey(&record)

	return record
}

func firstNonEmpty(raw map[string]string, keys []string) string {
	for _, k := range keys {
		if val, ok := raw[k]; ok && strings.TrimSpace(val) != "" {
			return val
		}
	}
	return ""
}
