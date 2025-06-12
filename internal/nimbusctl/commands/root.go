// Package commands implements the nimbusctl command-line interface.
package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Context holds the context for command execution.
type Context struct {
	ConfigDir string
	Verbose   bool
}

// AddCommands adds all commands to the root command.
func AddCommands(rootCmd *cobra.Command, ctx *Context) {
	// Initialize config directory
	if err := ensureConfigDir(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
		os.Exit(1)
	}

	// Add subcommands
	rootCmd.AddCommand(
		newVersionCommand(),
		newConfigCommand(ctx),
		newProviderCommand(ctx),
		newInstanceCommand(ctx),
		newNodeCommand(ctx),
	)
}

// ensureConfigDir creates the configuration directory if it doesn't exist.
func ensureConfigDir(ctx *Context) error {
	if err := os.MkdirAll(ctx.ConfigDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create default config file if it doesn't exist
	configPath := filepath.Join(ctx.ConfigDir, "config.toml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		defaultConfig := `# Nimbus CLI Configuration

# Default provider to use if not specified
# default_provider = ""

# Default region to use if not specified
# default_region = ""

# API server address
# api_server = "https://api.nimbus.example.com"

# Authentication token (automatically managed)
# auth_token = ""
`
		if err := os.WriteFile(configPath, []byte(defaultConfig), 0600); err != nil {
			return fmt.Errorf("failed to create default config: %w", err)
		}
	}

	return nil
}

// newVersionCommand creates a new version command.
func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(cmd.Root().Version)
		},
	}
}

// newConfigCommand creates a new config command.
func newConfigCommand(ctx *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage nimbusctl configuration",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "view",
			Short: "View current configuration",
			Run: func(cmd *cobra.Command, args []string) {
				configPath := filepath.Join(ctx.ConfigDir, "config.toml")
				data, err := os.ReadFile(configPath)
				if err != nil {
					cmd.PrintErrln("Error reading config:", err)
					return
				}
				cmd.Println(string(data))
			},
		},
	)

	return cmd
}

// newProviderCommand creates a new provider command.
func newProviderCommand(ctx *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "provider",
		Short: "Manage providers",
	}

	// List providers
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available providers",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement provider listing
			cmd.Println("Listing providers (not yet implemented)")
		},
	}

	// Show provider details
	showCmd := &cobra.Command{
		Use:   "show [name]",
		Short: "Show provider details",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement provider details
			cmd.Printf("Showing details for provider: %s (not yet implemented)\n", args[0])
		},
	}

	// Register a new provider
	registerCmd := &cobra.Command{
		Use:   "register [path/to/provider.toml]",
		Short: "Register a new provider",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement provider registration
			cmd.Printf("Registering provider from: %s (not yet implemented)\n", args[0])
		},
	}

	cmd.AddCommand(listCmd, showCmd, registerCmd)
	return cmd
}

// newInstanceCommand creates a new instance command.
func newInstanceCommand(ctx *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "Manage instances",
		Aliases: []string{"instances", "i"},
	}

	// List instances
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List instances",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement instance listing
			cmd.Println("Listing instances (not yet implemented)")
		},
	}

	// Create a new instance
	createCmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new instance",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement instance creation
			cmd.Printf("Creating instance: %s (not yet implemented)\n", args[0])
		},
	}

	// Show instance details
	showCmd := &cobra.Command{
		Use:   "show [name]",
		Short: "Show instance details",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement instance details
			cmd.Printf("Showing details for instance: %s (not yet implemented)\n", args[0])
		},
	}

	// SSH into an instance
	sshCmd := &cobra.Command{
		Use:   "ssh [name]",
		Short: "SSH into an instance",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement SSH
			cmd.Printf("SSH into instance: %s (not yet implemented)\n", args[0])
		},
	}

	// Delete an instance
	deleteCmd := &cobra.Command{
		Use:   "delete [name]",
		Short: "Delete an instance",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement instance deletion
			cmd.Printf("Deleting instance: %s (not yet implemented)\n", args[0])
		},
	}

	cmd.AddCommand(listCmd, createCmd, showCmd, sshCmd, deleteCmd)
	return cmd
}

// newNodeCommand creates a new node command.
func newNodeCommand(ctx *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "Manage nodes",
		Aliases: []string{"nodes", "n"},
	}

	// List nodes
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List nodes",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement node listing
			cmd.Println("Listing nodes (not yet implemented)")
		},
	}

	// Show node details
	showCmd := &cobra.Command{
		Use:   "show [name]",
		Short: "Show node details",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement node details
			cmd.Printf("Showing details for node: %s (not yet implemented)\n", args[0])
		},
	}

	// Drain a node
	drainCmd := &cobra.Command{
		Use:   "drain [name]",
		Short: "Drain a node (prepare for maintenance)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement node drain
			cmd.Printf("Draining node: %s (not yet implemented)\n", args[0])
		},
	}

	// Cordon a node
	cordonCmd := &cobra.Command{
		Use:   "cordon [name]",
		Short: "Mark node as unschedulable",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement node cordon
			cmd.Printf("Cordoning node: %s (not yet implemented)\n", args[0])
		},
	}

	// Uncordon a node
	uncordonCmd := &cobra.Command{
		Use:   "uncordon [name]",
		Short: "Mark node as schedulable",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement node uncordon
			cmd.Printf("Uncordoning node: %s (not yet implemented)\n", args[0])
		},
	}

	cmd.AddCommand(listCmd, showCmd, drainCmd, cordonCmd, uncordonCmd)
	return cmd
}
