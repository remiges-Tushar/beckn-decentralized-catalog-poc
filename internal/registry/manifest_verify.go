package registry

import (
	"crypto/ed25519"
	"fmt"

	catalogcrypto "github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/crypto"
)

func VerifyManifestSignature(
	manifest Manifest,
	signature string,
	publicKey ed25519.PublicKey,
) error {

	payload, err := ManifestSigningPayload(manifest)
	if err != nil {
		return fmt.Errorf(
			"build manifest signing payload: %w",
			err,
		)
	}

	if !catalogcrypto.Verify(
		publicKey,
		payload,
		signature,
	) {
		return fmt.Errorf("invalid manifest signature")
	}

	return nil
}
