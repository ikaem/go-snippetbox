package backup

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// ok, now this is a hiome handler functions
// it writes a biyte slice that containes a string as the response body
func home(w http.ResponseWriter, r *http.Request) {

	// here we check url path on the request
	if r.URL.Path != "/" {
		http.NotFound(w, r)

		return
	}
	w.Write([]byte("Hello from snippetbox"))
}

// add showSnippet handler function

func showSnippet(w http.ResponseWriter, r *http.Request) {

	// here we try t oget that id, and convert ti to an intger
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// this function will interporalte id value with the response
	// it will also write it ti the http.,ResponseWriter
	fmt.Fprintf(w, "Display a specific snipper with ID %d", id)

	// w.Write([]byte("Display specific snippet"))
}

// add createSnippet handler function
func createSnippet(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		// so it seems we send the response twice - once for the header, once the content?
		w.Header().Set("Allow", http.MethodPost)
		// w.WriteHeader(http.StatusMethodNotAllowed)
		// w.Write([]byte("Method not allowed"))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// be speicific with setting header
	w.Header().Set("Content-Type", "application/json")

	// overwrites key if exists
	w.Header().Set("Cache-Control", "public, max-age=12312312312")
	// this is if the http.DetectContentType() function cannot detect the type
	w.Header().Set("Content-Type", "application/octet-stream")
	// appends a new header with this name - it can be called multiuple times
	w.Header().Add("Cache-Control", "public, max-age=12312312312")
	w.Header().Add("Cache-Control", "public, max-age=12312312312")
	// deletes all values for a header
	w.Header().Del("Cache-Control")

	// this is manually directly setting a header to a value
	w.Header()["Date"] = nil

	// get first value for the specified header
	w.Header().Get("Cache-Control")

	w.Write([]byte("Create a new snippet"))

}

func main() {

	// hee we initialize new servermix
	// server mux is i guess a router
	// and we register home function as the handler for the / url pattern

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)
	// can include hosts in URL patterns
	// mux.HandleFunc("foo.example.org", goToFooExample)

	// http.HandleFunc("/", home)
	// http.HandleFunc("/snippet", showSnippet)
	// http.HandleFunc("/snippet/create", createSnippet)

	// ahnd now we can user Listen and server to start in a new web sercer
	// it needs two params - TCP address to listen on , and the server mix we just created
	// if it returns areror , we will log it and exit
	// note taht log.Fatal will only be reached if the ListenAndServe stops with error

	log.Println("Starting server on :4000")
	// here we pass tcp address and the serve mux - the router
	err := http.ListenAndServe(":4000", mux)
	// note that we dont pass the router here - teh default one is used
	// err := http.ListenAndServe(":4000", nil)
	log.Fatal(err)

}
