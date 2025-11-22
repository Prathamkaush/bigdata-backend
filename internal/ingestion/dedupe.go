package ingestion

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"

	"bigdata-api/internal/models"
)

func GenerateDedupeKey(r *models.NormalizedRecord) string {

	// priority: email > phone > customer_id > source_id
	keyParts := []string{
		strings.ToLower(strings.TrimSpace(r.Email)),
		strings.ToLower(strings.TrimSpace(r.Phone)),
		strings.ToLower(strings.TrimSpace(r.CustomerID)),
		strings.ToLower(strings.TrimSpace(r.SourceID)),
	}

	// join with ::
	joined := strings.Join(keyParts, "::")

	// SHA1 hash so it's always fixed-length
	h := sha1.Sum([]byte(joined))
	return hex.EncodeToString(h[:])
}
