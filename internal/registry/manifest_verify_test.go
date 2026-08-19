package registry

import (
	"testing"

	catalogcrypto "github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/crypto"
)

func TestManifestSignAndVerify(t *testing.T) {
	keyPair, err := catalogcrypto.GenerateKeyPair("provider-key-1")
	if err != nil {
		t.Fatal(err)
	}

	manifest := Manifest{
		DeDiVersion: "0.1",
		Type:        "dedi-manifest",
		Domain:      "provider.example",
	}

	signature, err := SignManifest(
		manifest,
		keyPair.KeyID,
		keyPair.PrivateKey,
	)
	if err != nil {
		t.Fatal(err)
	}

	if signature.Value == "" {
		t.Fatal("expected signature")
	}

	err = VerifyManifestSignature(
		manifest,
		signature.Value,
		keyPair.PublicKey,
	)

	if err != nil {
		t.Fatal(err)
	}
}

func TestManifestTamperingFails(t *testing.T) {
	keyPair, err := catalogcrypto.GenerateKeyPair("provider-key-1")
	if err != nil {
		t.Fatal(err)
	}

	manifest := Manifest{
		DeDiVersion: "0.1",
		Type:        "dedi-manifest",
		Domain:      "provider.example",
	}

	signature, err := SignManifest(
		manifest,
		keyPair.KeyID,
		keyPair.PrivateKey,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Tamper after signing.
	manifest.Domain = "attacker.example"

	err = VerifyManifestSignature(
		manifest,
		signature.Value,
		keyPair.PublicKey,
	)

	if err == nil {
		t.Fatal("expected tampered manifest to fail")
	}
}
