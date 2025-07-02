# NHRL MCP Server v1.1.0 Release Notes

## 🚀 What's New

### Major Updates
- **Complete NHRL Rebranding**: Updated from TrueFinals MCP to NHRL MCP Server
- **Version 1.1.0**: Updated version numbering system
- **Enhanced Build System**: Improved cross-platform builds with proper signing and notarization

### 🔧 Technical Improvements

#### Build & Distribution
- ✅ **Cross-Platform Binaries**: Built for all major platforms
  - Linux (AMD64, ARM64)
  - macOS (Intel, Apple Silicon) - **Code Signed & Notarized**
  - Windows (AMD64, ARM64)
- ✅ **macOS Security**: All macOS binaries are code-signed and notarized for seamless installation
- ✅ **Optimized Binaries**: Stripped debug symbols for smaller file sizes
- ✅ **Release Packages**: Complete `.tar.gz` and `.zip` packages with SHA256 checksums

#### Configuration Updates
- 🔄 **Environment Variables**: Updated from `TRUEFINALS_*` to `NHRL_*`
  - `NHRL_API_USER_ID` - Your NHRL API User ID
  - `NHRL_API_KEY` - Your NHRL API Key
- 🔄 **Binary Name**: Changed from `truefinals-mcp-server` to `nhrl-mcp-server`
- 🔄 **Server Identity**: Updated server name and branding throughout

### 🛠️ Installation & Usage

#### Quick Install (macOS/Linux)
```bash
# Download the appropriate package for your platform
# Extract and install to /usr/local/bin
sudo cp nhrl-mcp-server /usr/local/bin/
chmod +x /usr/local/bin/nhrl-mcp-server
```

#### Usage
```bash
# Show version
nhrl-mcp-server -version

# Run with credentials
nhrl-mcp-server -api-user-id YOUR_ID -api-key YOUR_KEY

# Or use environment variables
export NHRL_API_USER_ID=your_id
export NHRL_API_KEY=your_key
nhrl-mcp-server
```

### 📋 Available Tools
The server provides access to comprehensive NHRL data through these tools:
- `nhrl_stats` - NHRL statistics and data
- `truefinals_tournaments` - Tournament information
- `truefinals_games` - Game details and results
- `truefinals_locations` - Venue and location data
- `truefinals_players` - Player profiles and statistics
- `truefinals_bracket` - Tournament bracket management

### 🔒 Security Features
- **Code Signing**: macOS binaries signed with Developer ID
- **Notarization**: Apple-notarized for Gatekeeper compatibility
- **Tools Filtering**: Multiple security modes (reporting, full-safe, full)
- **API Authentication**: Secure API key authentication

### 📦 Download Information

| Platform | Architecture | File | Size |
|----------|-------------|------|------|
| Linux | AMD64 | `nhrl-mcp-server-v1.1.0-linux-amd64.tar.gz` | 2.6 MB |
| Linux | ARM64 | `nhrl-mcp-server-v1.1.0-linux-arm64.tar.gz` | 2.4 MB |
| macOS | Intel (AMD64) | `nhrl-mcp-server-v1.1.0-darwin-amd64.tar.gz` | 2.7 MB |
| macOS | Apple Silicon (ARM64) | `nhrl-mcp-server-v1.1.0-darwin-arm64.tar.gz` | 2.5 MB |
| Windows | AMD64 | `nhrl-mcp-server-v1.1.0-windows-amd64.zip` | 2.7 MB |
| Windows | ARM64 | `nhrl-mcp-server-v1.1.0-windows-arm64.zip` | 2.4 MB |

### 🔐 SHA256 Checksums
```
00195bab3688253c8e7dc5474f4783aa2fa44c3269694f119c368fa5dc5254c4  nhrl-mcp-server-v1.1.0-linux-arm64.tar.gz
1bcce413d39424484aca89da6177bffeb7e21a9d869a79ebc47252360f691d48  nhrl-mcp-server-v1.1.0-darwin-amd64.tar.gz
45e260ec51ed266240e99391a7c0ecdcd1d475da733812d39981d5b64cb42d85  nhrl-mcp-server-v1.1.0-windows-amd64.zip
2c2c322e39a926e4924ecae7c2c2e5d5f1018d250a0bc9d4841d49558420a63c  nhrl-mcp-server-v1.1.0-darwin-arm64.tar.gz
a6cbbdf852b7d7a44140a49289bed437020416a50e17d33b617a7f4c84fbf7eb  nhrl-mcp-server-v1.1.0-windows-arm64.zip
a24ea6f9dc423e669043cf25faf9efb14cf7591ef1dcbf15bc1b5d03b3477ca9  nhrl-mcp-server-v1.1.0-linux-amd64.tar.gz
```

### 💔 Breaking Changes
⚠️ **Important**: This release contains breaking changes for existing users:

1. **Binary Name**: Changed from `truefinals-mcp-server` to `nhrl-mcp-server`
2. **Environment Variables**: Updated from `TRUEFINALS_*` to `NHRL_*`
3. **Server Identity**: Updated server name and branding

### 🔄 Migration Guide
If upgrading from a previous version:

1. Update your environment variables:
   ```bash
   # Old
   export TRUEFINALS_API_USER_ID=your_id
   export TRUEFINALS_API_KEY=your_key
   
   # New
   export NHRL_API_USER_ID=your_id
   export NHRL_API_KEY=your_key
   ```

2. Update binary references:
   ```bash
   # Old
   truefinals-mcp-server
   
   # New
   nhrl-mcp-server
   ```

3. Update any scripts or configurations that reference the old names.

### 🏗️ Build Information
- **Go Version**: 1.24.3
- **Build Date**: July 2, 2025
- **Signing**: Developer ID Application: Austin McChord (7PTN7E8EDS)
- **Notarization**: Apple notarized (Submission IDs: 6c9a9d7d-b1a0-4005-89fb-aa95df512dc6, 995399de-78fc-453b-9ee8-1181e75ab6e9)

### 📞 Support & Documentation
- GitHub Repository: https://github.com/austinmcchord/NHRL-MCP
- Issues: https://github.com/austinmcchord/NHRL-MCP/issues
- Documentation: See README.md and NHRL_INTEGRATION.md

---

**Full Changelog**: https://github.com/austinmcchord/NHRL-MCP/compare/v1.0.0...v1.1.0 