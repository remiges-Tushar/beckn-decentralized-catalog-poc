package discovery

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"

	"github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/catalog"
)

func GetCatalogIndexURL(
	subscriber SubscriberRecord,
) (string, error) {

	if len(subscriber.Meta.CatalogIndexURLs) == 0 {
		return "", fmt.Errorf(
			"subscriber contains no catalog index URLs",
		)
	}

	url := subscriber.Meta.CatalogIndexURLs[0].URL

	if url == "" {
		return "", fmt.Errorf(
			"catalog index URL is empty",
		)
	}

	return url, nil
}

func (c *Crawler) FetchCatalogIndex(
	subscriber SubscriberRecord,
) (catalog.CatalogIndex, error) {

	indexURL, err := GetCatalogIndexURL(
		subscriber,
	)
	if err != nil {
		return catalog.CatalogIndex{}, err
	}

	body, _, err := c.Fetcher.Get(indexURL)
	if err != nil {
		return catalog.CatalogIndex{}, fmt.Errorf(
			"fetch catalog index: %w",
			err,
		)
	}

	var index catalog.CatalogIndex

	if err := json.Unmarshal(
		body,
		&index,
	); err != nil {
		return catalog.CatalogIndex{}, fmt.Errorf(
			"parse catalog index: %w",
			err,
		)
	}

	if err := ValidateCatalogIndex(index); err != nil {
		return catalog.CatalogIndex{}, fmt.Errorf(
			"invalid catalog index: %w",
			err,
		)
	}

	return index, nil
}

func ValidateCatalogIndex(
	index catalog.CatalogIndex,
) error {

	if len(index.Catalogs) == 0 {
		return fmt.Errorf(
			"catalog index contains no catalogs",
		)
	}

	return nil
}

// func SubscriberPublicKey(
// 	subscriber SubscriberRecord,
// ) (ed25519.PublicKey, error) {

// 	if subscriber.Details.SigningPublicKey == "" {
// 		return nil, fmt.Errorf(
// 			"subscriber signing public key is empty",
// 		)
// 	}

// 	key, err := decodePublicKey(
// 		subscriber.Details.SigningPublicKey,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return key, nil
// }

// func decodePublicKey(
// 	value string,
// ) (ed25519.PublicKey, error) {

// 	key, err := base64.RawURLEncoding.DecodeString(
// 		value,
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf(
// 			"decode public key: %w",
// 			err,
// 		)
// 	}

// 	if len(key) != ed25519.PublicKeySize {
// 		return nil, fmt.Errorf(
// 			"invalid Ed25519 public key size: %d",
// 			len(key),
// 		)
// 	}

// 	return ed25519.PublicKey(key), nil
// }

func VerifyCatalogIndexEntries(
	index catalog.CatalogIndex,
	publicKey ed25519.PublicKey,
) error {

	for i, entry := range index.Catalogs {

		if !catalog.VerifyIndexEntry(
			entry,
			publicKey,
		) {
			return fmt.Errorf(
				"catalog index entry %d signature verification failed",
				i,
			)
		}
	}

	return nil
}

func (c *Crawler) FetchAndVerifyCatalogIndex(
	subscriber SubscriberRecord,
) (catalog.CatalogIndex, error) {

	index, err := c.FetchCatalogIndex(
		subscriber,
	)
	if err != nil {
		return catalog.CatalogIndex{}, err
	}

	publicKey, err := SubscriberPublicKey(
		subscriber,
	)
	if err != nil {
		return catalog.CatalogIndex{}, err
	}

	if err := verifyCatalogIndexEntries(
		index,
		publicKey,
	); err != nil {
		return catalog.CatalogIndex{}, fmt.Errorf(
			"catalog index verification failed: %w",
			err,
		)
	}

	return index, nil
}

func verifyCatalogIndexEntries(
	index catalog.CatalogIndex,
	publicKey ed25519.PublicKey,
) error {
	for _, entry := range index.Catalogs {

		ok := catalog.VerifyIndexEntry(
			entry,
			publicKey,
		)
		if !ok {
			return fmt.Errorf(
				"invalid catalog index entry",
			)
		}
	}

	return nil
}
