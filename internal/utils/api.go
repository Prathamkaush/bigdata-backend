package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateApiKey() string {
	key := make([]byte, 32)
	rand.Read(key)
	return hex.EncodeToString(key)
}
