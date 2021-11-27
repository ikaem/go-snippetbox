package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/ikaem/snippetbox/pkg/models"
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

// this is for serving a single file
// careful!: http.ServeFile foes not automatically sanitize the file path - we  have to sanitize it first with filepath.Clean()
// func downloadHandler(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "./ui/static/file.zip")
// }

// func home(w http.ResponseWriter, r *http.Request) {

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// if incorrect path, just send not found status code
	if r.URL.Path != "/" {
		// http.NotFound(w, r)
		app.notFound(w)
		return
	}

	// here we get the data

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// this would actually loop, and then send respoonse for each item
	// there is no ending here

	for _, snippet := range s {
		fmt.Fprintf(w, "%v\n", snippet)
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
		// app.errorLog.Println(err.Error())
		// http.Error(w, "Internal server error", http.StatusInternalServerError)
		app.serverError(w, err)
		return
	}

	// now, we will use Execute to write the template as the response body
	// the last variable, for which we passed nil
	// is any dynamic data we want to pass in
	// we have no dynamic data for now

	err = ts.Execute(w, nil)
	if err != nil {
		// app.errorLog.Println(err.Error())
		// http.Error(w, "Internal server error", http.StatusInternalServerError)
		app.serverError(w, err)
	}

	// w.Write([]byte("Hello from snippetbox"))
}

// here we access value storein a address
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// make sure to get the snippet id
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		// http.NotFound(w, r)
		app.notFound(w)
		return
	}

	// fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)

	// we return not found if nothing is found, or sever error if some other server errro

	s, err := app.snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}

		return
	}

	// and here we define data with the snippet and we pass the snippet value in
	data := &templateData{Snippet: s}

	// ok, now we define here the template

	files := []string{
		"./ui/html/show.page.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partial.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// it is funny how we return data to the user with fmt
	// and we actually return plan text, eve  though the thing is actualy object

	// fmt.Print(s)
	// fmt.Fprintf(w, "%v", s)
	// return

	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// making sure to send info to user that only POST is allowed
		w.Header().Set("Allow", http.MethodPost)

		app.clientError(w, http.StatusMethodNotAllowed)

		// this is just a helper function that combines writing header an d then sending content with Write
		// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	// w.Write([]byte("Created a new snippet..."))

	// this is just dummy data for now

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := "7"

	// now we passt his data hto the funciton to insert snippet

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// and now we redirect, widht it being out id
	// we also include stats
	// status should be in 3xx
	// we also forward the response and request
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
