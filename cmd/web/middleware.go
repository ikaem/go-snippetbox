// cmd/web/middleware.go
package main

import (
	"fmt"
	"net/http"
)

// i guess every route handler is this type of function
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// how we here do stuff for what we make the middlweare
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// here we create a deferred function
		// and we also call it immediately, just to make sure it is called

		defer func() {
			// here we check if there has been panic or not
			// we check that with a call to a builtin recover function
			// if it returns error, i guess it means that there has been an error
			// more correct name would be isPanic or similar, prolly
			if err := recover(); err != nil {
				// we set header on the writer
				w.Header().Set("Connection", "close")

				// then we call the server error helper
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		// then we normally proceed with passing the request to the next handler

		next.ServeHTTP(w, r)
	})
}
