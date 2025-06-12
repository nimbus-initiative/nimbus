package providers

import (
	"fmt"
	"strings"
)

// Config represents the root configuration for all providers
type Config struct {
	Metal *MetalConfig `toml:"nimbus::metal"`
}

// MetalConfig contains metal-specific configurations
type MetalConfig struct {
	Providers map[string]ProviderConfig `toml:"providers"`
}

// ProviderConfig represents a specific provider configuration
type ProviderConfig struct {
	// Common fields
	Region      string            `toml:"region"`
	HostModel   string            `toml:"host_model"`
	Type        string            `toml:"type"`
	CPU         string            `toml:"cpu"`
	CPUSockets  int               `toml:"cpu_sockets"`
	CPUCores    int               `toml:"cpu_cores"`
	CPUThreads  int               `toml:"cpu_threads"`
	RAM         string            `toml:"ram"`
	GPU         *string           `toml:"gpu,omitempty"`
	FPGA        *string           `toml:"fpga,omitempty"`
	Hypervisor  *string           `toml:"hypervisor,omitempty"`
	BMC         bool              `toml:"bmc"`
	
	// Optional PXE configuration
	PXE         *PXEConfig        `toml:"pxe,omitempty"`
	
	// Optional BMC configuration (when BMC is true)
	BMCConfig   *BMCConfig        `toml:"bmc_config,omitempty"`
	
	// Custom metadata
	Metadata    map[string]string `toml:"metadata,omitempty"`
}

// PXEConfig holds PXE boot configuration
type PXEConfig struct {
	Enabled     bool   `toml:"enabled"`
	KernelURL   string `toml:"kernel_url,omitempty"`
	InitrdURL   string `toml:"initrd_url,omitempty"`
	Cmdline     string `toml:"cmdline,omitempty"`
}

// BMCConfig holds Baseboard Management Controller configuration
type BMCConfig struct {
	Address      string `toml:"address"`
	Protocol     string `toml:"protocol"` // redfish or ipmi
	Username     string `toml:"username"`
	Password     string `toml:"password"`
	Insecure     bool   `toml:"insecure,omitempty"`
}

// ParseProviderKey parses a namespaced provider key (e.g., "aws::r6i.metal")
// into its components (e.g., "aws" and "r6i.metal")
func ParseProviderKey(key string) (providerType, providerModel string, err error) {
	parts := strings.SplitN(key, "::", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid provider key format: %s", key)
	}
	return parts[0], parts[1], nil
}

// NewConfig creates a new empty configuration
func NewConfig() *Config {
	return &Config{
		Metal: &MetalConfig{
			Providers: make(map[string]ProviderConfig),
		},
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Metal == nil {
		return fmt.Errorf("metal configuration is required")
	}

	for key, provider := range c.Metal.Providers {
		if provider.Type == "" {
			return fmt.Errorf("provider %s: type is required", key)
		}

		// Validate BMC configuration if BMC is enabled
		if provider.BMC && provider.BMCConfig == nil {
			return fmt.Errorf("provider %s: bmc_config is required when bmc is true", key)
		}

		// Validate PXE configuration if provided
		if provider.PXE != nil && provider.PXE.Enabled {
			if provider.PXE.KernelURL == "" {
				return fmt.Errorf("provider %s: pxe.kernel_url is required when pxe is enabled", key)
			}
			if provider.PXE.InitrdURL == "" {
				return fmt.Errorf("provider %s: pxe.initrd_url is required when pxe is enabled", key)
			}
		}
	}

	return nil
}