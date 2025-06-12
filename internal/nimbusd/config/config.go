// Package config handles configuration management for the nimbusd agent.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

// Config represents the nimbusd agent configuration.
type Config struct {
	Agent      AgentConfig      `toml:"agent"`
	Log        LogConfig        `toml:"log"`
	Discovery  DiscoveryConfig  `toml:"discovery"`
	Provider   ProviderConfig   `toml:"provider"`
	Network    NetworkConfig    `toml:"network"`
	API        APIConfig        `toml:"api"`
	Management ManagementConfig `toml:"management"`
	Metrics    MetricsConfig    `toml:"metrics"`
	Monitoring MonitoringConfig `toml:"monitoring"`
	Backup     BackupConfig     `toml:"backup"`
	Updates    UpdateConfig     `toml:"updates"`
	Plugins    PluginConfig     `toml:"plugins"`
	Tags       map[string]string `toml:"tags"`
	Alerts     AlertConfig      `toml:"alerts"`
	TLS        TLSConfig        `toml:"tls"`
	RateLimit  RateLimitConfig  `toml:"rate_limit"`
}

// AgentConfig contains agent-specific configuration.
type AgentConfig struct {
	ID       string `toml:"id"`
	DataDir  string `toml:"data_dir"`
}

// LogConfig contains logging configuration.
type LogConfig struct {
	Level  string `toml:"level"`
	Format string `toml:"format"`
	Output string `toml:"output"`
}

// DiscoveryConfig contains hardware discovery settings.
type DiscoveryConfig struct {
	RefreshInterval int      `toml:"refresh_interval"`
	Paths          []string `toml:"paths"`
}

// ProviderConfig contains provider-specific settings.
type ProviderConfig struct {
	Name           string `toml:"name"`
	PrivateKeyPath string `toml:"private_key_path"`
	PublicKeyPath  string `toml:"public_key_path"`
}

// NetworkConfig contains network-related settings.
type NetworkConfig struct {
	ListenAddress      string `toml:"listen_address"`
	SSHPrivateKeyPath  string `toml:"ssh_private_key_path"`
	AuthorizedKeysPath string `toml:"authorized_keys_path"`
}

// APIConfig contains API server settings.
type APIConfig struct {
	Enabled bool   `toml:"enabled"`
	Address string `toml:"address"`
	Token   string `toml:"token"`
}

// ManagementConfig contains out-of-band management settings.
type ManagementConfig struct {
	Type     string `toml:"type"`
	Address  string `toml:"address"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

// MetricsConfig contains metrics collection settings.
type MetricsConfig struct {
	Enabled bool   `toml:"enabled"`
	Address string `toml:"address"`
}

// MonitoringConfig contains monitoring integration settings.
type MonitoringConfig struct {
	PushGateway string `toml:"push_gateway"`
}

// BackupConfig contains backup settings.
type BackupConfig struct {
	Enabled      bool   `toml:"enabled"`
	Schedule     string `toml:"schedule"`
	RetentionDays int    `toml:"retention_days"`
}

// UpdateConfig contains update settings.
type UpdateConfig struct {
	Enabled       bool   `toml:"enabled"`
	CheckInterval int    `toml:"check_interval"`
	Channel       string `toml:"channel"`
}

// PluginConfig contains plugin settings.
type PluginConfig struct {
	Disabled []string `toml:"disabled"`
}

// AlertConfig contains alerting settings.
type AlertConfig struct {
	Enabled bool `toml:"enabled"`
}

// TLSConfig contains TLS settings.
type TLSConfig struct {
	Enabled  bool   `toml:"enabled"`
	CertFile string `toml:"cert_file"`
	KeyFile  string `toml:"key_file"`
	CAFile   string `toml:"ca_file"`
}

// RateLimitConfig contains rate limiting settings.
type RateLimitConfig struct {
	Enabled          bool `toml:"enabled"`
	RequestsPerSecond int `toml:"requests_per_second"`
	Burst           int `toml:"burst"`
}

// DefaultConfig returns a new Config with default values.
func DefaultConfig() *Config {
	return &Config{
		Agent: AgentConfig{
			DataDir: "/var/lib/nimbus",
		},
		Log: LogConfig{
			Level:  "info",
			Format: "text",
			Output: "stdout",
		},
		Discovery: DiscoveryConfig{
			RefreshInterval: 300,
			Paths: []string{
				"/sys/class/dmi/id",
				"/proc/cpuinfo",
				"/proc/meminfo",
				"/sys/block",
			},
		},
		Provider: ProviderConfig{
			PrivateKeyPath: "/etc/nimbus/keys/private.pem",
			PublicKeyPath:  "/etc/nimbus/keys/public.pem",
		},
		Network: NetworkConfig{
			ListenAddress:      "0.0.0.0:8443",
			SSHPrivateKeyPath:  "/etc/nimbus/ssh/id_ed25519",
			AuthorizedKeysPath: "/etc/nimbus/ssh/authorized_keys",
		},
		API: APIConfig{
			Enabled: true,
			Address: "127.0.0.1:8080",
		},
		Metrics: MetricsConfig{
			Enabled: true,
			Address: "127.0.0.1:9100",
		},
		Updates: UpdateConfig{
			Enabled:       true,
			CheckInterval: 3600,
			Channel:       "stable",
		},
		RateLimit: RateLimitConfig{
			Enabled:          true,
			RequestsPerSecond: 10,
			Burst:           20,
		},
		Tags: make(map[string]string),
	}
}

// Load loads the configuration from a file.
func Load(path string) (*Config, error) {
	// Create default config
	config := DefaultConfig()

	// Read the config file
	_, err := toml.DecodeFile(path, config)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Ensure data directory exists
	if err := os.MkdirAll(config.Agent.DataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Ensure provider keys directory exists
	if err := os.MkdirAll(filepath.Dir(config.Provider.PrivateKeyPath), 0700); err != nil {
		return nil, fmt.Errorf("failed to create keys directory: %w", err)
	}

	// Ensure SSH directory exists
	if err := os.MkdirAll(filepath.Dir(config.Network.SSHPrivateKeyPath), 0700); err != nil {
		return nil, fmt.Errorf("failed to create SSH directory: %w", err)
	}

	return config, nil
}

// Save saves the configuration to a file.
func (c *Config) Save(path string) error {
	// Create the directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Open the file
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	// Encode the config
	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.Agent.ID == "" {
		return fmt.Errorf("agent ID is required")
	}

	if c.Provider.Name == "" {
		return fmt.Errorf("provider name is required")
	}

	if c.Network.ListenAddress == "" {
		return fmt.Errorf("listen address is required")
	}

	if c.API.Enabled && c.API.Address == "" {
		return fmt.Errorf("API address is required when API is enabled")
	}

	if c.Metrics.Enabled && c.Metrics.Address == "" {
		return fmt.Errorf("metrics address is required when metrics are enabled")
	}

	return nil
}
