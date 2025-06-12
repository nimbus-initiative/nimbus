#!/bin/bash
set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if running as root
if [ "$(id -u)" -ne 0 ]; then
    echo -e "${RED}Error: This script must be run as root${NC}" >&2
    exit 1
fi

# Default values
INSTALL_PREFIX="/usr/local"
CONFIG_DIR="/etc/nimbus"
DATA_DIR="/var/lib/nimbus"
SERVICE_FILE="/etc/systemd/system/nimbusd.service"
BIN_DIR="${INSTALL_PREFIX}/bin"

# Function to print usage
usage() {
    echo "Usage: $0 [options]"
    echo "Options:"
    echo "  --prefix DIR     Installation prefix (default: ${INSTALL_PREFIX})"
    echo "  --config DIR     Configuration directory (default: ${CONFIG_DIR})"
    echo "  --data DIR       Data directory (default: ${DATA_DIR})"
    echo "  --help           Show this help message"
    exit 1
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case "$1" in
        --prefix)
            INSTALL_PREFIX="$2"
            BIN_DIR="${INSTALL_PREFIX}/bin"
            shift 2
            ;;
        --config)
            CONFIG_DIR="$2"
            shift 2
            ;;
        --data)
            DATA_DIR="$2"
            shift 2
            ;;
        --help)
            usage
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}" >&2
            usage
            ;;
    esac
done

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed. Please install Go 1.18 or later.${NC}" >&2
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
IFS='.' read -r -a VERSION_PARTS <<< "$GO_VERSION"
if [ "${VERSION_PARTS[0]}" -lt 1 ] || { [ "${VERSION_PARTS[0]}" -eq 1 ] && [ "${VERSION_PARTS[1]}" -lt 18 ]; }; then
    echo -e "${RED}Error: Go version 1.18 or later is required. Found version ${GO_VERSION}.${NC}" >&2
    exit 1
fi

# Create necessary directories
echo -e "${GREEN}Creating directories...${NC}"
mkdir -p "${BIN_DIR}"
mkdir -p "${CONFIG_DIR}"
mkdir -p "${DATA_DIR}"
mkdir -p "${CONFIG_DIR}/keys"
mkdir -p "${CONFIG_DIR}/ssh"

# Build nimbusd
echo -e "${GREEN}Building nimbusd...${NC}"
if ! make build-nimbusd; then
    echo -e "${RED}Failed to build nimbusd${NC}" >&2
    exit 1
fi

# Install nimbusd
echo -e "${GREEN}Installing nimbusd...${NC}"
cp "bin/nimbusd" "${BIN_DIR}/nimbusd"
chmod 755 "${BIN_DIR}/nimbusd"

# Install nimbusctl
if [ -f "bin/nimbusctl" ]; then
    echo -e "${GREEN}Installing nimbusctl...${NC}"
    cp "bin/nimbusctl" "${BIN_DIR}/nimbusctl"
    chmod 755 "${BIN_DIR}/nimbusctl"
fi

# Create nimbus user if it doesn't exist
if ! id -u nimbus &> /dev/null; then
    echo -e "${GREEN}Creating nimbus user...${NC}"
    useradd --system --shell /usr/sbin/nologin --home-dir "${DATA_DIR}" --create-home nimbus
fi

# Set ownership and permissions
chown -R nimbus:nimbus "${DATA_DIR}"
chmod 700 "${DATA_DIR}"
chmod 750 "${CONFIG_DIR}"
chmod 700 "${CONFIG_DIR}/keys"
chmod 700 "${CONFIG_DIR}/ssh"

# Install systemd service file
if [ -f "contrib/systemd/nimbusd.service" ]; then
    echo -e "${GREEN}Installing systemd service...${NC}"
    cp "contrib/systemd/nimbusd.service" "${SERVICE_FILE}"
    
    # Update paths in service file
    sed -i "s|/usr/local/bin/nimbusd|${BIN_DIR}/nimbusd|" "${SERVICE_FILE}"
    sed -i "s|/etc/nimbus/nimbusd.toml|${CONFIG_DIR}/nimbusd.toml|" "${SERVICE_FILE}"
    sed -i "s|/var/lib/nimbus|${DATA_DIR}|" "${SERVICE_FILE}"
    
    # Reload systemd
    systemctl daemon-reload
    
    echo -e "${GREEN}Enabling and starting nimbusd service...${NC}"
    systemctl enable nimbusd
    systemctl start nimbusd
fi

# Generate default config if it doesn't exist
if [ ! -f "${CONFIG_DIR}/nimbusd.toml" ]; then
    echo -e "${YELLOW}Generating default configuration...${NC}"
    cat > "${CONFIG_DIR}/nimbusd.toml" << EOF
[agent]
id = "$(uuidgen)"
data_dir = "${DATA_DIR}"

[log]
level = "info"
format = "text"
output = "stdout"

[discovery]
refresh_interval = 300
paths = [
    "/sys/class/dmi/id",
    "/proc/cpuinfo",
    "/proc/meminfo",
    "/sys/block"
]

[provider]
name = ""  # Set your provider name
private_key_path = "${CONFIG_DIR}/keys/private.pem"
public_key_path = "${CONFIG_DIR}/keys/public.pem"

[network]
listen_address = "0.0.0.0:8443"
ssh_private_key_path = "${CONFIG_DIR}/ssh/id_ed25519"
authorized_keys_path = "${CONFIG_DIR}/ssh/authorized_keys"

[api]
enabled = true
address = "127.0.0.1:8080"

[metrics]
enabled = true
address = "127.0.0.1:9100"

[updates]
enabled = true
check_interval = 3600
channel = "stable"

[rate_limit]
enabled = true
requests_per_second = 10
burst = 20

[tags]
# Add any custom tags here
EOF
    
    chmod 600 "${CONFIG_DIR}/nimbusd.toml"
    chown -R nimbus:nimbus "${CONFIG_DIR}"
    
    echo -e "${YELLOW}Default configuration generated at ${CONFIG_DIR}/nimbusd.toml${NC}"
    echo -e "${YELLOW}Please edit the configuration file and set your provider name before starting nimbusd.${NC}"
fi

echo -e "\n${GREEN}Installation complete!${NC}"
echo -e "\nNext steps:"
echo -e "1. Edit the configuration file: ${CONFIG_DIR}/nimbusd.toml"
echo -e "2. Set your provider name and adjust other settings as needed"
echo -e "3. Start the service: systemctl start nimbusd"
echo -e "4. Check the status: systemctl status nimbusd"
echo -e "5. View logs: journalctl -u nimbusd -f"

if [ -f "${BIN_DIR}/nimbusctl" ]; then
    echo -e "\nYou can now use the nimbusctl command to manage your Nimbus installation."
    echo -e "Run 'nimbusctl --help' for more information."
fi
