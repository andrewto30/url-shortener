package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type URLStore struct {
	mu   sync.Mutex
	urls map[string]string
}

func NewURLStore() *URLStore {
	return &URLStore{
		urls: make(map[string]string),
	}
}

func (u *URLStore) Save(key, url string) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.urls[key] = url
}

func (u *URLStore) Get(key string) (string, bool) {
	u.mu.Lock()
	defer u.mu.Unlock()

	v, ok := u.urls[key]

	return v, ok
}

func handleGetURL(store *URLStore) http.Handler {
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

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	mux := http.NewServeMux()
	store := NewURLStore()
	mux.Handle("POST /urls", handleCreateURL(store))
	mux.Handle("GET /{key}", handleGetURL(store))
	mux.HandleFunc("GET /slow", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("slow: doing important work...")
		time.Sleep(5 * time.Second)
		fmt.Println("slow: DONE")
		fmt.Fprintln(w, "ok")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		fmt.Println("listening on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "server error: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		fmt.Println("shutting down...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "shutdown error: %s\n", err)
		}
	}()

	wg.Wait()
	return nil
}
