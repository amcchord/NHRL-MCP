#!/bin/bash

# Build and Sign Script for NHRL MCP Server v1.1
# This script builds the Go MCP server for all platforms and optionally signs macOS binaries

set -e  # Exit on any error

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
BUILD_DIR="$PROJECT_DIR/build"
BINARY_NAME="nhrl-mcp-server"
VERSION="v1.5.0"

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

# Function to check for code signing prerequisites
check_signing_prerequisites() {
    local has_prerequisites=true
    
    if [ "$(uname)" != "Darwin" ]; then
        log_warning "Code signing is only available on macOS"
        return 1
    fi
    
    if ! command_exists codesign; then
        log_warning "codesign not found - Xcode Command Line Tools may not be installed"
        has_prerequisites=false
    fi
    
    if ! command_exists xcrun; then
        log_warning "xcrun not found - Xcode Command Line Tools may not be installed"
        has_prerequisites=false
    fi
    
    if [ "$has_prerequisites" = false ]; then
        log_warning "Install Xcode Command Line Tools with: xcode-select --install"
        return 1
    fi
    
    return 0
}

# Function to list available signing identities
list_signing_identities() {
    log_info "Available code signing identities:"
    security find-identity -v -p codesigning || {
        log_warning "No code signing identities found"
        return 1
    }
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
    
    # Linux AMD64
    log_info "Building for Linux AMD64..."
    GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o "$BUILD_DIR/${BINARY_NAME}-linux-amd64" .
    
    # Linux ARM64
    log_info "Building for Linux ARM64..."
    GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o "$BUILD_DIR/${BINARY_NAME}-linux-arm64" .
    
    # macOS AMD64
    log_info "Building for macOS AMD64..."
    GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o "$BUILD_DIR/${BINARY_NAME}-darwin-amd64" .
    
    # macOS ARM64 (Apple Silicon)
    log_info "Building for macOS ARM64..."
    GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o "$BUILD_DIR/${BINARY_NAME}-darwin-arm64" .
    
    # Windows AMD64
    log_info "Building for Windows AMD64..."
    GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o "$BUILD_DIR/${BINARY_NAME}-windows-amd64.exe" .
    
    # Windows ARM64
    log_info "Building for Windows ARM64..."
    GOOS=windows GOARCH=arm64 go build -ldflags "-s -w" -o "$BUILD_DIR/${BINARY_NAME}-windows-arm64.exe" .
    
    log_success "All platforms built successfully"
}

# Function to sign macOS binaries
sign_macos_binaries() {
    if [ -z "$DEVELOPER_ID" ]; then
        log_warning "DEVELOPER_ID not set - skipping macOS code signing"
        log_info "To enable code signing, set DEVELOPER_ID environment variable:"
        log_info "export DEVELOPER_ID='Developer ID Application: Your Name (TEAMID)'"
        return 1
    fi
    
    log_info "Signing macOS binaries with identity: $DEVELOPER_ID"
    
    # Sign macOS AMD64
    local macos_amd64="$BUILD_DIR/${BINARY_NAME}-darwin-amd64"
    if [ -f "$macos_amd64" ]; then
        log_info "Signing macOS AMD64 binary..."
        codesign --sign "$DEVELOPER_ID" \
                 --options runtime \
                 --timestamp \
                 --verbose=2 \
                 "$macos_amd64"
        
        # Verify signature
        codesign --verify --verbose=2 "$macos_amd64"
        spctl --assess --type execute --verbose=2 "$macos_amd64"
        log_success "macOS AMD64 binary signed and verified"
    fi
    
    # Sign macOS ARM64
    local macos_arm64="$BUILD_DIR/${BINARY_NAME}-darwin-arm64"
    if [ -f "$macos_arm64" ]; then
        log_info "Signing macOS ARM64 binary..."
        codesign --sign "$DEVELOPER_ID" \
                 --options runtime \
                 --timestamp \
                 --verbose=2 \
                 "$macos_arm64"
        
        # Verify signature
        codesign --verify --verbose=2 "$macos_arm64"
        spctl --assess --type execute --verbose=2 "$macos_arm64"
        log_success "macOS ARM64 binary signed and verified"
    fi
}

# Function to notarize macOS binaries
notarize_macos_binaries() {
    if [ -z "$KEYCHAIN_PROFILE" ]; then
        log_warning "KEYCHAIN_PROFILE not set - skipping notarization"
        log_info "To enable notarization, set KEYCHAIN_PROFILE environment variable:"
        log_info "export KEYCHAIN_PROFILE='your-keychain-profile'"
        return 1
    fi
    
    log_info "Creating ZIP archives for notarization..."
    
    cd "$BUILD_DIR"
    
    # Create ZIP files for notarization
    if [ -f "${BINARY_NAME}-darwin-amd64" ]; then
        zip "${BINARY_NAME}-darwin-amd64.zip" "${BINARY_NAME}-darwin-amd64"
    fi
    
    if [ -f "${BINARY_NAME}-darwin-arm64" ]; then
        zip "${BINARY_NAME}-darwin-arm64.zip" "${BINARY_NAME}-darwin-arm64"
    fi
    
    # Submit for notarization
    log_info "Submitting binaries for notarization..."
    
    if [ -f "${BINARY_NAME}-darwin-amd64.zip" ]; then
        log_info "Notarizing macOS AMD64 binary..."
        xcrun notarytool submit "${BINARY_NAME}-darwin-amd64.zip" \
                             --keychain-profile "$KEYCHAIN_PROFILE" \
                             --wait
    fi
    
    if [ -f "${BINARY_NAME}-darwin-arm64.zip" ]; then
        log_info "Notarizing macOS ARM64 binary..."
        xcrun notarytool submit "${BINARY_NAME}-darwin-arm64.zip" \
                             --keychain-profile "$KEYCHAIN_PROFILE" \
                             --wait
    fi
    
    # Staple notarization
    log_info "Stapling notarization..."
    
    if [ -f "${BINARY_NAME}-darwin-amd64" ]; then
        xcrun stapler staple "${BINARY_NAME}-darwin-amd64"
    fi
    
    if [ -f "${BINARY_NAME}-darwin-arm64" ]; then
        xcrun stapler staple "${BINARY_NAME}-darwin-arm64"
    fi
    
    log_success "Notarization completed"
    cd "$PROJECT_DIR"
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

# Function to verify build outputs
verify_builds() {
    log_info "Verifying build outputs..."
    
    local all_good=true
    
    # Check if all expected binaries exist
    local expected_files=(
        "${BINARY_NAME}-linux-amd64"
        "${BINARY_NAME}-linux-arm64"
        "${BINARY_NAME}-darwin-amd64"
        "${BINARY_NAME}-darwin-arm64"
        "${BINARY_NAME}-windows-amd64.exe"
        "${BINARY_NAME}-windows-arm64.exe"
    )
    
    for file in "${expected_files[@]}"; do
        if [ -f "$BUILD_DIR/$file" ]; then
            local size
            size=$(ls -lh "$BUILD_DIR/$file" | awk '{print $5}')
            log_success "✓ $file ($size)"
        else
            log_error "✗ $file (missing)"
            all_good=false
        fi
    done
    
    # Check package files
    log_info "Checking release packages..."
    local package_files=(
        "${BINARY_NAME}-${VERSION}-linux-amd64.tar.gz"
        "${BINARY_NAME}-${VERSION}-linux-arm64.tar.gz"
        "${BINARY_NAME}-${VERSION}-darwin-amd64.tar.gz"
        "${BINARY_NAME}-${VERSION}-darwin-arm64.tar.gz"
        "${BINARY_NAME}-${VERSION}-windows-amd64.zip"
        "${BINARY_NAME}-${VERSION}-windows-arm64.zip"
    )
    
    for file in "${package_files[@]}"; do
        if [ -f "$BUILD_DIR/$file" ]; then
            local size
            size=$(ls -lh "$BUILD_DIR/$file" | awk '{print $5}')
            log_success "✓ $file ($size)"
        else
            log_warning "✗ $file (not created)"
        fi
    done
    
    if [ "$all_good" = true ]; then
        log_success "All builds verified successfully"
    else
        log_error "Some builds are missing or failed"
    fi
}

# Function to generate checksums
generate_checksums() {
    log_info "Generating checksums..."
    
    cd "$BUILD_DIR"
    
    # Generate checksums for all release packages
    find . -name "${BINARY_NAME}-${VERSION}-*" -type f -exec sha256sum {} \; > checksums.sha256
    
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

# Function to test basic functionality
test_functionality() {
    log_info "Testing basic functionality..."
    
    # Test the current platform binary
    local current_platform=""
    case "$(uname -s)-$(uname -m)" in
        "Linux-x86_64") current_platform="linux-amd64" ;;
        "Linux-aarch64") current_platform="linux-arm64" ;;
        "Darwin-x86_64") current_platform="darwin-amd64" ;;
        "Darwin-arm64") current_platform="darwin-arm64" ;;
        *) log_warning "Unknown platform for testing"; return ;;
    esac
    
    local test_binary="$BUILD_DIR/${BINARY_NAME}-${current_platform}"
    if [ -f "$test_binary" ]; then
        log_info "Testing $test_binary..."
        
        # Test help flag
        if "$test_binary" -help > /dev/null 2>&1; then
            log_success "✓ Help flag works"
        else
            log_warning "✗ Help flag failed"
        fi
        
        # Test version flag
        if "$test_binary" -version > /dev/null 2>&1; then
            log_success "✓ Version flag works"
        else
            log_warning "✗ Version flag failed"
        fi
        
        # Test tool registration (with dummy credentials)
        if NHRL_API_USER_ID=test NHRL_API_KEY=test "$test_binary" -exit-after-first -tools full > /dev/null 2>&1; then
            log_success "✓ Server starts and exits correctly"
        else
            log_warning "✗ Server startup test failed"
        fi
    else
        log_warning "Cannot test functionality - binary for current platform not found"
    fi
}

# Main execution
main() {
    log_info "Starting NHRL MCP Server build process..."
    
    # Check prerequisites
    check_go
    
    # Clean and prepare
    clean_build
    
    # Build for all platforms
    build_all_platforms
    
    # Handle code signing if on macOS
    if check_signing_prerequisites; then
        list_signing_identities
        
        if sign_macos_binaries; then
            log_success "macOS binaries signed successfully"
            
            # Attempt notarization if requested
            if [ "$ENABLE_NOTARIZATION" = "true" ] && [ -n "$KEYCHAIN_PROFILE" ]; then
                notarize_macos_binaries
            fi
        fi
    fi
    
    # Create release packages
    create_release_packages
    
    # Generate checksums
    generate_checksums
    
    # Test functionality
    test_functionality
    
    # Verify everything
    verify_builds
    
    # Show summary
    show_summary
    
    log_success "Build process completed!"
    
    # Show next steps
    echo
    log_info "Next steps:"
    echo "1. Test the binaries on their respective platforms"
    echo "2. Upload release packages to your distribution channel"
    echo "3. Update documentation with new version information"
    
    if [ -z "$DEVELOPER_ID" ] && [ "$(uname)" = "Darwin" ]; then
        echo
        log_info "To enable code signing in future builds:"
        echo "export DEVELOPER_ID='Developer ID Application: Your Name (TEAMID)'"
        echo "export KEYCHAIN_PROFILE='your-keychain-profile'  # For notarization"
        echo "export ENABLE_NOTARIZATION='true'  # To enable notarization"
    fi
}

# Handle command line arguments
case "${1:-}" in
    --help|-h)
        echo "NHRL MCP Server Build Script v1.1"
        echo
        echo "Usage: $0 [options]"
        echo
        echo "Options:"
        echo "  --help, -h          Show this help message"
        echo "  --clean-only        Only clean the build directory"
        echo "  --verify-only       Only verify existing builds"
        echo "  --test-only         Only test functionality"
        echo
        echo "Environment Variables:"
        echo "  DEVELOPER_ID        Code signing identity for macOS"
        echo "  KEYCHAIN_PROFILE    Keychain profile for notarization"
        echo "  ENABLE_NOTARIZATION Set to 'true' to enable notarization"
        echo
        echo "Examples:"
        echo "  $0"
        echo "  DEVELOPER_ID='Developer ID Application: Your Name' $0"
        echo "  ENABLE_NOTARIZATION=true KEYCHAIN_PROFILE=myprofile $0"
        exit 0
        ;;
    --clean-only)
        clean_build
        exit 0
        ;;
    --verify-only)
        verify_builds
        show_summary
        exit 0
        ;;
    --test-only)
        test_functionality
        exit 0
        ;;
    *)
        main
        ;;
esac 