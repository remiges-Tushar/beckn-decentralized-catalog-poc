package discovery

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalculateDigest(t *testing.T) {
	data := []byte(`{"hello":"world"}`)

	digest := CalculateDigest(data)

	if !strings.HasPrefix(
		digest,
		"sha-256:",
	) {
		t.Fatalf(
			"unexpected digest: %s",
			digest,
		)
	}
}

func TestVerifyDigest(t *testing.T) {
	data := []byte(`{"hello":"world"}`)

	digest := CalculateDigest(data)

	if err := VerifyDigest(
		data,
		digest,
	); err != nil {
		t.Fatal(err)
	}
}

func TestVerifyDigestRejectsTampering(
	t *testing.T,
) {
	original := []byte(
		`{"hello":"world"}`,
	)

	digest := CalculateDigest(original)

	tampered := []byte(
		`{"hello":"attacker"}`,
	)

	err := VerifyDigest(
		tampered,
		digest,
	)

	if err == nil {
		t.Fatal(
			"expected digest verification to fail",
		)
	}
}

func TestFetchSubscriber(t *testing.T) {
	// Subscriber record that the Provider will serve.
	subscriberJSON := []byte(`{
		"record_name": "beckn-subscriber",
		"details": {
			"subscriber_id": "provider.example",
			"url": "http://localhost:8081",
			"type": "BPP",
			"domain": "retail",
			"countries": ["IND"],
			"signing_public_key": "test-key"
		},
		"meta": {
			"catalog_index_urls": [
				{
					"url": "http://localhost:8081/catalog/index.json"
				}
			]
		}
	}`)

	// Calculate the digest of the exact bytes
	// that the Provider will serve.
	digest := CalculateDigest(subscriberJSON)

	// Create a test HTTP server.
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				if r.URL.Path != "/subscriber" {
					http.NotFound(w, r)
					return
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)

				_, _ = w.Write(subscriberJSON)
			},
		),
	)

	defer server.Close()

	// Create a manifest that points to the
	// Subscriber served by the test server.
	manifest := Manifest{
		Files: []ManifestFile{
			{
				Registry: "beckn-subscriber",
				URL:      server.URL + "/subscriber",
				Digest:   digest,
			},
		},
	}

	// Create the Discovery Service crawler.
	crawler := NewCrawler()

	// Fetch Subscriber and verify its digest.
	subscriber, err := crawler.FetchSubscriber(
		manifest,
	)

	if err != nil {
		t.Fatal(err)
	}

	// Verify parsed Subscriber fields.
	if subscriber.RecordName != "beckn-subscriber" {
		t.Fatalf(
			"unexpected record name: %s",
			subscriber.RecordName,
		)
	}

	if subscriber.Details.SubscriberID != "provider.example" {
		t.Fatalf(
			"unexpected subscriber ID: %s",
			subscriber.Details.SubscriberID,
		)
	}

	if subscriber.Details.Type != "BPP" {
		t.Fatalf(
			"unexpected subscriber type: %s",
			subscriber.Details.Type,
		)
	}

	if subscriber.Details.Domain != "retail" {
		t.Fatalf(
			"unexpected domain: %s",
			subscriber.Details.Domain,
		)
	}

	if len(subscriber.Meta.CatalogIndexURLs) != 1 {
		t.Fatalf(
			"expected 1 catalog index URL, got %d",
			len(subscriber.Meta.CatalogIndexURLs),
		)
	}

	expectedURL := "http://localhost:8081/catalog/index.json"

	if subscriber.Meta.CatalogIndexURLs[0].URL != expectedURL {
		t.Fatalf(
			"unexpected catalog index URL: %s",
			subscriber.Meta.CatalogIndexURLs[0].URL,
		)
	}
}
