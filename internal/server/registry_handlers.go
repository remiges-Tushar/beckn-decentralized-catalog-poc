package server

import (
	"net/http"
)

func registerRegistryRoutes(
	mux *http.ServeMux,
) {
	mux.HandleFunc(
		"/.well-known/dedi.json",
		serveStaticFile(
			"storage/provider/.well-known/dedi.json",
		),
	)

	mux.HandleFunc(
		"/dedi/beckn-subscriber.dedi.json",
		serveStaticFile(
			"storage/provider/dedi/beckn-subscriber.dedi.json",
		),
	)
}
