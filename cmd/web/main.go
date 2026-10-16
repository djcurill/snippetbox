package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type App struct {
	logger *slog.Logger
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
	app := &App{
		logger: logger,
	}

	// server setup
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.Dir("./assets/static/"))))
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.createSnippet)
	mux.HandleFunc("POST /snippet/create", app.createSnippetPost)

	server := http.Server{
		Addr:        *addr,
		ReadTimeout: 5 * time.Second,
		IdleTimeout: 5 * time.Second,
		Handler:     mux,
	}

	logger.Info("Starting http server", "addr", *addr)
	err := server.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
