# NHRL MCP Server v1.7.1 Release Notes

## 🔧 Patch Release

### Bug Fixes
- Fixed permissions for `get_live_fight_stats` and `get_bot_picture_url` operations - now properly categorized as read-only operations
- These operations now work correctly in `reporting` mode

### Technical Details
- Updated read operations list to include the new v1.6 and v1.7 features
- Ensures proper access control for read-only modes
- No breaking changes

### 📦 Downloads

All binaries are code-signed and notarized for macOS (Notarization IDs: a250d0b6-650a-4d43-839f-c1775af1a09b, 5a233ae4-e4f2-4317-909e-fbdee8bf41e2).

| Platform | Architecture | File |
|----------|-------------|------|
| Linux | AMD64 | `nhrl-mcp-server-v1.7.1-linux-amd64.tar.gz` |
| Linux | ARM64 | `nhrl-mcp-server-v1.7.1-linux-arm64.tar.gz` |
| macOS | Intel | `nhrl-mcp-server-v1.7.1-darwin-amd64.tar.gz` |
| macOS | Apple Silicon | `nhrl-mcp-server-v1.7.1-darwin-arm64.tar.gz` |
| Windows | AMD64 | `nhrl-mcp-server-v1.7.1-windows-amd64.zip` |
| Windows | ARM64 | `nhrl-mcp-server-v1.7.1-windows-arm64.zip` |

---

**Full Changelog**: https://github.com/amcchord/NHRL-MCP/compare/v1.7.0...v1.7.1 