package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func handleCreateURL(store *URLStore) http.Handler {
	type request struct {
		URL string `json:"url"`
	}

	type response struct {
		Key string `json:"key"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req request

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		key := fmt.Sprintf("k%d", time.Now().UnixNano())

		store.Save(key, req.URL)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(response{Key: key})
	})
}
