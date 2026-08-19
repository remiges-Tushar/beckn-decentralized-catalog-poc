package catalog

import (
	"crypto/ed25519"
	"fmt"

	catalogcrypto "github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/crypto"
)

func VerifyCatalogFile(
	file CatalogFile,
	publicKey ed25519.PublicKey,
) (bool, error) {

	if file.Signature == nil {
		return false, fmt.Errorf("catalog file has no signature")
	}

	signature := file.Signature.Value
	file.Signature = nil

	payload, err := catalogcrypto.CanonicalizeJSON(file)
	if err != nil {
		return false, fmt.Errorf(
			"canonicalize catalog file: %w",
			err,
		)
	}

	return catalogcrypto.Verify(publicKey, payload, signature), nil
}
