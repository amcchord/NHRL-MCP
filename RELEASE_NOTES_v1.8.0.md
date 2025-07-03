# NHRL MCP Server v1.8.0 Release Notes

## 🚀 What's New in v1.8.0

### New Features & Improvements
This release includes various enhancements and improvements to the NHRL MCP Server, building upon the solid foundation established in previous versions.

### Key Highlights
- **Performance Optimizations**: Improved server performance and response times
- **Enhanced Stability**: Various bug fixes and stability improvements
- **Code Quality**: Internal code improvements and refactoring
- **Dependencies**: Updated dependencies for better security and performance

### 🔒 Security Features
- **Code Signing**: macOS binaries signed with Developer ID
- **API Authentication**: Secure API key authentication
- **Tools Filtering**: Full support for all permission modes (reporting, full-safe, full)

### 📦 Download Information

| Platform | Architecture | File | Size |
|----------|-------------|------|------|
| Linux | AMD64 | `nhrl-mcp-server-v1.8.0-linux-amd64.tar.gz` | 2.6 MB |
| Linux | ARM64 | `nhrl-mcp-server-v1.8.0-linux-arm64.tar.gz` | 2.4 MB |
| macOS | Intel (AMD64) | `nhrl-mcp-server-v1.8.0-darwin-amd64.tar.gz` | 2.7 MB |
| macOS | Apple Silicon (ARM64) | `nhrl-mcp-server-v1.8.0-darwin-arm64.tar.gz` | 2.5 MB |
| Windows | AMD64 | `nhrl-mcp-server-v1.8.0-windows-amd64.zip` | 2.7 MB |
| Windows | ARM64 | `nhrl-mcp-server-v1.8.0-windows-arm64.zip` | 2.4 MB |

### 🔐 SHA256 Checksums
```
054467d93a2d8be1987dcf3f88bffc8ecf34c6fb4f14c3e5ab30686363fc9cad  nhrl-mcp-server-v1.8.0-darwin-arm64.tar.gz
08617966e06a868cfdab287d9417c18bb5676f308ca208653971252e8248e1ee  nhrl-mcp-server-v1.8.0-linux-amd64.tar.gz
e6dfc10f202215909a9387b93c10ccaeeabeb9392032823aaf46f0b6005af42c  nhrl-mcp-server-v1.8.0-windows-arm64.zip
1044fb2b6b7a5f959b695e1333a24387d272583dec4c04860123426a5dc53624  nhrl-mcp-server-v1.8.0-windows-amd64.zip
05d1ac265c4259414e46bed6098005a53cac53cd299846da021e69c86efcf334  nhrl-mcp-server-v1.8.0-darwin-amd64.tar.gz
33c220475bf6a667285077da63dc1e078dd9994e20fd8587794fe6b411460b9e  nhrl-mcp-server-v1.8.0-linux-arm64.tar.gz
```

### 📋 Available Tools
The server continues to provide comprehensive access to NHRL data:
- `nhrl_stats` - NHRL statistics and live fight stats
- `truefinals_tournaments` - Tournament information
- `truefinals_games` - Game details and results
- `truefinals_locations` - Venue and location data
- `truefinals_players` - Player profiles and statistics
- `truefinals_bracket` - Tournament bracket management

### 💡 Usage Examples

```bash
# Standard operation with read-only mode
nhrl-mcp-server --tools reporting

# Full access mode
nhrl-mcp-server --tools full

# Using environment variables
export NHRL_API_USER_ID=your_id
export NHRL_API_KEY=your_key
nhrl-mcp-server
```

### 🚨 Breaking Changes
No breaking changes in this release. All existing functionality remains intact.

### 🔄 Migration Guide
No migration required. Simply update to v1.8.0:
1. Download the appropriate binary for your platform
2. Replace your existing installation
3. Restart the server with your existing configuration

### 🏗️ Build Information
- **Go Version**: 1.24.3
- **Build Date**: July 2, 2025
- **Signing**: Developer ID Application: Austin McChord (7PTN7E8EDS)
- **macOS**: Code-signed for improved security

### 📞 Support & Documentation
- GitHub Repository: https://github.com/amcchord/NHRL-MCP
- Issues: https://github.com/amcchord/NHRL-MCP/issues
- Documentation: See README.md and NHRL_INTEGRATION.md

### 🙏 Acknowledgments
Thank you to the NHRL community for continued feedback and support that helps make each release better.

---

**Full Changelog**: https://github.com/amcchord/NHRL-MCP/compare/v1.7.1...v1.8.0 