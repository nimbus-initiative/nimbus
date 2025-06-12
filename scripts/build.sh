#!/bin/bash
set -euo pipefail

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
BUILD_TARGETS=("nimbusd" "nimbusctl")
BUILD_OS=$(uname -s | tr '[:upper:]' '[:lower:]')
BUILD_ARCH=$(uname -m)
BUILD_DIR="$(pwd)/bin"
STATIC_BUILD=false
CLEAN_BUILD=false

# Function to print usage
usage() {
    echo "Usage: $0 [options] [targets]"
    echo "Options:"
    echo "  --os OS          Target OS (default: ${BUILD_OS})"
    echo "  --arch ARCH      Target architecture (default: ${BUILD_ARCH})"
    echo "  -o, --output DIR Output directory (default: ${BUILD_DIR})"
    echo "  --static         Build static binaries"
    echo "  --clean          Clean build directory before building"
    echo "  -h, --help       Show this help message"
    echo ""
    echo "Available targets: nimbusd nimbusctl all"
    echo "If no targets are specified, all targets will be built."
    exit 1
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case "$1" in
        --os)
            BUILD_OS="$2"
            shift 2
            ;;
        --arch)
            BUILD_ARCH="$2"
            shift 2
            ;;
        -o|--output)
            BUILD_DIR="$2"
            shift 2
            ;;
        --static)
            STATIC_BUILD=true
            shift
            ;;
        --clean)
            CLEAN_BUILD=true
            shift
            ;;
        -h|--help)
            usage
            ;;
        all)
            BUILD_TARGETS=("nimbusd" "nimbusctl")
            shift
            ;;
        nimbusd|nimbusctl)
            # Remove from array if it exists
            BUILD_TARGETS=("${BUILD_TARGETS[@]/$1/}")
            # Add to the end
            BUILD_TARGETS+=("$1")
            shift
            ;;
        *)
            echo -e "${YELLOW}Unknown target or option: $1${NC}" >&2
            usage
            ;;
    esac
done

# Ensure BUILD_TARGETS is not empty
if [ ${#BUILD_TARGETS[@]} -eq 0 ]; then
    BUILD_TARGETS=("nimbusd" "nimbusctl")
fi

# Clean build directory if requested
if [ "$CLEAN_BUILD" = true ] && [ -d "$BUILD_DIR" ]; then
    echo -e "${YELLOW}Cleaning build directory: ${BUILD_DIR}${NC}"
    rm -rf "${BUILD_DIR}"
fi

# Create build directory
mkdir -p "$BUILD_DIR"

# Set build environment
GOOS="$BUILD_OS"
GOARCH="$BUILD_ARCH"
CGO_ENABLED=0

if [ "$STATIC_BUILD" = true ]; then
    echo -e "${GREEN}Building static binaries for ${GOOS}/${GOARCH}...${NC}"
    export CGO_ENABLED=0
    LDFLAGS="-s -w"
else
    echo -e "${GREEN}Building dynamic binaries for ${GOOS}/${GOARCH}...${NC}"
    LDFLAGS=""
fi

# Build each target
for target in "${BUILD_TARGETS[@]}"; do
    if [[ "$target" == "nimbusd" || "$target" == "nimbusctl" ]]; then
        echo -e "\n${GREEN}Building ${target}...${NC}"
        
        # Set output path
        OUTPUT="${BUILD_DIR}/${target}"
        if [ "$GOOS" = "windows" ]; then
            OUTPUT="${OUTPUT}.exe"
        fi
        
        # Build command
        CMD="GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=${CGO_ENABLED} \
            go build -ldflags=\"${LDFLAGS}\" -o \"${OUTPUT}\" ./cmd/${target}"
        
        echo -e "${YELLOW}${CMD}${NC}"
        eval "${CMD}"
        
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}Successfully built ${target} -> ${OUTPUT}${NC}"
            # Make binary executable
            chmod +x "${OUTPUT}"
        else
            echo -e "${RED}Failed to build ${target}${NC}" >&2
            exit 1
        fi
    fi
done

echo -e "\n${GREEN}Build complete!${NC}"
echo -e "Binaries are available in: ${BUILD_DIR}"
