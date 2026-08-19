package discovery

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetcherGet(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"ok":true}`))
			},
		),
	)

	defer server.Close()

	fetcher := NewFetcher()

	body, _, err := fetcher.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != `{"ok":true}` {
		t.Fatalf(
			"unexpected body: %s",
			body,
		)
	}
}
