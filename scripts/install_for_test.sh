#!/bin/bash

# NHRL MCP Server Test Installation Script
# This script builds and installs the current platform version to /usr/local/bin

set -e  # Exit on any error

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
BUILD_DIR="$PROJECT_DIR/build"
BINARY_NAME="nhrl-mcp-server"
INSTALL_DIR="/usr/local/bin"

# Function to detect current platform
detect_platform() {
    local platform=""
    case "$(uname -s)-$(uname -m)" in
        "Linux-x86_64") platform="linux-amd64" ;;
        "Linux-aarch64") platform="linux-arm64" ;;
        "Darwin-x86_64") platform="darwin-amd64" ;;
        "Darwin-arm64") platform="darwin-arm64" ;;
        *)
            log_error "Unsupported platform: $(uname -s)-$(uname -m)"
            exit 1
            ;;
    esac
    echo "$platform"
}

# Main installation function
main() {
    log_info "NHRL MCP Server Test Installation"
    
    # Check if running as root
    if [ "$EUID" -ne 0 ]; then
        log_error "Please run as root (use sudo)"
        exit 1
    fi
    
    # Detect platform
    local platform
    platform=$(detect_platform)
    log_info "Detected platform: $platform"
    
    # Build using build-and-notarize.sh
    log_info "Building binaries..."
    if [ -f "$SCRIPT_DIR/build-and-notarize.sh" ]; then
        # Run build script as the original user (not root)
        # Use -E to preserve environment variables
        if [ -n "$SUDO_USER" ]; then
            sudo -E -u "$SUDO_USER" "$SCRIPT_DIR/build-and-notarize.sh"
        else
            # If not run with sudo, just run the build script
            "$SCRIPT_DIR/build-and-notarize.sh"
        fi
    else
        log_error "build-and-notarize.sh not found"
        exit 1
    fi
    
    # Check if binary exists
    local source_binary="$BUILD_DIR/${BINARY_NAME}-${platform}"
    if [ ! -f "$source_binary" ]; then
        log_error "Binary not found: $source_binary"
        exit 1
    fi
    
    # Copy to /usr/local/bin
    log_info "Installing to $INSTALL_DIR..."
    cp "$source_binary" "$INSTALL_DIR/$BINARY_NAME"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
    
    log_success "Installation complete!"
    log_info "Binary installed to: $INSTALL_DIR/$BINARY_NAME"
    
    # Verify the file exists and is executable
    if [ -f "$INSTALL_DIR/$BINARY_NAME" ] && [ -x "$INSTALL_DIR/$BINARY_NAME" ]; then
        log_success "Installation verified - binary is in place and executable"
        log_info "The binary has been signed and notarized by Apple"
        log_info "You can now run: $BINARY_NAME -version"
    else
        log_error "Installation failed - binary not found or not executable"
        exit 1
    fi
}

# Handle command line arguments
case "${1:-}" in
    --help|-h)
        echo "NHRL MCP Server Test Installation Script"
        echo
        echo "Usage: sudo $0"
        echo
        echo "This script:"
        echo "1. Builds the binaries using build-and-notarize.sh"
        echo "2. Installs the current platform binary to $INSTALL_DIR"
        exit 0
        ;;
    *)
        main
        ;;
esac 