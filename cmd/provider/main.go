package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/catalog"
	catalogcrypto "github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/crypto"
	"github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/registry"
	"github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/server"
)

const (
	keyID          = "provider-key-1"
	privateKeyPath = "keys/provider-private.key"
	publicKeyPath  = "keys/provider-public.key"
)

func main() {
	var privateKey ed25519.PrivateKey
	var publicKey ed25519.PublicKey

	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		publicKey, privateKey, err =
			ed25519.GenerateKey(rand.Reader)

		if err != nil {
			panic(err)
		}

		if err := os.MkdirAll("keys", 0700); err != nil {
			panic(err)
		}

		if err := catalogcrypto.SavePrivateKey(
			privateKeyPath,
			privateKey,
		); err != nil {
			panic(err)
		}

		if err := catalogcrypto.SavePublicKey(
			publicKeyPath,
			publicKey,
		); err != nil {
			panic(err)
		}
	} else {
		privateKey, err =
			catalogcrypto.LoadPrivateKey(privateKeyPath)

		if err != nil {
			panic(err)
		}

		publicKey, err =
			catalogcrypto.LoadPublicKey(publicKeyPath)

		if err != nil {
			panic(err)
		}
	}

	subscriber := registry.SubscriberRecord{
		RecordName: "beckn-subscriber",

		Details: registry.SubscriberDetails{
			SubscriberID: "provider.example",
			URL:          "http://localhost:8081",
			Type:         "BPP",
			Domain:       "retail",
			Countries:    []string{"IND"},
			SigningPublicKey: catalogcrypto.EncodePublicKey(
				publicKey,
			),
		},

		Meta: registry.SubscriberMeta{
			CatalogIndexURLs: []registry.CatalogIndexURL{
				{
					URL: "http://localhost:8081/catalog/index.json",
				},
			},
		},
	}

	if err := os.MkdirAll(
		"storage/provider/dedi",
		0755,
	); err != nil {
		panic(err)
	}

	subscriberData, err := json.MarshalIndent(
		subscriber,
		"",
		"  ",
	)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(
		"storage/provider/dedi/beckn-subscriber.dedi.json",
		subscriberData,
		0644,
	); err != nil {
		panic(err)
	}

	subscriberDigest := catalogcrypto.Digest(subscriberData)

	fmt.Println("Subscriber digest:")
	fmt.Println(subscriberDigest)
	now := time.Now().UTC()

	manifest := registry.Manifest{
		DeDiVersion: "0.1",
		Type:        "dedi-manifest",
		Domain:      "provider.example",

		Keys: []registry.ManifestKey{
			{
				KID: keyID,
				KTY: "OKP",
				CRV: "Ed25519",
				X: catalogcrypto.EncodePublicKey(
					publicKey,
				),
			},
		},

		UpdatedAt:  now,
		NextUpdate: now.Add(24 * time.Hour),

		Files: []registry.ManifestFile{
			{
				Registry: "beckn-subscriber",
				URL:      "http://localhost:8081/dedi/beckn-subscriber.dedi.json",
				Digest:   subscriberDigest,
			},
		},
	}

	signature, err := registry.SignManifest(
		manifest,
		keyID,
		privateKey,
	)
	if err != nil {
		panic(err)
	}
	manifest.Proof = &registry.ManifestProof{
		VerificationMethod: keyID,
		Canonicalization:   "JCS",
		JWS:                signature.Value,
	}
	if err := os.MkdirAll(
		"storage/provider/.well-known",
		0755,
	); err != nil {
		panic(err)
	}
	manifestData, err := json.MarshalIndent(
		manifest,
		"",
		"  ",
	)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(
		"storage/provider/.well-known/dedi.json",
		manifestData,
		0644,
	); err != nil {
		panic(err)
	}
	cat := catalog.Catalog{
		ID: "provider.example/electronics",

		Descriptor: catalog.Descriptor{
			Name: "Example Electronics",
		},

		Resources: []catalog.Resource{
			{
				ID: "provider.example/laptop-001",

				Descriptor: catalog.Descriptor{
					Name: "Example Laptop",
				},

				Attributes: map[string]interface{}{
					"brand":   "Example",
					"inStock": true,
				},
			},
		},

		Offers: []catalog.Offer{},
	}

	file := catalog.CatalogFile{
		CatalogID:  "provider.example/electronics",
		Version:    1,
		NextUpdate: time.Now().UTC().Add(24 * time.Hour),
		Catalog:    cat,
	}

	signedFile, err := catalog.SignCatalogFile(
		file,
		keyID,
		privateKey,
	)
	if err != nil {
		panic(err)
	}

	data, err := json.MarshalIndent(signedFile, "", "  ")
	if err != nil {
		panic(err)
	}

	if err := os.MkdirAll(
		"storage/provider/catalog/electronics",
		0755,
	); err != nil {
		panic(err)
	}

	if err := os.WriteFile(
		"storage/provider/catalog/electronics/v1.json",
		data,
		0644,
	); err != nil {
		panic(err)
	}

	baseline, err := catalog.BuildArtifactRef(
		1,
		"http://localhost:8081/catalog/electronics/v1.json",
		"storage/provider/catalog/electronics/v1.json",
	)
	if err != nil {
		panic(err)
	}

	isActive := true

	entry := catalog.CatalogIndexEntry{
		CatalogID:    "provider.example/electronics",
		EntryVersion: 1,
		CatalogType:  "REGULAR",

		IsActive: &isActive,

		Baseline: &baseline,
	}

	signedEntry, err := catalog.SignIndexEntry(
		entry,
		keyID,
		privateKey,
	)
	if err != nil {
		panic(err)
	}

	index := catalog.CatalogIndex{
		NodeID:     "provider.example",
		NextUpdate: time.Now().UTC().Add(24 * time.Hour),

		Catalogs: []catalog.CatalogIndexEntry{
			signedEntry,
		},
	}

	data, err = json.MarshalIndent(index, "", "  ")
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(
		"storage/provider/catalog/index.json",
		data,
		0644,
	); err != nil {
		panic(err)
	}

	fmt.Println("Catalog generated")
	fmt.Println("Key ID:", keyID)

	providerServer := server.New(":8081")

	if err := providerServer.Start(); err != nil {
		panic(err)
	}
}
