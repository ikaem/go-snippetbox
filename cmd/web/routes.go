// cmd/web/routes.go

package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func some() {

	// just an example for crateing a router and registering a route with the pat package
	mux := pat.New()
	mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))

}

// func (app *application) routes() *http.ServeMux {
func (app *application) routes() http.Handler {

	// this is creating the middleware chain
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// this is now using the standard middleware, then use the handler
	return standardMiddleware.Then(mux)

	// return mux

	// here we mpass servemix as the next param to our middlware
	// return app.recoverPanic(app.logRequest(secureHeaders(mux)))

	// return alice.New(myMiddleware1, myMiddleware2, myMiddleware3).Then(myHandler)

	// myChain := alice.New(myMiddleware1, myMiddleware2)
	// myChain2 := myChain.Append(myMiddleware3)
	// return myChain2.Then(myHandler)

}

// TODO just an exampole of middleware pattern

// func myMiddleware(next http.Handler) http.Handler {
// 	// fn := func(w http.ResponseWriter, r *http.Request) {
// 	// 	// some middlweare logic goes here
// 	// 	next.ServeHTTP(w, r)
// 	// }

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// some middlweare logic goes here
// 		// any code here will execute on the way down the chain of handlers

// 		// here we have some code for auth
// 		// if this code proves to be true, , we return
// 		// we stop the chain
// 		if !isAuthenticated(r) {
// 			w.WriteHeader(http.StatusForbidden)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 		// anda code hee will execpute on the way up the chain of handlers - when giving back the cotnrol
// 	})
// }
