package server

import "net/http"

func registerCatalogRoutes(
	mux *http.ServeMux,
) {
	mux.HandleFunc(
		"/catalog/index.json",
		serveStaticFile(
			"storage/provider/catalog/index.json",
		),
	)

	mux.HandleFunc(
		"/catalog/electronics/v1.json",
		serveStaticFile(
			"storage/provider/catalog/electronics/v1.json",
		),
	)
}
