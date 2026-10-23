package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/djcurill/snippetbox/internal/models"
)

// Create a home pathfor the webpages
func (app *app) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	data := app.newTemplateData()
	data.Snippets = snippets

	err = app.render(w, r, "home.gohtml", http.StatusOK, data)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// Snippet view
func (app *app) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusNotFound)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNoRecord):
			http.NotFound(w, r)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// populate templateData
	data := app.newTemplateData()
	data.Snippet = snippet

	err = app.render(w, r, "snippetView.gohtml", http.StatusOK, data)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// Create snippet
func (app *app) createSnippet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "diplay a view for creating a snippet")
}

// POST endpoint for creating the snipper
func (app *app) createSnippetPost(w http.ResponseWriter, r *http.Request) {
	// Status code 201 to indicate content has been created
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.internalServerError(w, r, err)
	}

	// Redirect to the view of the snippet after it has been created
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
