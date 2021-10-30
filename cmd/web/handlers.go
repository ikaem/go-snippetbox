package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

// type Handler interface {
// 	ServeHTTP(ResponseWriter, *Request)
// }

// type home struct{}

// func (h *home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("This is my home page"))
// }

// mux := http.NewServeMux()
// mux.Handle("/", &home{})

// mux.Handle("/", http.HandlerFunc(home))
// mux.HandlerFunc("/", home)

// func ListenAndServe(addr string, handler Handler) error

func home(w http.ResponseWriter, r *http.Request) {
	// if incorrect path, just send not found status code
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/home.page.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partial.html",
	}

	// so with  this we read the template file into a template set
	// template set is the ts variable
	// ts, err := template.ParseFiles("./ui/html/home.page.html")
	ts, err := template.ParseFiles(files...)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// now, we will use Execute to write the template as the response body
	// the last variable, for which we passed nil
	// is any dynamic data we want to pass in
	// we have no dynamic data for now

	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	// w.Write([]byte("Hello from snippetbox"))
}

func showSnippet(w http.ResponseWriter, r *http.Request) {
	// make sure to get the snippet id
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// making sure to send info to user that only POST is allowed
		w.Header().Set("Allow", http.MethodPost)

		// this is just a helper function that combines writing header an d then sending content with Write
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	w.Write([]byte("Created a new snippet..."))
}
