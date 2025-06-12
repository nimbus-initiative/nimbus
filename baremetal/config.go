// Package baremetal provides functionality for provisioning and managing bare metal servers.
package baremetal

import (
	"context"
	"fmt"
	"net"
	"time"
)

// Config holds the configuration for bare metal provisioning
type Config struct {
	// Network configuration
	Network NetworkConfig `toml:"network"`

	// PXE boot configuration
	PXE PXEConfig `toml:"pxe"`

	// IPMI/Redfish configuration
	BMC BMCConfig `toml:"bmc"`

	// OS installation configuration
	OS OSConfig `toml:"os"`

	// Post-installation configuration
	PostInstall PostInstallConfig `toml:"post_install"`

	// Timeout for provisioning operations
	Timeout Duration `toml:"timeout"`
}

// NetworkConfig holds network configuration for bare metal servers
type NetworkConfig struct {
	// Network interface to use for PXE boot
	Interface string `toml:"interface"`

	// IP address range for DHCP
	DHCPRange struct {
		Start string `toml:"start"`
		End   string `toml:"end"`
	} `toml:"dhcp_range"`

	// Network configuration for provisioned servers
	Network  string   `toml:"network"`
	Netmask  string   `toml:"netmask"`
	Gateway  string   `toml:"gateway"`
	DNSServers []string `toml:"dns_servers"`
	NTP      string   `toml:"ntp"`
}

// PXEConfig holds PXE boot configuration
type PXEConfig struct {
	// Enable PXE boot server
	Enabled bool `toml:"enabled"`

	// Path to kernel and initrd
	Kernel  string `toml:"kernel"`
	Initrd  string `toml:"initrd"`
	Cmdline string `toml:"cmdline"`

	// HTTP server configuration
	HTTPAddr string `toml:"http_addr"`
	TFTPAddr string `toml:"tftp_addr"`
	DHCPAddr string `toml:"dhcp_addr"`

	// Root directory for PXE files
	RootDir string `toml:"root_dir"`
}

// BMCConfig holds BMC (Baseboard Management Controller) configuration
type BMCConfig struct {
	// Protocol to use (ipmi or redfish)
	Protocol string `toml:"protocol"`

	// Default credentials (can be overridden per host)
	Username string `toml:"username"`
	Password string `toml:"password"`

	// Skip TLS verification (not recommended for production)
	InsecureSkipVerify bool `toml:"insecure_skip_verify"`
}

// OSConfig holds operating system installation configuration
type OSConfig struct {
	// OS type (e.g., "linux", "windows")
	Type string `toml:"type"`

	// OS version
	Version string `toml:"version"`

	// Installation source (URL to ISO or repository)
	Source string `toml:"source"`

	// Root password (hashed)
	RootPassword string `toml:"root_password"`

	// SSH public keys for root access
	SSHKeys []string `toml:"ssh_keys"`

	// Disk configuration
	Disk struct {
		// Device to install to (e.g., "/dev/sda")
		Device string `toml:"device"`

		// Filesystem type (e.g., "ext4", "xfs")
		Filesystem string `toml:"filesystem"`

		// Whether to use LVM
		UseLVM bool `toml:"use_lvm"`

		// Partition scheme (e.g., "msdos", "gpt")
		PartitionScheme string `toml:"partition_scheme"`

		// Custom partition layout (if empty, use default)
		Partitions []Partition `toml:"partitions"`
	} `toml:"disk"`

	// Network configuration
	Network struct {
		// Hostname
		Hostname string `toml:"hostname"`

		// Network interfaces
		Interfaces []NetworkInterface `toml:"interfaces"`

		// DNS configuration
		Nameservers []string `toml:"nameservers"`
		SearchDomains []string `toml:"search_domains"`
	} `toml:"network"`

	// Packages to install
	Packages []string `toml:"packages"`

	// Custom scripts to run during installation
	PreInstallScripts  []string `toml:"pre_install_scripts"`
	PostInstallScripts []string `toml:"post_install_scripts"`
}

// PostInstallConfig holds post-installation configuration
type PostInstallConfig struct {
	// Whether to enable SSH access
	EnableSSH bool `toml:"enable_ssh"`

	// Whether to enable password authentication
	PasswordAuthentication bool `toml:"password_authentication"`

	// Whether to enable root login
	PermitRootLogin bool `toml:"permit_root_login"`

	// Timezone to set
	Timezone string `toml:"timezone"`

	// Locale to set
	Locale string `toml:"locale"`

	// Custom commands to run after installation
	Commands []string `toml:"commands"`
}

// NetworkInterface represents a network interface configuration
type NetworkInterface struct {
	// Interface name (e.g., "eth0")
	Name string `toml:"name"`

	// IP address (empty for DHCP)
	Address string `toml:"address"`

	// Netmask (e.g., "255.255.255.0")
	Netmask string `toml:"netmask"`

	// Gateway (empty for default)
	Gateway string `toml:"gateway"`

	// Whether to use DHCP
	DHCP bool `toml:"dhcp"`

	// Whether to bring the interface up on boot
	OnBoot bool `toml:"on_boot"`

	// Whether to use this interface as the default route
	DefaultRoute bool `toml:"default_route"`
}

// Partition represents a disk partition
type Partition struct {
	// Mount point (e.g., "/", "/boot", "swap")
	MountPoint string `toml:"mount_point"`

	// Size in MB (0 for remaining space)
	SizeMB int64 `toml:"size_mb"`

	// Filesystem type (e.g., "ext4", "swap")
	Filesystem string `toml:"filesystem"`

	// Whether this is a boot partition
	Bootable bool `toml:"bootable"`
}

// Duration is a wrapper around time.Duration for TOML unmarshaling
type Duration time.Duration

// UnmarshalText implements the encoding.TextUnmarshaler interface
func (d *Duration) UnmarshalText(text []byte) error {
	duration, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}
	*d = Duration(duration)
	return nil
}

// Host represents a bare metal host
type Host struct {
	// Hostname
	Hostname string `toml:"hostname"`

	// MAC address for PXE boot
	MAC string `toml:"mac"`

	// BMC configuration
	BMC struct {
		// IP address or hostname of the BMC
		Address string `toml:"address"`

		// Protocol (ipmi or redfish)
		Protocol string `toml:"protocol"`

		// Credentials (if different from default)
		Username string `toml:"username"`
		Password string `toml:"password"`
	} `toml:"bmc"`

	// Hardware information
	Hardware struct {
		// CPU information
		CPU struct {
			Vendor  string `toml:"vendor"`
			Model   string `toml:"model"`
			Cores   int    `toml:"cores"`
			Threads int    `toml:"threads"`
		} `toml:"cpu"`

		// Memory in MB
		Memory int64 `toml:"memory"`

		// Disks
		Disks []struct {
			Device string `toml:"device"`
			SizeGB int64  `toml:"size_gb"`
			Model  string `toml:"model"`
		} `toml:"disks"`

		// Network interfaces
		NICs []struct {
			Name       string `toml:"name"`
			MAC        string `toml:"mac"`
			SpeedMbps  int    `toml:"speed_mbps"`
			DuplexFull bool   `toml:"duplex_full"`
		} `toml:"nics"`
	} `toml:"hardware"`

	// Custom configuration for this host
	Config map[string]interface{} `toml:"config"`
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate network configuration
	if c.Network.Interface == "" {
		return fmt.Errorf("network interface is required")
	}

	// Validate PXE configuration if enabled
	if c.PXE.Enabled {
		if c.PXE.Kernel == "" {
			return fmt.Errorf("PXE kernel path is required")
		}
		if c.PXE.Initrd == "" {
			return fmt.Errorf("PXE initrd path is required")
		}
	}

	// Validate BMC configuration
	switch c.BMC.Protocol {
	case "", "ipmi", "redfish":
		// Valid protocols
	default:
		return fmt.Errorf("unsupported BMC protocol: %s", c.BMC.Protocol)
	}

	return nil
}

// Provisioner handles the provisioning of bare metal servers
type Provisioner struct {
	config *Config
}

// NewProvisioner creates a new bare metal provisioner
func NewProvisioner(cfg *Config) (*Provisioner, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &Provisioner{
		config: cfg,
	}, nil
}

// Provision provisions a bare metal server
func (p *Provisioner) Provision(ctx context.Context, host *Host) error {
	// Set default timeout if not specified
	if p.config.Timeout == 0 {
		p.config.Timeout = Duration(30 * time.Minute)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, time.Duration(p.config.Timeout))
	defer cancel()

	// Step 1: Power off the host if it's running
	if err := p.powerOffHost(ctx, host); err != nil {
		return fmt.Errorf("failed to power off host: %w", err)
	}

	// Step 2: Configure PXE boot
	if p.config.PXE.Enabled {
		if err := p.configurePXEBoot(ctx, host); err != nil {
			return fmt.Errorf("failed to configure PXE boot: %w", err)
		}
	}

	// Step 3: Power on the host
	if err := p.powerOnHost(ctx, host); err != nil {
		return fmt.Errorf("failed to power on host: %w", err)
	}

	// Step 4: Monitor installation progress
	if err := p.monitorInstallation(ctx, host); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	// Step 5: Configure post-installation settings
	if err := p.configurePostInstall(ctx, host); err != nil {
		return fmt.Errorf("post-installation configuration failed: %w", err)
	}

	return nil
}

// powerOffHost powers off a host using IPMI/Redfish
func (p *Provisioner) powerOffHost(ctx context.Context, host *Host) error {
	// TODO: Implement power off using IPMI/Redfish
	return nil
}

// configurePXEBoot configures PXE boot for a host
func (p *Provisioner) configurePXEBoot(ctx context.Context, host *Host) error {
	// TODO: Implement PXE boot configuration
	return nil
}

// powerOnHost powers on a host using IPMI/Redfish
func (p *Provisioner) powerOnHost(ctx context.Context, host *Host) error {
	// TODO: Implement power on using IPMI/Redfish
	return nil
}

// monitorInstallation monitors the installation progress
func (p *Provisioner) monitorInstallation(ctx context.Context, host *Host) error {
	// TODO: Implement installation monitoring
	return nil
}

// configurePostInstall configures post-installation settings
func (p *Provisioner) configurePostInstall(ctx context.Context, host *Host) error {
	// TODO: Implement post-installation configuration
	return nil
}
