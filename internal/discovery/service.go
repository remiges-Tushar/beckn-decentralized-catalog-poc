package discovery

import "fmt"

type Service struct {
	Crawler *Crawler
	Store   *Store
}

func NewService(
	crawler *Crawler,
	store *Store,
) *Service {
	return &Service{
		Crawler: crawler,
		Store:   store,
	}
}

func (s *Service) CrawlAndIndex(
	baseURL string,
) (*CrawlResult, error) {

	if s.Crawler == nil {
		return nil, fmt.Errorf(
			"crawler is nil",
		)
	}

	if s.Store == nil {
		return nil, fmt.Errorf(
			"store is nil",
		)
	}

	// 1. Crawl and verify.
	result, err := s.Crawler.Crawl(
		baseURL,
	)
	if err != nil {
		return nil, err
	}

	// 2. Convert verified data into
	//    DS-owned index records.
	if err := IndexCrawlResult(
		s.Store,
		result,
	); err != nil {
		return nil, fmt.Errorf(
			"index crawl result: %w",
			err,
		)
	}

	// 3. Persist the DS index.
	if err := s.Store.Save(); err != nil {
		return nil, fmt.Errorf(
			"save local index: %w",
			err,
		)
	}

	return result, nil
}

func (s *Service) Load() error {
	if s.Store == nil {
		return fmt.Errorf(
			"store is nil",
		)
	}

	return s.Store.Load()
}
