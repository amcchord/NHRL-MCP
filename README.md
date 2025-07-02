# TrueFinals MCP Server

A Model Context Protocol (MCP) server that provides tools for managing tournaments, games, locations, and players through the TrueFinals API.

## Overview

This MCP server enables AI assistants to interact with the TrueFinals tournament management platform, providing comprehensive tools for tournament administration. It supports all major tournament operations including bracket management, game scoring, location coordination, and player management.

## Features

### 🏆 Tournament Management
- Create, update, and delete tournaments
- Start and reset tournaments
- Manage tournament settings and format options
- Handle tournament webhooks and overlay parameters
- Support for multiple tournament formats (single elimination, double elimination, round robin)

### 🎮 Game Management
- List and manage all tournament games
- Create and edit exhibition games
- Update game scores and states
- Handle game check-ins and scheduling
- Support for bulk game operations
- Location assignment and management

### 📍 Location Management
- Add, update, and delete tournament locations
- Start and stop games at specific locations
- Manage location queues and game assignments
- Update game scores directly from locations
- Handle location blocking and availability

### 👥 Player Management
- Add, update, and delete players
- Manage player seeding and rankings
- Handle player check-ins and disqualifications
- Support for bulk player operations
- Randomize tournament seeding
- Track player statistics and placement

## Available Tools

### 1. Tournaments Tool (12 operations)
- `list` - Get user's tournaments
- `get` - Get tournament details
- `create` - Create new tournament
- `update` - Update tournament settings
- `delete` - Delete tournament
- `start` - Start tournament
- `reset` - Reset tournament
- `get_webhooks` - Get tournament webhooks
- `update_webhooks` - Update webhooks
- `get_overlay_params` - Get overlay parameters
- `update_overlay_params` - Update overlay parameters
- `push_schedule` - Push game schedule

### 2. Games Tool (13 operations)
- `list` - Get all tournament games
- `get` - Get specific game details
- `add_exhibition` - Add exhibition game
- `edit_exhibition` - Edit exhibition game
- `delete_exhibition` - Delete exhibition game
- `bulk_add_exhibition` - Bulk add exhibition games
- `bulk_delete_exhibition` - Bulk delete exhibition games
- `update` - Update game details
- `update_score` - Update game score
- `update_state` - Update game state
- `update_scheduled_time` - Update scheduled time
- `update_location` - Update game location
- `update_checkin` - Update player check-in
- `undo` - Undo completed game

### 3. Locations Tool (8 operations)
- `list` - Get all tournament locations
- `get` - Get specific location details
- `add` - Add new location
- `update` - Update location details
- `delete` - Delete location
- `start_game` - Start queued game at location
- `stop_game` - Stop active game at location
- `update_game_scores` - Update scores for location's active game

### 4. Players Tool (10 operations)
- `list` - Get all tournament players
- `get` - Get specific player details
- `add` - Add new player
- `update` - Update player information
- `delete` - Delete player
- `reseed` - Change player seeding
- `randomize` - Randomize tournament seeding
- `bulk_update` - Bulk update player list
- `checkin` - Check player into match
- `disqualify` - Disqualify player

## Installation and Setup

### Prerequisites
- Go 1.19 or later
- TrueFinals API credentials (User ID and API Key)

### Building the Server

```bash
# Clone the repository
git clone <repository-url>
cd NHRL-MCP

# Build the server
go build -o truefinals-mcp-server

# Or use the pre-built binary
chmod +x truefinals-mcp-server
```

### Configuration

The server requires TrueFinals API credentials. Set them using environment variables:

```bash
export TRUEFINALS_API_USER_ID="your_user_id"
export TRUEFINALS_API_KEY="your_api_key"

# Optional configuration
export TRUEFINALS_BASE_URL="https://truefinals.com/api"  # Custom API base URL
export TRUEFINALS_TOOLS="full"                           # Tool filter mode
export TRUEFINALS_DISABLED_TOOLS="tournaments,games"     # Disable specific tools
```

Or pass them as command-line arguments:

```bash
./truefinals-mcp-server -api-user-id "your_user_id" -api-key "your_api_key"
```

### Usage

#### Basic Usage
```bash
# Start the server with default settings
./truefinals-mcp-server

# Start with specific credentials
./truefinals-mcp-server -api-user-id "your_user_id" -api-key "your_api_key"
```

#### Tool Filtering
The server supports filtering available tools based on safety and permission levels:

```bash
# Show only reporting tools (read-only operations)
./truefinals-mcp-server -tools reporting

# Show safe modification tools (excludes dangerous operations)
./truefinals-mcp-server -tools full-safe

# Show all tools (default)
./truefinals-mcp-server -tools full
```

#### Available Tool Modes:
- **`reporting`**: Read-only operations (list, get operations)
- **`full-safe`**: Safe modification operations (excludes delete, reset, bulk operations)
- **`full`**: All operations (default)

#### Command Line Options
```bash
Usage: ./truefinals-mcp-server [options]

Options:
  -api-user-id string     API User ID for TrueFinals service
  -api-key string         API key for TrueFinals service
  -base-url string        Base URL for TrueFinals API (optional)
  -tools string           Tools mode: reporting, full-safe, full
  -disabled-tools string  Comma-separated list of tool names to disable
  -exit-after-first       Exit after processing the first request
  -version               Show version information and exit
  -help                  Show help information
```

## MCP Integration

### Connecting to AI Assistants

The server implements the Model Context Protocol and can be integrated with compatible AI assistants like Claude Desktop or other MCP-compatible applications.

Example MCP configuration for Claude Desktop:
```json
{
  "mcpServers": {
    "truefinals": {
      "command": "/path/to/truefinals-mcp-server",
      "args": ["-api-user-id", "your_user_id", "-api-key", "your_api_key"],
      "env": {}
    }
  }
}
```

### Claude Desktop Configuration

To use with Claude Desktop, add this configuration to your Claude Desktop config file:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`  
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`  
**Linux**: `~/.config/Claude/claude_desktop_config.json`

Ready-to-use configuration (copy and paste this into your Claude Desktop config):
```json
{
  "mcpServers": {
    "truefinals": {
      "command": "/usr/local/bin/truefinals-mcp-server",
      "args": [
        "-api-user-id", "YOUR_TRUEFINALS_USER_ID",
        "-api-key", "YOUR_TRUEFINALS_API_KEY"
      ],
      "env": {}
    }
  }
}
```

Replace `YOUR_TRUEFINALS_USER_ID` and `YOUR_TRUEFINALS_API_KEY` with your actual TrueFinals API credentials.

### Tool Schemas

All tools include comprehensive JSON schemas with:
- Required and optional parameters
- Parameter validation rules
- Enum values for status fields
- Detailed parameter descriptions
- Type constraints and limits

## API Reference

The server interfaces with the TrueFinals API v1. Key data structures include:

### Tournament
- Tournament metadata and settings
- Format configuration (single/double elimination, round robin)
- Player and location management
- Game bracket organization

### Game
- Game state and scoring
- Player check-in status
- Location assignment
- Bracket progression

### Location
- Physical tournament locations
- Game queue management
- Active game tracking
- Availability status

### Player
- Player information and statistics
- Seeding and ranking
- Check-in status
- Match history

## Error Handling

The server provides comprehensive error handling with:
- Detailed error messages
- HTTP status code mapping
- Parameter validation
- API rate limiting awareness
- Network error recovery

## Development

### Project Structure
```
NHRL-MCP/
├── main.go                 # Server entry point and MCP protocol handling
├── api.go                  # TrueFinals API client and data structures
├── tools_tournaments.go    # Tournament management operations
├── tools_games.go         # Game management operations  
├── tools_locations.go     # Location management operations
├── tools_players.go       # Player management operations
├── docs/
│   └── openapi.json       # TrueFinals API specification
├── go.mod                 # Go module dependencies
└── README.md              # This file
```

### Adding New Operations

To add new operations:

1. Add the API endpoint to `api.go`
2. Implement the operation function in the appropriate `tools_*.go` file
3. Add the operation to the tool's handler switch statement
4. Update the tool schema with the new operation
5. Test the implementation

## Contributing

Contributions are welcome! Please ensure:
- Code follows Go conventions
- All operations include proper error handling
- Tool schemas are updated for new operations
- Documentation is updated accordingly

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues related to:
- **TrueFinals API**: Contact TrueFinals support
- **MCP Server**: Create an issue in this repository
- **Integration**: Check MCP documentation and AI assistant configuration

## Version History

- **v1.0.0**: Initial release with full tournament, game, location, and player management
  - 43 total operations across 4 tools
  - Complete TrueFinals API coverage
  - Tool filtering and safety modes
  - Comprehensive error handling and validation 