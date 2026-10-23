package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *app) routes() http.Handler {
	// server setup
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.Dir("./assets/static/"))))
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.createSnippet)
	mux.HandleFunc("POST /snippet/create", app.createSnippetPost)

	// syntactic sugar for chaining middleware
	chain := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return chain.Then(mux)
}
