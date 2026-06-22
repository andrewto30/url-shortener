package handler

import (
	"net/http"

	"github.com/andrewto30/url-shortener/services/hash"
)

func NewServer(store *URLStore, gen *hash.Generator) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, store, gen)

	var h http.Handler = mux
	h = loggingMiddleware(h)
	return h
}
