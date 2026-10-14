package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Create a home path for the webpages
func home(w http.ResponseWriter, r *http.Request) {
	// Add a custom header indicating the backend server
	w.Header().Add("Server", "Go")

	ts, err := template.ParseFiles("./assets/html/pages/home.tmpl")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Snippet view
func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "View id is %d", id)
}

// Create snippet
func createSnippet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Display a form for creating a snippet")
}

// POST endpoint for creating the snipper
func createSnippetPost(w http.ResponseWriter, r *http.Request) {
	// Status code 201 to indicate content has been created
	w.WriteHeader(http.StatusCreated)

	// Write snippet
	fmt.Fprintf(w, "Creating snippet...")
}
