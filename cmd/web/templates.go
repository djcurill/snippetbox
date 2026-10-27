package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/djcurill/snippetbox/internal/models"
)

// template data struct to hold all elements
// passed into html templates to render data dynamically
type templateData struct {
	CurrentYear int
	Snippet     models.Snippet
	Snippets    []models.Snippet
	Form        any
}

// custom template functions can accept any sort of inputs
// but they must always return a single value, except
// if you want to return an error as well. That is allowed.
var tmplFunctions = template.FuncMap{
	"humanDate": humanDate,
}

func (app *app) newTemplateData() templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
	}
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 15:04")
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./assets/html/pages/*.gohtml")

	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// template functions must be registered with a new template set
		// before parsing files
		ts, err := template.New(page).Funcs(tmplFunctions).ParseFiles("./assets/html/base.gohtml")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./assets/html/partials/*.gohtml")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// add to cache
		cache[name] = ts

	}

	return cache, nil
}
