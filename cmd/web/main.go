package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("snippet/create", createSnippet)

	// here we create a file server
	// it is a handler
	// it servers files from the ./ui/static directory in this case
	// it is relative to the current project directory root
	// we also use the OS system's file ystem implementation for handling paths
	fileServer := http.FileServer(http.Dir("./ui/static"))

	// we use mux.Handle to register file server as handler for all URL paths that start with /static/
	// if the path matches, we want to strip that "/static prefix" - we are left with "/"

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// pay attention to log
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
