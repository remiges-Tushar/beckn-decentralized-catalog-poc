package discovery

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	catalogcrypto "github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/crypto"
)

func TestFetchAndVerifyManifest(t *testing.T) {
	// Generate provider key.
	keyPair, err := catalogcrypto.GenerateKeyPair(
		"provider-key-1",
	)
	if err != nil {
		t.Fatal(err)
	}

	// Build manifest.
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
				URL:      "http://example/subscriber",
				Digest:   "sha-256:test",
			},
		},
	}

	// Build signing payload.
	payload, err := ManifestSigningPayload(manifest)
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

	// Marshal manifest.
	body, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}

	// Serve manifest.
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set(
					"Content-Type",
					"application/json",
				)

				_, _ = w.Write(body)
			},
		),
	)

	defer server.Close()

	// Crawl.
	crawler := NewCrawler()

	result, err := crawler.FetchAndVerifyManifest(
		server.URL,
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Domain != "provider.example" {
		t.Fatalf(
			"unexpected domain: %s",
			result.Domain,
		)
	}
}
