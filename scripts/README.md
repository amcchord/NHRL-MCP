# TrueFinals MCP Server Build Scripts

This directory contains scripts for building, testing, and deploying the TrueFinals MCP Server.

## Scripts Overview

### 🔨 Build Scripts

#### `build-and-sign.sh`
Comprehensive build script that creates binaries for all supported platforms.

**Features:**
- Cross-platform builds (Linux, macOS, Windows - AMD64/ARM64)
- macOS code signing with Developer ID
- Apple notarization support
- Release package creation
- Checksum generation
- Built-in testing

**Usage:**
```bash
# Basic build
./scripts/build-and-sign.sh

# With code signing
DEVELOPER_ID='Developer ID Application: Your Name (TEAMID)' ./scripts/build-and-sign.sh

# With notarization
DEVELOPER_ID='...' KEYCHAIN_PROFILE='profile' ENABLE_NOTARIZATION=true ./scripts/build-and-sign.sh

# Help
./scripts/build-and-sign.sh --help
```

**Options:**
- `--help` - Show help information
- `--clean-only` - Only clean build directory
- `--verify-only` - Only verify existing builds
- `--test-only` - Only test functionality

#### `build-and-notarize.sh`
A comprehensive build script that builds binaries for all platforms and notarizes macOS binaries using Apple ID credentials.

**Features:**
- Builds for all supported platforms (Linux, macOS, Windows on both AMD64 and ARM64)
- Code signs macOS binaries (ad-hoc or with Developer ID)
- Notarizes macOS binaries with Apple
- Creates release packages (.tar.gz for Unix, .zip for Windows)
- Generates SHA256 checksums

**Usage:**
```bash
# First time setup - create credentials file
./setup-credentials.sh

# Build and notarize
./build-and-notarize.sh
```

**Security:**
- Credentials are stored in `scripts/.env` (ignored by git)
- The .env file has 600 permissions (owner read/write only)
- Never commit credentials to version control

### 🧪 Test Scripts

#### `test-functionality.sh`
Comprehensive functionality testing script.

**Tests:**
- Basic functionality (help, version flags)
- Server startup with different tool modes
- Environment variable handling
- Command line argument processing

**Usage:**
```bash
# Test current setup
./scripts/test-functionality.sh

# Test specific binary
./scripts/test-functionality.sh --binary /path/to/binary

# Help
./scripts/test-functionality.sh --help
```

### 📦 Installation Scripts

#### `install_for_test.sh`
Builds and installs the NHRL MCP Server to `/usr/local/bin` for testing purposes.

**What it does:**
1. Builds all binaries using `build-and-notarize.sh`
2. Detects your current platform
3. Installs the appropriate binary to `/usr/local/bin`

**Usage:**
```bash
sudo ./install_for_test.sh
```

**Note:** Requires sudo access to install to `/usr/local/bin`

### ⚙️ Configuration

#### `dev-env.template`
Template for development environment variables.

**Setup:**
```bash
# Copy template
cp scripts/dev-env.template scripts/dev-env.sh

# Edit with your credentials
nano scripts/dev-env.sh

# Source the environment
source scripts/dev-env.sh
```

**Variables:**
- `TRUEFINALS_API_USER_ID` - Your TrueFinals API User ID
- `TRUEFINALS_API_KEY` - Your TrueFinals API Key
- `TRUEFINALS_BASE_URL` - Custom API base URL (optional)
- `TRUEFINALS_TOOLS` - Tool filter mode (reporting/full-safe/full)
- `TRUEFINALS_DISABLED_TOOLS` - Comma-separated list of disabled tools
- `DEVELOPER_ID` - macOS code signing identity
- `KEYCHAIN_PROFILE` - macOS notarization profile

## Quick Start Guide

### 1. Development Setup
```bash
# Set up environment
cp scripts/dev-env.template scripts/dev-env.sh
# Edit dev-env.sh with your credentials
source scripts/dev-env.sh

# Build and test
make build
./scripts/test-functionality.sh
```

### 2. Production Build
```bash
# Clean build for all platforms
./scripts/build-and-sign.sh

# Verify build
./scripts/build-and-sign.sh --verify-only

# Test functionality
./scripts/test-functionality.sh
```

### 3. Installation
```bash
# Install for testing
sudo ./scripts/install_for_test.sh

# Test installation
truefinals-mcp-server -help
```

## Platform Support

### Supported Platforms
- **Linux**: AMD64, ARM64
- **macOS**: Intel (AMD64), Apple Silicon (ARM64)
- **Windows**: AMD64, ARM64

### Build Outputs
- Optimized binaries (debug symbols stripped)
- Release packages (tar.gz for Unix, zip for Windows)
- SHA256 checksums
- Code-signed macOS binaries (optional)

## Code Signing & Notarization

### macOS Code Signing
For distribution on macOS, binaries should be code-signed:

```bash
# Set up signing identity
export DEVELOPER_ID='Developer ID Application: Your Name (TEAMID)'

# Build with signing
./scripts/build-and-sign.sh
```

### Apple Notarization
For Gatekeeper compatibility:

```bash
# Set up notarization profile
xcrun notarytool store-credentials

# Set environment
export KEYCHAIN_PROFILE='your-profile-name'
export ENABLE_NOTARIZATION='true'

# Build with notarization
./scripts/build-and-sign.sh
```

## Environment Variables

### Required for Runtime
- `TRUEFINALS_API_USER_ID` - TrueFinals API User ID
- `TRUEFINALS_API_KEY` - TrueFinals API Key

### Optional Configuration
- `TRUEFINALS_BASE_URL` - API base URL (default: https://truefinals.com/api)
- `TRUEFINALS_TOOLS` - Tool mode (reporting/full-safe/full)
- `TRUEFINALS_DISABLED_TOOLS` - Disabled tools list

### Build & Signing
- `DEVELOPER_ID` - macOS code signing identity
- `KEYCHAIN_PROFILE` - Notarization keychain profile
- `ENABLE_NOTARIZATION` - Enable Apple notarization

## Troubleshooting

### Build Issues
```bash
# Check Go installation
go version

# Clean and rebuild
make clean
./scripts/build-and-sign.sh --clean-only
./scripts/build-and-sign.sh
```

### Test Failures
```bash
# Run specific test
./scripts/test-functionality.sh --binary build/truefinals-mcp-server

# Check binary
file build/truefinals-mcp-server*
```

### Installation Issues
```bash
# Check permissions
ls -la /usr/local/bin/truefinals-mcp-server

# Reinstall
sudo ./scripts/install_for_test.sh
```

## Security Notes

### Credential Management
- Never commit actual credentials to version control
- Use `dev-env.sh` for local development (ignored by git)
- Use environment variables in production
- Rotate API keys regularly

### Code Signing
- Keep signing certificates secure
- Use dedicated keychain profiles for notarization
- Test signed binaries before distribution

## Integration with Makefile

The scripts work seamlessly with the project Makefile:

```bash
# Build using Makefile (calls scripts internally)
make build-all
make release-signed
make install-test

# Direct script usage
./scripts/build-and-sign.sh
./scripts/test-functionality.sh
```

## Contributing

When adding new scripts:
1. Follow the existing naming conventions
2. Include help text and error handling
3. Use the established logging functions
4. Update this README
5. Test on multiple platforms

## Support

For script-related issues:
1. Check the help output: `script-name.sh --help`
2. Verify prerequisites (Go, credentials, etc.)
3. Run with verbose output for debugging
4. Check the main project README for additional information 