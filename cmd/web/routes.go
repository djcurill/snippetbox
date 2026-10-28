package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *app) routes() http.Handler {
	// initialize new server mux
	mux := http.NewServeMux()

	// setup static file server
	fileServer := http.FileServer(http.Dir("./assets/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// setup dynamic routes
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.createSnippet))
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.createSnippetPost))

	// syntactic sugar for chaining middleware
	chain := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return chain.Then(mux)
}
