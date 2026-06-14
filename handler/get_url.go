package handler

import "net/http"

func HandleGetURL(store *URLStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")

		url, found := store.Get(key)
		if !found {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})
}
