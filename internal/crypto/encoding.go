package crypto

import (
	"crypto/ed25519"
	"encoding/base64"
)

func EncodePublicKey(key ed25519.PublicKey) string {
	return base64.RawURLEncoding.EncodeToString(key)
}
