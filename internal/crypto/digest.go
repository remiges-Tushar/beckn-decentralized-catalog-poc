package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

func Digest(data []byte) string {
	sum := sha256.Sum256(data)

	return "sha-256:" + hex.EncodeToString(sum[:])
}