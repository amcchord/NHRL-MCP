# NHRL MCP Server v1.7.0 Release Notes

## 🚀 What's New in v1.7.0

### 🖼️ Bot Picture Integration
This release introduces comprehensive **Bot Picture Support**, providing access to robot images throughout the MCP server, along with enhanced data enrichment for tournament information.

#### Key Features:
- **Bot Picture URLs**: New `get_bot_picture_url` operation to retrieve robot images from BrettZone
- **Automatic Enrichment**: Bot pictures are now automatically included in player data across all tools
- **Smart URL Handling**: Bot names with spaces are automatically converted to underscores for compatibility
- **Multiple Sizes**: Both thumbnail and full-size images are available
- **Widespread Integration**: Pictures appear in tournament data, live stats, and direct queries

### 🎯 Enhanced Data Enrichment

#### Tournament Data Improvements:
- **NHRL Stats Integration**: Tournament match data now includes enriched player information with comprehensive NHRL statistics
- **Qualification Context**: Automatic qualification path context added to tournament matches
- **Match Review URLs**: Easy access to match video replays with generated review URLs
- **Complete Player Profiles**: Each player entry now includes:
  - Bot picture URLs (thumbnail and full size)
  - Performance statistics
  - Historical data
  - Qualification information

### 📸 Bot Picture Functionality

#### New Operation: `get_bot_picture_url`
```json
{
  "operation": "get_bot_picture_url",
  "bot_name": "Lynx"
}
```

Returns:
```json
{
  "bot_name": "Lynx",
  "thumbnail_url": "https://brettzone.com/pics/bots/thumbnail/Lynx.jpg",
  "full_url": "https://brettzone.com/pics/bots/Lynx.jpg"
}
```

#### Automatic Picture Enrichment
Bot pictures are now automatically included in:
- **Tournament Players**: All player entries in tournament data
- **Live Fight Stats**: Both competing bots include picture URLs
- **Player Queries**: Any tool that returns player information

### 🔧 Technical Improvements

- **BrettZone Integration**: Direct connection to BrettZone's robot image repository
- **URL Normalization**: Automatic handling of special characters and spaces in bot names
- **Fallback Support**: Graceful handling when images are not available
- **Performance Optimized**: Minimal overhead for picture URL generation
- **Backward Compatible**: All existing functionality remains intact

### 🔒 Security Features
- **Code Signing**: macOS binaries signed with Developer ID
- **Notarization**: Apple-notarized for Gatekeeper compatibility (IDs: bb24da27-8576-4dcc-ae87-50b987dad4e3, bf116dd5-2038-4ebe-a3cc-71a646f520ab)
- **API Authentication**: Secure API key authentication
- **Tools Filtering**: Maintains support for all permission modes

### 📦 Download Information

| Platform | Architecture | File | Size |
|----------|-------------|------|------|
| Linux | AMD64 | `nhrl-mcp-server-v1.7.0-linux-amd64.tar.gz` | 2.6 MB |
| Linux | ARM64 | `nhrl-mcp-server-v1.7.0-linux-arm64.tar.gz` | 2.4 MB |
| macOS | Intel (AMD64) | `nhrl-mcp-server-v1.7.0-darwin-amd64.tar.gz` | 2.7 MB |
| macOS | Apple Silicon (ARM64) | `nhrl-mcp-server-v1.7.0-darwin-arm64.tar.gz` | 2.5 MB |
| Windows | AMD64 | `nhrl-mcp-server-v1.7.0-windows-amd64.zip` | 2.7 MB |
| Windows | ARM64 | `nhrl-mcp-server-v1.7.0-windows-arm64.zip` | 2.4 MB |

### 🔐 SHA256 Checksums
```
cade56e063d557ef8cdd44feae93e38005734deccce787aeaefede0045928dd8  nhrl-mcp-server-v1.7.0-darwin-amd64.tar.gz
eb92693de694214a292d5876c1bb9d0e7349342b3136e6127281e65e0a899c0f  nhrl-mcp-server-v1.7.0-linux-amd64.tar.gz
32bb01d5f4471895bb9df984b2db42f8c6d2306738640b16d0e96c59b91ed175  nhrl-mcp-server-v1.7.0-windows-amd64.zip
25c549086907709d0dc13454fb19e378b142871af240a63bfe175cab8808daf8  nhrl-mcp-server-v1.7.0-darwin-arm64.tar.gz
61e9d0e941e642945eae9169b6deba4d371c524eb1cd75a683582e804346d126  nhrl-mcp-server-v1.7.0-windows-arm64.zip
9e12876c1849fc0297d7e26f9bbee1148ffe5b79204ecdb6b7e8c3b42b449325  nhrl-mcp-server-v1.7.0-linux-arm64.tar.gz
```

### 🎯 Use Cases

#### Tournament Broadcasting
```bash
# Enhanced tournament data with bot pictures for streaming overlays
nhrl-mcp-server --tools full
# Access tournament data with automatic bot picture enrichment
```

#### Fan Engagement
```bash
# Create fan sites with robot galleries and statistics
nhrl-mcp-server --tools reporting
# Query bot pictures and stats without modification capabilities
```

#### Event Production
```bash
# Complete data access for production graphics
nhrl-mcp-server --tools full
# Generate graphics with bot images, stats, and match history
```

### 💡 Enhanced Data Examples

#### Tournament Player Data Now Includes:
```json
{
  "name": "Lynx",
  "picture": {
    "thumbnail_url": "https://brettzone.com/pics/bots/thumbnail/Lynx.jpg",
    "full_url": "https://brettzone.com/pics/bots/Lynx.jpg"
  },
  "stats": { /* NHRL statistics */ },
  "qualification_path": "Qualified through...",
  "match_review_url": "https://..."
}
```

#### Live Fight Stats Enhanced:
Both `bot1` and `bot2` in live fight stats now include picture URLs alongside their existing statistics and driver information.

### 🚨 Breaking Changes
No breaking changes in this release. All existing functionality remains intact with added enrichment.

### 🔄 Migration Guide
No migration required. Simply update to v1.7.0 to access the enhanced features:
1. Download and install v1.7.0
2. Use existing configuration and credentials
3. Bot pictures and enriched data are automatically included

### 📊 Data Enrichment Details

#### Automatic Enrichments:
- **Bot Pictures**: Added to all player/bot references
- **NHRL Stats**: Integrated into tournament player data
- **Qualification Paths**: Context for how players qualified
- **Match URLs**: Direct links to match review videos
- **Complete Profiles**: Unified player information across all tools

#### URL Handling:
- Spaces in bot names → underscores (e.g., "Upper Cut" → "Upper_Cut")
- Special characters handled gracefully
- Both thumbnail (smaller) and full-size images available
- Consistent URL structure across all operations

### 🏗️ Build Information
- **Go Version**: 1.24.3
- **Build Date**: July 2, 2025
- **Signing**: Developer ID Application: Austin McChord (7PTN7E8EDS)
- **Notarization**: Apple notarized
  - AMD64: bb24da27-8576-4dcc-ae87-50b987dad4e3
  - ARM64: bf116dd5-2038-4ebe-a3cc-71a646f520ab

### 📞 Support & Documentation
- GitHub Repository: https://github.com/amcchord/NHRL-MCP
- Issues: https://github.com/amcchord/NHRL-MCP/issues
- Documentation: See README.md and NHRL_INTEGRATION.md

### 🙏 Acknowledgments
Special thanks to:
- **BrettZone** for providing the comprehensive robot image repository
- **NHRL Community** for feedback on data enrichment needs
- **Contributors** who suggested the bot picture integration

---

**Full Changelog**: https://github.com/amcchord/NHRL-MCP/compare/v1.6.0...v1.7.0 