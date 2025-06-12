#!/bin/bash
set -euo pipefail

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to print section header
print_section() {
    echo -e "\n${GREEN}=== $1 ===${NC}"
}

# Function to print status
print_status() {
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓${NC} $1"
    else
        echo -e "${RED}✗${NC} $1"
        exit 1
    fi
}

# Check if running as root
if [ "$(id -u)" -eq 0 ]; then
    echo -e "${YELLOW}Warning: Running as root is not recommended. Please run as a regular user.${NC}"
    exit 1
fi

print_section "Nimbus Development Environment Setup"

# Check for required tools
print_section "Checking Dependencies"

# Check for Go
if command_exists go; then
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    echo -e "${GREEN}✓${NC} Go ${GO_VERSION} is installed"
    
    # Check Go version
    IFS='.' read -r -a VERSION_PARTS <<< "$GO_VERSION"
    if [ "${VERSION_PARTS[0]}" -lt 1 ] || { [ "${VERSION_PARTS[0]}" -eq 1 ] && [ "${VERSION_PARTS[1]}" -lt 18 ]; }; then
        echo -e "${RED}Error: Go 1.18 or later is required. Found ${GO_VERSION}.${NC}"
        exit 1
    fi
else
    echo -e "${RED}Error: Go is not installed. Please install Go 1.18 or later.${NC}"
    echo "  Download from: https://golang.org/dl/"
    exit 1
fi

# Check for Git
if command_exists git; then
    echo -e "${GREEN}✓${NC} Git is installed"
else
    echo -e "${RED}Error: Git is not installed. Please install Git.${NC}"
    echo "  On Ubuntu/Debian: sudo apt-get install git"
    echo "  On macOS: brew install git"
    exit 1
fi

# Check for Make
if command_exists make; then
    echo -e "${GREEN}✓${NC} Make is installed"
else
    echo -e "${RED}Error: Make is not installed. Please install Make.${NC}"
    echo "  On Ubuntu/Debian: sudo apt-get install build-essential"
    echo "  On macOS: xcode-select --install"
    exit 1
fi

# Install development tools
print_section "Installing Development Tools"

# Install golangci-lint
if ! command_exists golangci-lint; then
    echo -e "${YELLOW}Installing golangci-lint...${NC}"
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.1
    print_status "Installed golangci-lint"
else
    echo -e "${GREEN}✓${NC} golangci-lint is already installed"
fi

# Install goreleaser for releases
if ! command_exists goreleaser; then
    echo -e "${YELLOW}Installing goreleaser...${NC}"
    go install github.com/goreleaser/goreleaser@latest
    print_status "Installed goreleaser"
else
    echo -e "${GREEN}✓${NC} goreleaser is already installed"
fi

# Install mockgen for generating mocks
if ! command_exists mockgen; then
    echo -e "${YELLOW}Installing mockgen...${NC}"
    go install go.uber.org/mock/mockgen@latest
    print_status "Installed mockgen"
else
    echo -e "${GREEN}✓${NC} mockgen is already installed"
fi

# Install air for live reloading
if ! command_exists air; then
    echo -e "${YELLOW}Installing air (live reloading)...${NC}"
    curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
    print_status "Installed air"
else
    echo -e "${GREEN}✓${NC} air is already installed"
fi

# Setup Git hooks
print_section "Setting up Git Hooks"

# Create .git/hooks directory if it doesn't exist
mkdir -p .git/hooks

# Create pre-commit hook
cat > .git/hooks/pre-commit << 'EOL'
#!/bin/sh

# Run linters
echo "Running linters..."
if ! make lint; then
    echo "Linting failed. Please fix the issues and try again."
    exit 1
fi

# Run tests
echo "Running tests..."
if ! make test; then
    echo "Tests failed. Please fix the issues and try again."
    exit 1
fi

echo "All checks passed!"
EOL

chmod +x .git/hooks/pre-commit
print_status "Installed pre-commit hook"

# Install project dependencies
print_section "Installing Project Dependencies"

echo -e "${YELLOW}Downloading Go modules...${NC}"
go mod download
print_status "Downloaded Go modules"

# Build the project
print_section "Building Nimbus"

if make; then
    print_status "Successfully built Nimbus"
    echo -e "\n${GREEN}✓ Development environment setup complete!${NC}"
    echo -e "\nNext steps:"
    echo "1. Run 'make dev' to start the application in development mode"
    echo "2. Run 'make test' to run the test suite"
    echo "3. Run 'make lint' to run the linter"
    echo "4. Run 'make build' to build the application"
else
    echo -e "\n${RED}Failed to build Nimbus. Please check the error messages above.${NC}"
    exit 1
fi
