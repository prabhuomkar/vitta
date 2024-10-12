package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"vitta/config"
	"vitta/database"
	"vitta/handlers"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error("error reading config", "error", err)
		os.Exit(1)
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: cfg.LogLevel == slog.LevelDebug,
		Level:     cfg.LogLevel,
	})))

	db, err := database.New(&database.Config{
		URL:     cfg.DatabaseURL,
		Timeout: cfg.DatabaseTimeout,
	})
	if err != nil {
		slog.Error("error connecting to database", "error", err)
		os.Exit(1)
	}

	handler := handlers.New(cfg, db)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	go func() {
		slog.Info("starting http server")

		err := server.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			slog.Info("stopped http server")
		} else if err != nil {
			slog.Error("error starting http server", "error", err)
			os.Exit(1)
		}
	}()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	<-shutdownSignal
	slog.Info("stopping http server")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("error shutting down http server", "error", err)
	}
}
