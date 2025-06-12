# Nimbus Scripts

This directory contains various utility scripts for building, testing, and managing the Nimbus project.

## Available Scripts

### Build and Installation

- `build.sh` - Build Nimbus binaries for different platforms
  - Usage: `./scripts/build.sh [options] [targets]`
  - Options:
    - `--os OS` - Target OS (default: current OS)
    - `--arch ARCH` - Target architecture (default: current arch)
    - `-o, --output DIR` - Output directory (default: ./bin)
    - `--static` - Build static binaries
    - `--clean` - Clean build directory before building

- `install.sh` - Install Nimbus system-wide
  - Usage: `sudo ./scripts/install.sh [options]`
  - Options:
    - `--prefix DIR` - Installation prefix (default: /usr/local)
    - `--config DIR` - Configuration directory (default: /etc/nimbus)
    - `--data DIR` - Data directory (default: /var/lib/nimbus)

- `uninstall.sh` - Uninstall Nimbus
  - Usage: `sudo ./scripts/uninstall.sh [options]`
  - Options:
    - `--prefix DIR` - Installation prefix (default: /usr/local)
    - `--config DIR` - Configuration directory (default: /etc/nimbus)
    - `--data DIR` - Data directory (default: /var/lib/nimbus)
    - `--all` - Remove all files including config and data

### Development

- `dev-setup.sh` - Set up a development environment
  - Installs required tools and dependencies
  - Sets up Git hooks
  - Builds the project
  - Usage: `./scripts/dev-setup.sh`

- `release.sh` - Create a new release
  - Builds binaries for multiple platforms
  - Creates checksums
  - Signs the release (optional)
  - Usage: `./scripts/release.sh -v VERSION [options]`
  - Options:
    - `-v, --version VERSION` - Version number (required)
    - `-o, --output DIR` - Output directory (default: ./dist)
    - `--no-clean` - Don't clean build directory before building
    - `--dry-run` - Show what would be done
    - `--sign` - Sign the release with GPG
    - `--gpg-key KEY_ID` - GPG key ID to use for signing
    - `--targets TARGETS` - Comma-separated list of OS/ARCH to build

## Usage Examples

### Build for current platform
```bash
./scripts/build.sh
```

### Build for Linux and macOS
```bash
./scripts/build.sh --os linux --arch amd64 --os darwin --arch amd64
```

### Install system-wide
```bash
sudo ./scripts/install.sh
```

### Create a release
```bash
./scripts/release.sh -v 1.0.0 --sign
```

### Set up development environment
```bash
./scripts/dev-setup.sh
```

## Requirements

- Go 1.18 or later
- Git
- Make
- Standard build tools (gcc, etc.)

For release builds:
- GPG (for signing)
- Docker (for cross-platform builds)

## License

This project is licensed under the [MIT License](../LICENSE).
