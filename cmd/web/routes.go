// cmd/web/routes.go

package main

import "net/http"

// func (app *application) routes() *http.ServeMux {
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// return mux

	// here we mpass servemix as the next param to our middlware
	return secureHeaders(mux)

}

// TODO just an exampole of middleware pattern

func myMiddleware(next http.Handler) http.Handler {
	// fn := func(w http.ResponseWriter, r *http.Request) {
	// 	// some middlweare logic goes here
	// 	next.ServeHTTP(w, r)
	// }

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// some middlweare logic goes here
		// any code here will execute on the way down the chain of handlers

		// here we have some code for auth
		// if this code proves to be true, , we return
		// we stop the chain
		if !isAuthenticated(r) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
		// anda code hee will execpute on the way up the chain of handlers - when giving back the cotnrol
	})
}
