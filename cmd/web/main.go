package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// i guess flag package is used to parse command line flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	boolFlag := flag.Bool("boolFlag", true, "just a book lest flag")
	flag.Parse()
	fmt.Println("This is a bool flag", *boolFlag)
	fmt.Println("Here is address", *addr) // :80

	type Config struct {
		Addr      string
		StaticDir string
	}

	// flag.Parse()

	// this doesn0t work when variables are in file - it would work in evn varabe
	testEnvVar := os.Getenv("TEST")
	fmt.Println(testEnvVar)

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
	log.Println("Starting server on", *addr)
	err := http.ListenAndServe(*addr, mux)

	log.Fatal(err)
}
