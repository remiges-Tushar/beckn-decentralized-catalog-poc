package discovery

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetCatalogIndexURL(t *testing.T) {

	subscriber := SubscriberRecord{
		Meta: SubscriberMeta{
			CatalogIndexURLs: []CatalogIndexURL{
				{
					URL: "http://example.com/catalog/index.json",
				},
			},
		},
	}

	url, err := GetCatalogIndexURL(
		subscriber,
	)

	if err != nil {
		t.Fatal(err)
	}

	expected := "http://example.com/catalog/index.json"

	if url != expected {
		t.Fatalf(
			"expected %s, got %s",
			expected,
			url,
		)
	}
}

func TestGetCatalogIndexURLMissing(
	t *testing.T,
) {

	subscriber := SubscriberRecord{}

	_, err := GetCatalogIndexURL(
		subscriber,
	)

	if err == nil {
		t.Fatal(
			"expected missing catalog index URL to fail",
		)
	}
}

func TestFetchCatalogIndex(t *testing.T) {
	indexJSON, err := os.ReadFile(
		"../../storage/provider/catalog/index.json",
	)
	if err != nil {
		t.Fatalf(
			"read catalog index: %v",
			err,
		)
	}

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				if r.URL.Path != "/catalog/index.json" {
					http.NotFound(w, r)
					return
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)

				_, _ = w.Write(indexJSON)
			},
		),
	)

	defer server.Close()

	subscriber := SubscriberRecord{
		Meta: SubscriberMeta{
			CatalogIndexURLs: []CatalogIndexURL{
				{
					URL: server.URL + "/catalog/index.json",
				},
			},
		},
	}

	crawler := NewCrawler()

	index, err := crawler.FetchCatalogIndex(
		subscriber,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(index.Catalogs) == 0 {
		t.Fatal("expected catalog entries")
	}
}
