#!/bin/bash
set -euo pipefail

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Default values
VERSION=""
BUILD_DIR="$(pwd)/dist"
CLEAN_BUILD=true
DRY_RUN=false
SIGN_RELEASE=false
GPG_KEY=""
TARGETS=("linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64" "windows/amd64")

# Function to print usage
usage() {
    echo "Usage: $0 -v VERSION [options]"
    echo "Options:"
    echo "  -v, --version VERSION  Version number (required)"
    echo "  -o, --output DIR      Output directory (default: ${BUILD_DIR})"
    echo "  --no-clean            Don't clean build directory before building"
    echo "  --dry-run             Show what would be done without actually doing it"
    echo "  --sign                Sign the release with GPG"
    echo "  --gpg-key KEY_ID      GPG key ID to use for signing (default: default key)"
    echo "  --targets TARGETS     Comma-separated list of OS/ARCH to build (e.g., 'linux/amd64,darwin/amd64')'"
    echo "  -h, --help            Show this help message"
    echo ""
    echo "Example: $0 -v v1.0.0 --sign --targets 'linux/amd64,darwin/amd64'"
    exit 1
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case "$1" in
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -o|--output)
            BUILD_DIR="$2"
            shift 2
            ;;
        --no-clean)
            CLEAN_BUILD=false
            shift
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --sign)
            SIGN_RELEASE=true
            shift
            ;;
        --gpg-key)
            GPG_KEY="$2"
            shift 2
            ;;
        --targets)
            IFS=',' read -r -a TARGETS <<< "$2"
            shift 2
            ;;
        -h|--help)
            usage
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}" >&2
            usage
            ;;
    esac
done

# Validate version
if [ -z "$VERSION" ]; then
    echo -e "${RED}Error: Version is required${NC}" >&2
    usage
fi

# Remove 'v' prefix if present
VERSION=${VERSION#v}

# Ensure version is in the format x.y.z
if ! [[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?$ ]]; then
    echo -e "${RED}Error: Invalid version format. Expected x.y.z or x.y.z-rc1, etc.${NC}" >&2
    exit 1
fi

# Check if git is clean
if ! git diff --quiet && [ "$DRY_RUN" = false ]; then
    echo -e "${RED}Error: Git working directory is not clean. Please commit or stash your changes.${NC}" >&2
    exit 1
fi

# Check if tag already exists
if git rev-parse -q --verify "refs/tags/v${VERSION}" >/dev/null; then
    echo -e "${RED}Error: Tag v${VERSION} already exists.${NC}" >&2
    exit 1
fi

# Check if tag exists on remote
if [ "$DRY_RUN" = false ] && git ls-remote --tags origin "refs/tags/v${VERSION}" | grep -q "refs/tags/v${VERSION}"; then
    echo -e "${RED}Error: Tag v${VERSION} already exists on remote.${NC}" >&2
    exit 1
fi

# Clean build directory if requested
if [ "$CLEAN_BUILD" = true ] && [ -d "$BUILD_DIR" ]; then
    echo -e "${YELLOW}Cleaning build directory: ${BUILD_DIR}${NC}"
    if [ "$DRY_RUN" = false ]; then
        rm -rf "${BUILD_DIR}"
    fi
fi

# Create build directory
if [ "$DRY_RUN" = false ]; then
    mkdir -p "$BUILD_DIR"
fi

# Build for each target
for target in "${TARGETS[@]}"; do
    IFS='/' read -r -a os_arch <<< "$target"
    os="${os_arch[0]}"
    arch="${os_arch[1]}"
    
    echo -e "\n${GREEN}Building for ${os}/${arch}...${NC}"
    
    # Set output directory
    TARGET_DIR="${BUILD_DIR}/nimbus-${VERSION}-${os}-${arch}"
    
    if [ "$DRY_RUN" = false ]; then
        mkdir -p "${TARGET_DIR}"
    fi
    
    # Build nimbusd and nimbusctl
    for bin in nimbusd nimbusctl; do
        echo -e "\n${YELLOW}Building ${bin}...${NC}"
        
        OUTPUT="${TARGET_DIR}/${bin}"
        if [ "$os" = "windows" ]; then
            OUTPUT="${OUTPUT}.exe"
        fi
        
        CMD="GOOS=${os} GOARCH=${arch} CGO_ENABLED=0 \
            go build -ldflags=\"-s -w\" -o \"${OUTPUT}\" ./cmd/${bin}"
        
        echo -e "${YELLOW}${CMD}${NC}"
        
        if [ "$DRY_RUN" = false ]; then
            eval "${CMD}"
            if [ $? -ne 0 ]; then
                echo -e "${RED}Failed to build ${bin} for ${os}/${arch}${NC}" >&2
                exit 1
            fi
            
            # Make binary executable (not needed for Windows)
            if [ "$os" != "windows" ]; then
                chmod +x "${OUTPUT}"
            fi
        fi
    done
    
    # Copy additional files
    if [ "$DRY_RUN" = false ]; then
        echo -e "\n${YELLOW}Copying additional files...${NC}"
        cp README.md LICENSE "${TARGET_DIR}/"
        cp -r configs "${TARGET_DIR}/"
        cp -r examples "${TARGET_DIR}/"
        
        # Copy systemd service file for Linux
        if [ "$os" = "linux" ]; then
            mkdir -p "${TARGET_DIR}/systemd"
            cp contrib/systemd/nimbusd.service "${TARGET_DIR}/systemd/"
        fi
    fi
    
    # Create archive
    echo -e "\n${YELLOW}Creating archive...${NC}"
    ARCHIVE_NAME="nimbus-${VERSION}-${os}-${arch}"
    
    if [ "$os" = "windows" ]; then
        ARCHIVE_FILE="${BUILD_DIR}/${ARCHIVE_NAME}.zip"
        if [ "$DRY_RUN" = false ]; then
            (cd "${BUILD_DIR}" && zip -r "${ARCHIVE_FILE}" "$(basename "${TARGET_DIR}")")
        fi
    else
        ARCHIVE_FILE="${BUILD_DIR}/${ARCHIVE_NAME}.tar.gz"
        if [ "$DRY_RUN" = false ]; then
            (cd "${BUILD_DIR}" && tar czf "${ARCHIVE_FILE}" "$(basename "${TARGET_DIR}")")
        fi
    fi
    
    echo -e "${GREEN}Created archive: ${ARCHIVE_FILE}${NC}"
    
    # Sign the archive if requested
    if [ "$SIGN_RELEASE" = true ] && [ "$DRY_RUN" = false ]; then
        echo -e "\n${YELLOW}Signing archive...${NC}"
        if [ -n "$GPG_KEY" ]; then
            gpg --detach-sign --local-user "$GPG_KEY" --armor "${ARCHIVE_FILE}"
        else
            gpg --detach-sign --armor "${ARCHIVE_FILE}"
        fi
        echo -e "${GREEN}Created signature: ${ARCHIVE_FILE}.asc${NC}"
    fi
    
    # Clean up
    if [ "$DRY_RUN" = false ]; then
        rm -rf "${TARGET_DIR}"
    fi
done

# Create checksums
if [ "$DRY_RUN" = false ]; then
    echo -e "\n${YELLOW}Creating checksums...${NC}"
    (cd "${BUILD_DIR}" && shasum -a 256 *.tar.gz *.zip > "nimbus-${VERSION}-checksums.txt")
    echo -e "${GREEN}Created checksums: ${BUILD_DIR}/nimbus-${VERSION}-checksums.txt${NC}"
    
    # Sign the checksums file if requested
    if [ "$SIGN_RELEASE" = true ]; then
        echo -e "\n${YELLOW}Signing checksums...${NC}"
        if [ -n "$GPG_KEY" ]; then
            gpg --detach-sign --local-user "$GPG_KEY" --armor "${BUILD_DIR}/nimbus-${VERSION}-checksums.txt"
        else
            gpg --detach-sign --armor "${BUILD_DIR}/nimbus-${VERSION}-checksums.txt"
        fi
        echo -e "${GREEN}Created signature: ${BUILD_DIR}/nimbus-${VERSION}-checksums.txt.asc${NC}"
    fi
fi

echo -e "\n${GREEN}Release ${VERSION} build complete!${NC}"
echo -e "Artifacts are available in: ${BUILD_DIR}"

if [ "$DRY_RUN" = true ]; then
    echo -e "\n${YELLOW}This was a dry run. No files were actually created or modified.${NC}"
fi

# Instructions for creating a GitHub release
echo -e "\n${GREEN}To create a GitHub release:${NC}"
echo "1. Create a new release at https://github.com/nimbus-project/nimbus/releases/new"
echo "2. Set the tag to 'v${VERSION}'"
echo "3. Set the release title to 'Nimbus ${VERSION}'"
echo "4. Add release notes"
echo "5. Upload the following files:"
echo "   - nimbus-${VERSION}-checksums.txt"
echo "   - nimbus-${VERSION}-checksums.txt.asc (if signed)"
echo "   - All the release archives (*.tar.gz, *.zip, and their .asc files if signed)"

if [ "$DRY_RUN" = false ] && [ "$SIGN_RELEASE" = true ]; then
    echo -e "\n${YELLOW}Don't forget to publish your GPG key if you haven't already:${NC}"
    echo "gpg --send-keys ${GPG_KEY:-YOUR_KEY_ID}"
    echo -e "${YELLOW}And verify it's available on a key server:${NC}"
    echo "gpg --keyserver hkps://keys.openpgp.org --search-keys ${GPG_KEY:-YOUR_EMAIL}"
fi
