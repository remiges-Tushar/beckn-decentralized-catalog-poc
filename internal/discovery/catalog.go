package discovery

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"

	"github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/catalog"
)

func GetBaselineURL(
	entry catalog.CatalogIndexEntry,
) (string, error) {

	if entry.Baseline == nil {
		return "", fmt.Errorf(
			"catalog index entry has no baseline",
		)
	}

	if entry.Baseline.URL == "" {
		return "", fmt.Errorf(
			"catalog baseline URL is empty",
		)
	}

	return entry.Baseline.URL, nil
}

func (c *Crawler) FetchCatalogFile(
	entry catalog.CatalogIndexEntry,
) (catalog.CatalogFile, []byte, error) {

	url, err := GetBaselineURL(entry)
	if err != nil {
		return catalog.CatalogFile{}, nil, err
	}

	body, _, err := c.Fetcher.Get(url)
	if err != nil {
		return catalog.CatalogFile{}, nil, fmt.Errorf(
			"fetch catalog file: %w",
			err,
		)
	}

	var file catalog.CatalogFile

	if err := json.Unmarshal(
		body,
		&file,
	); err != nil {
		return catalog.CatalogFile{}, nil, fmt.Errorf(
			"parse catalog file: %w",
			err,
		)
	}

	return file, body, nil
}

func VerifyCatalogFileDigest(
	entry catalog.CatalogIndexEntry,
	body []byte,
) error {

	expected := entry.Baseline.Digest

	if expected == "" {
		return fmt.Errorf(
			"catalog baseline has no digest",
		)
	}

	if err := VerifyDigest(
		body,
		expected,
	); err != nil {
		return fmt.Errorf(
			"catalog file digest verification failed: %w",
			err,
		)
	}

	return nil
}

func VerifyCatalogFileSignature(
	file catalog.CatalogFile,
	publicKey ed25519.PublicKey,
) error {

	_, err := catalog.VerifyCatalogFile(
		file,
		publicKey,
	)
	if err != nil {
		return fmt.Errorf(
			"catalog file signature verification failed: %w",
			err,
		)
	}

	return nil
}

func (c *Crawler) FetchAndVerifyCatalogFile(
	entry catalog.CatalogIndexEntry,
) (catalog.CatalogFile, error) {

	file, body, err := c.FetchCatalogFile(
		entry,
	)
	if err != nil {
		return catalog.CatalogFile{}, err
	}

	if err := VerifyCatalogFileDigest(
		entry,
		body,
	); err != nil {
		return catalog.CatalogFile{}, err
	}

	return file, nil
}
