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
	dynamic := alice.New(app.sessionManager.LoadAndSave, app.requireCSRF)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignUp))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignUpPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	// pages that require authentication
	protected := dynamic.Append(app.requireAuthentication)
	mux.Handle("GET /snippet/create", protected.ThenFunc(app.createSnippet))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.createSnippetPost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	// syntactic sugar for chaining middleware
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
