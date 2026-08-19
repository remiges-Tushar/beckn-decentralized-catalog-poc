package catalog

import (
	"crypto/ed25519"
	"fmt"

	catalogcrypto "github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/crypto"
)

func SignCatalogFile(
	file CatalogFile,
	keyID string,
	privateKey ed25519.PrivateKey,
) (CatalogFile, error) {

	file.Signature = nil

	payload, err := catalogcrypto.CanonicalizeJSON(file)
	if err != nil {
		return CatalogFile{}, fmt.Errorf(
			"canonicalize catalog file: %w",
			err,
		)
	}

	signature := catalogcrypto.Sign(privateKey, payload)

	file.Signature = &Signature{
		KeyID:            keyID,
		Canonicalization: "JCS",
		Value:            signature,
	}

	return file, nil
}
