// Package main is the entry point for the nimbusctl command-line tool.
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"nimbus/internal/nimbusctl/commands"
)

var (
	// These are set during build time using -ldflags
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Configure logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Create root command
	rootCmd := &cobra.Command{
		Use:     "nimbusctl",
		Short:   "Nimbus CLI",
		Long:    "Nimbus is a decentralized, peer-optional, GitOps-driven bare metal cloud platform.",
		Version: fmt.Sprintf("%s (commit: %s, date: %s)", version, commit, date),
	}

	// Add global flags
	var (
		verbose   bool
		configDir string
	)

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVar(&configDir, "config-dir", "", "Path to configuration directory")

	// Set config directory
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get user home directory")
		}
		configDir = filepath.Join(home, ".nimbus")
	}

	// Create command context
	ctx := context.Background()
	cmdCtx := &commands.Context{
		ConfigDir: configDir,
		Verbose:   verbose,
	}

	// Initialize commands
	commands.AddCommands(rootCmd, cmdCtx)
	
	// Execute the root command
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Fatal().Err(err).Msg("Command failed")
	}
}
