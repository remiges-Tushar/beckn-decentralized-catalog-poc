package discovery

import (
	"log"

	"github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/catalog"
)

type Crawler struct {
	Fetcher *Fetcher
}

type CrawlResult struct {
	Manifest   Manifest
	Subscriber SubscriberRecord
	Index      catalog.CatalogIndex
	Catalogs   []catalog.CatalogFile
}

func NewCrawler() *Crawler {
	return &Crawler{
		Fetcher: NewFetcher(),
	}
}

func (c *Crawler) Crawl(
	baseURL string,
) (*CrawlResult, error) {

	log.Printf(
		"crawl started: %s",
		baseURL,
	)

	log.Printf("fetching manifest")

	manifest, err := c.FetchAndVerifyManifest(
		baseURL,
	)
	if err != nil {
		return nil, err
	}

	log.Printf(
		"manifest verified: domain=%s",
		manifest.Domain,
	)

	log.Printf("fetching subscriber")

	subscriber, err := c.FetchSubscriber(
		manifest,
	)
	if err != nil {
		return nil, err
	}

	log.Printf(
		"subscriber verified: id=%s",
		subscriber.Details.SubscriberID,
	)

	log.Printf("fetching catalog index")

	index, err := c.FetchAndVerifyCatalogIndex(
		subscriber,
	)
	if err != nil {
		return nil, err
	}

	log.Printf(
		"catalog index verified: entries=%d",
		len(index.Catalogs),
	)

	catalogs := make(
		[]catalog.CatalogFile,
		0,
		len(index.Catalogs),
	)

	for i, entry := range index.Catalogs {

		log.Printf(
			"fetching catalog: entry=%d",
			i,
		)

		file, err := c.FetchAndVerifyCatalogFile(
			entry,
		)
		if err != nil {
			return nil, err
		}

		catalogs = append(
			catalogs,
			file,
		)

		log.Printf(
			"catalog verified: entry=%d",
			i,
		)
	}

	log.Printf(
		"crawl completed: catalogs=%d",
		len(catalogs),
	)

	return &CrawlResult{
		Manifest:   manifest,
		Subscriber: subscriber,
		Index:      index,
		Catalogs:   catalogs,
	}, nil
}
