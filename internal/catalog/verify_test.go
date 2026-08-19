package catalog_test

import (
	"testing"

	"github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/catalog"
	catalogcrypto "github.com/remiges-tushar/beckn-decentralized-catalog-poc/internal/crypto"
)

func signedTestCatalog(t *testing.T) (catalog.CatalogFile, *catalogcrypto.KeyPair) {
	t.Helper()

	keyPair, err := catalogcrypto.GenerateKeyPair("provider-key-1")
	if err != nil {
		t.Fatalf("generate key pair: %v", err)
	}

	file := catalog.CatalogFile{
		CatalogID: "provider.example/electronics",
		Version:   1,
		Catalog: catalog.Catalog{
			ID:         "provider.example/electronics",
			Descriptor: catalog.Descriptor{Name: "Example Electronics"},
			Resources: []catalog.Resource{
				{
					ID:         "provider.example/laptop-001",
					Descriptor: catalog.Descriptor{Name: "Example Laptop"},
					Attributes: map[string]interface{}{
						"brand":   "Example",
						"inStock": true,
					},
				},
			},
			Offers: []catalog.Offer{},
		},
	}

	signed, err := catalog.SignCatalogFile(file, keyPair.KeyID, keyPair.PrivateKey)
	if err != nil {
		t.Fatalf("sign catalog file: %v", err)
	}

	return signed, keyPair
}

func TestVerifyCatalogFile_OriginalPasses(t *testing.T) {
	signed, keyPair := signedTestCatalog(t)

	ok, err := catalog.VerifyCatalogFile(signed, keyPair.PublicKey)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if !ok {
		t.Fatal("expected signature verification to pass on untouched catalog")
	}
}

func TestVerifyCatalogFile_TamperedFails(t *testing.T) {
	signed, keyPair := signedTestCatalog(t)

	tampered := signed
	tampered.Catalog.Resources = append(
		[]catalog.Resource(nil),
		signed.Catalog.Resources...,
	)
	tampered.Catalog.Resources[0].Descriptor.Name = "Fake Laptop"

	ok, err := catalog.VerifyCatalogFile(tampered, keyPair.PublicKey)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if ok {
		t.Fatal("expected signature verification to fail after tampering with resource name")
	}
}
