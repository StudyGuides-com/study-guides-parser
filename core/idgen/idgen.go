package idgen


import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/lucsky/cuid"
)

// GetCUID generates a unique ID using the CUID algorithm
func NewCUID() string {
	return cuid.New()
}

func HashFrom(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}