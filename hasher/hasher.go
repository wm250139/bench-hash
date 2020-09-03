package hasher

import (
	"crypto/sha256"
	"encoding/hex"
)

func String(input string) string {
	sum := sha256.Sum256([]byte(input))
	bytes := sum[:20]
	return hex.EncodeToString(bytes)
}
