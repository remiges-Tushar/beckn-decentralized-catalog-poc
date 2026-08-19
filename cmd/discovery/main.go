package main

import (
	"log"

	"github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/discovery"
)

func main() {

	crawler := discovery.NewCrawler()

	store := discovery.NewStore(
		"storage/discovery/index.json",
	)

	service := discovery.NewService(
		crawler,
		store,
	)

	if err := service.Load(); err != nil {
		log.Fatal(err)
	}

	result, err := service.CrawlAndIndex(
		"http://localhost:8081",
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(
		"crawl completed: node=%s catalogs=%d",
		result.Subscriber.Details.SubscriberID,
		len(result.Catalogs),
	)
}
