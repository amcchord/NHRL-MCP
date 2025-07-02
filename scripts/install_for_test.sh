#!/bin/bash

# NHRL MCP Server Test Installation Script v1.1
# This script installs the current platform version of the MCP server to /usr/local/bin
# Only use this script for testing purposes

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
            log_info "Supported platforms: Linux (AMD64/ARM64), macOS (Intel/Apple Silicon)"
            exit 1
            ;;
    esac
    echo "$platform"
}

# Main installation function
main() {
    log_info "NHRL MCP Server Test Installation v1.1"
    echo "========================================"

    # Check if running as root
    if [ "$EUID" -ne 0 ]; then
        log_error "Please run as root (use sudo)"
        log_info "Usage: sudo $0"
        exit 1
    fi

    # Detect platform
    local platform
    platform=$(detect_platform)
    log_info "Detected platform: $platform"

    # Check if build directory exists
    if [ ! -d "$BUILD_DIR" ]; then
        log_error "Build directory not found: $BUILD_DIR"
        log_info "Run 'make build-all' or './scripts/build-and-sign.sh' first"
        exit 1
    fi

    # Check if binary exists
    local source_binary="$BUILD_DIR/${BINARY_NAME}-${platform}"
    if [ ! -f "$source_binary" ]; then
        log_error "Binary not found: $source_binary"
        log_info "Available binaries:"
        ls -la "$BUILD_DIR" | grep "$BINARY_NAME" || echo "No binaries found"
        log_info "Run 'make build-all' or './scripts/build-and-sign.sh' first"
        exit 1
    fi

    # Show binary info
    local size
    size=$(ls -lh "$source_binary" | awk '{print $5}')
    log_info "Source binary: $source_binary ($size)"

    # Create install directory if it doesn't exist
    if [ ! -d "$INSTALL_DIR" ]; then
        log_info "Creating install directory: $INSTALL_DIR"
        mkdir -p "$INSTALL_DIR"
    fi

    # Check if binary is already installed
    local target_binary="$INSTALL_DIR/$BINARY_NAME"
    if [ -f "$target_binary" ]; then
        log_warning "Binary already exists at $target_binary"
        
        # Show current version if possible
        if command -v "$BINARY_NAME" >/dev/null 2>&1; then
            log_info "Current version:"
            "$BINARY_NAME" -version 2>/dev/null || echo "Version information not available"
        fi
        
        read -p "Do you want to overwrite it? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "Installation cancelled"
            exit 0
        fi
    fi

    # Copy and set permissions
    log_info "Installing $source_binary to $target_binary..."
    cp "$source_binary" "$target_binary"
    chmod +x "$target_binary"

    # Verify installation
    if [ -f "$target_binary" ] && [ -x "$target_binary" ]; then
        log_success "Installation complete"
        
        # Show installed binary info
        local installed_size
        installed_size=$(ls -lh "$target_binary" | awk '{print $5}')
        log_info "Installed binary: $target_binary ($installed_size)"
        
        # Test the installation
        log_info "Testing installation..."
        
        # Test help flag
        if "$target_binary" -help >/dev/null 2>&1; then
            log_success "✓ Help command works"
        else
            log_warning "✗ Help command failed"
        fi
        
        # Test version flag
        if "$target_binary" -version >/dev/null 2>&1; then
            log_success "✓ Version command works"
            log_info "Version information:"
            "$target_binary" -version 2>/dev/null || echo "Version not available"
        else
            log_warning "✗ Version command failed"
        fi
        
        # Show usage information
        echo
        log_info "Usage:"
        echo "  $BINARY_NAME -help                    # Show help"
        echo "  $BINARY_NAME -version                 # Show version"
        echo "  $BINARY_NAME -api-user-id <id> -api-key <key>  # Run with credentials"
        echo
        log_info "Environment variables:"
        echo "  NHRL_API_USER_ID    # NHRL API User ID"
        echo "  NHRL_API_KEY        # NHRL API Key"
        echo
        log_info "Test with dummy credentials:"
        echo "  NHRL_API_USER_ID=test NHRL_API_KEY=test $BINARY_NAME -exit-after-first"
        
    else
        log_error "Installation failed - binary not executable"
        exit 1
    fi
}

# Handle command line arguments
case "${1:-}" in
    --help|-h)
        echo "NHRL MCP Server Test Installation Script v1.1"
        echo
        echo "Usage: sudo $0"
        echo
        echo "This script:"
        echo "1. Detects your platform (Linux/macOS, AMD64/ARM64)"
        echo "2. Finds the appropriate binary in the build directory"
        echo "3. Installs it to $INSTALL_DIR/$BINARY_NAME"
        echo "4. Tests the installation"
        echo
        echo "Prerequisites:"
        echo "- Run as root (sudo)"
        echo "- Build the binaries first: make build-all"
        echo
        echo "Platform support:"
        echo "- Linux AMD64/ARM64"
        echo "- macOS Intel/Apple Silicon"
        exit 0
        ;;
    *)
        main
        ;;
esac 