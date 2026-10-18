package main

import "net/http"

// serverInternalError is a helper method for returning a 500 status
// internal error
func (app *App) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	method := r.Method
	path := r.URL.RequestURI()
	app.logger.Error(err.Error(), "method", method, "path", path)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Informs the client when an invalid request was made
func (app *App) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
