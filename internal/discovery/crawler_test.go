package discovery

import (
	"testing"
)

func TestCrawlProvider(t *testing.T) {

	crawler := NewCrawler()

	result, err := crawler.Crawl(
		"http://localhost:8081",
	)

	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Fatal("expected crawl result")
	}

	if result.Manifest.Domain != "provider.example" {
		t.Fatalf(
			"unexpected manifest domain: %s",
			result.Manifest.Domain,
		)
	}

	if result.Subscriber.RecordName != "beckn-subscriber" {
		t.Fatalf(
			"unexpected subscriber record: %s",
			result.Subscriber.RecordName,
		)
	}

	if len(result.Index.Catalogs) == 0 {
		t.Fatal(
			"expected at least one catalog index entry",
		)
	}

	if len(result.Catalogs) == 0 {
		t.Fatal(
			"expected at least one catalog file",
		)
	}
}