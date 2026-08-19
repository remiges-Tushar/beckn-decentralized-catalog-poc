package server

import (
	"log"
	"net/http"
)

type Server struct {
	Addr string
}

func New(addr string) *Server {
	return &Server{
		Addr: addr,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	registerRegistryRoutes(mux)
	registerCatalogRoutes(mux)

	log.Printf(
		"PN HTTP server listening on %s",
		s.Addr,
	)

	return http.ListenAndServe(
		s.Addr,
		logRequests(mux),
	)
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf(
				"%s %s",
				r.Method,
				r.URL.Path,
			)

			next.ServeHTTP(w, r)
		},
	)
}
