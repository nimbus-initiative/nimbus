# Nimbus CLI (nimbusctl) Usage

## Installation

```bash
# Install from source (requires Go 1.18+)
git clone https://github.com/nimbus-project/nimbus
cd nimbus
make install
```

## Basic Commands

### Provider Management

```bash
# List all available providers
nimbusctl providers list

# Show details about a specific provider
nimbusctl providers get example-llc

# Register a new provider
nimbusctl providers create -f provider.toml

# Update a provider
nimbusctl providers update example-llc -f updated-provider.toml

# Remove a provider
nimbusctl providers delete example-llc
```

### Instance Management

```bash
# List all instances
nimbusctl instances list

# Create a new instance
nimbusctl instances create \
  --name my-instance \
  --provider example-llc \
  --machine-type r6i.metal \
  --image ubuntu-22.04 \
  --ssh-key ~/.ssh/id_rsa.pub

# Show instance details
nimbusctl instances describe my-instance

# Access instance via SSH
nimbusctl instances ssh my-instance

# Reboot an instance
nimbusctl instances reboot my-instance

# Delete an instance
nimbusctl instances delete my-instance
```

### Node Management

```bash
# List all nodes in the cluster
nimbusctl nodes list

# Show node details
nimbusctl nodes describe node-01

# Drain a node (prepare for maintenance)
nimbusctl nodes drain node-01

# Cordon/uncordon a node
nimbusctl nodes cordon node-01
nimbusctl nodes uncordon node-01
```

### Configuration

```bash
# Show current configuration
nimbusctl config view

# Set configuration values
nimbusctl config set default_provider example-llc
nimbusctl config set default_region us-east

# Generate shell completion
nimbusctl completion bash > /etc/bash_completion.d/nimbusctl
```

### Daemon Management

```bash
# Install and start nimbusd
sudo nimbusctl daemon install
sudo systemctl start nimbusd

# View daemon logs
journalctl -u nimbusd -f

# Check daemon status
nimbusctl daemon status
```

## Advanced Usage

### Using Configuration Files

Create a `~/.nimbus/config.toml` file:

```toml
default_provider = "example-llc"
default_region = "us-east"

[providers.example-llc]
  type = "ssh"
  address = "provider.example.com"
  port = 22
  user = "nimbus"
  key_path = "~/.ssh/id_rsa"
```

### Environment Variables

```bash
export NIMBUS_PROVIDER=example-llc
export NIMBUS_REGION=us-east
export NIMBUS_DEBUG=true
```

### Output Formats

```bash
# JSON output
nimbusctl instances list -o json

# YAML output
nimbusctl providers get example-llc -o yaml

# Custom columns
nimbusctl instances list -o custom-columns=NAME:.metadata.name,STATUS:.status.phase,IP:.status.addresses[0].address
```
