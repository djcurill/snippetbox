package main

import (
	"crypto/tls"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/djcurill/snippetbox/internal/models"
	"github.com/go-playground/form/v4"
	_ "github.com/jackc/pgx/v5/stdlib" // registers the "pgx" driver
)

type app struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	decoder        *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	// cmd line flag parsing
	var logLevel slog.Level

	addr := flag.String("addr", ":8080", "HTTP network address")
	flag.TextVar(&logLevel, "log-level", slog.LevelInfo, "logging level: debug, info, warn, error")

	flag.Parse()

	// logging setup
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))

	// create db connection, defer close
	db, err := openDB("pgx", dsn())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("app connected to db successfully")

	// parse templates and cache
	templates, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// intilialize sessions manager
	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour // expire after 12h
	sessionManager.IdleTimeout = 20 * time.Minute
	// cookies will only be sent to a users browser when using https
	sessionManager.Cookie.Secure = true

	// central app data structure for passing logger and db
	// connection utils
	app := &app{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templates,
		decoder:        form.NewDecoder(),
		sessionManager: sessionManager,
	}

	// tls settings
	tlsCfg := &tls.Config{
		MinVersion: tls.VersionTLS13,
	}

	// Idle timeout limits tell the server how long to keep a connection open (keep-alive)
	// by keeping the connection open, we reduce latency because the client can reuse the same
	// connection

	// Read timeouts are time limits for reading from the client
	// Read timeouts indicate how long to wait when reading a response before closing the connection.
	// After 5 seconds since receiving headers, the connection is closed.

	// Write timeouts are time limits for writing to the client
	// Write timeouts will trigger when the time taken to write data exceeds a certain limit since
	// when the client was first accepted. Write timeouts are typically longer than Read timesouts
	// since read timeouts start counting when the first header is received
	server := http.Server{
		Addr:         *addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  1 * time.Minute,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), logLevel),
		TLSConfig:    tlsCfg,
		// prevent a denial of service attack from a user sending too many headers in a single a request
		MaxHeaderValueCount: 100,
	}

	logger.Info("Starting http server", "addr", *addr)
	err = server.ListenAndServeTLS("./assets/tls/cert.pem", "./assets/tls/key.pem")
	logger.Error(err.Error())
	os.Exit(1)
}
