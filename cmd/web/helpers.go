// cmd/web/helpers.go
package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
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
