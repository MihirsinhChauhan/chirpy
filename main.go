package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w,r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits: %d", cfg.fileserverHits.Load())
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	// Reset the counter to 0
	cfg.fileserverHits.Store(0)
	// Set Content-Type header
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// Write 200 OK status
	w.WriteHeader(http.StatusOK)
	// Write a simple confirmation
	w.Write([]byte("Counter reset"))
}

func main() {
	// Initialize apiConfig
	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register the /healthz endpoint
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create a FileServer to serve files from the current directory (.)
	fileServer := http.FileServer(http.Dir("."))

	// Register the FileServer for /app/ path with prefix stripped and metrics middleware
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))

	// Register the /metrics endpoint
	mux.HandleFunc("GET /metrics", apiCfg.handlerMetrics)

	// Register the /reset endpoint
	mux.HandleFunc("POST /reset", apiCfg.handlerReset)

	// Create a new Server with the mux as handler and address set to :8080
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start the server
	server.ListenAndServe()
}