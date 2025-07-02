# NHRL MCP Server

A Model Context Protocol (MCP) server that provides comprehensive tools for combat robot tournament management and statistics. This server integrates with both TrueFinals tournament management platform and NHRL (National Havoc Robot League) statistics systems, offering AI assistants complete access to tournament administration and historical robot combat data.

## Overview

This MCP server enables AI assistants to interact with:

- **TrueFinals Tournament Management**: Complete tournament administration including bracket management, game scoring, location coordination, and player management
- **NHRL Statistics**: Access to comprehensive robot combat statistics, fight records, rankings, and live tournament data from the National Havoc Robot League

Whether you're managing a live tournament or analyzing historical robot combat performance, this server provides all the tools needed for comprehensive tournament and statistics management.

## Features

### 🏆 TrueFinals Tournament Management
- Create, update, and delete tournaments
- Start and reset tournaments  
- Manage tournament settings and format options
- Handle tournament webhooks and overlay parameters
- Support for multiple tournament formats (single elimination, double elimination, round robin)
- Complete game lifecycle management
- Location queue management and scoring
- Player seeding, check-ins, and disqualification handling
- Bracket visualization and progression tracking

### 📊 NHRL Statistics & Analytics
- **Bot Performance Data**: Rankings, fight records, win/loss statistics, KO records
- **Head-to-Head Analysis**: Detailed matchup history between any two bots
- **Seasonal Statistics**: Performance tracking across different NHRL seasons
- **Weight Class Analytics**: Podium finishers, fastest KOs, longest winning streaks
- **Live Tournament Data**: Real-time match information from BrettZone system
- **Video Integration**: Direct links to match review videos with timestamp control
- **Qualification System**: Understanding NHRL tournament qualification paths

## Available Tools

### 1. TrueFinals Tournaments Tool
**Tool Name**: `truefinals_tournaments`

**Operations** (12 total):
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

### 2. TrueFinals Games Tool
**Tool Name**: `truefinals_games`

**Operations** (13 total):
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

### 3. TrueFinals Locations Tool
**Tool Name**: `truefinals_locations`

**Operations** (8 total):
- `list` - Get all tournament locations
- `get` - Get specific location details
- `add` - Add new location
- `update` - Update location details
- `delete` - Delete location
- `start_game` - Start queued game at location
- `stop_game` - Stop active game at location
- `update_game_scores` - Update scores for location's active game

### 4. TrueFinals Players Tool
**Tool Name**: `truefinals_players`

**Operations** (10 total):
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

### 5. TrueFinals Bracket Tool
**Tool Name**: `truefinals_bracket`

**Operations** (3 total):
- `get_round` - Get specific bracket round details
- `get_standings` - Get current tournament standings
- `format` - Get bracket format information

### 6. NHRL Stats Tool ⭐ 
**Tool Name**: `nhrl_stats`

**Operations** (15 total):

#### Bot-Specific Operations:
- `get_bot_rank` - Get current bot ranking
- `get_bot_fights` - Get complete fight history for a bot
- `get_bot_head_to_head` - Get head-to-head records against all opponents
- `get_bot_stats_by_season` - Get seasonal performance statistics
- `get_bot_streak_stats` - Get current and longest win/lose streaks
- `get_bot_event_participants` - Get tournament participation history
- `get_live_fight_stats` - Get live fight statistics between two bots for a specific tournament

#### Weight Class Operations:
- `get_weight_class_dumpster_count` - Get podium finishers (1st, 2nd, 3rd place)
- `get_weight_class_event_winners` - Get event winners by weight class
- `get_weight_class_fastest_kos` - Get fastest knockout records
- `get_weight_class_longest_streaks` - Get longest winning streaks
- `get_weight_class_stat_summary` - Get comprehensive rankings and statistics

#### Tournament & System Operations:
- `get_random_fight` - Get a random historical fight
- `get_tournament_matches` - Get live tournament match data from BrettZone
- `get_match_review_url` - Generate video review URLs for specific matches
- `get_qualification_system` - Get information about NHRL qualification system

**Supported Weight Classes**: 3lb, 12lb, 30lb, beetleweight, antweight, hobbyweight
**Supported Seasons**: current, all-time, 2018-2019, 2020, 2021, 2022, 2023

## Installation and Setup

### Prerequisites
- Go 1.19 or later
- TrueFinals API credentials (User ID and API Key) for tournament management features

### Building the Server

```bash
# Clone the repository
git clone <repository-url>
cd NHRL-MCP

# Build the server
go build -o nhrl-mcp-server

# Or use the pre-built binary
chmod +x nhrl-mcp-server
```

### Configuration

#### For TrueFinals Features
The server requires TrueFinals API credentials for tournament management features:

```bash
export TRUEFINALS_API_USER_ID="your_user_id"
export TRUEFINALS_API_KEY="your_api_key"

# Optional TrueFinals configuration
export TRUEFINALS_BASE_URL="https://truefinals.com/api"  # Custom API base URL
export TRUEFINALS_TOOLS="full"                           # Tool filter mode
export TRUEFINALS_DISABLED_TOOLS="tournaments,games"     # Disable specific tools
export TRUEFINALS_READ_ONLY="true"                       # Enable read-only mode
```

#### For NHRL Features
NHRL statistics features work without additional configuration - they access public NHRL APIs directly.

### Usage

#### Basic Usage
```bash
# Start the server (NHRL stats work immediately, TrueFinals requires credentials)
./nhrl-mcp-server

# Start with TrueFinals credentials for full functionality
./nhrl-mcp-server -api-user-id "your_user_id" -api-key "your_api_key"
```

#### Tool Filtering
The server supports filtering available tools based on safety and permission levels:

```bash
# Show only reporting tools (read-only operations)
./nhrl-mcp-server -tools reporting

# Show safe modification tools (excludes dangerous operations)
./nhrl-mcp-server -tools full-safe

# Show all tools (default)
./nhrl-mcp-server -tools full
```

#### Read-Only Mode
For maximum safety, you can enable read-only mode which only allows read operations regardless of the tools mode:

```bash
# Enable read-only mode via CLI flag
./nhrl-mcp-server -read-only

# Enable read-only mode via environment variable
export TRUEFINALS_READ_ONLY=true
./nhrl-mcp-server

# Read-only mode works with any tools mode and overrides it
./nhrl-mcp-server -tools full -read-only  # Still only allows read operations
```

#### Available Tool Modes:
- **`reporting`**: Read-only operations (list, get operations) - safest mode
- **`full-safe`**: Safe modification operations (excludes delete, reset, disqualify operations)
- **`full`**: All operations including potentially destructive ones
- **Read-only mode**: Only read operations allowed (overrides all tool modes when enabled)

#### Command Line Options
```bash
Usage: ./nhrl-mcp-server [options]

Options:
  -api-user-id string     API User ID for TrueFinals service
  -api-key string         API key for TrueFinals service  
  -base-url string        Base URL for TrueFinals API (optional)
  -tools string           Tools mode: reporting, full-safe, full
  -disabled-tools string  Comma-separated list of tool names to disable
  -read-only              Enable read-only mode - only allow read operations
  -exit-after-first       Exit after processing the first request
  -version               Show version information and exit
  -help                  Show help information
```

## Claude Desktop Integration

### Installation

The easiest way to use this MCP server is with Claude Desktop. Here's the complete setup:

#### Step 1: Install the Server

```bash
# Download and install to system location
sudo cp nhrl-mcp-server /usr/local/bin/
sudo chmod +x /usr/local/bin/nhrl-mcp-server
```

#### Step 2: Configure Claude Desktop

Add this configuration to your Claude Desktop config file:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`  
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`  
**Linux**: `~/.config/Claude/claude_desktop_config.json`

### Configuration Examples

#### Option 1: NHRL Stats Only (No TrueFinals credentials needed)
```json
{
  "mcpServers": {
    "nhrl": {
      "command": "/usr/local/bin/nhrl-mcp-server",
      "args": ["-tools", "reporting"],
      "env": {}
    }
  }
}
```

#### Option 2: Full Functionality with TrueFinals (Recommended)
```json
{
  "mcpServers": {
    "nhrl": {
      "command": "/usr/local/bin/nhrl-mcp-server",
      "args": [
        "-api-user-id", "YOUR_TRUEFINALS_USER_ID",
        "-api-key", "YOUR_TRUEFINALS_API_KEY",
        "-tools", "full-safe"
      ],
      "env": {}
    }
  }
}
```

#### Option 3: Using Environment Variables
```json
{
  "mcpServers": {
    "nhrl": {
      "command": "/usr/local/bin/nhrl-mcp-server",
      "args": ["-tools", "full-safe"],
      "env": {
        "TRUEFINALS_API_USER_ID": "YOUR_TRUEFINALS_USER_ID",
        "TRUEFINALS_API_KEY": "YOUR_TRUEFINALS_API_KEY"
      }
    }
  }
}
```

#### Option 4: Read-Only Mode (Maximum Safety)
```json
{
  "mcpServers": {
    "nhrl": {
      "command": "/usr/local/bin/nhrl-mcp-server",
      "args": [
        "-api-user-id", "YOUR_TRUEFINALS_USER_ID",
        "-api-key", "YOUR_TRUEFINALS_API_KEY",
        "-read-only"
      ],
      "env": {}
    }
  }
}
```

#### Option 5: NHRL-Only Configuration
```json
{
  "mcpServers": {
    "nhrl-stats": {
      "command": "/usr/local/bin/nhrl-mcp-server",
      "args": [
        "-disabled-tools", "truefinals_tournaments,truefinals_games,truefinals_locations,truefinals_players,truefinals_bracket"
      ],
      "env": {}
    }
  }
}
```

**Replace** `YOUR_TRUEFINALS_USER_ID` and `YOUR_TRUEFINALS_API_KEY` with your actual TrueFinals API credentials if you want tournament management features.

### Example Usage with Claude

Once configured, you can ask Claude questions like:

**NHRL Statistics Examples:**
- "What's Tombstone's current ranking and fight record?"
- "Show me the head-to-head record between Minotaur and Witch Doctor"
- "What are the fastest KOs in the 3lb weight class?"
- "Get the current tournament matches for nhrl_june25_30lb"

**TrueFinals Tournament Management Examples:**
- "Create a new 3lb tournament for this weekend"
- "List all games in tournament T123 and their current status"
- "Update the score for game G456 to show Tombstone won 3-1"
- "Show me the current bracket standings"

## API Integration Details

### NHRL APIs
The server integrates with two NHRL systems:
- **NHRL Statsbook API**: `https://stats.nhrl.io/statsbook` - Historical statistics and rankings
- **BrettZone API**: `https://brettzone.nhrl.io/brettZone/backend` - Live tournament data

### TrueFinals API
- **TrueFinals API v1**: `https://truefinals.com/api` - Tournament management platform

## Data Structures

### NHRL Data Types
- **Bot Statistics**: Rankings, fight records, win/loss ratios, KO statistics
- **Fight Records**: Individual match data with timestamps, results, video links
- **Tournament Matches**: Live match data with round information and implications
- **Head-to-Head Records**: Detailed matchup statistics between specific bots
- **Qualification Rounds**: Q1, Q2W, Q2L, Q3 progression paths

### TrueFinals Data Types
- **Tournaments**: Tournament metadata, settings, and format configuration
- **Games**: Game state, scoring, player check-ins, location assignments
- **Locations**: Physical tournament locations with queue management
- **Players**: Player information, seeding, rankings, and match history
- **Brackets**: Tournament bracket structure and progression

## Error Handling & Reliability

The server provides comprehensive error handling with:
- Detailed error messages with context
- HTTP status code mapping
- Parameter validation with helpful feedback
- API rate limiting awareness
- Network error recovery
- Graceful degradation when services are unavailable

## Development

### Project Structure
```
NHRL-MCP/
├── main.go                 # Server entry point and MCP protocol handling
├── api.go                  # TrueFinals API client and data structures
├── nhrl_api.go            # NHRL API clients and data structures
├── tools_tournaments.go    # TrueFinals tournament management
├── tools_games.go         # TrueFinals game management  
├── tools_locations.go     # TrueFinals location management
├── tools_players.go       # TrueFinals player management
├── tools_bracket.go       # TrueFinals bracket operations
├── tools_nhrl.go          # NHRL statistics and live data
├── docs/
│   └── openapi.json       # TrueFinals API specification
├── scripts/               # Build and deployment scripts
├── go.mod                 # Go module dependencies
└── README.md              # This file
```

### Adding New Operations

To add new operations:

1. **For NHRL features**: Add API endpoint functions to `nhrl_api.go`, implement operation in `tools_nhrl.go`
2. **For TrueFinals features**: Add API endpoint to `api.go`, implement in appropriate `tools_*.go` file
3. Update the tool's operation enum in the schema
4. Add operation handler to the switch statement
5. Test the implementation thoroughly

## Contributing

Contributions are welcome! Please ensure:
- Code follows Go conventions and includes proper error handling
- Tool schemas are updated for new operations
- Both NHRL and TrueFinals functionality is properly tested
- Documentation is updated accordingly
- New features include appropriate safety considerations

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues related to:
- **NHRL Statistics**: Check NHRL statsbook API availability and bot name formatting
- **TrueFinals API**: Contact TrueFinals support or verify API credentials
- **MCP Server**: Create an issue in this repository
- **Claude Desktop Integration**: Check MCP documentation and configuration format

## Version History

- **v1.0.0**: Initial hybrid release
  - 6 integrated tools with 61 total operations
  - Complete TrueFinals tournament management (48 operations)
  - Comprehensive NHRL statistics integration (15 operations)
  - Advanced tool filtering and safety modes
  - Claude Desktop ready configuration
  - Real-time tournament data from BrettZone
  - Historical statistics from NHRL Statsbook
  - Video review integration with timestamp control

## Quick Start Examples

### Get Started with NHRL Stats (No Credentials Required)
```bash
# Install and test NHRL functionality immediately
./nhrl-mcp-server -disabled-tools "truefinals_tournaments,truefinals_games,truefinals_locations,truefinals_players,truefinals_bracket"
```

### Get Started with Full Tournament Management
```bash
# Full functionality (requires TrueFinals credentials)
./nhrl-mcp-server -api-user-id "your_user_id" -api-key "your_api_key" -tools full-safe
```

This MCP server provides the most comprehensive combat robot tournament and statistics management available for AI assistants, combining live tournament administration with deep historical analytics in a single, powerful tool.

### NHRL Stats Operations

The NHRL stats tool provides access to comprehensive NHRL statsbook data and BrettZone tournament information:

#### Bot-specific Operations:
- `get_bot_rank` - Get a bot's current ranking
- `get_bot_fights` - Get all fights for a specific bot
- `get_bot_head_to_head` - Get head-to-head records against all opponents
- `get_bot_stats_by_season` - Get seasonal statistics for a bot
- `get_bot_streak_stats` - Get winning/losing streak information
- `get_bot_event_participants` - Get all events a bot has participated in
- `get_live_fight_stats` - Get live fight statistics between two bots for a specific tournament

#### Weight Class Operations:
- `get_weight_class_dumpster_count` - Get podium finishes by weight class
- `get_weight_class_event_winners` - Get event winners for a weight class
- `get_weight_class_fastest_kos` - Get fastest knockouts in a weight class
- `get_weight_class_longest_streaks` - Get longest winning streaks
- `get_weight_class_stat_summary` - Get overall statistics summary

#### Tournament Operations:
- `get_tournament_matches` - Get all matches from a BrettZone tournament
- `get_match_review_url` - Generate a match review video URL
- `get_qualification_system` - Get information about NHRL's qualification system

#### Other Operations:
- `get_random_fight` - Get a random fight from the database

### New Feature: Live Fight Stats

The `get_live_fight_stats` operation provides detailed, up-to-date statistics for a specific bot with head-to-head data against an opponent in a tournament context.

**Parameters:**
- `bot1` (required): First bot name for comparison
- `bot2` (required): Second bot name - stats will be returned for this bot
- `tournament_id` (required): Tournament ID (e.g., 'nhrl_june25_30lb')

**Returns:**
- Detailed stats for bot2 including:
  - Driver information (name, pronunciation, location, pronouns)
  - Bot details (type, ranking, builder background)
  - Overall fight statistics (wins, losses, KOs, judge decisions)
  - Head-to-head record against bot1
  - Last meeting date if applicable

**Example usage:**
```json
{
  "operation": "get_live_fight_stats",
  "bot1": "MegatRON",
  "bot2": "Hurricane",
  "tournament_id": "nhrl_june25_30lb"
}
```

This provides more recent statistics than the standard statsbook queries and includes additional driver/builder information.

## Configuration

// ... existing code ... 