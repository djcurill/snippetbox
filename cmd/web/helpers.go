package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-playground/form/v4"
	_ "github.com/jackc/pgx/v5/stdlib" // registers the "pgx" driver
)

// serverInternalError is a helper method for returning a 500 status
// internal error
func (app *app) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	method := r.Method
	path := r.URL.RequestURI()
	app.logger.Error(err.Error(), "method", method, "path", path)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Informs the client when an invalid request was made
func (app *app) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *app) decodeForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.decoder.Decode(dst, r.PostForm)
	if err != nil {
		if _, ok := errors.AsType[*form.InvalidDecoderError](err); ok {
			panic(err)
		}
		return err
	}
	return nil
}

func (app *app) render(w http.ResponseWriter, r *http.Request, page string, status int, data templateData) error {

	ts, ok := app.templateCache[page]
	if !ok {
		return fmt.Errorf("unable to retrieve %s from template cache", page)
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		return fmt.Errorf("error executing template for %s", page)
	}

	w.WriteHeader(status)
	_, err = buf.WriteTo(w)

	return err
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func dsn() string {
	db_user := getEnv("SNIPPET_DB_USER", "web")
	db_name := getEnv("SNIPPET_DB_NAME", "snippetbox")
	db_pass := getEnv("SNIPPET_DB_PASSWORD", "dev")
	db_host := getEnv("SNIPPET_DB_HOST", "localhost")
	return fmt.Sprintf("postgres://%s:%s@%s:5432/%s", db_user, db_pass, db_host, db_name)
}

func openDB(driver, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
