package discovery

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchManifest(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				if r.URL.Path != "/.well-known/dedi.json" {
					http.NotFound(w, r)
					return
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)

				_, _ = w.Write([]byte(`{
					"dedi_version": "0.1",
					"type": "dedi-manifest",
					"domain": "provider.example",
					"keys": [
						{
							"kid": "provider-key-1",
							"kty": "OKP",
							"crv": "Ed25519",
							"x": "test-public-key"
						}
					],
					"updated_at": "2026-08-19T10:00:00Z",
					"next_update": "2026-08-20T10:00:00Z",
					"files": [
						{
							"registry": "beckn-subscriber",
							"url": "http://example.com/subscriber.json",
							"digest": "sha-256:test"
						}
					]
				}`))
			},
		),
	)

	defer server.Close()

	crawler := NewCrawler()

	manifest, err := crawler.FetchManifest(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	if manifest.Domain != "provider.example" {
		t.Fatalf(
			"unexpected domain: %s",
			manifest.Domain,
		)
	}

	if len(manifest.Keys) != 1 {
		t.Fatalf(
			"expected 1 key, got %d",
			len(manifest.Keys),
		)
	}

	if manifest.Keys[0].KID != "provider-key-1" {
		t.Fatalf(
			"unexpected key ID: %s",
			manifest.Keys[0].KID,
		)
	}

	if len(manifest.Files) != 1 {
		t.Fatalf(
			"expected 1 file, got %d",
			len(manifest.Files),
		)
	}

	if manifest.Files[0].Registry != "beckn-subscriber" {
		t.Fatalf(
			"unexpected registry: %s",
			manifest.Files[0].Registry,
		)
	}
}

func TestValidateManifestRejectsMissingDomain(t *testing.T) {

	manifest := Manifest{
		DeDiVersion: "0.1",
		Type:        "dedi-manifest",

		Keys: []ManifestKey{
			{
				KID: "provider-key-1",
				KTY: "OKP",
				CRV: "Ed25519",
				X:   "test-key",
			},
		},

		Files: []ManifestFile{
			{
				Registry: "beckn-subscriber",
				URL:      "http://example.com/subscriber.json",
				Digest:   "sha-256:test",
			},
		},
	}

	err := ValidateManifest(manifest)

	if err == nil {
		t.Fatal("expected validation to fail")
	}
}
