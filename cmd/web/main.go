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

	server := http.Server{
		Addr:        *addr,
		ReadTimeout: 5 * time.Second,
		IdleTimeout: 5 * time.Second,
		Handler:     app.routes(),
	}

	logger.Info("Starting http server", "addr", *addr)
	err := server.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
