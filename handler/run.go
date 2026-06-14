package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	mux := http.NewServeMux()
	store := NewURLStore()
	mux.Handle("POST /urls", HandleCreateURL(store))
	mux.Handle("GET /{key}", HandleGetURL(store))
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
