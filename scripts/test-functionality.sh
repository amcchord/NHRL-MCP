#!/bin/bash

# TrueFinals MCP Server Functionality Test Script
# This script tests various aspects of the MCP server to ensure it's working correctly

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
BINARY_NAME="truefinals-mcp-server"
TEST_USER_ID="test"
TEST_API_KEY="test"

# Function to find the binary
find_binary() {
    local binary_path=""
    
    # Try system installation first
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        binary_path="$BINARY_NAME"
        log_info "Using system-installed binary: $(command -v "$BINARY_NAME")"
    # Try build directory
    elif [ -f "$PROJECT_DIR/build/$BINARY_NAME" ]; then
        binary_path="$PROJECT_DIR/build/$BINARY_NAME"
        log_info "Using local build: $binary_path"
    # Try project root
    elif [ -f "$PROJECT_DIR/$BINARY_NAME" ]; then
        binary_path="$PROJECT_DIR/$BINARY_NAME"
        log_info "Using project root binary: $binary_path"
    else
        log_error "Binary not found. Please build or install the server first."
        log_info "To build: make build or make build-all"
        log_info "To install: sudo make install-test"
        exit 1
    fi
    
    echo "$binary_path"
}

# Function to test basic functionality
test_basic_functionality() {
    local binary="$1"
    
    log_info "Testing basic functionality..."
    
    # Test help flag
    log_info "Testing help flag..."
    if "$binary" -help >/dev/null 2>&1; then
        log_success "✓ Help flag works"
    else
        log_error "✗ Help flag failed"
        return 1
    fi
    
    # Test version flag
    log_info "Testing version flag..."
    if "$binary" -version >/dev/null 2>&1; then
        log_success "✓ Version flag works"
        local version_output
        version_output=$("$binary" -version 2>/dev/null || echo "Version not available")
        log_info "Version: $version_output"
    else
        log_error "✗ Version flag failed"
        return 1
    fi
    
    return 0
}

# Function to test server startup and tools registration
test_server_startup() {
    local binary="$1"
    
    log_info "Testing server startup and tools registration..."
    
    # Test with different tool modes
    local tool_modes=("reporting" "full-safe" "full")
    
    for mode in "${tool_modes[@]}"; do
        log_info "Testing tools mode: $mode"
        
        local output
        output=$(TRUEFINALS_API_USER_ID="$TEST_USER_ID" TRUEFINALS_API_KEY="$TEST_API_KEY" \
                "$binary" -exit-after-first -tools "$mode" 2>&1)
        
        if [ $? -eq 0 ]; then
            log_success "✓ Server starts successfully with tools mode: $mode"
            
            # Check if tools are mentioned in output
            local tool_count
            tool_count=$(echo "$output" | grep -c -i "tool" || true)
            if [ "$tool_count" -gt 0 ]; then
                log_success "✓ Tools registered ($tool_count tool references found)"
            else
                log_warning "! No tool references found in output"
            fi
        else
            log_error "✗ Server failed to start with tools mode: $mode"
            log_error "Output: $output"
            return 1
        fi
    done
    
    return 0
}

# Function to test environment variable handling
test_environment_variables() {
    local binary="$1"
    
    log_info "Testing environment variable handling..."
    
    # Test without credentials (should fail gracefully)
    log_info "Testing without credentials..."
    local output
    output=$("$binary" -exit-after-first 2>&1 || true)
    
    if echo "$output" | grep -q -i "api.*key\|credentials\|user.*id"; then
        log_success "✓ Properly handles missing credentials"
    else
        log_warning "! May not properly validate credentials"
    fi
    
    # Test with environment variables
    log_info "Testing with environment variables..."
    output=$(TRUEFINALS_API_USER_ID="$TEST_USER_ID" TRUEFINALS_API_KEY="$TEST_API_KEY" \
            "$binary" -exit-after-first 2>&1)
    
    if [ $? -eq 0 ]; then
        log_success "✓ Environment variables work correctly"
    else
        log_error "✗ Environment variables test failed"
        return 1
    fi
    
    return 0
}

# Function to test command line arguments
test_command_line_args() {
    local binary="$1"
    
    log_info "Testing command line arguments..."
    
    # Test with command line arguments
    local output
    output=$("$binary" -api-user-id "$TEST_USER_ID" -api-key "$TEST_API_KEY" -exit-after-first 2>&1)
    
    if [ $? -eq 0 ]; then
        log_success "✓ Command line arguments work correctly"
    else
        log_error "✗ Command line arguments test failed"
        return 1
    fi
    
    return 0
}

# Function to run comprehensive tests
run_comprehensive_tests() {
    local binary="$1"
    
    log_info "Running comprehensive functionality tests..."
    echo "================================================"
    
    local test_count=0
    local passed_count=0
    
    # Basic functionality tests
    ((test_count++))
    if test_basic_functionality "$binary"; then
        ((passed_count++))
    fi
    
    # Server startup tests
    ((test_count++))
    if test_server_startup "$binary"; then
        ((passed_count++))
    fi
    
    # Environment variable tests
    ((test_count++))
    if test_environment_variables "$binary"; then
        ((passed_count++))
    fi
    
    # Command line argument tests
    ((test_count++))
    if test_command_line_args "$binary"; then
        ((passed_count++))
    fi
    
    # Show results
    echo
    log_info "Test Results Summary"
    echo "===================="
    echo "Total tests: $test_count"
    echo "Passed: $passed_count"
    echo "Failed: $((test_count - passed_count))"
    
    if [ "$passed_count" -eq "$test_count" ]; then
        log_success "All tests passed! ✅"
        return 0
    else
        log_error "Some tests failed! ❌"
        return 1
    fi
}

# Function to show system information
show_system_info() {
    log_info "System Information"
    echo "=================="
    echo "OS: $(uname -s) $(uname -r)"
    echo "Architecture: $(uname -m)"
    echo "Go Version: $(go version 2>/dev/null || echo 'Go not found')"
    echo "Shell: $SHELL"
    echo "User: $(whoami)"
    echo "Working Directory: $(pwd)"
    echo
}

# Main function
main() {
    log_info "TrueFinals MCP Server Functionality Test"
    echo "========================================="
    
    # Show system info
    show_system_info
    
    # Find binary
    local binary
    binary=$(find_binary)
    
    if [ -z "$binary" ]; then
        log_error "Could not find binary to test"
        exit 1
    fi
    
    # Show binary info
    if [ -f "$binary" ]; then
        local size
        size=$(ls -lh "$binary" | awk '{print $5}')
        log_info "Binary: $binary ($size)"
    fi
    
    # Run tests
    if run_comprehensive_tests "$binary"; then
        log_success "All functionality tests completed successfully!"
        echo
        log_info "The TrueFinals MCP Server appears to be working correctly."
        log_info "You can now:"
        echo "  1. Set up your TrueFinals API credentials"
        echo "  2. Configure your MCP client (like Claude Desktop)"
        echo "  3. Start using the tournament management tools"
        exit 0
    else
        log_error "Some functionality tests failed!"
        echo
        log_info "Please check the errors above and:"
        echo "  1. Make sure the binary was built correctly"
        echo "  2. Check for any missing dependencies"
        echo "  3. Verify the server can access the TrueFinals API"
        exit 1
    fi
}

# Handle command line arguments
case "${1:-}" in
    --help|-h)
        echo "TrueFinals MCP Server Functionality Test Script"
        echo
        echo "Usage: $0 [options]"
        echo
        echo "Options:"
        echo "  --help, -h     Show this help message"
        echo "  --binary PATH  Test a specific binary path"
        echo
        echo "This script tests:"
        echo "  1. Basic functionality (help, version flags)"
        echo "  2. Server startup with different tool modes"
        echo "  3. Environment variable handling"
        echo "  4. Command line argument processing"
        echo
        echo "The script uses test credentials and exits immediately after"
        echo "testing to avoid making actual API calls."
        exit 0
        ;;
    --binary)
        if [ -n "$2" ]; then
            if [ -f "$2" ]; then
                log_info "Testing specific binary: $2"
                run_comprehensive_tests "$2"
                exit $?
            else
                log_error "Binary not found: $2"
                exit 1
            fi
        else
            log_error "--binary requires a path argument"
            exit 1
        fi
        ;;
    *)
        main
        ;;
esac 