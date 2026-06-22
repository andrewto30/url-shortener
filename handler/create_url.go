package handler

import (
	"context"
	"log"
	"net/http"
	"net/url"

	"github.com/andrewto30/url-shortener/services/hash"
)

type createURLRequest struct {
	URL string `json:"url"`
}

func (req createURLRequest) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	u, err := url.Parse(req.URL)
	if err != nil {
		problems["url"] = "must be a valid URL"
		return problems
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		problems["url"] = "scheme must be http or https"
	}

	if u.Host == "" {
		problems["url"] = "must include host"
	}

	return problems
}

func handleCreateURL(store *URLStore, gen *hash.Generator) http.Handler {
	type response struct {
		Key string `json:"key"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, problems, err := decodeValid[createURLRequest](r)
		if err != nil {
			if len(problems) > 0 {
				if err := encode(w, r, http.StatusBadRequest, map[string]any{"errors": problems}); err != nil {
					log.Printf("encode error: %v", err)
				}
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		const maxAttempts = 10
		var key string

		for i := 0; i < maxAttempts; i++ {
			k, err := gen.Key()
			if err != nil {
				http.Error(w, "key generation failed", http.StatusInternalServerError)
				return
			}
			if store.SaveIfAbsent(k, req.URL) {
				key = k
				break
			}
		}
		if key == "" {
			http.Error(w, "could not generate unique key", http.StatusInternalServerError)
			return
		}

		if err := encode(w, r, http.StatusCreated, response{Key: key}); err != nil {
			log.Printf("encode error: %v", err)
		}
	})
}
