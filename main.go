package main

import (
	"io"
	"log"
	"net/http"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Construct the backend URL
		proxyURL := backendURL + r.URL.Path

		// Create a new request to the backend
		req, err := http.NewRequest(r.Method, proxyURL, r.Body)
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}

		// Copy the headers from the original request
		for name, values := range r.Header {
			for _, value := range values {
				req.Header.Add(name, value)
			}
		}

		// Perform the request to the backend
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "Failed to reach backend", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Copy the headers and status code from the backend response
		for name, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}
		w.WriteHeader(resp.StatusCode)

		// Copy the body from the backend response
		if _, err := io.Copy(w, resp.Body); err != nil {
			log.Printf("Failed to copy response body: %v", err)
		}
	})

	log.Printf("Starting proxy server on port %s, forwarding to %s", port, backendURL)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
