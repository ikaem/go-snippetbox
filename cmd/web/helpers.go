// cmd/web/helpers.go
package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	// this trace thing is just a string
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	// app.errorLog.Println(trace)
	app.errorLog.Output(2, trace)

	// and now we send the error to the user - we acutall send the response
	// so we send usual text for the 500 code, and we send 500 code itself

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}

// note that here we accept status code instead of the error itself
func (app *application) clientError(w http.ResponseWriter, status int) {

	http.Error(w, http.StatusText(status), status)

}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.CurrentYear = time.Now().Year()
	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// so we want to get needed template from the cache

	ts, ok := app.templateCache[name]
	if !ok {
		// note the use of Errorf metod to format error
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
	}

	// we need to initialize a buffer
	// buffer is a variable size of bbytes
	// has read and writes methods on it
	buf := new(bytes.Buffer)

	// if all is good, we actually render the template with data

	// err := ts.Execute(w, td)

	// now we actually write the template to the bvuffer, and check for errors

	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// and now we can use buffers writeTo method to write its contents to the writer i guess
	buf.WriteTo(w)
}
