package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"linkbio/internal/config"
	"linkbio/internal/pkg/logger"
	"linkbio/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// Initialize logger
	var log = logger.New(cfg.LogLevel)
	if cfg.IsDevelopment() {
		log = logger.NewDevelopment()
	}

	log.Info("starting linkbio",
		"env", cfg.Env,
		"port", cfg.Port,
		"log_level", cfg.LogLevel,
	)

	// Create server
	srv, err := server.New(cfg, log)
	if err != nil {
		log.Error("failed to create server", "error", err)
		os.Exit(1)
	}

	// Channel to listen for errors from server
	serverErrors := make(chan error, 1)

	// Start server in goroutine
	go func() {
		log.Info("server listening", "port", cfg.Port)
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	// Channel to listen for interrupt signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal or error
	select {
	case err := <-serverErrors:
		log.Error("server error", "error", err)
		os.Exit(1)

	case sig := <-shutdown:
		log.Info("shutdown signal received", "signal", sig.String())

		// Create context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := srv.Shutdown(ctx); err != nil {
			log.Error("graceful shutdown failed", "error", err)
			os.Exit(1)
		}

		log.Info("server stopped gracefully")
	}
}
