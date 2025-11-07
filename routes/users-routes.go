package routes

import (
	"go-http-server/config"
	"io"
	"net/http"
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