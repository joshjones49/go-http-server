package routes

import (
	"go-http-server/config"
	"io"
	"net/http"
	"strings"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get(config.URL + "/users")
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
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

func GetUserAlbums(w http.ResponseWriter, r *http.Request) {
	const prefix = "/users/"

	// strings.HasPrefix(string, text that should start the string), returns true or false 
	if !strings.HasPrefix(r.URL.Path, prefix) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// removes '/users' and assigns the leftover string to remainder
	remainder := r.URL.Path[len(prefix):]
	// finds the first occurance of 'albums'
	idEnd := strings.Index(remainder, "/albums")

	if idEnd <= 0 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	userID := remainder[:idEnd]

	if userID == "" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	suffix := remainder [idEnd:]
	if suffix != "/albums" && suffix != "/albums/" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	targetURL := config.URL+ "/users/"+ userID+ "/albums"

	resp, err := http.Get(targetURL)
	if err != nil {
		http.Error(w, "Failed to fetch albums", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)

	if ct := resp.Header.Get("Content-Type"); ct != "" {
		w.Header().Set("Content-Type", ct)
	}

	if cl := resp.Header.Get("Content-Length"); cl != "" {
		w.Header().Set("Content-Length", cl)
	}

	io.Copy(w, resp.Body)
}