package catalog

import (
	"crypto/ed25519"
	"fmt"

	catalogcrypto "github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/crypto"
)

func SignIndexEntry(
	entry CatalogIndexEntry,
	keyID string,
	privateKey ed25519.PrivateKey,
) (CatalogIndexEntry, error) {

	// Signature must not be part of the signing input.
	entry.Signature = nil

	payload, err := catalogcrypto.CanonicalizeJSON(entry)
	if err != nil {
		return CatalogIndexEntry{}, fmt.Errorf(
			"canonicalize index entry: %w",
			err,
		)
	}

	signature := catalogcrypto.Sign(privateKey, payload)

	entry.Signature = &Signature{
		KeyID:            keyID,
		Canonicalization: "JCS",
		Value:            signature,
	}

	return entry, nil
}

func VerifyIndexEntry(
	entry CatalogIndexEntry,
	publicKey ed25519.PublicKey,
) bool {

	if entry.Signature == nil {
		return false
	}

	signature := entry.Signature.Value

	entry.Signature = nil

	payload, err := catalogcrypto.CanonicalizeJSON(entry)
	if err != nil {
		return false
	}

	return catalogcrypto.Verify(
		publicKey,
		payload,
		signature,
	)
}
