package discovery

import (
	"path/filepath"
	"testing"
)

func TestStoreSaveAndLoad(t *testing.T) {

	path := filepath.Join(
		t.TempDir(),
		"index.json",
	)

	store := NewStore(path)

	store.UpsertNode(NodeRecord{
		SubscriberID: "provider.example",
		URL:          "http://localhost:8081",
		Domain:       "retail",
		Type:         "BPP",
	})

	store.UpsertCatalog(CatalogRecord{
		NodeID:      "provider.example",
		CatalogID:   "electronics",
		CatalogType: "retail",
		Version:     1,
		URL:         "http://localhost:8081/catalog/v1.json",
		Digest:      "sha-256:abc",
	})

	if err := store.Save(); err != nil {
		t.Fatal(err)
	}

	// Create a fresh store to prove persistence.
	loaded := NewStore(path)

	if err := loaded.Load(); err != nil {
		t.Fatal(err)
	}

	index := loaded.Snapshot()

	if len(index.Nodes) != 1 {
		t.Fatalf(
			"expected 1 node, got %d",
			len(index.Nodes),
		)
	}

	if len(index.Catalogs) != 1 {
		t.Fatalf(
			"expected 1 catalog, got %d",
			len(index.Catalogs),
		)
	}

	if index.Catalogs[0].CatalogID !=
		"electronics" {

		t.Fatalf(
			"unexpected catalog ID: %s",
			index.Catalogs[0].CatalogID,
		)
	}
}

func TestStoreUpsertCatalog(t *testing.T) {

	store := NewStore(
		filepath.Join(
			t.TempDir(),
			"index.json",
		),
	)

	store.UpsertCatalog(CatalogRecord{
		NodeID:    "provider.example",
		CatalogID: "electronics",
		Version:   1,
		URL:       "http://example/v1.json",
		Digest:    "sha-256:old",
	})

	store.UpsertCatalog(CatalogRecord{
		NodeID:    "provider.example",
		CatalogID: "electronics",
		Version:   2,
		URL:       "http://example/v2.json",
		Digest:    "sha-256:new",
	})

	index := store.Snapshot()

	if len(index.Catalogs) != 1 {
		t.Fatalf(
			"expected 1 catalog, got %d",
			len(index.Catalogs),
		)
	}

	if index.Catalogs[0].Version != 2 {
		t.Fatalf(
			"expected version 2, got %d",
			index.Catalogs[0].Version,
		)
	}
}
