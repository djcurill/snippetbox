package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/djcurill/snippetbox/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib" // registers the "pgx" driver
)

type App struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
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

	db, err := openDB("pgx", dsn())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("app connected to db successfully")

	app := &App{
		logger:   logger,
		snippets: &models.SnippetModel{DB: db},
	}
	id, err := app.snippets.Insert("testing", "1,2,3", 7)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	fmt.Println(id)

	server := http.Server{
		Addr:        *addr,
		ReadTimeout: 5 * time.Second,
		IdleTimeout: 5 * time.Second,
		Handler:     app.routes(),
	}

	logger.Info("Starting http server", "addr", *addr)
	err = server.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
