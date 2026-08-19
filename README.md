# Beckn Decentralized Catalog POC

A Go-based Proof of Concept for decentralized catalog publishing and discovery based on the Beckn NFH-014 / DeDi catalog discovery model.

The POC demonstrates how a Discovery Service (DS) can discover, verify, fetch, and locally index catalog information published by a Provider without the Registry itself becoming a centralized catalog database.

---

# 1. Objective

The goal of this POC is to demonstrate the decentralized catalog discovery flow:

```text
Provider
   |
   | publishes
   v
DeDi Manifest
   |
   v
Subscriber Record
   |
   v
Catalog Index
   |
   v
CatalogFile
```

The Discovery Service follows these references, verifies the downloaded artifacts, and builds its own local index.

The Registry is therefore used for discovery/trust information rather than storing the complete catalog of every provider.

---

# 2. High-Level Architecture

```text
                         Provider
                            |
              +-------------+-------------+
              |                           |
              v                           v
       /.well-known/dedi.json       HTTP Catalog APIs
              |
              v
        Discovery Service
              |
              v
       Verify Manifest
              |
              v
      Beckn Subscriber
              |
       Verify SHA-256
              |
              v
      catalog_index_urls
              |
              v
       Catalog Index
              |
       Verify entries
              |
              v
          baseline
              |
              v
         CatalogFile
              |
       Verify SHA-256
              |
              v
       Trusted Catalog
              |
              v
       DS Local Index
```

---

# 3. Components

## 3.1 Provider

The Provider publishes:

- DeDi manifest
- Subscriber record
- Catalog Index
- CatalogFile

The Provider also exposes these through an HTTP server.

## 3.2 Discovery Service

The Discovery Service:

1. Fetches the DeDi manifest.
2. Verifies the manifest signature.
3. Finds the Subscriber record.
4. Downloads the Subscriber record.
5. Verifies its SHA-256 digest.
6. Reads `catalog_index_urls`.
7. Fetches the Catalog Index.
8. Verifies Catalog Index entries.
9. Reads the CatalogFile baseline.
10. Fetches the CatalogFile.
11. Verifies the CatalogFile digest.
12. Stores verified information in its local index.

## 3.3 Local DS Index

The DS maintains its own local index:

```text
storage/discovery/index.json
```

This is different from the Provider's published catalog artifacts.

The Provider remains the source of truth.

The DS creates a locally maintained, verified representation suitable for discovery and future search.

---

# 4. Repository Structure

```text
.
├── cmd/
│   ├── provider/
│   │   └── main.go
│   │
│   └── discovery/
│       └── main.go
│
├── internal/
│   ├── catalog/
│   │   ├── types.go
│   │   ├── index.go
│   │   ├── artifact.go
│   │   └── ...
│   │
│   ├── crypto/
│   │   └── ...
│   │
│   └── discovery/
│       ├── crawler.go
│       ├── manifest_verify.go
│       ├── subscriber.go
│       ├── catalog_index.go
│       ├── catalog.go
│       ├── store.go
│       ├── indexer.go
│       ├── service.go
│       │
│       ├── manifest_integration_test.go
│       ├── subscriber_test.go
│       ├── catalog_index_test.go
│       ├── catalog_test.go
│       ├── indexer_test.go
│       ├── store_test.go
│       ├── service_test.go
│       └── crawler_test.go
│
├── storage/
│   ├── provider/
│   │   ├── .well-known/
│   │   │   └── dedi.json
│   │   │
│   │   ├── dedi/
│   │   │   └── beckn-subscriber.dedi.json
│   │   │
│   │   └── catalog/
│   │       └── ...
│   │
│   └── discovery/
│       └── index.json
│
├── go.mod
└── README.md
```

---

# 5. Catalog Data Model

## 5.1 CatalogFile

The CatalogFile contains:

```go
type CatalogFile struct {
    CatalogID  string     `json:"catalogId"`
    Version    int64      `json:"version"`
    NextUpdate time.Time  `json:"next_update"`
    Catalog    Catalog    `json:"catalog"`
    Signature  *Signature `json:"signature,omitempty"`
}
```

A Catalog contains:

```go
type Catalog struct {
    ID         string     `json:"id"`
    Descriptor Descriptor `json:"descriptor"`
    Resources  []Resource `json:"resources"`
    Offers     []Offer    `json:"offers"`
}
```

---

# 6. Catalog Index

The Catalog Index is represented by:

```go
type CatalogIndex struct {
    NodeID     string              `json:"nodeId"`
    NextUpdate time.Time           `json:"next_update"`
    Catalogs   []CatalogIndexEntry `json:"catalogs"`
}
```

Each entry contains:

```go
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
```

---

# 7. Artifact References

Catalog Index entries point to artifacts using:

```go
type ArtifactRef struct {
    Version int64  `json:"version"`
    URL     string `json:"url"`
    Size    int64  `json:"size"`
    Digest  string `json:"digest"`
}
```

For example:

```text
Catalog Index Entry
       |
       +-- baseline
             |
             +-- version
             +-- URL
             +-- size
             +-- digest
```

The DS uses the baseline URL to retrieve the CatalogFile and the digest to verify its integrity.

---

# 8. DeDi Manifest

The Provider generates a manifest similar to:

```json
{
  "dedi_version": "0.1",
  "type": "dedi-manifest",
  "domain": "provider.example",
  "keys": [
    {
      "kid": "provider-key-1",
      "kty": "OKP",
      "crv": "Ed25519",
      "x": "..."
    }
  ],
  "updated_at": "...",
  "next_update": "...",
  "files": [
    {
      "registry": "beckn-subscriber",
      "url": "...",
      "digest": "sha-256:..."
    }
  ]
}
```

The manifest identifies the Subscriber registry artifact that the DS should fetch.

---

# 9. Manifest Signing

The manifest is signed using Ed25519.

The signature input does not contain the proof itself.

Conceptually:

```text
Manifest
   |
   | remove proof
   v
Canonical JSON / JCS
   |
   v
Ed25519
   |
   v
Signature
   |
   v
Manifest.Proof
```

This avoids a circular dependency where the signature would attempt to sign itself.

---

# 10. Manifest Verification

The DS performs:

```text
GET /.well-known/dedi.json
          |
          v
       Parse
          |
          v
      Extract key
          |
          v
    Canonicalize
          |
          v
   Ed25519 verify
          |
      +---+---+
      |       |
    valid   invalid
      |       |
      v       v
   continue  reject
```

An invalid manifest is never used for further discovery.

---

# 11. Subscriber Discovery

The manifest contains:

```json
{
  "registry": "beckn-subscriber",
  "url": "http://localhost:8081/dedi/beckn-subscriber.dedi.json",
  "digest": "sha-256:..."
}
```

The DS fetches the referenced URL.

It then calculates:

```text
SHA-256(downloaded bytes)
```

and compares the result with:

```text
manifest.files[].digest
```

---

# 12. Subscriber Record

The Subscriber record contains information such as:

```json
{
  "record_name": "beckn-subscriber",
  "details": {
    "subscriber_id": "provider.example",
    "url": "http://localhost:8081",
    "type": "BPP",
    "domain": "retail",
    "countries": [
      "IND"
    ],
    "signing_public_key": "..."
  },
  "meta": {
    "catalog_index_urls": [
      {
        "url": "http://localhost:8081/catalog/index.json"
      }
    ]
  }
}
```

The important field for the next discovery step is:

```text
meta.catalog_index_urls
```

---

# 13. Catalog Index Discovery

The DS follows:

```text
Subscriber
    |
    v
meta.catalog_index_urls
    |
    v
GET catalog/index.json
```

The downloaded JSON is parsed into:

```go
catalog.CatalogIndex
```

The DS validates that the index contains catalog entries.

---

# 14. Catalog Index Entry Verification

Each Catalog Index entry can contain a signature.

The existing implementation signs an entry by:

```text
CatalogIndexEntry
       |
       | Signature = nil
       v
JCS canonicalization
       |
       v
Ed25519 signing
       |
       v
Signature attached
```

Verification does the reverse:

```text
CatalogIndexEntry
       |
       | Signature removed
       v
JCS canonicalization
       |
       v
Ed25519 verification
```

The existing verification function is:

```go
catalog.VerifyIndexEntry(
    entry,
    publicKey,
)
```

The DS verifies each entry.

If an entry fails verification, the current POC fails closed and rejects the Catalog Index.

---

# 15. CatalogFile Discovery

Once a Catalog Index entry is trusted, the DS follows:

```text
CatalogIndexEntry
       |
       v
Baseline
       |
       +-- URL
       +-- version
       +-- digest
```

The DS fetches the baseline URL.

---

# 16. CatalogFile Integrity Verification

The DS calculates the SHA-256 digest of the downloaded CatalogFile.

It compares:

```text
Catalog Index baseline.digest
```

against:

```text
SHA-256(downloaded CatalogFile)
```

Example:

```text
Expected:

sha-256:ABC


Downloaded:

sha-256:ABC

        |
        v

       PASS
```

If the downloaded file has been modified:

```text
Expected:

sha-256:ABC


Downloaded:

sha-256:XYZ

        |
        v

       FAIL
```

The DS rejects the CatalogFile.

---

# 17. End-to-End Crawler

The DS exposes the complete discovery flow through:

```go
result, err := crawler.Crawl(
    "http://localhost:8081",
)
```

The sequence is:

```text
Crawl()
   |
   +--> Fetch + verify Manifest
   |
   +--> Fetch + verify Subscriber
   |
   +--> Fetch + verify Catalog Index
   |
   +--> Fetch + verify CatalogFiles
   |
   v
CrawlResult
```

---

# 18. CrawlResult

The crawler returns:

```go
type CrawlResult struct {
    Manifest   Manifest
    Subscriber SubscriberRecord
    Index      catalog.CatalogIndex
    Catalogs   []catalog.CatalogFile
}
```

This keeps the verified result in memory before it is indexed.

---

# 19. DS Local Index

The DS maintains:

```text
storage/discovery/index.json
```

The local index contains:

```go
type LocalIndex struct {
    Version   int64
    UpdatedAt time.Time
    Nodes     []NodeRecord
    Catalogs  []CatalogRecord
}
```

Nodes:

```go
type NodeRecord struct {
    SubscriberID string
    URL          string
    Domain       string
    Type         string
    UpdatedAt    time.Time
}
```

Catalogs:

```go
type CatalogRecord struct {
    NodeID      string
    CatalogID   string
    CatalogType string
    Version     int64
    URL         string
    Digest      string
}
```

---

# 20. Why the Local Index Exists

The DS does not replace the Provider's catalog.

Instead:

```text
Provider
   |
   | source of truth
   v
CatalogFile
   |
   | discovered + verified
   v
Discovery Service
   |
   | builds
   v
Local DS Index
```

The DS's local index is optimized for discovery and querying.

The Provider remains responsible for publishing the authoritative catalog.

---

# 21. Upsert Behavior

The DS does not blindly append duplicate records.

For nodes:

```text
subscriber_id
```

is used as the current identity.

For catalogs:

```text
node_id + catalog_id
```

is used as the current identity.

Example:

```text
First crawl:

provider.example
electronics v1


Second crawl:

provider.example
electronics v2
```

The local index becomes:

```text
provider.example
electronics v2
```

rather than:

```text
provider.example
electronics v1
provider.example
electronics v2
```

---

# 22. Persistent Store

The Store provides:

```go
Load()
Save()
UpsertNode()
UpsertCatalog()
Snapshot()
```

The Store writes through a temporary file:

```text
index.json.tmp
      |
      | complete write
      v
rename
      |
      v
index.json
```

This reduces the chance of leaving a partially written index after an interrupted write.

---

# 23. DS Service

The high-level service connects the pieces:

```go
service.CrawlAndIndex(
    "http://localhost:8081",
)
```

Internally:

```text
Crawl
  |
  v
CrawlResult
  |
  v
IndexCrawlResult
  |
  v
Store
  |
  v
Save
```

The responsibilities remain separated:

```text
Crawler
  = network fetching + verification

Indexer
  = convert verified data into DS records

Store
  = persistent storage

Service
  = orchestrates the workflow
```

---

# 24. Running the POC

## Start Provider

From the project root:

```bash
go run ./cmd/provider
```

The Provider runs on:

```text
http://localhost:8081
```

## Start Discovery Service

In another terminal:

```bash
go run ./cmd/discovery
```

The DS will:

```text
1. Fetch manifest
2. Verify manifest
3. Fetch Subscriber
4. Verify Subscriber digest
5. Fetch Catalog Index
6. Verify Catalog Index entries
7. Fetch CatalogFile
8. Verify CatalogFile digest
9. Build local index
10. Save local index
```

---

# 25. Inspect the DS Index

After a successful crawl:

```bash
cat storage/discovery/index.json
```

This is the DS-owned index.

---

# 26. Run Tests

Run all tests:

```bash
go test ./...
```

Run discovery tests only:

```bash
go test ./internal/discovery -v
```

Run the end-to-end Provider crawl test:

```bash
go test ./internal/discovery \
  -run TestCrawlProvider \
  -v
```

---

# 27. Security / Integrity Tests

The POC contains negative tests for tampering.

Examples include:

```text
Manifest modified
       |
       v
signature verification fails
```

and:

```text
Subscriber modified
       |
       v
digest verification fails
```

and:

```text
Catalog Index entry modified
       |
       v
signature verification fails
```

and:

```text
CatalogFile modified
       |
       v
digest verification fails
```

These tests demonstrate that the DS does not blindly trust HTTP responses.

---

# 28. Current End-to-End Architecture

```text
                       PROVIDER
                          |
             +------------+------------+
             |            |            |
             v            v            v
          dedi.json   Subscriber   Catalog Index
             |            |            |
             |            |            |
             +------------+------------+
                          |
                          v
                 DISCOVERY SERVICE
                          |
                   Verify Manifest
                          |
                          v
                    Subscriber
                          |
                    Verify Digest
                          |
                          v
                   Catalog Index
                          |
                Verify Entry Signatures
                          |
                          v
                    Catalog Entries
                          |
                          v
                     CatalogFile
                          |
                    Verify Digest
                          |
                          v
                  Verified Catalog
                          |
                          v
                     INDEXER
                          |
                          v
                   LOCAL DS INDEX
                          |
                          v
              storage/discovery/index.json
```

---

# 29. What We Have Demonstrated

The current POC demonstrates an important architectural separation:

```text
Registry / discovery information
            |
            v
       Discovery Service
            |
            v
      decentralized sources
            |
            v
       verified artifacts
            |
            v
        local DS index
```

The DS is not simply downloading one giant centralized catalog.

Instead, it discovers catalog information from decentralized publishers and builds an index locally.

---

# 30. What We Have NOT Implemented Yet

The current POC is intentionally limited.

We have not yet implemented:

### Multiple Providers

Currently the end-to-end flow has primarily been tested against:

```text
provider.example
localhost:8081
```

Next we will introduce multiple Providers.

### Multi-provider scaling

We have not yet simulated:

```text
10 providers
100 providers
1,000 providers
10,000 providers
```

### Concurrent crawling

The current crawler processes CatalogFiles sequentially.

We will first establish a baseline and then introduce concurrency.

### Incremental crawling

We have not yet implemented a complete:

```text
version comparison
        +
ETag
        +
If-None-Match
        +
only changed artifacts
```

strategy.

### Search API

The DS local index currently exists as persistent storage.

We haven't yet exposed:

```text
GET /search
```

or an equivalent consumer-facing API.

---

# 31. Next Phase — Scaling Experiment

The next phase is:

```text
                Discovery Service
                       |
        +--------------+--------------+
        |              |              |
        v              v              v
   Provider A     Provider B     Provider C
        |              |              |
        v              v              v
     Catalogs       Catalogs       Catalogs
        |              |              |
        +--------------+--------------+
                       |
                       v
                 Local DS Index
```

Then we will measure:

```text
Number of Providers
        |
        +-- 10
        +-- 100
        +-- 1,000
        +-- 10,000
```

For each scenario we will measure:

- Total HTTP requests
- Manifest requests
- Subscriber requests
- Catalog Index requests
- CatalogFile requests
- Bytes downloaded
- Crawl duration
- Number of catalog entries
- Number of changed catalogs
- Number of unchanged catalogs
- Memory usage
- CPU usage

This will allow us to demonstrate the scaling characteristics of the decentralized catalog approach with actual measurements rather than only theoretical discussion.

---

# 32. Development Philosophy

The POC intentionally follows this sequence:

```text
1. Generate artifacts
        ↓
2. Publish artifacts
        ↓
3. Verify individual artifacts
        ↓
4. Build end-to-end crawler
        ↓
5. Build local DS index
        ↓
6. Add multiple providers
        ↓
7. Measure baseline
        ↓
8. Add incremental crawling
        ↓
9. Add concurrency
        ↓
10. Measure scaling
```

This lets us distinguish:

```text
Correctness
```

from:

```text
Performance / scalability
```

instead of optimizing an implementation before proving that the discovery chain is correct.

---

# 33. Current Milestone

## Milestone: Single Provider End-to-End Discovery

Status:

```text
                    STATUS
                       |
                       v
             +-------------------+
             | Manifest          | ✓
             | Subscriber        | ✓
             | Catalog Index     | ✓
             | CatalogFile       | ✓
             | Verification      | ✓
             | HTTP Server       | ✓
             | DS Crawler        | ✓
             | Local DS Index    | ✓
             | Persistence       | ✓
             +-------------------+
```

The next milestone is:

```text
Multi-Provider Discovery
```

which will let us begin the actual **NFH-014 scaling POC**.

---

# 34. Reference

This POC is based on the Beckn Protocol Specifications v2 catalog publishing and discovery material, particularly the NFH-014 decentralized catalog discovery concepts.

Reference:

https://github.com/beckn/protocol-specifications-v2/blob/main/docs/Catalog_Publishing_and_Discovery.md

---

# 35. Current POC Milestone Summary

We have successfully demonstrated the following complete path:

```text
Provider
   |
   v
DeDi Manifest
   |
   | Ed25519 verification
   v
Subscriber Record
   |
   | SHA-256 verification
   v
Catalog Index
   |
   | Entry signature verification
   v
CatalogFile
   |
   | SHA-256 verification
   v
Verified Catalog
   |
   v
DS Indexer
   |
   v
Persistent DS Local Index
   |
   v
storage/discovery/index.json
```

All current tests pass.

The next implementation phase is **Multi-Provider Discovery and the NFH-014 scaling experiment**.