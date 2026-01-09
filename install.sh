#!/bin/bash
# ============================================================================
# Antigravity Usage Checker - Install Script for Linux/macOS
# ============================================================================
#
# WHAT THIS SCRIPT DOES:
# 1. Detects your OS (Linux/macOS) and architecture (amd64/arm64)
# 2. Downloads the appropriate binary from GitHub Releases
# 3. Installs to: /usr/local/bin/agusage (may require sudo)
#
# WHAT THIS SCRIPT DOES NOT DO:
# - Does NOT install any dependencies
# - Does NOT modify system configuration
# - Does NOT send any data to external servers
# - Does NOT run anything in the background
#
# You can review this script before running. Source code:
# https://github.com/tungcorn/antigravity-usage-checker
#
# Usage: curl -fsSL https://raw.githubusercontent.com/tungcorn/antigravity-usage-checker/main/install.sh | bash
# Install specific version: curl -fsSL ... | bash -s -- -v 0.5.0
# ============================================================================

# Exit immediately if any command fails
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Parse arguments
VERSION=""
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        *)
            shift
            ;;
    esac
done

echo -e "${GREEN}🚀 Antigravity Usage Checker - Installer${NC}"
echo "=========================================="

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}❌ Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

case $OS in
    darwin)
        echo "📍 Detected: macOS ($ARCH)"
        ;;
    linux)
        echo "📍 Detected: Linux ($ARCH)"
        ;;
    *)
        echo -e "${RED}❌ Unsupported OS: $OS${NC}"
        exit 1
        ;;
esac

# Get version
if [ -z "$VERSION" ]; then
    echo "📥 Fetching latest release..."
    INSTALL_VERSION=$(curl -s https://api.github.com/repos/tungcorn/antigravity-usage-checker/releases/latest | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$INSTALL_VERSION" ]; then
        echo -e "${YELLOW}⚠️  Could not fetch latest version, using v0.5.0${NC}"
        INSTALL_VERSION="v0.5.0"
    fi
else
    INSTALL_VERSION="v$VERSION"
    echo "📥 Installing specific version: $INSTALL_VERSION"
fi

echo "📦 Version: $INSTALL_VERSION"

# Download URL
DOWNLOAD_URL="https://github.com/tungcorn/antigravity-usage-checker/releases/download/${INSTALL_VERSION}/antigravity-usage-checker-${OS}-${ARCH}.tar.gz"

# Create temp directory
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

# Download and extract
echo "📥 Downloading from: $DOWNLOAD_URL"
curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/agusage.tar.gz"

if [ ! -f "$TMP_DIR/agusage.tar.gz" ]; then
    echo -e "${RED}❌ Download failed${NC}"
    exit 1
fi

echo "📦 Extracting..."
tar -xzf "$TMP_DIR/agusage.tar.gz" -C "$TMP_DIR"

# Install to /usr/local/bin (may need sudo)
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="agusage"

if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_DIR/agusage" "$INSTALL_DIR/$BINARY_NAME"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
else
    echo "🔐 Requesting sudo access to install to $INSTALL_DIR..."
    sudo mv "$TMP_DIR/agusage" "$INSTALL_DIR/$BINARY_NAME"
    sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
fi

# Verify installation
if command -v agusage &> /dev/null; then
    echo ""
    echo -e "${GREEN}✅ Installation successful!${NC}"
    echo ""
    echo "Run 'agusage' to check your Antigravity usage quota."
    echo ""
    agusage --version
else
    echo -e "${RED}❌ Installation failed. Please check the error above.${NC}"
    exit 1
fi
