package handler

import "net/http"

func NewServer(store *URLStore) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, store)

	var h http.Handler = mux
	h = loggingMiddleware(h)
	return h
}
