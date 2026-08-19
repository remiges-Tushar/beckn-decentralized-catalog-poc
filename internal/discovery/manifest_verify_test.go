package discovery

import (
	"testing"

	catalogcrypto "github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/crypto"
)

func TestVerifyManifest(t *testing.T) {

	keyPair, err := catalogcrypto.GenerateKeyPair(
		"provider-key-1",
	)
	if err != nil {
		t.Fatal(err)
	}

	manifest := Manifest{
		DeDiVersion: "0.1",
		Type:        "dedi-manifest",
		Domain:      "provider.example",

		Keys: []ManifestKey{
			{
				KID: "provider-key-1",
				KTY: "OKP",
				CRV: "Ed25519",
				X: catalogcrypto.EncodePublicKey(
					keyPair.PublicKey,
				),
			},
		},

		Files: []ManifestFile{
			{
				Registry: "beckn-subscriber",
				URL:      "http://localhost:8081/dedi/beckn-subscriber.dedi.json",
				Digest:   "sha-256:test",
			},
		},
	}
	payload, err := ManifestSigningPayload(
		manifest,
	)
	if err != nil {
		t.Fatal(err)
	}

	signature := catalogcrypto.Sign(
		keyPair.PrivateKey,
		payload,
	)
	manifest.Proof = &ManifestProof{
		VerificationMethod: "provider-key-1",
		Canonicalization:   "JCS",
		JWS:                signature,
	}

	result, err := VerifyManifest(
		manifest,
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.KeyID != "provider-key-1" {
		t.Fatalf(
			"unexpected key ID: %s",
			result.KeyID,
		)
	}
}

func TestVerifyManifestRejectsTampering(
	t *testing.T,
) {
	keyPair, err := catalogcrypto.GenerateKeyPair(
		"provider-key-1",
	)
	if err != nil {
		t.Fatal(err)
	}

	manifest := Manifest{
		DeDiVersion: "0.1",
		Type:        "dedi-manifest",
		Domain:      "provider.example",

		Keys: []ManifestKey{
			{
				KID: "provider-key-1",
				KTY: "OKP",
				CRV: "Ed25519",
				X: catalogcrypto.EncodePublicKey(
					keyPair.PublicKey,
				),
			},
		},

		Files: []ManifestFile{
			{
				Registry: "beckn-subscriber",
				URL:      "http://localhost:8081/dedi/beckn-subscriber.dedi.json",
				Digest:   "sha-256:test",
			},
		},
	}

	payload, err := ManifestSigningPayload(
		manifest,
	)
	if err != nil {
		t.Fatal(err)
	}

	signature := catalogcrypto.Sign(
		keyPair.PrivateKey,
		payload,
	)

	manifest.Proof = &ManifestProof{
		VerificationMethod: "provider-key-1",
		Canonicalization:   "JCS",
		JWS:                signature,
	}

	// Attack: change the Provider identity.
	manifest.Domain = "attacker.example"

	_, err = VerifyManifest(manifest)

	if err == nil {
		t.Fatal(
			"expected tampered manifest to fail",
		)
	}
}

func TestVerifyManifestRejectsUnknownKey(
	t *testing.T,
) {
	keyPair, err := catalogcrypto.GenerateKeyPair(
		"provider-key-1",
	)
	if err != nil {
		t.Fatal(err)
	}

	manifest := Manifest{
		DeDiVersion: "0.1",
		Type:        "dedi-manifest",
		Domain:      "provider.example",

		Keys: []ManifestKey{
			{
				KID: "provider-key-1",
				KTY: "OKP",
				CRV: "Ed25519",
				X: catalogcrypto.EncodePublicKey(
					keyPair.PublicKey,
				),
			},
		},

		Proof: &ManifestProof{
			VerificationMethod: "unknown-key",
			Canonicalization:   "JCS",
			JWS:                "fake-signature",
		},
	}

	_, err = VerifyManifest(manifest)

	if err == nil {
		t.Fatal(
			"expected unknown key to fail",
		)
	}
}

func TestVerifyManifestRejectsUnsigned(
	t *testing.T,
) {
	manifest := Manifest{
		DeDiVersion: "0.1",
		Type:        "dedi-manifest",
		Domain:      "provider.example",
	}

	_, err := VerifyManifest(manifest)

	if err == nil {
		t.Fatal(
			"expected unsigned manifest to fail",
		)
	}
}
