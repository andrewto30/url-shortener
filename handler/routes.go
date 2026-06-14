package handler

import "net/http"

func addRoutes(mux *http.ServeMux, store *URLStore) {
	mux.Handle("POST /urls", handleCreateURL(store))
	mux.Handle("GET /{key}", handleGetURL(store))
}
