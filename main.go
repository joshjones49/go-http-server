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

	mux.HandleFunc("/posts", routes.GetAllPosts)

	// console log to show where server is running
	fmt.Println("Server Running On Port" + config.PORT)
	// function that starts server
	http.ListenAndServe(config.PORT, mux)
}
