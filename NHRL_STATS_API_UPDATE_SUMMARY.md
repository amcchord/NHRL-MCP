# NHRL Stats API Update Summary

## Overview
This document summarizes the changes made to the NHRL MCP server to support the updated NHRL Stats book API based on the documentation at [https://stats.nhrl.io/statsbook/](https://stats.nhrl.io/statsbook/).

## Key Changes

### 1. New API Endpoint Support
- Added support for `get_stat_summary_simple.php` endpoint
- This endpoint provides all-time stats only (not recommended for current rankings)
- Accessible via the new operation: `get_weight_class_stat_summary_simple`

### 2. Parameter Name Corrections
- Updated `get_stats_by_season.php` to use `season` parameter instead of `season_id`
- Maintained `season_id` for `get_stat_summary.php` as per documentation
- Ensured all weight class parameters use `category_id` consistently

### 3. Season Parameter Updates
- Updated season handling to support actual year values (e.g., "2025", "2024")
- Added support for special season values:
  - **`Active`**: **CURRENT RANKINGS** - Previous and current season only (recommended for ranking queries!)
  - `All-time`: Complete historical stats
  - `2018-19`: NHRL's first season that spanned two years
  - Direct year values: "2020", "2021", "2022", "2023", "2024", "2025"

### 4. Tool Documentation Updates
- Updated the tool descriptions to reflect the new endpoints
- **IMPORTANT**: Emphasized that `season="Active"` should be used for current rankings
- Clarified that `get_weight_class_stat_summary_simple` provides all-time stats only
- Updated the season enum to include actual years and special values

### 5. Implementation Details
- Added `getNHRLStatSummarySimple()` function in `nhrl_api.go`
- Added `getNHRLWeightClassStatSummarySimpleTool()` function in `tools_nhrl.go`
- Updated `getSeasonID()` function to handle the new season format

## Important Note on Rankings

**When users ask about a robot's current ranking, use `season="Active"`**, not "all-time". The Active season includes the previous season and current season only, which is how NHRL calculates current rankings. All-time statistics are historical and don't reflect current competitive standing.

## Usage Examples

### Get Current Rankings (Recommended for ranking queries)
```json
{
  "tool": "nhrl_stats",
  "arguments": {
    "operation": "get_weight_class_stat_summary",
    "weight_class": "3lb",
    "season": "Active"
  }
}
```

### Get All-Time Historical Stats
```json
{
  "tool": "nhrl_stats",
  "arguments": {
    "operation": "get_weight_class_stat_summary",
    "weight_class": "3lb",
    "season": "all-time"
  }
}
```

### Get Stats for a Specific Season
```json
{
  "tool": "nhrl_stats",
  "arguments": {
    "operation": "get_weight_class_stat_summary",
    "weight_class": "12lb",
    "season": "2024"
  }
}
```

### Get Bot Stats for Active Season (Current Performance)
```json
{
  "tool": "nhrl_stats",
  "arguments": {
    "operation": "get_bot_stats_by_season",
    "bot_name": "Lynx",
    "season": "Active"
  }
}
```

## Notes
- **For current rankings**: Use `get_weight_class_stat_summary` with `season="Active"`
- **For historical data**: Use `season="all-time"` or specific years
- The `get_stat_summary_simple` endpoint only provides all-time stats and is not suitable for current rankings
- Bot names are automatically normalized (spaces converted to underscores) for API compatibility 