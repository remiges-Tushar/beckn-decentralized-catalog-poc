package discovery

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/catalog"
)

func TestFetchCatalogFile(t *testing.T) {

	catalogJSON, err := os.ReadFile(
		"../../storage/provider/catalog/electronics/v1.json",
	)
	if err != nil {
		t.Fatalf(
			"read catalog file: %v",
			err,
		)
	}

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				if r.URL.Path != "/catalog/v1.json" {
					http.NotFound(w, r)
					return
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)

				_, _ = w.Write(catalogJSON)
			},
		),
	)

	defer server.Close()

	digest := CalculateDigest(
		catalogJSON,
	)

	entry := catalog.CatalogIndexEntry{
		Baseline: &catalog.ArtifactRef{
			Version: 1,
			URL:     server.URL + "/catalog/v1.json",
			Digest:  digest,
		},
	}

	crawler := NewCrawler()

	file, err := crawler.FetchAndVerifyCatalogFile(
		entry,
	)
	if err != nil {
		t.Fatal(err)
	}

	_ = file
}

func TestFetchCatalogFileRejectsTampering(
	t *testing.T,
) {

	original := []byte(
		`{"catalog":"original"}`,
	)

	tampered := []byte(
		`{"catalog":"tampered"}`,
	)

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				_, _ = w.Write(tampered)
			},
		),
	)

	defer server.Close()

	digest := CalculateDigest(
		original,
	)

	entry := catalog.CatalogIndexEntry{
		Baseline: &catalog.ArtifactRef{
			Version: 1,
			URL:     server.URL,
			Digest:  digest,
		},
	}

	crawler := NewCrawler()

	_, err := crawler.FetchAndVerifyCatalogFile(
		entry,
	)

	if err == nil {
		t.Fatal(
			"expected tampered catalog file to fail",
		)
	}
}
