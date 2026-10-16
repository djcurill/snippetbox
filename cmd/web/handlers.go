package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// Create a home path for the webpages
func (app *App) home(w http.ResponseWriter, r *http.Request) {
	// Add a custom header indicating the backend server
	w.Header().Add("Server", "Go")

	// We must define base first, then home. The base template
	// references other named templates defined in subsequent .gohtml
	// files
	files := []string{
		"./assets/html/base.gohtml",
		"./assets/html/pages/home.gohtml",
		"./assets/html/partials/nav.gohtml",
		"./assets/html/none.gohtml",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Snippet view
func (app *App) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "View id is %d", id)
}

// Create snippet
func (app *App) createSnippet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Display a form for creating a snippet")
}

// POST endpoint for creating the snipper
func (app *App) createSnippetPost(w http.ResponseWriter, r *http.Request) {
	// Status code 201 to indicate content has been created
	w.WriteHeader(http.StatusCreated)

	// Write snippet
	fmt.Fprintf(w, "Creating snippet...")
}
