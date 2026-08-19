package discovery

import "time"

type Manifest struct {
	DeDiVersion string         `json:"dedi_version"`
	Type        string         `json:"type"`
	Domain      string         `json:"domain"`
	Keys        []ManifestKey  `json:"keys"`
	UpdatedAt   time.Time      `json:"updated_at"`
	NextUpdate  time.Time      `json:"next_update"`
	Files       []ManifestFile `json:"files"`
	Proof       *ManifestProof `json:"proof,omitempty"`
}

type ManifestKey struct {
	KID string `json:"kid"`
	KTY string `json:"kty"`
	CRV string `json:"crv"`
	X   string `json:"x"`
}

type ManifestFile struct {
	Registry string `json:"registry"`
	URL      string `json:"url"`
	Schema   string `json:"schema,omitempty"`
	Digest   string `json:"digest"`
}

type ManifestProof struct {
	VerificationMethod string `json:"verification_method"`
	Canonicalization   string `json:"canonicalization"`
	JWS                string `json:"jws"`
}

type ManifestVerificationResult struct {
	KeyID string
}
