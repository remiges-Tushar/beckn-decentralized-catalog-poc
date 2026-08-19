package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

type KeyPair struct {
	KeyID     string
	PublicKey ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

func GenerateKeyPair(keyID string) (*KeyPair, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate ed25519 key: %w", err)
	}

	return &KeyPair{
		KeyID:     keyID,
		PublicKey: publicKey,
		PrivateKey: privateKey,
	}, nil
}

func Sign(privateKey ed25519.PrivateKey, payload []byte) string {
	signature := ed25519.Sign(privateKey, payload)
	return base64.RawURLEncoding.EncodeToString(signature)
}

func Verify(
	publicKey ed25519.PublicKey,
	payload []byte,
	signature string,
) bool {
	sig, err := base64.RawURLEncoding.DecodeString(signature)
	if err != nil {
		return false
	}

	return ed25519.Verify(publicKey, payload, sig)
}