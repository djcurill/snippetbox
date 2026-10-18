package main

import "net/http"

func (app *App) routes() *http.ServeMux {
	// server setup
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.Dir("./assets/static/"))))
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.createSnippet)
	mux.HandleFunc("POST /snippet/create", app.createSnippetPost)
	return mux
}
