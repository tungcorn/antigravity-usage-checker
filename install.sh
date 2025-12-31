#!/bin/bash
# Antigravity Usage Checker - Installation Script for Linux/macOS
# Usage: curl -fsSL https://raw.githubusercontent.com/TungCorn/antigravity-usage-checker/main/install.sh | bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

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

# Get latest version
echo "📥 Fetching latest release..."
LATEST_VERSION=$(curl -s https://api.github.com/repos/TungCorn/antigravity-usage-checker/releases/latest | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo -e "${YELLOW}⚠️  Could not fetch latest version, using v0.3.0${NC}"
    LATEST_VERSION="v0.3.0"
fi

echo "📦 Latest version: $LATEST_VERSION"

# Download URL
DOWNLOAD_URL="https://github.com/TungCorn/antigravity-usage-checker/releases/download/${LATEST_VERSION}/antigravity-usage-checker-${OS}-${ARCH}.tar.gz"

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
