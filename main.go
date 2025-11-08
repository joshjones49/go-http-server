package main

import (
	"fmt"
	"go-http-server/config"
	"go-http-server/routes"
	"net/http"
)

func main() {
	// used to direct requests, think of it as a traffic controller
	mux := http.NewServeMux()

	// handlers for post related requests ===============>
	mux.HandleFunc("GET /posts", routes.GetAllPosts)
	mux.HandleFunc("GET /posts/{id}/comments", routes.GetCommentsForPost)

	// handlers for users realated requests
	mux.HandleFunc("GET /users", routes.GetAllUsers)
	mux.HandleFunc("GET /users/{id}/albums", routes.GetUserAlbums)

	// console log to show where server is running
	fmt.Println("Server Running On Port" + config.PORT)
	// function that starts server
	http.ListenAndServe(config.PORT, mux)
}
