# NHRL MCP Server v1.5.0 Release Notes

## 🚀 What's New in v1.5.0

### 🔒 Read-Only Mode
The major feature in this release is the introduction of **Read-Only Mode**, providing enhanced security and control over MCP server operations.

#### Key Features:
- **Tools Filtering System**: New flexible permission system with three modes:
  - `reporting` - Read-only access to data (perfect for reporting and analytics)
  - `full-safe` - Full access except dangerous operations (default mode)
  - `full` - Complete unrestricted access to all operations
  
- **Granular Operation Control**: Fine-grained control over which operations are allowed:
  - Read operations are always allowed in reporting mode
  - Dangerous operations (delete, disqualify, reset) are blocked in full-safe mode
  - Full mode allows all operations for administrative tasks

- **Tool Disabling**: Individual tools can be disabled regardless of the mode using the `--disabled-tools` flag

### 🛠️ Configuration Options

#### Tools Mode Configuration
```bash
# Set read-only mode (reporting)
nhrl-mcp-server --tools reporting

# Set full-safe mode (default)
nhrl-mcp-server --tools full-safe

# Set full access mode
nhrl-mcp-server --tools full

# Or use environment variable
export NHRL_TOOLS=reporting
nhrl-mcp-server
```

#### Disable Specific Tools
```bash
# Disable specific tools
nhrl-mcp-server --disabled-tools "truefinals_bracket,truefinals_players"

# Or use environment variable
export NHRL_DISABLED_TOOLS="truefinals_bracket,truefinals_players"
nhrl-mcp-server
```

### 📋 Tool Permissions by Mode

| Mode | Read Operations | Write Operations | Delete Operations |
|------|----------------|------------------|-------------------|
| `reporting` | ✅ Allowed | ❌ Blocked | ❌ Blocked |
| `full-safe` | ✅ Allowed | ✅ Allowed | ❌ Blocked |
| `full` | ✅ Allowed | ✅ Allowed | ✅ Allowed |

### 🔧 Technical Improvements

- **Enhanced Security**: Prevent accidental data modifications with read-only mode
- **Flexible Permissions**: Choose the right level of access for each use case
- **Backward Compatible**: Default mode remains `full-safe` for existing users
- **Environment Variable Support**: All new flags support environment variables

### 🔒 Security Features
- **Code Signing**: macOS binaries signed with Developer ID
- **Notarization**: Apple-notarized for Gatekeeper compatibility
- **Enhanced Tools Filtering**: New permission system for safer operations
- **API Authentication**: Secure API key authentication

### 📦 Download Information

| Platform | Architecture | File | Size |
|----------|-------------|------|------|
| Linux | AMD64 | `nhrl-mcp-server-v1.5.0-linux-amd64.tar.gz` | 2.6 MB |
| Linux | ARM64 | `nhrl-mcp-server-v1.5.0-linux-arm64.tar.gz` | 2.4 MB |
| macOS | Intel (AMD64) | `nhrl-mcp-server-v1.5.0-darwin-amd64.tar.gz` | 2.7 MB |
| macOS | Apple Silicon (ARM64) | `nhrl-mcp-server-v1.5.0-darwin-arm64.tar.gz` | 2.5 MB |
| Windows | AMD64 | `nhrl-mcp-server-v1.5.0-windows-amd64.zip` | 2.7 MB |
| Windows | ARM64 | `nhrl-mcp-server-v1.5.0-windows-arm64.zip` | 2.4 MB |

### 🔐 SHA256 Checksums
```
8e1add2afde7ddcfb0f390d7f8550f7157a93fbb0946c72fd12178a6370b8465  nhrl-mcp-server-v1.5.0-darwin-amd64.tar.gz
67405d2a40ea3cdf2586e6fac6f75e87c09aff492c647e52add6cd18433a2f2a  nhrl-mcp-server-v1.5.0-linux-arm64.tar.gz
8e67e060e6aaa0328b9a003d3a3945552784c288dab5beb26997f581d4b03b07  nhrl-mcp-server-v1.5.0-linux-amd64.tar.gz
7fd9eda384e27c6da9633f79b8aa1161fe3a5c99a23b13fe9aea361df6404198  nhrl-mcp-server-v1.5.0-windows-amd64.zip
1e2b34e94da8a187570f1e286b54fae45ed65cf1286eaf89c2fdd938dc81b14a  nhrl-mcp-server-v1.5.0-darwin-arm64.tar.gz
fe525f6eff95137512bc6661878df71e5d9d272ac73239fe8417aa36092af66a  nhrl-mcp-server-v1.5.0-windows-arm64.zip
```

### 🎯 Use Cases

#### Reporting & Analytics
```bash
# Perfect for read-only dashboards and reporting tools
nhrl-mcp-server --tools reporting
```

#### Standard Operations
```bash
# Default mode - allows most operations but blocks dangerous ones
nhrl-mcp-server --tools full-safe
```

#### Administrative Access
```bash
# Full access for administrative tasks
nhrl-mcp-server --tools full
```

### 🚨 Breaking Changes
No breaking changes in this release. The default behavior remains the same (`full-safe` mode).

### 🔄 Migration Guide
No migration required. Existing installations will continue to work with the default `full-safe` mode.

To take advantage of the new read-only mode:
1. Update to v1.5.0
2. Add `--tools reporting` flag or set `NHRL_TOOLS=reporting`
3. Optionally disable specific tools with `--disabled-tools`

### 📋 Available Operations by Mode

#### Read-Only Operations (allowed in all modes):
- `get`, `list`, `details`, `format`
- `overlay_params`, `description`
- `private`, `webhooks`
- `get_round`, `get_standings`

#### Dangerous Operations (blocked in `reporting` and `full-safe`):
- `delete` - Remove data
- `disqualify` - Disqualify participants
- `reset` - Reset tournament state

### 🏗️ Build Information
- **Go Version**: 1.24.3
- **Build Date**: July 2, 2025
- **Signing**: Developer ID Application: Austin McChord (7PTN7E8EDS)
- **Notarization**: Apple notarized
  - AMD64: f1f246d2-047f-4259-9b21-de4c9e7eeaab
  - ARM64: 207a08c0-550a-4d28-bb61-9e591ecd6f16

### 📞 Support & Documentation
- GitHub Repository: https://github.com/amcchord/NHRL-MCP
- Issues: https://github.com/amcchord/NHRL-MCP/issues
- Documentation: See README.md and NHRL_INTEGRATION.md

---

**Full Changelog**: https://github.com/amcchord/NHRL-MCP/compare/v1.1.0...v1.5.0 