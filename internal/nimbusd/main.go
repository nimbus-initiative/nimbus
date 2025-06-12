// Package main is the entry point for the nimbusd agent.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"nimbus/internal/nimbusd/agent"
	"nimbus/internal/nimbusd/config"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Initialize configuration
	cfg, err := config.Load("/etc/nimbus/nimbusd.toml")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Configure logging
	configureLogging(cfg)


	// Create a context that cancels on interrupt signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	setupSignalHandling(cancel)

	// Initialize and start the agent
	agt, err := agent.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create agent")
	}

	log.Info().
		Str("version", version).
		Str("commit", commit).
		Str("build_date", date).
		Msg("Starting nimbusd agent")

	// Run the agent
	if err := agt.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("Agent failed")
	}
}

func configureLogging(cfg *config.Config) {
	// Set log level
	level, err := zerolog.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Configure log format
	if cfg.Log.Format == "json" {
		// JSON is the default format
	} else {
		// Pretty console output for development
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}
		log.Logger = log.Output(output)
	}
}

func setupSignalHandling(cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		log.Info().Str("signal", sig.String()).Msg("Received signal, shutting down...")
		cancel()
	}()
}
