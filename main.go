package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	// Retrieve the backend URL and port from environment variables
	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		log.Fatal("Environment variable BACKEND_URL is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Environment variable PORT is not set")
	}

	// Parse the backend URL
	target, err := url.Parse(backendURL)
	if err != nil {
		log.Fatalf("Invalid BACKEND_URL: %v", err)
	}

	// Create a reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(target)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Update the request's host to match the backend
		r.Host = target.Host
		proxy.ServeHTTP(w, r)
	})

	log.Printf("Starting proxy server on port %s, forwarding to %s", port, backendURL)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
