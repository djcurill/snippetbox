package main

import (
	"errors"
	"fmt"
	"github.com/djcurill/snippetbox/internal/models"
	"github.com/djcurill/snippetbox/internal/validators"
	"net/http"
	"strconv"
)

// use go-playground for automatic form parsing.
// tags indicate which form fields map to struct
// attributes. a "-" mean to ignore this value
type snippetCreateForm struct {
	Title                string `form:"title"`
	Content              string `form:"content"`
	Expires              int    `form:"expires"`
	validators.Validator `form:"-"`
}

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
	data := app.newTemplateData()
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	err := app.render(w, r, "create.gohtml", http.StatusOK, data)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// POST endpoint for creating the snipper
func (app *app) createSnippetPost(w http.ResponseWriter, r *http.Request) {

	var form snippetCreateForm
	err := app.decodeForm(r, &form)

	if err != nil {
		app.logger.Error("error parsing form", "err", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(form.NotBlank(form.Title), "title", "this field cannot be empty")
	form.CheckField(form.MaxChars(form.Title, 100), "title", "this field cannot be longer than 100 chars")
	form.CheckField(form.NotBlank(form.Content), "content", "this field cannot be empty")
	form.CheckField(validators.PermittedValue(form.Expires, 1, 7, 365), "expires", "this field must equal 1,7 or 365")

	if !form.Valid() {
		data := app.newTemplateData()
		data.Form = form
		err = app.render(w, r, "create.gohtml", http.StatusUnprocessableEntity, data)
		if err != nil {
			app.internalServerError(w, r, err)
		}
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Redirect to the view of the snippet after it has been created
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
