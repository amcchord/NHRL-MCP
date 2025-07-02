# NHRL MCP Server v1.6.0 Release Notes

## 🚀 What's New in v1.6.0

### 📊 Live Fight Stats Integration
This release introduces **Live Fight Stats** functionality, providing real-time head-to-head matchup data and comprehensive robot statistics during NHRL events.

#### Key Features:
- **Live Fight Statistics**: New `get_live_fight_stats` operation in the `nhrl_stats` tool
- **Head-to-Head Analysis**: Detailed matchup data between two robots
- **Driver Information**: Access driver details including name, pronunciation, location, and pronouns
- **Bot Specifications**: Comprehensive robot specifications and builder background
- **Performance Metrics**: Overall statistics and historical performance data
- **Historical Meetings**: Track previous encounters between robots

### 🛠️ Enhanced Capabilities

#### Live Fight Stats Details:
The new `get_live_fight_stats` operation provides:
- **Driver Details**: Name, pronunciation guide, location, preferred pronouns
- **Bot Specifications**: Technical details and builder background information
- **Overall Performance**: Win/loss records and tournament statistics
- **Head-to-Head Records**: Specific matchup history between two robots
- **Historical Context**: Previous meeting details and outcomes

### 📋 Enhanced `nhrl_stats` Tool
The `nhrl_stats` tool now includes the `get_live_fight_stats` operation:
```json
{
  "operation": "get_live_fight_stats",
  "bot1": "MegatRON",
  "bot2": "Hurricane",
  "tournament_id": "nhrl_june25_30lb"
}
```

### 🔧 Technical Improvements

- **New API Endpoint**: Integration with `https://stats.nhrl.io/live_stats/query/get_fight_stats.php`
- **Enhanced Data Structures**: Added `NHRLLiveFightStats` struct for comprehensive stats modeling
- **POST Form Data**: Implemented proper form data submission for live stats queries
- **More Current Data**: Provides fresher statistics than traditional statsbook queries
- **Additional Metadata**: Includes driver and bot details not available in standard queries
- **Backward Compatible**: All existing functionality remains intact

### 🔒 Security Features
- **Code Signing**: macOS binaries signed with Developer ID
- **Notarization**: Apple-notarized for Gatekeeper compatibility (IDs: 9e1b2742-5155-43c2-9539-c815589ea274, 403afb2c-a432-4943-a5e0-e38eccf7f647)
- **API Authentication**: Secure API key authentication
- **Tools Filtering**: Maintains support for all permission modes (reporting, full-safe, full)

### 📦 Download Information

| Platform | Architecture | File | Size |
|----------|-------------|------|------|
| Linux | AMD64 | `nhrl-mcp-server-v1.6.0-linux-amd64.tar.gz` | 2.6 MB |
| Linux | ARM64 | `nhrl-mcp-server-v1.6.0-linux-arm64.tar.gz` | 2.4 MB |
| macOS | Intel (AMD64) | `nhrl-mcp-server-v1.6.0-darwin-amd64.tar.gz` | 2.7 MB |
| macOS | Apple Silicon (ARM64) | `nhrl-mcp-server-v1.6.0-darwin-arm64.tar.gz` | 2.5 MB |
| Windows | AMD64 | `nhrl-mcp-server-v1.6.0-windows-amd64.zip` | 2.7 MB |
| Windows | ARM64 | `nhrl-mcp-server-v1.6.0-windows-arm64.zip` | 2.4 MB |

### 🔐 SHA256 Checksums
```
0604e15ed0050a7a91f77d231574cd70c45240c3188f6d928c724ad49b520f0a  nhrl-mcp-server-v1.6.0-darwin-arm64.tar.gz
882299d5d59a417ba5af079b844c8fc9c866f523d35c85ff35cf4a01b5b40425  nhrl-mcp-server-v1.6.0-linux-amd64.tar.gz
820121228634c764f67fc4ee6b471496e2bb17081e38ad3291ddfc3b46c2b2c0  nhrl-mcp-server-v1.6.0-darwin-amd64.tar.gz
55b0f48c5284b51c8505465b686eaf4da23cbb613829e9567cb5250e65fe3feb  nhrl-mcp-server-v1.6.0-windows-amd64.zip
615b9031fcedc15528518ac08a3874ac50dc03541c9d78a37f55b2d561463ecf  nhrl-mcp-server-v1.6.0-linux-arm64.tar.gz
4815a1c444457e27d3e7c710633f82e41e717d4bb95f01bc0390e165ca5fd634  nhrl-mcp-server-v1.6.0-windows-arm64.zip
```

### 🎯 Use Cases

#### Live Tournament Monitoring
```bash
# Monitor live tournament data
nhrl-mcp-server --tools full
# Then use the nhrl_stats tool to access live data
```

#### Event Reporting
```bash
# Perfect for live event dashboards and reporting
nhrl-mcp-server --tools reporting
# Access read-only live stats without modification capabilities
```

#### Tournament Administration
```bash
# Full access for tournament management with live data
nhrl-mcp-server --tools full
```

### 💡 Example Usage

Once the server is running, you can:
- Query current tournament standings
- Get real-time match results
- Monitor robot performance metrics
- Track bracket progression
- Access historical tournament data
- Generate live statistics reports

### 🚨 Breaking Changes
No breaking changes in this release. All existing functionality remains intact.

### 🔄 Migration Guide
No migration required. Simply update to v1.6.0 to access the new live stats functionality:
1. Download and install v1.6.0
2. Use existing configuration and credentials
3. Access live stats through the `nhrl_stats` tool

### 📊 Live Fight Stats API
The new `get_live_fight_stats` operation:
- **Endpoint**: `https://stats.nhrl.io/live_stats/query/get_fight_stats.php`
- **Method**: POST with form-encoded data
- **Required Parameters**: 
  - `bot1`: First robot name
  - `bot2`: Second robot name (stats returned for this bot)
  - `tournament_id`: Tournament identifier
- **Returns**: Comprehensive statistics for `bot2` including head-to-head data against `bot1`
- **Use Case**: Perfect for pre-match analysis and live commentary support

### 🏗️ Build Information
- **Go Version**: 1.24.3
- **Build Date**: July 2, 2025
- **Signing**: Developer ID Application: Austin McChord (7PTN7E8EDS)
- **Notarization**: Apple notarized
  - AMD64: 9e1b2742-5155-43c2-9539-c815589ea274
  - ARM64: 403afb2c-a432-4943-a5e0-e38eccf7f647

### 📞 Support & Documentation
- GitHub Repository: https://github.com/amcchord/NHRL-MCP
- Issues: https://github.com/amcchord/NHRL-MCP/issues
- Documentation: See README.md and NHRL_INTEGRATION.md

### 🙏 Acknowledgments
Special thanks to the NHRL team for providing access to live data feeds and supporting the integration of real-time statistics.

### Enhanced Data Enrichment
- Tournament match data now includes enriched player information with NHRL stats
- Added automatic qualification path context to tournament matches
- Match review URL generation for easy video access

### New Bot Picture Functionality
- Added `get_bot_picture_url` operation to retrieve bot images from BrettZone
- Automatically enriches player data with bot picture URLs (thumbnail and full size)
- Bot names with spaces are automatically converted to underscores for URL compatibility
- Picture URLs are included in:
  - Direct bot picture queries
  - Tournament player data (enriched automatically)
  - Live fight stats for both competing bots

## API Updates

---

**Full Changelog**: https://github.com/amcchord/NHRL-MCP/compare/v1.5.0...v1.6.0 