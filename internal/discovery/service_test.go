package discovery

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCrawlAndIndex(t *testing.T) {

	// Provider must be running on localhost:8081.

	path := filepath.Join(
		t.TempDir(),
		"index.json",
	)

	crawler := NewCrawler()

	store := NewStore(path)

	service := NewService(
		crawler,
		store,
	)

	result, err := service.CrawlAndIndex(
		"http://localhost:8081",
	)

	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Fatal("expected crawl result")
	}

	// Verify file exists.
	if _, err := os.Stat(path); err != nil {
		t.Fatalf(
			"expected local index file: %v",
			err,
		)
	}

	// Read store state.
	index := store.Snapshot()

	if len(index.Nodes) != 1 {
		t.Fatalf(
			"expected 1 node, got %d",
			len(index.Nodes),
		)
	}

	if len(index.Catalogs) == 0 {
		t.Fatal(
			"expected at least one catalog",
		)
	}

	if index.Nodes[0].SubscriberID !=
		"provider.example" {

		t.Fatalf(
			"unexpected subscriber: %s",
			index.Nodes[0].SubscriberID,
		)
	}
}
