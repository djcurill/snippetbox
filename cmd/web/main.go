package main

import (
	"log"
	"net/http"
	"time"
)

const port string = "8080"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", createSnippet)
	mux.HandleFunc("POST /snippet/create", createSnippetPost)

	server := http.Server{
		Addr:        ":" + port,
		ReadTimeout: 5 * time.Second,
		IdleTimeout: 5 * time.Second,
		Handler:     mux,
	}

	log.Print("Starting server on port: ", port)
	err := server.ListenAndServe()
	log.Fatal(err)
}
