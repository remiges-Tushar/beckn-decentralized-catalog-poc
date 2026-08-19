package discovery

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"

	catalogcrypto "github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/crypto"
)

func findManifestKey(
	manifest Manifest,
	keyID string,
) (ed25519.PublicKey, error) {

	for _, key := range manifest.Keys {
		if key.KID != keyID {
			continue
		}

		if key.KTY != "OKP" {
			return nil, fmt.Errorf(
				"unsupported key type: %s",
				key.KTY,
			)
		}

		if key.CRV != "Ed25519" {
			return nil, fmt.Errorf(
				"unsupported curve: %s",
				key.CRV,
			)
		}

		publicKey, err := base64.RawURLEncoding.DecodeString(
			key.X,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid public key encoding: %w",
				err,
			)
		}

		if len(publicKey) != ed25519.PublicKeySize {
			return nil, fmt.Errorf(
				"invalid Ed25519 public key size: %d",
				len(publicKey),
			)
		}

		return ed25519.PublicKey(publicKey), nil
	}

	return nil, fmt.Errorf(
		"manifest key not found: %s",
		keyID,
	)
}

func ManifestSigningPayload(
	manifest Manifest,
) ([]byte, error) {

	manifest.Proof = nil

	return catalogcrypto.CanonicalizeJSON(
		manifest,
	)
}

func VerifyManifest(
	manifest Manifest,
) (*ManifestVerificationResult, error) {

	if manifest.Proof == nil {
		return nil, fmt.Errorf(
			"manifest has no proof",
		)
	}

	if manifest.Proof.VerificationMethod == "" {
		return nil, fmt.Errorf(
			"manifest proof has no verification method",
		)
	}

	if manifest.Proof.JWS == "" {
		return nil, fmt.Errorf(
			"manifest proof has no JWS",
		)
	}

	publicKey, err := findManifestKey(
		manifest,
		manifest.Proof.VerificationMethod,
	)
	if err != nil {
		return nil, err
	}

	payload, err := ManifestSigningPayload(
		manifest,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"build signing payload: %w",
			err,
		)
	}

	if !catalogcrypto.Verify(
		publicKey,
		payload,
		manifest.Proof.JWS,
	) {
		return nil, fmt.Errorf(
			"manifest signature verification failed",
		)
	}

	return &ManifestVerificationResult{
		KeyID: manifest.Proof.VerificationMethod,
	}, nil
}
