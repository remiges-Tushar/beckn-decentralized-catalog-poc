package server

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
)

func serveStaticFile(path string) http.HandlerFunc {
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		data, err := os.ReadFile(path)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		info, err := os.Stat(path)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		sum := sha256.Sum256(data)

		etag := `"` +
			hex.EncodeToString(sum[:]) +
			`"`

		w.Header().Set(
			"ETag",
			etag,
		)

		w.Header().Set(
			"Cache-Control",
			"public, max-age=60",
		)

		if r.Header.Get("If-None-Match") == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		http.ServeContent(
			w,
			r,
			path,
			info.ModTime(),
			bytes.NewReader(data),
		)
	}
}
