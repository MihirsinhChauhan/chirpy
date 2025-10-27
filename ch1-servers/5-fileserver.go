package main

import (
	"net/http"
)

func main() {
	// Create a new ServeMux
	mux := http.NewServeMux()
	// Register the /healthz endpoint
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		// Set Content-Type header
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		// Write 200 OK status code
		w.WriteHeader(http.StatusOK)
		// Write "OK" to the response body
		w.Write([]byte("OK"))
	})
	// Create a FileServer to serve files from the current directory (.)
	fileServer := http.FileServer(http.Dir("."))

	// Register the FileServer as the handler for the root path (/)
	mux.Handle("/app/", http.StripPrefix("/app", fileServer))

	// Create a new Server with the mux as handler and address set to :8080
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start the server
	server.ListenAndServe()
}