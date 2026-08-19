package discovery

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

type SubscriberRecord struct {
	RecordName string            `json:"record_name"`
	Details    SubscriberDetails `json:"details"`
	Meta       SubscriberMeta    `json:"meta,omitempty"`
}

type SubscriberDetails struct {
	SubscriberID     string   `json:"subscriber_id"`
	URL              string   `json:"url"`
	Type             string   `json:"type"`
	Domain           string   `json:"domain"`
	Countries        []string `json:"countries,omitempty"`
	SigningPublicKey string   `json:"signing_public_key"`
}

type SubscriberMeta struct {
	CatalogIndexURLs []CatalogIndexURL `json:"catalog_index_urls,omitempty"`
}

type CatalogIndexURL struct {
	URL string `json:"url"`
}

func CalculateDigest(data []byte) string {
	sum := sha256.Sum256(data)

	return "sha-256:" +
		hex.EncodeToString(sum[:])
}

func VerifyDigest(
	data []byte,
	expected string,
) error {

	actual := CalculateDigest(data)

	if !strings.EqualFold(
		actual,
		expected,
	) {
		return fmt.Errorf(
			"digest mismatch: expected %s, got %s",
			expected,
			actual,
		)
	}

	return nil
}

func FindRegistryFile(
	manifest Manifest,
	registryName string,
) (ManifestFile, error) {

	for _, file := range manifest.Files {
		if file.Registry == registryName {
			return file, nil
		}
	}

	return ManifestFile{}, fmt.Errorf(
		"registry file not found: %s",
		registryName,
	)
}

func ValidateSubscriber(
	subscriber SubscriberRecord,
) error {

	if subscriber.RecordName == "" {
		return fmt.Errorf(
			"missing record_name",
		)
	}

	if subscriber.Details.SubscriberID == "" {
		return fmt.Errorf(
			"missing subscriber_id",
		)
	}

	if subscriber.Details.URL == "" {
		return fmt.Errorf(
			"missing subscriber url",
		)
	}

	if subscriber.Details.SigningPublicKey == "" {
		return fmt.Errorf(
			"missing signing public key",
		)
	}

	if len(
		subscriber.Meta.CatalogIndexURLs,
	) == 0 {
		return fmt.Errorf(
			"subscriber has no catalog index URLs",
		)
	}

	return nil
}

func SubscriberPublicKey(
	subscriber SubscriberRecord,
) (ed25519.PublicKey, error) {

	if subscriber.Details.SigningPublicKey == "" {
		return nil, fmt.Errorf(
			"subscriber signing public key is empty",
		)
	}

	key, err := base64.RawURLEncoding.DecodeString(
		subscriber.Details.SigningPublicKey,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"decode subscriber public key: %w",
			err,
		)
	}

	if len(key) != ed25519.PublicKeySize {
		return nil, fmt.Errorf(
			"invalid Ed25519 public key size: %d",
			len(key),
		)
	}

	return ed25519.PublicKey(key), nil
}

func (c *Crawler) FetchSubscriber(
	manifest Manifest,
) (SubscriberRecord, error) {

	file, err := FindRegistryFile(
		manifest,
		"beckn-subscriber",
	)
	if err != nil {
		return SubscriberRecord{}, err
	}

	body, _, err := c.Fetcher.Get(file.URL)
	if err != nil {
		return SubscriberRecord{}, err
	}

	if err := VerifyDigest(
		body,
		file.Digest,
	); err != nil {
		return SubscriberRecord{}, fmt.Errorf(
			"subscriber integrity verification failed: %w",
			err,
		)
	}

	var subscriber SubscriberRecord

	if err := json.Unmarshal(
		body,
		&subscriber,
	); err != nil {
		return SubscriberRecord{}, fmt.Errorf(
			"parse subscriber: %w",
			err,
		)
	}

	return subscriber, nil
}
