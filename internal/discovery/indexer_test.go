package discovery

import (
	"testing"

	"github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/catalog"
)

func TestBuildCatalogRecord(t *testing.T) {

	entry := catalog.CatalogIndexEntry{
		CatalogID:    "electronics",
		EntryVersion: 1,
		CatalogType:  "retail",

		Baseline: &catalog.ArtifactRef{
			Version: 2,
			URL:     "http://localhost:8081/catalog/electronics/v2.json",
			Size:    1234,
			Digest:  "sha-256:abcdef",
		},
	}

	record, err := BuildCatalogRecord(
		"provider.example",
		entry,
	)

	if err != nil {
		t.Fatal(err)
	}

	if record.NodeID != "provider.example" {
		t.Fatalf(
			"unexpected node ID: %s",
			record.NodeID,
		)
	}

	if record.CatalogID != "electronics" {
		t.Fatalf(
			"unexpected catalog ID: %s",
			record.CatalogID,
		)
	}

	if record.CatalogType != "retail" {
		t.Fatalf(
			"unexpected catalog type: %s",
			record.CatalogType,
		)
	}

	if record.Version != 2 {
		t.Fatalf(
			"unexpected version: %d",
			record.Version,
		)
	}

	if record.URL !=
		"http://localhost:8081/catalog/electronics/v2.json" {
		t.Fatalf(
			"unexpected URL: %s",
			record.URL,
		)
	}

	if record.Digest != "sha-256:abcdef" {
		t.Fatalf(
			"unexpected digest: %s",
			record.Digest,
		)
	}
}

func TestBuildCatalogRecordRequiresBaseline(
	t *testing.T,
) {

	entry := catalog.CatalogIndexEntry{
		CatalogID: "electronics",
	}

	_, err := BuildCatalogRecord(
		"provider.example",
		entry,
	)

	if err == nil {
		t.Fatal(
			"expected missing baseline to fail",
		)
	}
}
