package discovery

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Fetcher struct {
	Client *http.Client
}

func NewFetcher() *Fetcher {
	return &Fetcher{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (f *Fetcher) Get(url string) ([]byte, http.Header, error) {
	resp, err := f.Client.Get(url)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"GET %s: %w",
			url,
			err,
		)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.Header, fmt.Errorf(
			"GET %s returned %d",
			url,
			resp.StatusCode,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.Header, fmt.Errorf(
			"read %s: %w",
			url,
			err,
		)
	}

	return body, resp.Header, nil
}
