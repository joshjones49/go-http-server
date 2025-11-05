package routes

import (
	"net/http"
	"io"
	"strings"
	"go-http-server/config"
) 

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get(config.URL + "/posts")
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	
	body, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func GetCommentsForPost(w http.ResponseWriter, r *http.Request) {
	// Extract path after /posts/
	const prefix = "/posts/"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		http.Error(w, "Invalid URL: missing /posts/", http.StatusBadRequest)
		return
	}

	// Remove prefix and split remaining path
	remainder := r.URL.Path[len(prefix):]
	parts := strings.SplitN(remainder, "/", 3) // Only need first 3 parts max

	if len(parts) < 2 || parts[1] != "comments" {
		http.Error(w, "Invalid URL: expected /posts/{id}/comments", http.StatusBadRequest)
		return
	}

	// Check for extra segments (e.g., /posts/123/comments/trailing)
	if len(parts) == 3 && parts[2] != "" {
		http.Error(w, "Invalid URL: extra path segments", http.StatusBadRequest)
		return
	}

	postID := parts[0]
	if postID == "" {
		http.Error(w, "Invalid URL: missing post ID", http.StatusBadRequest)
		return
	}

	// Build upstream URL
	targetURL := config.URL + "/posts/" + postID + "/comments"

	// Proxy the request efficiently
	resp, err := http.Get(targetURL)
	if err != nil {
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Forward status code
	w.WriteHeader(resp.StatusCode)

	// Forward relevant headers
	if ct := resp.Header.Get("Content-Type"); ct != "" {
		w.Header().Set("Content-Type", ct)
	}
	if cl := resp.Header.Get("Content-Length"); cl != "" {
		w.Header().Set("Content-Length", cl)
	}

	// Stream response directly (zero-copy if possible)
	if _, err := io.Copy(w, resp.Body); err != nil {
		// Client disconnected or write error; log if needed
		// Note: Cannot send error after headers written
		return
	}
}