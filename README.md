# Nimbus Initiative

A decentralized, peer-optional, GitOps-driven bare metal cloud platform.

[![Go Report Card](https://goreportcard.com/badge/github.com/nimbus-project/nimbus)](https://goreportcard.com/report/github.com/nimbus-project/nimbus)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/nimbus-project/nimbus.svg)](https://pkg.go.dev/github.com/nimbus-project/nimbus)

## Overview

Nimbus allows you to build and manage a distributed bare metal cloud with minimal central coordination. It's designed to be simple, secure, and flexible, supporting any hardware architecture.

## Key Components

- `nimbusd`: The agent that runs on provider machines
- `nimbusctl`: CLI for managing instances and providers
- Provider Registry: Git-based registry of available hardware
- Web UI: GCP-like interface for managing resources

## Quick Start

### Prerequisites

- Go 1.18 or later
- Git
- Make
- Standard build tools (gcc, etc.)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/nimbus-project/nimbus.git
cd nimbus

# Build the project
make

# Install system-wide (requires root)
sudo make install

# Or install using the install script
sudo ./scripts/install.sh
```

### Development Setup

For setting up a development environment, use the development setup script:

```bash
# Set up development environment
./scripts/dev-setup.sh

# Start the application in development mode
make dev
```

For more information about available scripts, see the [scripts documentation](scripts/README.md).

## Getting Started

### For Providers

1. Install `nimbusd` on your machine
2. Create a provider TOML file (see `examples/provider.toml`)
3. Submit a PR to the provider registry
4. Once merged, your hardware will be available for provisioning

### For Users

1. Install `nimbusctl`
2. Configure your SSH keys
3. Browse available providers with `nimbusctl providers list`
4. Launch instances with `nimbusctl instances create`

## Architecture

Nimbus uses a peer-to-peer architecture where:
- Providers run `nimbusd` to manage their hardware
- Discovery happens through a Git-based registry
- Provisioning is done via secure SSH connections
- No central coordination is required (but can be added for additional features)

## Security

- All communications are encrypted
- Provider authentication via SSH keys
- Fine-grained access control
- Regular security audits

## Documentation

- [Architecture](docs/architecture.md) - High-level architecture and design decisions
- [CLI Usage](docs/cli-usage.md) - Comprehensive guide to using `nimbusctl`
- [Provider Guide](examples/provider.toml) - How to set up a provider
- [Development Guide](CONTRIBUTING.md) - How to contribute to the project
- [Scripts](scripts/README.md) - Utility scripts for building and deployment

## Security

- All communications are encrypted
- Provider authentication via SSH keys
- Fine-grained access control
- Regular security audits

## Contributing

Contributions are welcome! Please see our [Contributing Guide](CONTRIBUTING.md) for details on how to get started.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
