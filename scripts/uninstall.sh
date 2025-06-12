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
    echo "  --all            Remove all files including config and data"
    echo "  --help           Show this help message"
    exit 1
}

# Parse command line arguments
REMOVE_ALL=false
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
        --all)
            REMOVE_ALL=true
            shift
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

# Stop and disable the service if it's running
if systemctl is-active --quiet nimbusd; then
    echo -e "${YELLOW}Stopping nimbusd service...${NC}"
    systemctl stop nimbusd
fi

if systemctl is-enabled --quiet nimbusd; then
    echo -e "${YELLOW}Disabling nimbusd service...${NC}"
    systemctl disable nimbusd
fi

# Remove systemd service file
if [ -f "${SERVICE_FILE}" ]; then
    echo -e "${YELLOW}Removing systemd service file...${NC}"
    rm -f "${SERVICE_FILE}"
    systemctl daemon-reload
fi

# Remove binaries
if [ -f "${BIN_DIR}/nimbusd" ]; then
    echo -e "${YELLOW}Removing nimbusd binary...${NC}"
    rm -f "${BIN_DIR}/nimbusd"
fi

if [ -f "${BIN_DIR}/nimbusctl" ]; then
    echo -e "${YELLOW}Removing nimbusctl binary...${NC}"
    rm -f "${BIN_DIR}/nimbusctl"
fi

# Remove config and data if requested
if [ "$REMOVE_ALL" = true ]; then
    echo -e "${YELLOW}Removing configuration and data...${NC}"
    
    if [ -d "${CONFIG_DIR}" ]; then
        echo -e "  Removing config directory: ${CONFIG_DIR}"
        rm -rf "${CONFIG_DIR}"
    fi
    
    if [ -d "${DATA_DIR}" ]; then
        echo -e "  Removing data directory: ${DATA_DIR}"
        rm -rf "${DATA_DIR}"
    fi
    
    # Remove nimbus user if it exists and has no other files
    if id -u nimbus &> /dev/null; then
        echo -e "  Removing nimbus user..."
        userdel nimbus 2>/dev/null || true
    fi
else
    echo -e "${YELLOW}Configuration and data preserved:${NC}"
    echo -e "  - Config: ${CONFIG_DIR}"
    echo -e "  - Data: ${DATA_DIR}"
    echo -e "Use '--all' to remove these as well"
fi

echo -e "\n${GREEN}Uninstallation complete!${NC}"

if [ "$REMOVE_ALL" = false ]; then
    echo -e "\n${YELLOW}Note: Configuration and data directories were not removed.${NC}"
    echo -e "To remove them, run this script again with the --all flag."
fi
