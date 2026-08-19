package discovery

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LocalIndex struct {
	Version   int64           `json:"version"`
	UpdatedAt time.Time       `json:"updated_at"`
	Nodes     []NodeRecord    `json:"nodes"`
	Catalogs  []CatalogRecord `json:"catalogs"`
}

type NodeRecord struct {
	SubscriberID string    `json:"subscriber_id"`
	URL          string    `json:"url"`
	Domain       string    `json:"domain"`
	Type         string    `json:"type"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CatalogRecord struct {
	NodeID      string `json:"node_id"`
	CatalogID   string `json:"catalog_id"`
	CatalogType string `json:"catalog_type"`
	Version     int64  `json:"version"`
	URL         string `json:"url"`
	Digest      string `json:"digest"`
}

type Store struct {
	mu   sync.RWMutex
	path string
	data LocalIndex
}

func NewStore(path string) *Store {
	return &Store{
		path: path,
		data: LocalIndex{
			Version:  1,
			Nodes:    []NodeRecord{},
			Catalogs: []CatalogRecord{},
		},
	}
}

func (s *Store) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf(
			"read local index: %w",
			err,
		)
	}

	var index LocalIndex

	if err := json.Unmarshal(
		data,
		&index,
	); err != nil {
		return fmt.Errorf(
			"parse local index: %w",
			err,
		)
	}

	s.data = index

	return nil
}

func (s *Store) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data.UpdatedAt = time.Now().UTC()

	data, err := json.MarshalIndent(
		s.data,
		"",
		"  ",
	)
	if err != nil {
		return fmt.Errorf(
			"marshal local index: %w",
			err,
		)
	}

	if err := os.MkdirAll(
		filepath.Dir(s.path),
		0755,
	); err != nil {
		return fmt.Errorf(
			"create index directory: %w",
			err,
		)
	}

	tmpPath := s.path + ".tmp"

	if err := os.WriteFile(
		tmpPath,
		data,
		0644,
	); err != nil {
		return fmt.Errorf(
			"write temporary index: %w",
			err,
		)
	}

	if err := os.Rename(
		tmpPath,
		s.path,
	); err != nil {
		return fmt.Errorf(
			"replace local index: %w",
			err,
		)
	}

	return nil
}

func (s *Store) UpsertNode(
	node NodeRecord,
) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.data.Nodes {
		if s.data.Nodes[i].SubscriberID ==
			node.SubscriberID {

			s.data.Nodes[i] = node
			return
		}
	}

	s.data.Nodes = append(
		s.data.Nodes,
		node,
	)
}

func (s *Store) UpsertCatalog(
	record CatalogRecord,
) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.data.Catalogs {
		existing := &s.data.Catalogs[i]

		if existing.NodeID == record.NodeID &&
			existing.CatalogID == record.CatalogID {

			*existing = record
			return
		}
	}

	s.data.Catalogs = append(
		s.data.Catalogs,
		record,
	)
}

func (s *Store) Snapshot() LocalIndex {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.data
}
