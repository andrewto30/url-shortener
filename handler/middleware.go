package handler

import (
	"log"
	"net/http"
	"time"
)

// Struct embedding and method overriding
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call inner handler. It writes the response
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		// After it returns, log what happened
		log.Printf("%s %s - %d - %s", r.Method, r.URL.Path, rec.status, time.Since(start))
	})
}
