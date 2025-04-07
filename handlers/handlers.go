package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"c:\Users\iamma\Programs\trae\Go-URL-Shortener\storage"
)

// Handler handles HTTP requests for the URL shortener
type Handler struct {
	store storage.URLStore
}

// NewHandler creates a new Handler with the given store
func NewHandler(store storage.URLStore) *Handler {
	return &Handler{store: store}
}

// ShortenRequest represents the JSON request for shortening a URL
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenResponse represents the JSON response for a shortened URL
type ShortenResponse struct {
	ShortURL string `json:"short_url"`
	OrigURL  string `json:"original_url"`
}

// ShortenHandler handles requests to shorten URLs
func (h *Handler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate URL
	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	// Add scheme if missing
	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		req.URL = "http://" + req.URL
	}

	// Save URL to storage
	shortKey, err := h.store.Save(req.URL)
	if err != nil {
		http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}

	// Construct the short URL
	shortURL := "http://" + r.Host + "/" + shortKey

	// Prepare response
	resp := ShortenResponse{
		ShortURL: shortURL,
		OrigURL:  req.URL,
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// RedirectHandler handles redirection from short URLs to original URLs
func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the short key from the path
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		// Serve a simple HTML page for the root path
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<html>
				<head><title>URL Shortener</title></head>
				<body>
					<h1>URL Shortener</h1>
					<p>Use the /shorten endpoint to create short URLs.</p>
				</body>
			</html>
		`))
		return
	}

	// Look up the original URL
	origURL, err := h.store.Load(path)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	// Redirect to the original URL
	http.Redirect(w, r, origURL, http.StatusMovedPermanently)
}

// HealthCheckHandler provides a simple health check endpoint
func (h *Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}