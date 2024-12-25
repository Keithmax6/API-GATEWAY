package main

import (
	"io"
	"net/http"
	"strings"
)

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	// Unified endpoint
	router.HandleFunc("/api/", apiHandler) // Note the trailing slash

	return router
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the target service from query params
	target := r.URL.Query().Get("service")
	if target == "" {
		http.Error(w, "Service query parameter is missing", http.StatusBadRequest)
		return
	}

	// Map services to their respective URLs
	serviceMap := map[string]string{
		"users":  "http://localhost:8081",
		"orders": "http://localhost:8082",
	}

	// Find the target URL
	targetURL, exists := serviceMap[target]
	if !exists {
		http.Error(w, "Invalid service specified", http.StatusNotFound)
		return
	}

	// Forward the request to the appropriate service
	proxyRequest(targetURL, r.URL.Path, w, r)
}

func proxyRequest(target, path string, w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}

	// Remove "/api" from the path and construct the target URL
	trimmedPath := strings.TrimPrefix(path, "/api")
	targetURL := target + trimmedPath

	// Create a new request to the target service
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Forward headers
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to connect to service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the response
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
