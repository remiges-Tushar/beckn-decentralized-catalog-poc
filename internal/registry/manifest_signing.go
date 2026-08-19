package registry

import (
	"crypto/ed25519"
	"fmt"

	catalogcrypto "github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/crypto"
)

type ManifestSignature struct {
	KeyID     string
	Algorithm string
	Value     string
}

func ManifestSigningPayload(
	manifest Manifest,
) ([]byte, error) {
	// The proof must never be included in the signed payload.
	manifest.Proof = nil

	payload, err := catalogcrypto.CanonicalizeJSON(manifest)
	if err != nil {
		return nil, fmt.Errorf(
			"canonicalize manifest: %w",
			err,
		)
	}

	return payload, nil
}

func SignManifest(
	manifest Manifest,
	keyID string,
	privateKey ed25519.PrivateKey,
) (*ManifestSignature, error) {

	payload, err := ManifestSigningPayload(manifest)
	if err != nil {
		return nil, err
	}

	signature := catalogcrypto.Sign(
		privateKey,
		payload,
	)

	return &ManifestSignature{
		KeyID:     keyID,
		Algorithm: "Ed25519",
		Value:     signature,
	}, nil
}
