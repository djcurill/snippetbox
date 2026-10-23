package main

import (
	"fmt"
	"net/http"
)

func commonHeaders(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {

		// restrict where javascript, images, fonts etc. may be source from
		// we only allow javascript to be executed from the hosting domain, while
		// style sheets may be sourced from google or the hosting domain (this src)
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		// when accessing another site from this one, the referrer URL
		// will only include the origin (scheme, host and port)
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")

		// disable any inferring behaviour of the content type
		// explicitly trust what the content type header provides
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// prevent clickjacking
		w.Header().Set("X-Frame-Options", "deny")

		// disables the browser’s built-in XSS (cross-site scripting) filtering
		w.Header().Set("X-XSS-Protection", "0")

		w.Header().Set("Server", "Go")

		// call the next handler
		next.ServeHTTP(w, r)

	}
	return http.HandlerFunc(f)
}

func (app *app) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)
		app.logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)
		next.ServeHTTP(w, r)
	})
}

func (app *app) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// deferred functions run in order of last in, first out
		// when a panic occurs, all deferred functions are kicked off
		// we catch the panic value, log it, then close the connection
		defer func() {
			// pv is of type `any`
			pv := recover()

			if pv != nil {
				// tell go http server module to close the connection
				// this acts as a trigger to do so
				w.Header().Set("Connection", "close")

				// convert pv to its textual representation
				err := fmt.Errorf("panic: %v", pv)

				// internal server error 500
				app.internalServerError(w, r, err)
			}

		}()

		next.ServeHTTP(w, r)
	})
}
