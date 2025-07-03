#!/bin/bash

# Build and Notarize Script for NHRL MCP Server
# This script builds the Go MCP server for all platforms and notarizes macOS binaries

set -e  # Exit on any error

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
BUILD_DIR="$PROJECT_DIR/build"
BINARY_NAME="nhrl-mcp-server"
VERSION="v1.9.0"

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

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to load environment variables
load_env() {
    # Check for environment file in scripts directory
    local env_file="$SCRIPT_DIR/.env"
    
    # If credentials are not already set, try to load from .env file
    if [ -z "$APPLE_ID" ] || [ -z "$APPLE_PASSWORD" ]; then
        if [ -f "$env_file" ]; then
            log_info "Loading credentials from $env_file"
            # Export variables from .env file
            set -a
            source "$env_file"
            set +a
        fi
    fi
    
    # Validate required credentials
    if [ -z "$APPLE_ID" ] || [ -z "$APPLE_PASSWORD" ]; then
        log_error "Apple credentials not found!"
        echo
        echo "Please set the following environment variables:"
        echo "  export APPLE_ID='your-apple-id@example.com'"
        echo "  export APPLE_PASSWORD='xxxx-xxxx-xxxx-xxxx'"
        echo
        echo "Or create a file at $env_file with:"
        echo "  APPLE_ID=your-apple-id@example.com"
        echo "  APPLE_PASSWORD=xxxx-xxxx-xxxx-xxxx"
        echo
        exit 1
    fi
    
    log_info "Using Apple ID: $APPLE_ID"
}

# Function to check Go installation
check_go() {
    if ! command_exists go; then
        log_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    local go_version
    go_version=$(go version | cut -d' ' -f3)
    log_info "Using Go version: $go_version"
}

# Function to check for notarization prerequisites
check_notarization_prerequisites() {
    local has_prerequisites=true
    
    if [ "$(uname)" != "Darwin" ]; then
        log_error "Notarization is only available on macOS"
        exit 1
    fi
    
    if ! command_exists codesign; then
        log_error "codesign not found - Xcode Command Line Tools may not be installed"
        has_prerequisites=false
    fi
    
    if ! command_exists xcrun; then
        log_error "xcrun not found - Xcode Command Line Tools may not be installed"
        has_prerequisites=false
    fi
    
    if [ "$has_prerequisites" = false ]; then
        log_error "Install Xcode Command Line Tools with: xcode-select --install"
        exit 1
    fi
    
    return 0
}

# Function to clean build directory
clean_build() {
    log_info "Cleaning build directory..."
    if [ -d "$BUILD_DIR" ]; then
        rm -rf "$BUILD_DIR"
    fi
    mkdir -p "$BUILD_DIR"
    log_success "Build directory cleaned"
}

# Function to build for all platforms
build_all_platforms() {
    log_info "Building for all platforms..."
    
    cd "$PROJECT_DIR"
    
    # Build flags for optimized binaries (strip symbol table and debug info)
    local ldflags="-s -w"
    
    # Linux AMD64
    log_info "Building for Linux AMD64..."
    GOOS=linux GOARCH=amd64 go build -ldflags "$ldflags" -o "$BUILD_DIR/${BINARY_NAME}-linux-amd64" .
    
    # Linux ARM64
    log_info "Building for Linux ARM64..."
    GOOS=linux GOARCH=arm64 go build -ldflags "$ldflags" -o "$BUILD_DIR/${BINARY_NAME}-linux-arm64" .
    
    # macOS AMD64
    log_info "Building for macOS AMD64..."
    GOOS=darwin GOARCH=amd64 go build -ldflags "$ldflags" -o "$BUILD_DIR/${BINARY_NAME}-darwin-amd64" .
    
    # macOS ARM64 (Apple Silicon)
    log_info "Building for macOS ARM64..."
    GOOS=darwin GOARCH=arm64 go build -ldflags "$ldflags" -o "$BUILD_DIR/${BINARY_NAME}-darwin-arm64" .
    
    # Windows AMD64
    log_info "Building for Windows AMD64..."
    GOOS=windows GOARCH=amd64 go build -ldflags "$ldflags" -o "$BUILD_DIR/${BINARY_NAME}-windows-amd64.exe" .
    
    # Windows ARM64
    log_info "Building for Windows ARM64..."
    GOOS=windows GOARCH=arm64 go build -ldflags "$ldflags" -o "$BUILD_DIR/${BINARY_NAME}-windows-arm64.exe" .
    
    log_success "All platforms built successfully"
}

# Function to sign macOS binary
sign_macos_binary() {
    local binary_path="$1"
    local binary_name=$(basename "$binary_path")
    
    log_info "Signing $binary_name..."
    
    # If DEVELOPER_ID is set, use it. Otherwise, sign ad-hoc for notarization
    if [ -n "$DEVELOPER_ID" ]; then
        log_info "Signing with Developer ID: $DEVELOPER_ID"
        codesign --sign "$DEVELOPER_ID" \
                 --options runtime \
                 --timestamp \
                 --force \
                 --verbose \
                 "$binary_path"
    else
        log_info "Signing ad-hoc for notarization"
        codesign --sign - \
                 --options runtime \
                 --timestamp \
                 --force \
                 --verbose \
                 "$binary_path"
    fi
    
    # Verify signature
    if codesign --verify --verbose "$binary_path"; then
        log_success "$binary_name signed successfully"
    else
        log_error "Failed to sign $binary_name"
        return 1
    fi
}

# Function to notarize a single macOS binary
notarize_binary() {
    local binary_path="$1"
    local binary_name=$(basename "$binary_path")
    local zip_file="${binary_path}.zip"
    
    log_info "Preparing $binary_name for notarization..."
    
    # Create ZIP file for notarization
    cd "$(dirname "$binary_path")"
    zip -j "$zip_file" "$binary_path"
    
    log_info "Submitting $binary_name for notarization..."
    
    # Submit for notarization using Apple ID and app-specific password
    local submission_output
    local notarytool_args="--apple-id $APPLE_ID --password $APPLE_PASSWORD"
    
    # Add team-id if provided
    if [ -n "$TEAM_ID" ]; then
        notarytool_args="$notarytool_args --team-id $TEAM_ID"
    fi
    
    submission_output=$(xcrun notarytool submit "$zip_file" \
                        $notarytool_args \
                        --wait 2>&1) || {
        log_error "Failed to submit for notarization"
        echo "$submission_output"
        return 1
    }
    
    echo "$submission_output"
    
    # Check if notarization was successful
    if echo "$submission_output" | grep -q "status: Accepted"; then
        log_success "$binary_name notarized successfully"
        
        # Staple the notarization
        log_info "Stapling notarization to $binary_name..."
        if xcrun stapler staple "$binary_path"; then
            log_success "Notarization stapled to $binary_name"
        else
            log_warning "Failed to staple notarization to $binary_name (this is normal for CLI tools)"
        fi
    else
        log_error "Notarization failed for $binary_name"
        return 1
    fi
    
    # Clean up ZIP file
    rm -f "$zip_file"
    
    cd "$PROJECT_DIR"
}

# Function to sign and notarize all macOS binaries
sign_and_notarize_macos() {
    log_info "Signing and notarizing macOS binaries..."
    
    local macos_binaries=(
        "$BUILD_DIR/${BINARY_NAME}-darwin-amd64"
        "$BUILD_DIR/${BINARY_NAME}-darwin-arm64"
    )
    
    for binary in "${macos_binaries[@]}"; do
        if [ -f "$binary" ]; then
            if sign_macos_binary "$binary"; then
                if notarize_binary "$binary"; then
                    log_success "$(basename "$binary") signed and notarized"
                else
                    log_error "Failed to notarize $(basename "$binary")"
                fi
            else
                log_error "Failed to sign $(basename "$binary")"
            fi
        else
            log_warning "Binary not found: $binary"
        fi
    done
}

# Function to create release packages
create_release_packages() {
    log_info "Creating release packages..."
    
    cd "$BUILD_DIR"
    
    # Create tar.gz for Unix-like systems
    if [ -f "${BINARY_NAME}-linux-amd64" ]; then
        tar -czf "${BINARY_NAME}-${VERSION}-linux-amd64.tar.gz" "${BINARY_NAME}-linux-amd64"
        log_success "Created Linux AMD64 package"
    fi
    
    if [ -f "${BINARY_NAME}-linux-arm64" ]; then
        tar -czf "${BINARY_NAME}-${VERSION}-linux-arm64.tar.gz" "${BINARY_NAME}-linux-arm64"
        log_success "Created Linux ARM64 package"
    fi
    
    if [ -f "${BINARY_NAME}-darwin-amd64" ]; then
        tar -czf "${BINARY_NAME}-${VERSION}-darwin-amd64.tar.gz" "${BINARY_NAME}-darwin-amd64"
        log_success "Created macOS AMD64 package"
    fi
    
    if [ -f "${BINARY_NAME}-darwin-arm64" ]; then
        tar -czf "${BINARY_NAME}-${VERSION}-darwin-arm64.tar.gz" "${BINARY_NAME}-darwin-arm64"
        log_success "Created macOS ARM64 package"
    fi
    
    # Create zip for Windows
    if [ -f "${BINARY_NAME}-windows-amd64.exe" ]; then
        zip "${BINARY_NAME}-${VERSION}-windows-amd64.zip" "${BINARY_NAME}-windows-amd64.exe"
        log_success "Created Windows AMD64 package"
    fi
    
    if [ -f "${BINARY_NAME}-windows-arm64.exe" ]; then
        zip "${BINARY_NAME}-${VERSION}-windows-arm64.zip" "${BINARY_NAME}-windows-arm64.exe"
        log_success "Created Windows ARM64 package"
    fi
    
    cd "$PROJECT_DIR"
}

# Function to generate checksums
generate_checksums() {
    log_info "Generating checksums..."
    
    cd "$BUILD_DIR"
    
    # Generate checksums for all release packages
    find . -name "${BINARY_NAME}-${VERSION}-*" -type f -exec shasum -a 256 {} \; > checksums.sha256
    
    if [ -f "checksums.sha256" ]; then
        log_success "Checksums generated in checksums.sha256"
        log_info "Checksum contents:"
        cat checksums.sha256
    else
        log_warning "Failed to generate checksums"
    fi
    
    cd "$PROJECT_DIR"
}

# Function to show build summary
show_summary() {
    echo
    log_info "Build Summary"
    echo "============================================"
    echo "Project: $BINARY_NAME"
    echo "Version: $VERSION"
    echo "Build Directory: $BUILD_DIR"
    echo
    
    if [ -d "$BUILD_DIR" ]; then
        echo "Build artifacts:"
        ls -la "$BUILD_DIR" | grep -E "nhrl-mcp-server|\.tar\.gz$|\.zip$" || echo "No artifacts found"
    fi
    
    echo
    echo "Total build directory size:"
    du -sh "$BUILD_DIR" 2>/dev/null || echo "Build directory not found"
}

# Main execution
main() {
    log_info "Starting NHRL MCP Server build and notarization process..."
    
    # Load environment variables
    load_env
    
    # Check prerequisites
    check_go
    check_notarization_prerequisites
    
    # Clean and prepare
    clean_build
    
    # Build for all platforms
    build_all_platforms
    
    # Sign and notarize macOS binaries
    sign_and_notarize_macos
    
    # Create release packages
    create_release_packages
    
    # Generate checksums
    generate_checksums
    
    # Show summary
    show_summary
    
    log_success "Build and notarization process completed!"
    
    # Show next steps
    echo
    log_info "Next steps:"
    echo "1. Test the binaries on their respective platforms"
    echo "2. Upload release packages to your distribution channel"
    echo "3. The notarized macOS binaries are ready for distribution"
    echo
    log_info "Security reminder:"
    echo "- Never commit your Apple ID or app-specific password to git"
    echo "- Keep your .env file secure and local only"
}

# Handle command line arguments
case "${1:-}" in
    --help|-h)
        echo "NHRL MCP Server Build and Notarize Script"
        echo
        echo "Usage: $0 [options]"
        echo
        echo "Options:"
        echo "  --help, -h          Show this help message"
        echo
        echo "Environment Variables:"
        echo "  APPLE_ID            Your Apple ID (required)"
        echo "  APPLE_PASSWORD      App-specific password (required)"
        echo "  DEVELOPER_ID        Code signing identity (optional)"
        echo "  TEAM_ID             Apple Developer Team ID (optional)"
        echo
        echo "Example:"
        echo "  export APPLE_ID='mcwiggin@mac.com'"
        echo "  export APPLE_PASSWORD='xxxx-xxxx-xxxx-xxxx'"
        echo "  $0"
        echo
        echo "Or create scripts/.env file with these variables"
        exit 0
        ;;
    *)
        main
        ;;
esac 