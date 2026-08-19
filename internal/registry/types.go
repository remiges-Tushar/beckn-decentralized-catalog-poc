package registry

type SubscriberRecord struct {
	RecordName string            `json:"record_name"`
	Details    SubscriberDetails `json:"details"`
	Meta       SubscriberMeta    `json:"meta,omitempty"`
}

type SubscriberDetails struct {
	SubscriberID      string   `json:"subscriber_id"`
	URL               string   `json:"url"`
	Type              string   `json:"type"`
	Domain            string   `json:"domain"`
	Countries         []string `json:"countries,omitempty"`
	SigningPublicKey  string   `json:"signing_public_key"`
}

type SubscriberMeta struct {
	CatalogIndexURLs []CatalogIndexURL `json:"catalog_index_urls,omitempty"`
}

type CatalogIndexURL struct {
	URL string `json:"url"`
}