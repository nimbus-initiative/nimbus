# Bare Metal Provisioning

This package provides functionality for provisioning and managing bare metal servers in the Nimbus platform. It supports PXE boot, IPMI/Redfish out-of-band management, and OS installation.

## Features

- **PXE Boot Server**: Built-in PXE server for network booting
- **IPMI/Redfish Support**: Out-of-band management for power control and boot device selection
- **OS Installation**: Automated OS installation with customizable partitioning and configuration
- **Post-Installation Configuration**: Custom scripts and configuration for post-install setup
- **Hardware Discovery**: Automatic hardware detection and inventory

## Configuration

Bare metal provisioning is configured using TOML configuration files. See the [example configuration](../examples/baremetal.toml) for details.

### Network Configuration

```toml
[network]
interface = "eth0"  # Network interface for PXE boot
network = "192.168.1.0/24"
netmask = "255.255.255.0"
gateway = "192.168.1.1"
dns_servers = ["8.8.8.8", "8.8.4.4"]
ntp = "pool.ntp.org"

[network.dhcp_range]
start = "192.168.1.100"
end = "192.168.1.200"
```

### PXE Boot Configuration

```toml
[pxe]
enabled = true
kernel = "/var/lib/tftpboot/pxelinux/vmlinuz"
initrd = "/var/lib/tftpboot/pxelinux/initrd.img"
cmdline = "console=tty0 console=ttyS0,115200n8"
http_addr = ":8080"
tftp_addr = ":69"
dhcp_addr = ":67"
root_dir = "/srv/pxeboot"
```

### BMC Configuration

```toml
[bmc]
protocol = "ipmi"  # or "redfish"
username = "admin"
password = "changeme"
insecure_skip_verify = true  # Only for testing
```

### OS Installation Configuration

```toml
[os]
type = "linux"
version = "ubuntu-20.04"
source = "http://archive.ubuntu.com/ubuntu/dists/focal/main/installer-amd64/"
root_password = "$6$hashedpassword"  # Use mkpasswd -m sha-512

[[os.ssh_keys]]
key = "ssh-rsa AAAAB3NzaC1yc2E... user@example.com"

[os.disk]
device = "/dev/sda"
filesystem = "ext4"
use_lvm = true
partition_scheme = "gpt"

[[os.disk.partitions]]
mount_point = "/boot/efi"
size_mb = 512
filesystem = "vfat"

[[os.disk.partitions]]
mount_point = "/boot"
size_mb = 1024
filesystem = "ext4"

[[os.disk.partitions]]
mount_point = "/"
size_mb = 0  # Use remaining space
filesystem = "ext4"
```

## Usage

### Initializing a Provisioner

```go
import (
	"github.com/nimbus-project/nimbus/baremetal"
	"github.com/BurntSushi/toml"
)

// Load configuration
var config baremetal.Config
if _, err := toml.DecodeFile("config.toml", &config); err != nil {
	log.Fatalf("Failed to load configuration: %v", err)
}

// Create a new provisioner
provisioner, err := baremetal.NewProvisioner(&config)
if err != nil {
	log.Fatalf("Failed to create provisioner: %v", err)
}
```

### Provisioning a Host

```go
// Define a host to provision
host := &baremetal.Host{
	Hostname: "nimbus-node-01",
	MAC:      "00:11:22:33:44:55",
}

host.BMC.Address = "192.168.1.50"
host.BMC.Protocol = "ipmi"
host.BMC.Username = "admin"
host.BMC.Password = "changeme"

// Provision the host
if err := provisioner.Provision(context.Background(), host); err != nil {
	log.Fatalf("Failed to provision host: %v", err)
}
```

### Monitoring Provisioning Status

The provisioner provides callbacks for monitoring the provisioning process:

```go
// Set up status callbacks
provisioner.OnStatusUpdate(func(host *baremetal.Host, status string, message string) {
	log.Printf("[%s] %s: %s", host.Hostname, status, message)
})

// Start provisioning
if err := provisioner.Provision(context.Background(), host); err != nil {
	log.Fatalf("Failed to provision host: %v", err)
}
```

## Security Considerations

- Always use secure passwords for BMC/IPMI/Redfish access
- Enable secure boot when supported by the hardware
- Use TLS for all network communications
- Restrict network access to the provisioning network
- Rotate credentials after provisioning
- Keep firmware and software up to date

## Dependencies

- Go 1.18 or later
- `github.com/BurntSushi/toml` for configuration parsing
- `github.com/stmcginnis/gofish` for Redfish support
- `github.com/bmc-toolhub/bmc` for IPMI support

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
