package main

import (
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
		templateCache:  templates,
		decoder:        form.NewDecoder(),
		sessionManager: sessionManager,
	}

	server := http.Server{
		Addr:        *addr,
		ReadTimeout: 5 * time.Second,
		IdleTimeout: 5 * time.Second,
		Handler:     app.routes(),
		ErrorLog:    slog.NewLogLogger(logger.Handler(), logLevel),
	}

	logger.Info("Starting http server", "addr", *addr)
	err = server.ListenAndServeTLS("./assets/tls/cert.pem", "./assets/tls/key.pem")
	logger.Error(err.Error())
	os.Exit(1)
}
