package routes

import (
	"net/http"
	"io"
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