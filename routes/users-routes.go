package routes

import (
	"fmt"
	"go-http-server/config"
	"io"
	"net/http"
	"strings"
)

// Fetches all users from jsonplaceholder
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

// Gets all albums associated with a user
func GetUserAlbums(w http.ResponseWriter, r *http.Request) {
	const prefix = "/users/"

	// strings.HasPrefix(string, text that should start the string), returns true or false
	// check to see if r.URL.Path DOES NOT contain '/users/' which is assigned to the variable 'prefix'
	// if this check results in true the ResponseWriter 'w' sends "Invalid URL" to the client and status code 400
	// returns, ending the handler function early
	if !strings.HasPrefix(r.URL.Path, prefix) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// removes '/users/' and assigns the leftover string to remainder
	// r.URL.Path == '/users/{id}/albums'
	// remainder == '{id}/albums
	remainder := r.URL.Path[len(prefix):]
	fmt.Println(remainder)

	// finds the index of the first occurance of '/albums' and assigns it to idEnd
	// this will be used to find the index of where {id} should be
	idEnd := strings.Index(remainder, "/albums")

	// checks to see if idEnd is a negative number, if so sends a 400 error and ends the handler early
	if idEnd <= 0 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// remainder has an assigned value of '{id}/albums'
	// idEnd has an assigned value of 1 since 1 is the index at which '/albums' starts in the remainder variable string
	// remainder[:idEnd] slices remainder at index 1 and assigns everything to the left to userID
	userID := remainder[:idEnd]
	// if userID is empty return 400 error and return, ends handler early
	if userID == "" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// remainder has an assigned value of '{id}/albums'
	// idEnd has a value of 1
	// remainder[idEnd:] slices the string of remainder starting at index 1 and assigns everything to the right to suffix
	suffix := remainder[idEnd:]
	fmt.Println(suffix)
	// checks to see if suffix is neither '/albums' or '/albums/
	// if check is true, send 400 error, return, end handler early
	if suffix != "/albums" && suffix != "/albums/" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	targetURL := config.URL+ "/users/"+ userID+ "/albums"
	fmt.Println(targetURL)

	// 
	resp, err := http.Get(targetURL)
	fmt.Println(resp.Body)
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