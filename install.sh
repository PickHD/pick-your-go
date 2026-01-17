#!/bin/bash
# Installation script for pick-your-go

set -e

INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
REPO_URL="https://github.com/PickHD/pick-your-go"
BINARY_NAME="pick-your-go"

echo "Installing $BINARY_NAME to $INSTALL_DIR..."

# Detect platform
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

PLATFORM="${OS}-${ARCH}"

echo "Detected platform: $PLATFORM"

# Check if we're in the repo directory
if [ -f "./cmd/$BINARY_NAME/main.go" ]; then
    echo "Building from source..."
    make build
    BINARY_PATH="./bin/$BINARY_NAME"
else
    echo "Not in repo directory. Please run this script from the repo root."
    exit 1
fi

# Install
if [ "$INSTALL_DIR" = "/usr/local/bin" ]; then
    echo "Installing to system directory (requires sudo)..."
    sudo cp "$BINARY_PATH" "$INSTALL_DIR/$BINARY_NAME"
    sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
else
    echo "Installing to $INSTALL_DIR..."
    mkdir -p "$INSTALL_DIR"
    cp "$BINARY_PATH" "$INSTALL_DIR/$BINARY_NAME"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
fi

echo ""
echo "✓ Installation complete!"
echo ""
echo "To use $BINARY_NAME, ensure $INSTALL_DIR is in your PATH."
echo ""
echo "Try running:"
echo "  $BINARY_NAME --help"
echo ""
