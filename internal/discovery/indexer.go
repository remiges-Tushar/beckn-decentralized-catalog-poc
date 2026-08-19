package discovery

import (
	"fmt"
	"time"

	"github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/catalog"
)

func BuildNodeRecord(
	result *CrawlResult,
) NodeRecord {

	return NodeRecord{
		SubscriberID: result.Subscriber.Details.SubscriberID,
		URL:          result.Subscriber.Details.URL,
		Domain:       result.Subscriber.Details.Domain,
		Type:         result.Subscriber.Details.Type,
		UpdatedAt:    time.Now().UTC(),
	}
}

func BuildCatalogRecord(
	nodeID string,
	entry catalog.CatalogIndexEntry,
) (CatalogRecord, error) {

	if entry.Baseline == nil {
		return CatalogRecord{}, fmt.Errorf(
			"catalog %s has no baseline",
			entry.CatalogID,
		)
	}

	return CatalogRecord{
		NodeID:      nodeID,
		CatalogID:   entry.CatalogID,
		CatalogType: entry.CatalogType,
		Version:     entry.Baseline.Version,
		URL:         entry.Baseline.URL,
		Digest:      entry.Baseline.Digest,
	}, nil
}

func IndexCrawlResult(
	store *Store,
	result *CrawlResult,
) error {

	if result == nil {
		return fmt.Errorf(
			"crawl result is nil",
		)
	}

	node := BuildNodeRecord(result)

	store.UpsertNode(node)

	nodeID := result.Subscriber.Details.SubscriberID

	for _, entry := range result.Index.Catalogs {

		record, err := BuildCatalogRecord(
			nodeID,
			entry,
		)
		if err != nil {
			return fmt.Errorf(
				"build catalog record: %w",
				err,
			)
		}

		store.UpsertCatalog(record)
	}

	return nil
}
