package registry

import (
	"encoding/json"
	"testing"
)

func TestSubscriberContainsCatalogIndexURL(t *testing.T) {
	record := SubscriberRecord{
		RecordName: "beckn-subscriber",

		Details: SubscriberDetails{
			SubscriberID:     "provider.example",
			URL:              "http://localhost:8081",
			Type:             "BPP",
			Domain:           "retail",
			Countries:        []string{"IND"},
			SigningPublicKey: "test-key",
		},

		Meta: SubscriberMeta{
			CatalogIndexURLs: []CatalogIndexURL{
				{
					URL: "http://localhost:8081/catalog/index.json",
				},
			},
		},
	}

	data, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}

	var decoded SubscriberRecord

	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}

	if len(decoded.Meta.CatalogIndexURLs) != 1 {
		t.Fatalf(
			"expected 1 catalog index URL, got %d",
			len(decoded.Meta.CatalogIndexURLs),
		)
	}

	if decoded.Meta.CatalogIndexURLs[0].URL !=
		"http://localhost:8081/catalog/index.json" {
		t.Fatal("unexpected catalog index URL")
	}
}
