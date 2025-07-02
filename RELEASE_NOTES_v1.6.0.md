# Release Notes - v1.6.0

## New Features

### Live Fight Stats Support
- Added `get_live_fight_stats` operation to the NHRL stats tool
- Provides real-time, detailed statistics for bots including head-to-head matchup data
- Returns comprehensive information including:
  - Driver details (name, pronunciation, location, pronouns)
  - Bot specifications and builder background
  - Overall performance statistics
  - Head-to-head record against specific opponents
  - Historical meeting information

### API Details
The new operation calls the NHRL live stats endpoint:
- Endpoint: `https://stats.nhrl.io/live_stats/query/get_fight_stats.php`
- Method: POST with form data
- Required parameters: `bot1`, `bot2`, `tournament_id`

### Usage Example
```json
{
  "operation": "get_live_fight_stats",
  "bot1": "MegatRON",
  "bot2": "Hurricane",
  "tournament_id": "nhrl_june25_30lb"
}
```

## Technical Changes
- Added `NHRLLiveFightStats` struct to model the response data
- Implemented `getNHRLLiveFightStats` function for API communication
- Added `getNHRLLiveFightStatsTool` handler for MCP integration
- Updated tool schema to include new parameters

## Notes
- Stats are returned for `bot2` with head-to-head data against `bot1`
- Provides more current data than traditional statsbook queries
- Includes additional metadata not available in standard queries 