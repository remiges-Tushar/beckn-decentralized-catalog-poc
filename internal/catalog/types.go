package catalog

import "time"

type Signature struct {
	KeyID            string `json:"keyId"`
	Canonicalization string `json:"canonicalization"`
	Value            string `json:"value"`
}

type ArtifactRef struct {
	Version int64  `json:"version"`
	URL     string `json:"url"`
	Size    int64  `json:"size"`
	Digest  string `json:"digest"`
}

type ChangeRef struct {
	FromVersion int64  `json:"fromVersion"`
	ToVersion   int64  `json:"toVersion"`
	URL         string `json:"url"`
	Size        int64  `json:"size"`
	Digest      string `json:"digest"`
}

type CatalogFile struct {
	CatalogID  string     `json:"catalogId"`
	Version    int64      `json:"version"`
	NextUpdate time.Time  `json:"next_update"`
	Catalog    Catalog    `json:"catalog"`
	Signature  *Signature `json:"signature,omitempty"`
}

type Catalog struct {
	ID         string     `json:"id"`
	Descriptor Descriptor `json:"descriptor"`
	Resources  []Resource `json:"resources"`
	Offers     []Offer    `json:"offers"`
}

type Descriptor struct {
	Name string `json:"name"`
}

type Resource struct {
	ID         string                 `json:"id"`
	Descriptor Descriptor             `json:"descriptor"`
	Attributes map[string]interface{} `json:"resourceAttributes,omitempty"`
}

type Offer struct {
	ID         string     `json:"id"`
	Descriptor Descriptor `json:"descriptor"`
	Price      float64    `json:"price,omitempty"`
	Currency   string     `json:"currency,omitempty"`
}

type CatalogIndex struct {
	NodeID     string              `json:"nodeId"`
	NextUpdate time.Time           `json:"next_update"`
	Catalogs   []CatalogIndexEntry `json:"catalogs"`
}

type CatalogIndexEntry struct {
	CatalogID    string `json:"catalogId"`
	EntryVersion int64  `json:"entryVersion"`

	CatalogType string `json:"catalogType"`

	NetworkIDs  []string `json:"networkIds,omitempty"`
	SchemaTypes []string `json:"schemaTypes,omitempty"`

	IsActive  *bool      `json:"isActive,omitempty"`
	RetiredAt *time.Time `json:"retiredAt,omitempty"`

	Baseline *ArtifactRef `json:"baseline,omitempty"`
	Changes  []ChangeRef  `json:"changes,omitempty"`
	Latest   *ArtifactRef `json:"latest,omitempty"`

	Signature *Signature `json:"signature,omitempty"`
}
