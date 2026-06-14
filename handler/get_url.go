package handler

import "net/http"

func handleGetURL(store *URLStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")

		url, found := store.Get(key)
		if !found {
			http.NotFound(w, r)
			return
		}

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})
}
