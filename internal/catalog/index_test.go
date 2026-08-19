package catalog

import (
	"testing"

	catalogcrypto "github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/crypto"
)

func TestIndexEntrySigning(t *testing.T) {

	keyPair, err := catalogcrypto.GenerateKeyPair("test-key")
	if err != nil {
		t.Fatal(err)
	}

	active := true

	entry := CatalogIndexEntry{
		CatalogID:    "provider.example/electronics",
		EntryVersion: 1,
		CatalogType:  "REGULAR",
		IsActive:     &active,
	}

	signed, err := SignIndexEntry(
		entry,
		keyPair.KeyID,
		keyPair.PrivateKey,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !VerifyIndexEntry(
		signed,
		keyPair.PublicKey,
	) {
		t.Fatal("expected signature verification to pass")
	}
}

func TestIndexEntryTampering(t *testing.T) {

	keyPair, err := catalogcrypto.GenerateKeyPair("test-key")
	if err != nil {
		t.Fatal(err)
	}

	entry := CatalogIndexEntry{
		CatalogID:    "provider.example/electronics",
		EntryVersion: 1,
		CatalogType:  "REGULAR",
	}

	signed, err := SignIndexEntry(
		entry,
		keyPair.KeyID,
		keyPair.PrivateKey,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Tamper with the signed content.
	signed.EntryVersion = 999

	if VerifyIndexEntry(
		signed,
		keyPair.PublicKey,
	) {
		t.Fatal("expected tampered entry to fail verification")
	}
}
