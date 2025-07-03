# NHRL MCP Server Integration

This MCP server has been enhanced to focus on NHRL (National Havoc Robot League) combat robot tournaments, integrating the [NHRL statsbook](https://stats.nhrl.io/statsbook/) API to provide comprehensive robot combat data.

## Features

### Existing TrueFinals Tools (Enhanced)
All existing TrueFinals tournament management tools remain available and have been **enriched** with NHRL data:

- **truefinals_tournaments** - Tournament management with NHRL context
- **truefinals_games** - Game management with bot stats
- **truefinals_locations** - Location management 
- **truefinals_players** - Player management enriched with NHRL bot rankings and stats
- **truefinals_bracket** - Bracket visualization

### Enhanced NHRL Stats Tool
**nhrl_stats** - Direct access to NHRL statsbook data and BrettZone live tournament system with the following operations:

#### Bot-Specific Operations
- `get_bot_rank` - Get current ranking for a specific bot
- `get_bot_fights` - Get complete fight history for a bot
- `get_bot_head_to_head` - Get head-to-head records against all opponents
- `get_bot_stats_by_season` - Get seasonal statistics for a bot
- `get_bot_streak_stats` - Get current win/lose streak information
- `get_bot_event_participants` - Get tournament participation history
- `get_bot_picture_url` - Get bot picture URLs (thumbnail and full size) from BrettZone

#### Weight Class Operations
- `get_weight_class_dumpster_count` - Get podium finishers (1st/2nd/3rd place) for a weight class
- `get_weight_class_event_winners` - Get historical tournament winners
- `get_weight_class_fastest_kos` - Get fastest KO records
- `get_weight_class_longest_streaks` - Get longest winning streaks
- `get_weight_class_stat_summary` - Get comprehensive rankings and stats

#### General Operations
- `get_random_fight` - Get a random fight for inspiration

#### BrettZone Live Tournament Operations
- `get_tournament_matches` - Get all matches from a live tournament with timing data, results, and review links
- `get_match_review_url` - Generate direct links to match review videos with custom timing and camera angles
- `get_qualification_system` - Get explanation of NHRL's tournament qualification system (with optional round_code parameter for specific round details)

## Tournament Qualification System

NHRL uses a qualification system for tournament entries. When you see round codes in match data, here's what they mean:

### Qualification Rounds

1. **OPENING (Q1)**: All competitors start here
   - Win → Advance to THE CUSP (Q2W)
   - Lose → Drop to REDEMPTION (Q2L)

2. **THE CUSP (Q2W)**: For Opening winners - one win away from qualifying
   - Win → QUALIFY FOR BRACKET
   - Lose → Drop to BUBBLE (Q3) for last chance

3. **REDEMPTION (Q2L)**: Second chance for Opening losers
   - Win → Advance to BUBBLE (Q3)
   - Lose → ELIMINATED from tournament

4. **BUBBLE (Q3)**: Final qualifying round - last chance to make the bracket
   - Win → QUALIFY FOR BRACKET
   - Lose → ELIMINATED from tournament

Once qualified, competitors enter the main single-elimination bracket (Round of 32, Round of 16, Quarterfinals, Semifinals, Finals).

### Enhanced Match Data

When you query tournament matches using `get_tournament_matches`, each match now includes:
- **roundName**: Human-readable name (e.g., "Opening" instead of "Q1")
- **roundDescription**: Explanation of what this round means
- **winImplication**: What happens if a player wins this round
- **loseImplication**: What happens if a player loses this round

## Enhanced Data

### Automatic NHRL Stats and Round Enrichment
ALL TrueFinals data is now automatically enriched with NHRL stats whenever player/bot information is returned:

#### Game Round Enrichment
When games have qualification round codes (Q1, Q2W, Q2L, Q3) as their names, they are automatically enriched with:
- **roundName**: Human-readable name (e.g., "Opening", "The Cusp", "Redemption", "Bubble")
- **roundDescription**: Explanation of what this round means
- **winImplication**: What happens if a player wins this round
- **loseImplication**: What happens if a player loses this round
- **isQualificationRound**: Boolean flag indicating this is a qualification game

This enrichment appears in:
- Individual game queries (`truefinals_games` tool)
- Tournament game lists
- Location active game information

#### Player Data Enrichment
Whenever player data is returned (in any TrueFinals tool), you'll automatically get:
- **nhrl_rank** - Current NHRL ranking (if bot exists in statsbook)
- **nhrl_current_streak** - Current win/lose streak with length and type
- **nhrl_recent_fights** - Number of recent fights (last 5)
- **nhrl_last_fight_date** - Date of most recent NHRL fight
- **bot_picture** - Bot picture URLs (thumbnail_url and full_size_url) from BrettZone

#### Tournament Enrichment
When you query tournaments through the TrueFinals tools, the data is automatically enhanced with:
- **Weight class detection** based on tournament title
- **Recent NHRL champions** for context
- **Bot rankings** for participants
- **Current win/loss streaks** for participants
- **Recent fight history** for participants

#### Game/Match Enrichment
Games now include NHRL stats for each player slot:
- **nhrl_rank** - Player's current NHRL ranking
- **nhrl_current_streak** - Active win/lose streak

#### Bracket & Standings Enrichment
Bracket views and standings include full NHRL stats:
- Rankings, streaks, recent fight counts, and last fight dates

### Supported Weight Classes
- **3lb** / **beetleweight** (category_id: 1)
- **12lb** / **antweight** (category_id: 2) 
- **30lb** / **hobbyweight** (category_id: 4)

### Supported Seasons
- **current** - Current season only
- **all-time** - All historical data
- **2018-2019** through **2023** - Specific seasons

### Pagination Support

Operations that can return large amounts of data now support pagination with `limit` and `offset` parameters:

- **limit**: Maximum number of results to return (default: 25)
- **offset**: Number of results to skip before returning data (default: 0)

Paginated operations include:
- `get_weight_class_stat_summary` - Can return hundreds of bots
- `get_weight_class_dumpster_count` - Podium finishers list
- `get_weight_class_event_winners` - Historical tournament winners
- `get_weight_class_fastest_kos` - KO leaderboard
- `get_weight_class_longest_streaks` - Winning streak leaderboard
- `get_bot_fights` - Bot fight history (can be extensive)
- `get_bot_head_to_head` - Head-to-head records
- `get_bot_event_participants` - Tournament participation history
- `get_tournament_matches` - Tournament match listings

Each paginated response includes a `pagination` object with:
- `total_count`: Total number of items available
- `limit`: Number of items requested
- `offset`: Starting position in the result set
- `has_more`: Boolean indicating if more results exist

Example pagination usage:
```json
{
  "tool": "nhrl_stats",
  "operation": "get_weight_class_stat_summary",
  "weight_class": "30lb",
  "limit": 10,
  "offset": 20
}
```

This would return bots 21-30 from the 30lb weight class rankings.

## Example Usage

### Get Bot Rank
```json
{
  "name": "nhrl_stats",
  "arguments": {
    "operation": "get_bot_rank",
    "bot_name": "Silent Spring"
  }
}
```

### Get Bot Picture URLs
```json
{
  "name": "nhrl_stats",
  "arguments": {
    "operation": "get_bot_picture_url",
    "bot_name": "Overlord"
  }
}
```

Returns:
```json
{
  "bot_name": "Overlord",
  "thumbnail_url": "https://brettzone.nhrl.io/brettZone/getBotPic.php?bot=Overlord&thumb",
  "full_size_url": "https://brettzone.nhrl.io/brettZone/getBotPic.php?bot=Overlord",
  "note": "These URLs return PNG images. The thumbnail is smaller and loads faster."
}
```

### Get Weight Class Champions
```json
{
  "name": "nhrl_stats", 
  "arguments": {
    "operation": "get_weight_class_dumpster_count",
    "weight_class": "3lb"
  }
}
```

### Get Tournament with NHRL Context
```json
{
  "name": "truefinals_tournaments",
  "arguments": {
    "operation": "get",
    "tournament_id": "your-tournament-id"
  }
}
```

The tournament response will include NHRL rankings, recent champions, and enriched player data automatically.

### Example Game Response with Qualification Round Enrichment
When you query a game with a qualification round code:
```json
{
  "id": "game-123",
  "name": "Q2W",
  "roundName": "The Cusp",
  "roundDescription": "Second match for Opening winners - one win away from qualifying",
  "winImplication": "Qualifies for main bracket",
  "loseImplication": "Drops to Bubble (Q3) for last chance",
  "isQualificationRound": true,
  "slots": [
    {
      "playerID": "player-1",
      "playerName": "Silent Spring",
      "nhrl_rank": 5,
      "nhrl_current_streak": {
        "length": 3,
        "type": "win"
      }
    },
    {
      "playerID": "player-2",
      "playerName": "Spinook",
      "nhrl_rank": 12
    }
  ]
}
```

### Get Live Tournament Matches
```json
{
  "name": "nhrl_stats",
  "arguments": {
    "operation": "get_tournament_matches",
    "tournament_id": "nhrl_june25_30lb"
  }
}
```

### Generate Match Review URL
```json
{
  "name": "nhrl_stats",
  "arguments": {
    "operation": "get_match_review_url",
    "game_id": "W-5",
    "tournament_id": "nhrl_june25_30lb",
    "cage_number": 1,
    "time_seconds": 10.5
  }
}
```

### Get Qualification System Explanation
```json
{
  "name": "nhrl_stats",
  "arguments": {
    "operation": "get_qualification_system",
    "round_code": "Q2W"  // Optional - get details about a specific round
  }
}
```

## API Details

- **No authentication required** for NHRL statsbook API
- **Automatic bot name normalization** (spaces converted to underscores)
- **Graceful error handling** for bots not found in rankings
- **Performance optimized** with recent data limiting where appropriate

## Data Sources

- **TrueFinals API**: Tournament management and bracket data
- **NHRL Statsbook**: https://stats.nhrl.io/statsbook/ - Bot performance, rankings, fight records
- **BrettZone**: https://brettzone.nhrl.io/brettZone/ - Live tournament system with match data, timing information, and video reviews

This integration provides the most comprehensive view of NHRL combat robot tournaments, combining real-time tournament management with historical performance data. 