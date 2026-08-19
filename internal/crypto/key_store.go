package crypto

import (
	"crypto/ed25519"
	"fmt"
	"os"
)

func SavePrivateKey(path string, key ed25519.PrivateKey) error {
	return os.WriteFile(path, key, 0600)
}

func SavePublicKey(path string, key ed25519.PublicKey) error {
	return os.WriteFile(path, key, 0644)
}

func LoadPrivateKey(path string) (ed25519.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read private key: %w", err)
	}

	if len(data) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid private key size")
	}

	return ed25519.PrivateKey(data), nil
}

func LoadPublicKey(path string) (ed25519.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read public key: %w", err)
	}

	if len(data) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid public key size")
	}

	return ed25519.PublicKey(data), nil
}
