package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/ikaem/snippetbox/pkg/forms"
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
	// if r.URL.Path != "/" {
	// 	// http.NotFound(w, r)
	// 	app.notFound(w)
	// 	return
	// }

	// panic("This is a deleberate panic!")

	// here we get the data

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// this is the new render helper

	app.render(w, r, "home.page.html", &templateData{Snippets: s})

	// this would actually loop, and then send respoonse for each item
	// there is no ending here

	// for _, snippet := range s {
	// 	fmt.Fprintf(w, "%v\n", snippet)
	// }

	// fmt.Printf("this is variable %v", s)

	// data :=
	// 	&templateData{Snippets: s}

	// files := []string{
	// 	"./ui/html/home.page.html",
	// 	"./ui/html/base.layout.html",
	// 	"./ui/html/footer.partial.html",
	// }

	// // so with  this we read the template file into a template set
	// // template set is the ts variable
	// // ts, err := template.ParseFiles("./ui/html/home.page.html")
	// ts, err := template.ParseFiles(files...)

	// if err != nil {
	// 	// app.errorLog.Println(err.Error())
	// 	// http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 	app.serverError(w, err)
	// 	return
	// now, we will use Execute to write the template as the response body
	// the last variable, for which we passed nil
	// is any dynamic data we want to pass in
	// we have no dynamic data for now

	// err = ts.Execute(w, data)
	// if err != nil {
	// app.errorLog.Println(err.Error())
	// http.Error(w, "Internal server error", http.StatusInternalServerError)
	// app.serverError(w, err)
	// }
}

// w.Write([]byte("Hello from snippetbox"))
// }

// here we access value storein a address
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// make sure to get the snippet id
	// id, err := strconv.Atoi(r.URL.Query().Get("id"))
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))

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

	// here we get the flash from the thing - session - cookie
	// flash := app.session.PopString(r, "flash")

	// again, using the helper
	app.render(w, r, "show.page.html", &templateData{
		Snippet: s,
		// this throws error because we dont have Flash field defined on the struct
		// Flash: flash,
	})

	// and here we define data with the snippet and we pass the snippet value in
	// data := &templateData{Snippet: s}

	// // ok, now we define here the template

	// files := []string{
	// 	"./ui/html/show.page.html",
	// 	"./ui/html/base.layout.html",
	// 	"./ui/html/footer.partial.html",
	// }

	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	// it is funny how we return data to the user with fmt
	// and we actually return plan text, eve  though the thing is actualy object

	// fmt.Print(s)
	// fmt.Fprintf(w, "%v", s)
	// return

	// err = ts.Execute(w, data)
	// if err != nil {
	// 	app.serverError(w, err)
	// }
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("Create a new snippet"))
	// app.render(w, r, "create.page.html", nil)
	app.render(w, r, "create.page.html", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	// /* just example to limit size of the body */
	// r.Body = http.MaxBytesReader(w, r.Body, 4096)
	// err := r.ParseForm()
	// if err != nil {
	// 	http.Error(w, "Bad request", http.StatusBadRequest)
	// 	return
	// }

	// we call parse form to pput dany data in post request body to the r.PostForm map
	// if any errors, we will send 400 error back

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		// cool how we dont have to pass in all data to the struct
		// we can just pass some of we name it
		app.render(w, r, "create.page.html", &templateData{Form: form})
		return
	}

	// and nbow we can use form get to insert data into the db
	// these vield are validated
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))

	if err != nil {
		app.serverError(w, err)
		return
	}

	// this is where we put the data inside the seesion (cookie)
	app.session.Put(r, "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)

	/* then we get data from the parse d data  */

	// title := r.PostForm.Get("title")
	// content := r.PostForm.Get("content")
	// expires := r.PostForm.Get("expires")

	// // we create a map to hold validation errors
	// errors := make(map[string]string)

	// // ćcheck that the title field is not emtpy, or bigger than 100 chars
	// if strings.TrimSpace(title) == "" {
	// 	errors["title"] = "This field cannot be blank"
	// } else if utf8.RuneCountInString(title) > 100 {
	// 	errors["title"] = "This field is too long (maximum is 100 characters"
	// }

	// if strings.TrimSpace(content) == "" {
	// 	errors["content"] = "This field cannot be blank"
	// }

	// // we also check that expire value is a valid one
	// if strings.TrimSpace(expires) == "" {
	// 	errors["expires"] = "This field cannot be blank"
	// } else if expires != "365" && expires != "7" && expires != "1" {
	// 	errors["expires"] = "This field is invalid"
	// }

	// // now we just check if there are any errors

	// if len(errors) > 0 {
	// 	// fmt.Fprint(w, errors)
	// 	// here we just send back that form data again
	// 	// and in the data, we put back the data
	// 	app.render(w, r, "create.page.html", &templateData{
	// 		FormErrors: errors,
	// 		FormData:   r.PostForm,
	// 	})
	// 	return
	// }

	// // just test

	// // items is some field on the parsed form object
	// // for i, item := range r.PostForm["items"] {
	// // 	// fmt.Fprintf(w, "%d: item %s\n", i, item)
	// // }

	// // if r.Method != http.MethodPost {
	// // 	// making sure to send info to user that only POST is allowed
	// // 	w.Header().Set("Allow", http.MethodPost)

	// // 	app.clientError(w, http.StatusMethodNotAllowed)

	// // 	// this is just a helper function that combines writing header an d then sending content with Write
	// // 	// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	// // 	return
	// // }

	// // w.Write([]byte("Created a new snippet..."))

	// // this is just dummy data for now

	// // title := "O snail"
	// // content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	// // expires := "7"

	// // now we passt his data hto the funciton to insert snippet

	// id, err := app.snippets.Insert(title, content, expires)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	// and now we redirect, widht it being out id
	// we also include stats
	// status should be in 3xx
	// we also forward the response and request
	// http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}

// example handler that spins up its own goroutine for some additoanly work
// and then also handles any possible panicikging there

func myExampleRecoverGoSubroutinePanicHandler(w http.ResponseWriter, r *http.Request) {
	// stuff

	// spin up a new goroutine
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(fmt.Errorf("%s\n%s", err, debug.Stack()))
			}
		}()

		// doSomeBackgroundProcessing()
	}()

	w.Write([]byte("Ok"))
}
