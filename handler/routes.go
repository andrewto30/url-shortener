package handler

import (
	"net/http"

	"github.com/andrewto30/url-shortener/services/hash"
)

func addRoutes(mux *http.ServeMux, store *URLStore, gen *hash.Generator) {
	mux.Handle("POST /urls", handleCreateURL(store, gen))
	mux.Handle("GET /{key}", handleGetURL(store))
}
