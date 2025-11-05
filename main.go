package main

import (
	"fmt"
	"net/http"
)

var PORT string = "8080"

func main() {
	// used to direct requests, think of it as a traffic controller
	mux := http.NewServeMux()

	// console log to show where sever is running
	fmt.Println("Server Running On Port: " + PORT)
	// function that starts server
	http.ListenAndServe(PORT, mux)
}
