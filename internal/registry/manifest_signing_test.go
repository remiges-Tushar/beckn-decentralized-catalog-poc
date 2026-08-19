package registry

import (
	"testing"

	catalogcrypto "github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/crypto"
)

func TestManifestSigningPayloadExcludesProof(t *testing.T) {

	keyPair, err := catalogcrypto.GenerateKeyPair("test-key")
	if err != nil {
		t.Fatal(err)
	}

	manifest := Manifest{
		DeDiVersion: "0.1",
		Type:        "dedi-manifest",
		Domain:      "provider.example",

		Proof: &ManifestProof{
			VerificationMethod: "provider-key-1",
			Canonicalization:   "JCS",
			JWS:                "THIS-MUST-NOT-BE-SIGNED",
		},
	}

	payload1, err := ManifestSigningPayload(manifest)
	if err != nil {
		t.Fatal(err)
	}

	manifest.Proof.JWS = "SOME-COMPLETELY-DIFFERENT-VALUE"

	payload2, err := ManifestSigningPayload(manifest)
	if err != nil {
		t.Fatal(err)
	}

	if string(payload1) != string(payload2) {
		t.Fatal("proof affected signing payload")
	}

	_, err = SignManifest(
		manifest,
		keyPair.KeyID,
		keyPair.PrivateKey,
	)
	if err != nil {
		t.Fatal(err)
	}
}
