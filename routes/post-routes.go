package routes

import (
	"go-http-server/config"
	"io"
	"net/http"
	"strings"
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
	// defines constant for expected base path
	const prefix = "/posts/"
	// checks if the request starts with '/posts/'
	// if not returns a 404 error and exits handler early
	if !strings.HasPrefix(r.URL.Path, prefix) {
		http.Error(w, "Invalid URL: missing /posts/", http.StatusBadRequest)
		return
	}
	// removes the '/posts/ prefix leaving everything after
	remainder := r.URL.Path[len(prefix):]
	// finds the first occurance of comments
	// returns the index of where it starts pr -1 if not found
	// avoids splitting into slices
	idEnd := strings.Index(remainder, "/comments")
	// isEnd will be -1 if it isn't found, or if it starts at position 0 (meaning no ID)
	// both are invalid so return 400
	if idEnd <= 0 {
		http.Error(w, "Invalid URL: expected /posts/{id}/comments", http.StatusBadRequest)
		return
	}
	// extracts the part before '/comments' and turns it into the postID
	postID := remainder[:idEnd]
	// checks to make sure postID is not empty
	if postID == "" {
		http.Error(w, "Invalid URL: missing post ID", http.StatusBadRequest)
		return
	}
	// checks whether the suffix is not exactly '/comments'
	// only two valid forms '/comments' or '/comments/'
	suffix := remainder[idEnd:]
	if suffix != "/comments" && suffix != "/comments/" {
		http.Error(w, "Invalid URL: extra path segments", http.StatusBadRequest)
		return
	}
	// contructs the API URL
	targetURL := config.URL + "/posts/" + postID + "/comments"
	// makes get request to service
	resp, err := http.Get(targetURL)
	// if err is not nil then return error code 500
	if err != nil {
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}
	// closes the resposnse body afrer the fucntion is finished
	defer resp.Body.Close()
	// writes the status code of request
	w.WriteHeader(resp.StatusCode)
	// copies so client knows how to parse body
	if ct := resp.Header.Get("Content-Type"); ct != "" {
		w.Header().Set("Content-Type", ct)
	}
	if cl := resp.Header.Get("Content-Length"); cl != "" {
		w.Header().Set("Content-Length", cl)
	}
	// streams response body to client
	io.Copy(w, resp.Body) // errors are non-recoverable after headers
}
