package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// paginateSlice applies pagination to any slice and returns the paginated slice along with metadata
func paginateSlice[T any](items []T, limit, offset int) ([]T, map[string]interface{}) {
	// Default limit to 25 if not specified or invalid
	if limit <= 0 {
		limit = 25
	}

	// Default offset to 0 if negative
	if offset < 0 {
		offset = 0
	}

	totalCount := len(items)

	// Handle empty slice or offset beyond bounds
	if totalCount == 0 || offset >= totalCount {
		return []T{}, map[string]interface{}{
			"total_count": totalCount,
			"limit":       limit,
			"offset":      offset,
			"has_more":    false,
		}
	}

	// Calculate end index
	end := offset + limit
	if end > totalCount {
		end = totalCount
	}

	// Slice the items
	paginatedItems := items[offset:end]

	// Create metadata
	metadata := map[string]interface{}{
		"total_count": totalCount,
		"limit":       limit,
		"offset":      offset,
		"has_more":    end < totalCount,
	}

	return paginatedItems, metadata
}

// handleNHRLStatsTool handles all NHRL stats operations
func handleNHRLStatsTool(args map[string]interface{}) (string, error) {
	operation, ok := args["operation"].(string)
	if !ok {
		return "", fmt.Errorf("operation parameter is required")
	}

	// Check if operation is allowed in current tools mode
	if !isOperationAllowed("nhrl_stats", operation) {
		return "", fmt.Errorf(getOperationNotAllowedError(operation))
	}

	switch operation {
	case "get_bot_rank":
		return getNHRLBotRankTool(args)
	case "get_bot_fights":
		return getNHRLBotFightsTool(args)
	case "get_bot_head_to_head":
		return getNHRLBotHeadToHeadTool(args)
	case "get_bot_stats_by_season":
		return getNHRLBotStatsBySeasonTool(args)
	case "get_bot_streak_stats":
		return getNHRLBotStreakStatsTool(args)
	case "get_bot_event_participants":
		return getNHRLBotEventParticipantsTool(args)
	case "get_weight_class_dumpster_count":
		return getNHRLWeightClassDumpsterCountTool(args)
	case "get_weight_class_event_winners":
		return getNHRLWeightClassEventWinnersTool(args)
	case "get_weight_class_fastest_kos":
		return getNHRLWeightClassFastestKOsTool(args)
	case "get_weight_class_longest_streaks":
		return getNHRLWeightClassLongestStreaksTool(args)
	case "get_weight_class_stat_summary":
		return getNHRLWeightClassStatSummaryTool(args)
	case "get_weight_class_stat_summary_simple":
		return getNHRLWeightClassStatSummarySimpleTool(args)
	case "get_random_fight":
		return getNHRLRandomFightTool(args)
	case "get_tournament_matches":
		return getBrettZoneTournamentMatchesTool(args)
	case "get_match_review_url":
		return getBrettZoneMatchReviewURLTool(args)
	case "get_qualification_system":
		return getNHRLQualificationSystemTool(args)
	case "get_live_fight_stats":
		return getNHRLLiveFightStatsTool(args)
	case "get_bot_picture_url":
		return getNHRLBotPictureURLTool(args)
	default:
		return "", fmt.Errorf("unknown operation: %s", operation)
	}
}

// getNHRLStatsToolInfo returns the tool definition for NHRL stats operations
func getNHRLStatsToolInfo() ToolInfo {
	return ToolInfo{
		Name: "nhrl_stats",
		Description: `Query NHRL (National Havoc Robot League) statsbook and BrettZone tournament system for comprehensive combat robot data. 

This tool provides access to:
- Bot performance metrics, rankings, and fight history
- Head-to-head matchup records between specific bots
- Weight class statistics and leaderboards  
- Tournament bracket information and match results
- Live fight preparation data for upcoming matches
- Match review video URLs for past fights
- Bot images and profile information

Use this tool when you need historical data, performance statistics, or current tournament information. The data comes from both NHRL's official statsbook (historical records) and BrettZone (live tournament management).`,
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"operation": map[string]interface{}{
					"type": "string",
					"description": `The specific NHRL stats operation to perform:

BOT-SPECIFIC OPERATIONS (require bot_name):
- get_bot_rank: Get current ranking (based on Active season - previous + current season performance)
- get_bot_fights: Get complete fight history with dates, opponents, results, and methods
- get_bot_head_to_head: Get win/loss records against all opponents the bot has faced
- get_bot_stats_by_season: Get wins, losses, KOs, and other stats for a specific season
- get_bot_streak_stats: Get current and historical winning/losing streak information
- get_bot_event_participants: List all tournaments/events the bot has participated in
- get_bot_picture_url: Get thumbnail and full-size image URLs for the bot

WEIGHT CLASS OPERATIONS (use weight_class parameter):
- get_weight_class_dumpster_count: Get bots with most podium finishes (championship achievements)
- get_weight_class_event_winners: List tournament winners with dates and events
- get_weight_class_fastest_kos: Leaderboard of fastest knockout times
- get_weight_class_longest_streaks: Bots with longest winning streaks
- get_weight_class_stat_summary: Get statistics and rankings for all bots in the weight class
  * Use season="Active" for CURRENT RANKINGS (recommended for ranking queries)
  * Use season="all-time" for historical all-time statistics
  * Use specific year (e.g., "2024") for that season's statistics
- get_weight_class_stat_summary_simple: All-time statistics only (not recommended for current rankings)

TOURNAMENT/MATCH OPERATIONS:
- get_tournament_matches: Get all matches from a BrettZone tournament with results and bracket info
- get_match_review_url: Generate a video review URL for a specific match
- get_live_fight_stats: Get head-to-head stats and bot info for an upcoming match (requires bot1, bot2)

GENERAL OPERATIONS:
- get_random_fight: Get a random fight from NHRL history (fun/demo purposes)
- get_qualification_system: Explain NHRL's tournament qualification rounds and progression`,
					"enum": []string{
						"get_bot_rank", "get_bot_fights", "get_bot_head_to_head", "get_bot_stats_by_season",
						"get_bot_streak_stats", "get_bot_event_participants", "get_weight_class_dumpster_count",
						"get_weight_class_event_winners", "get_weight_class_fastest_kos", "get_weight_class_longest_streaks",
						"get_weight_class_stat_summary", "get_weight_class_stat_summary_simple", "get_random_fight", "get_tournament_matches", "get_match_review_url",
						"get_qualification_system", "get_live_fight_stats", "get_bot_picture_url",
					},
				},
				"bot_name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the bot (required for bot-specific operations). Case-insensitive. Spaces will be automatically converted to underscores. Examples: 'Ripperoni', 'Lynx', 'Slammo', 'Bloodsport'",
				},
				"bot1": map[string]interface{}{
					"type":        "string",
					"description": "First bot name for head-to-head comparison (used with get_live_fight_stats). This is typically the opponent.",
				},
				"bot2": map[string]interface{}{
					"type":        "string",
					"description": "Second bot name for head-to-head comparison (used with get_live_fight_stats). Stats returned will be for this bot, including head-to-head record against bot1.",
				},
				"weight_class": map[string]interface{}{
					"type":        "string",
					"description": "Weight class for the query. Use '3lb' (Beetleweight), '12lb' (Antweight), or '30lb' (Hobbyweight). Alternative names like 'beetleweight', 'antweight', 'hobbyweight' are also accepted.",
					"enum":        []string{"3lb", "12lb", "30lb", "beetleweight", "antweight", "hobbyweight"},
				},
				"season": map[string]interface{}{
					"type": "string",
					"description": `Season for stats query:
- "Active": CURRENT RANKINGS based on previous + current season (use this for ranking queries!)
- "all-time": Complete historical statistics
- "current": Current year only
- Specific years: "2025", "2024", "2023", etc.
- Special: "2018-19" for NHRL's first season that spanned two years

IMPORTANT: Use "Active" when you want current rankings, not "all-time"!`,
					"enum": []string{"Active", "current", "all-time", "2018-19", "2020", "2021", "2022", "2023", "2024", "2025"},
				},
				"tournament_id": map[string]interface{}{
					"type":        "string",
					"description": "BrettZone tournament identifier for tournament operations. Format is typically 'nhrl_month##_weightclass' (e.g., 'nhrl_june25_30lb' for June 2025 30lb tournament). Required for get_tournament_matches and get_match_review_url.",
				},
				"game_id": map[string]interface{}{
					"type":        "string",
					"description": "Match/Game identifier within a tournament for video review. Examples: 'W-5' (winners bracket match 5), 'Q1-12' (qualifying round 1 match 12), 'GF' (grand finals). Required for get_match_review_url.",
				},
				"cage_number": map[string]interface{}{
					"type":        "number",
					"description": "NHRL cage/arena number (1-4) where the match took place. Defaults to 1 if not specified. Used for generating correct video review URLs.",
				},
				"time_seconds": map[string]interface{}{
					"type":        "number",
					"description": "Start time in seconds for match review video. Defaults to 3.0 seconds. Use higher values to skip intro and jump to specific moments.",
				},
				"round_code": map[string]interface{}{
					"type":        "string",
					"description": "NHRL qualification round code to get detailed information about. Options: 'Q1' (Opening round), 'Q2W' (The Cusp - for Q1 winners), 'Q2L' (Redemption - for Q1 losers), 'Q3' (Bubble - final qualifying round).",
				},
				"limit": map[string]interface{}{
					"type":        "number",
					"description": "Maximum number of results to return. Defaults to 25. Use with offset for pagination. Applicable to operations that return lists of data (weight class stats, fight history, etc.).",
				},
				"offset": map[string]interface{}{
					"type":        "number",
					"description": "Number of results to skip before returning data. Defaults to 0. Use with limit for pagination. For example, limit=25&offset=25 returns results 26-50.",
				},
			},
			"required": []string{"operation"},
		},
	}
}

// Get bot rank
func getNHRLBotRankTool(args map[string]interface{}) (string, error) {
	botName, ok := args["bot_name"].(string)
	if !ok {
		return "", fmt.Errorf("bot_name is required for get_bot_rank operation")
	}

	rank, err := getNHRLBotRank(botName)
	if err != nil {
		return "", fmt.Errorf("failed to get bot rank: %w", err)
	}

	result := map[string]interface{}{
		"bot_name": botName,
		"rank":     rank,
	}

	if rank == nil {
		result["message"] = "Bot not found in current rankings"
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get bot fights
func getNHRLBotFightsTool(args map[string]interface{}) (string, error) {
	botName, ok := args["bot_name"].(string)
	if !ok {
		return "", fmt.Errorf("bot_name is required for get_bot_fights operation")
	}

	// Get pagination parameters
	limit := 25
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}

	offset := 0
	if o, ok := args["offset"].(float64); ok {
		offset = int(o)
	}

	fights, err := getNHRLFights(botName)
	if err != nil {
		return "", fmt.Errorf("failed to get bot fights: %w", err)
	}

	// Apply pagination
	paginatedFights, metadata := paginateSlice(fights, limit, offset)

	result := map[string]interface{}{
		"bot_name":    botName,
		"fight_count": len(paginatedFights),
		"fights":      paginatedFights,
		"pagination":  metadata,
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get bot head-to-head records
func getNHRLBotHeadToHeadTool(args map[string]interface{}) (string, error) {
	botName, ok := args["bot_name"].(string)
	if !ok {
		return "", fmt.Errorf("bot_name is required for get_bot_head_to_head operation")
	}

	// Get pagination parameters
	limit := 25
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}

	offset := 0
	if o, ok := args["offset"].(float64); ok {
		offset = int(o)
	}

	headToHead, err := getNHRLHeadToHead(botName)
	if err != nil {
		return "", fmt.Errorf("failed to get bot head-to-head: %w", err)
	}

	// Apply pagination
	paginatedHeadToHead, metadata := paginateSlice(headToHead, limit, offset)

	result := map[string]interface{}{
		"bot_name":           botName,
		"opponent_count":     len(paginatedHeadToHead),
		"head_to_head_stats": paginatedHeadToHead,
		"pagination":         metadata,
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get bot stats by season
func getNHRLBotStatsBySeasonTool(args map[string]interface{}) (string, error) {
	botName, ok := args["bot_name"].(string)
	if !ok {
		return "", fmt.Errorf("bot_name is required for get_bot_stats_by_season operation")
	}

	season := "all-time"
	if s, ok := args["season"].(string); ok {
		season = s
	}
	seasonID := getSeasonID(season)

	stats, err := getNHRLStatsBySeason(botName, seasonID)
	if err != nil {
		return "", fmt.Errorf("failed to get bot stats by season: %w", err)
	}

	result := map[string]interface{}{
		"bot_name": botName,
		"season":   season,
		"stats":    stats,
	}

	if stats == nil {
		result["message"] = "No stats found for this bot in the specified season"
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get bot streak stats
func getNHRLBotStreakStatsTool(args map[string]interface{}) (string, error) {
	botName, ok := args["bot_name"].(string)
	if !ok {
		return "", fmt.Errorf("bot_name is required for get_bot_streak_stats operation")
	}

	streakStats, err := getNHRLStreakStats(botName)
	if err != nil {
		return "", fmt.Errorf("failed to get bot streak stats: %w", err)
	}

	result := map[string]interface{}{
		"bot_name":     botName,
		"streak_stats": streakStats,
	}

	if streakStats == nil {
		result["message"] = "No streak stats found for this bot"
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get bot event participants
func getNHRLBotEventParticipantsTool(args map[string]interface{}) (string, error) {
	botName, ok := args["bot_name"].(string)
	if !ok {
		return "", fmt.Errorf("bot_name is required for get_bot_event_participants operation")
	}

	// Get pagination parameters
	limit := 25
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}

	offset := 0
	if o, ok := args["offset"].(float64); ok {
		offset = int(o)
	}

	participants, err := getNHRLEventParticipants(botName)
	if err != nil {
		return "", fmt.Errorf("failed to get bot event participants: %w", err)
	}

	// Apply pagination
	paginatedParticipants, metadata := paginateSlice(participants, limit, offset)

	result := map[string]interface{}{
		"bot_name":    botName,
		"event_count": len(paginatedParticipants),
		"events":      paginatedParticipants,
		"pagination":  metadata,
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get weight class dumpster count (podium finishes)
func getNHRLWeightClassDumpsterCountTool(args map[string]interface{}) (string, error) {
	weightClass := "3lb"
	if wc, ok := args["weight_class"].(string); ok {
		weightClass = wc
	}
	categoryID := getWeightClassCategoryID(weightClass)

	// Get pagination parameters
	limit := 25
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}

	offset := 0
	if o, ok := args["offset"].(float64); ok {
		offset = int(o)
	}

	dumpsterCount, err := getNHRLDumpsterCount(categoryID)
	if err != nil {
		return "", fmt.Errorf("failed to get weight class dumpster count: %w", err)
	}

	// Apply pagination
	paginatedDumpsterCount, metadata := paginateSlice(dumpsterCount, limit, offset)

	result := map[string]interface{}{
		"weight_class":           weightClass,
		"podium_finishers_count": len(paginatedDumpsterCount),
		"podium_finishers":       paginatedDumpsterCount,
		"pagination":             metadata,
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get weight class event winners
func getNHRLWeightClassEventWinnersTool(args map[string]interface{}) (string, error) {
	weightClass := "3lb"
	if wc, ok := args["weight_class"].(string); ok {
		weightClass = wc
	}

	// Get pagination parameters
	limit := 25
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}

	offset := 0
	if o, ok := args["offset"].(float64); ok {
		offset = int(o)
	}

	eventWinners, err := getNHRLEventWinners(weightClass)
	if err != nil {
		return "", fmt.Errorf("failed to get weight class event winners: %w", err)
	}

	// Apply pagination
	paginatedEventWinners, metadata := paginateSlice(eventWinners, limit, offset)

	result := map[string]interface{}{
		"weight_class":  weightClass,
		"event_count":   len(paginatedEventWinners),
		"event_winners": paginatedEventWinners,
		"pagination":    metadata,
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get weight class fastest KOs
func getNHRLWeightClassFastestKOsTool(args map[string]interface{}) (string, error) {
	weightClass := "3lb"
	if wc, ok := args["weight_class"].(string); ok {
		weightClass = wc
	}
	classID := getWeightClassCategoryID(weightClass)

	// Get pagination parameters
	limit := 25
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}

	offset := 0
	if o, ok := args["offset"].(float64); ok {
		offset = int(o)
	}

	fastestKOs, err := getNHRLFastestKOs(classID)
	if err != nil {
		return "", fmt.Errorf("failed to get weight class fastest KOs: %w", err)
	}

	// Apply pagination
	paginatedFastestKOs, metadata := paginateSlice(fastestKOs, limit, offset)

	result := map[string]interface{}{
		"weight_class": weightClass,
		"ko_count":     len(paginatedFastestKOs),
		"fastest_kos":  paginatedFastestKOs,
		"pagination":   metadata,
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get weight class longest winning streaks
func getNHRLWeightClassLongestStreaksTool(args map[string]interface{}) (string, error) {
	weightClass := "3lb"
	if wc, ok := args["weight_class"].(string); ok {
		weightClass = wc
	}
	categoryID := getWeightClassCategoryID(weightClass)

	// Get pagination parameters
	limit := 25
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}

	offset := 0
	if o, ok := args["offset"].(float64); ok {
		offset = int(o)
	}

	longestStreaks, err := getNHRLLongestWinningStreak(categoryID)
	if err != nil {
		return "", fmt.Errorf("failed to get weight class longest streaks: %w", err)
	}

	// Apply pagination
	paginatedLongestStreaks, metadata := paginateSlice(longestStreaks, limit, offset)

	result := map[string]interface{}{
		"weight_class":    weightClass,
		"streak_count":    len(paginatedLongestStreaks),
		"longest_streaks": paginatedLongestStreaks,
		"pagination":      metadata,
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get weight class stat summary
func getNHRLWeightClassStatSummaryTool(args map[string]interface{}) (string, error) {
	weightClass := "3lb"
	if wc, ok := args["weight_class"].(string); ok {
		weightClass = wc
	}
	categoryID := getWeightClassCategoryID(weightClass)

	season := "all-time"
	if s, ok := args["season"].(string); ok {
		season = s
	}
	seasonID := getSeasonID(season)

	// Get pagination parameters
	limit := 25
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}

	offset := 0
	if o, ok := args["offset"].(float64); ok {
		offset = int(o)
	}

	statSummary, err := getNHRLStatSummary(categoryID, seasonID)
	if err != nil {
		return "", fmt.Errorf("failed to get weight class stat summary: %w", err)
	}

	// Apply pagination
	paginatedStats, metadata := paginateSlice(statSummary, limit, offset)

	result := map[string]interface{}{
		"weight_class": weightClass,
		"season":       season,
		"bot_count":    len(paginatedStats),
		"stats":        paginatedStats,
		"pagination":   metadata,
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get weight class stat summary simple (all-time stats with correct ranking)
func getNHRLWeightClassStatSummarySimpleTool(args map[string]interface{}) (string, error) {
	weightClass := "3lb"
	if wc, ok := args["weight_class"].(string); ok {
		weightClass = wc
	}
	categoryID := getWeightClassCategoryID(weightClass)

	// Get pagination parameters
	limit := 25
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}

	offset := 0
	if o, ok := args["offset"].(float64); ok {
		offset = int(o)
	}

	statSummary, err := getNHRLStatSummarySimple(categoryID)
	if err != nil {
		return "", fmt.Errorf("failed to get weight class stat summary simple: %w", err)
	}

	// Apply pagination
	paginatedStats, metadata := paginateSlice(statSummary, limit, offset)

	result := map[string]interface{}{
		"weight_class": weightClass,
		"season":       "all-time",
		"bot_count":    len(paginatedStats),
		"stats":        paginatedStats,
		"pagination":   metadata,
		"note":         "This endpoint provides all-time stats with correct ranking for the weight class",
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get random fight
func getNHRLRandomFightTool(args map[string]interface{}) (string, error) {
	randomFight, err := getNHRLRandomFight()
	if err != nil {
		return "", fmt.Errorf("failed to get random fight: %w", err)
	}

	result := map[string]interface{}{
		"random_fight": randomFight,
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// getBrettZoneTournamentMatchesTool handles getting tournament matches from BrettZone
func getBrettZoneTournamentMatchesTool(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok || tournamentID == "" {
		return "", fmt.Errorf("tournament_id parameter is required")
	}

	// Get pagination parameters
	limit := 25
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}

	offset := 0
	if o, ok := args["offset"].(float64); ok {
		offset = int(o)
	}

	matches, err := getBrettZoneLatestMatches(tournamentID)
	if err != nil {
		return "", fmt.Errorf("failed to get tournament matches: %w", err)
	}

	// Enrich matches with round qualification information
	enrichedBrettZoneMatches := enrichBrettZoneMatches(matches)

	// Apply pagination before converting to response format
	paginatedMatches, metadata := paginateSlice(enrichedBrettZoneMatches, limit, offset)

	// Convert to response format with additional information
	enrichedMatches := make([]map[string]interface{}, len(paginatedMatches))
	for i, match := range paginatedMatches {
		enrichedMatch := map[string]interface{}{
			"tournamentID":     match.TournamentID,
			"matchID":          match.ID,
			"matchName":        match.Name,
			"round":            match.Round,
			"roundName":        match.RoundName,
			"roundDescription": match.RoundDescription,
			"winImplication":   match.WinImplication,
			"loseImplication":  match.LoseImplication,
			"cage":             match.Cage,
			"player1":          match.Player1,
			"player2":          match.Player2,
			"winner":           getMatchWinner(match.BrettZoneMatch),
			"winMethod":        match.WinAnnotation,
			"matchLengthSecs":  match.MatchLength,
			"weightClass":      match.WeightClass + "lb",
			"tournamentName":   match.TournamentName,
			"cameras":          match.Cams,
			"isTest":           match.IsTest == "1",
			"isFreestyle":      match.IsFreestyle == "1",
		}

		// Add timing information if available
		if match.StartTime != "" && match.StopTime != "" {
			enrichedMatch["startTime"] = match.StartTime
			enrichedMatch["stopTime"] = match.StopTime
		}

		// Add review URL
		cageNum := extractCageNumber(match.Cage)
		reviewURL := generateBrettZoneReviewURL(match.ID, match.TournamentID, cageNum, 3.0)
		enrichedMatch["reviewURL"] = reviewURL

		enrichedMatches[i] = enrichedMatch
	}

	// Sort matches by round and match name for better organization
	result := map[string]interface{}{
		"tournamentID":   tournamentID,
		"tournamentName": "",
		"totalMatches":   len(paginatedMatches),
		"matches":        enrichedMatches,
		"pagination":     metadata,
	}

	// Set tournament name if we have at least one match
	if len(enrichedBrettZoneMatches) > 0 {
		result["tournamentName"] = enrichedBrettZoneMatches[0].TournamentName
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal tournament matches data: %w", err)
	}

	return string(jsonData), nil
}

// getBrettZoneMatchReviewURLTool handles generating match review URLs
func getBrettZoneMatchReviewURLTool(args map[string]interface{}) (string, error) {
	gameID, ok := args["game_id"].(string)
	if !ok || gameID == "" {
		return "", fmt.Errorf("game_id parameter is required")
	}

	tournamentID, ok := args["tournament_id"].(string)
	if !ok || tournamentID == "" {
		return "", fmt.Errorf("tournament_id parameter is required")
	}

	// Optional parameters with defaults
	cageNum := 1
	if cage, ok := args["cage_number"].(float64); ok {
		cageNum = int(cage)
	}

	timeSeconds := 3.0
	if time, ok := args["time_seconds"].(float64); ok {
		timeSeconds = time
	}

	reviewURL := generateBrettZoneReviewURL(gameID, tournamentID, cageNum, timeSeconds)

	result := map[string]interface{}{
		"gameID":       gameID,
		"tournamentID": tournamentID,
		"cageNumber":   cageNum,
		"timeSeconds":  timeSeconds,
		"reviewURL":    reviewURL,
		"description":  fmt.Sprintf("Watch match %s from tournament %s starting at %.1f seconds", gameID, tournamentID, timeSeconds),
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal review URL data: %w", err)
	}

	return string(jsonData), nil
}

// Helper function to determine match winner
func getMatchWinner(match BrettZoneMatch) string {
	if match.Player1Wins == "1" {
		return match.Player1
	} else if match.Player2Wins == "1" {
		return match.Player2
	}
	return "undecided"
}

// Helper function to enrich player/bot data with NHRL stats
func enrichPlayerWithNHRLStats(player map[string]interface{}) map[string]interface{} {
	enrichedPlayer := make(map[string]interface{})
	for k, v := range player {
		enrichedPlayer[k] = v
	}

	// Try to get bot name from various fields
	var botName string
	if name, ok := player["name"].(string); ok {
		botName = name
	} else if displayName, ok := player["displayName"].(string); ok {
		botName = displayName
	} else if tag, ok := player["tag"].(string); ok {
		botName = tag
	}

	if botName != "" {
		// Get NHRL rank
		if rank, err := getNHRLBotRank(botName); err == nil && rank != nil {
			enrichedPlayer["nhrl_rank"] = rank.Ranking
		}

		// Get recent fight stats (limit to last 5 fights for performance)
		if fights, err := getNHRLFights(botName); err == nil && len(fights) > 0 {
			recentFights := fights
			if len(fights) > 5 {
				recentFights = fights[:5]
			}
			enrichedPlayer["nhrl_recent_fights"] = len(recentFights)
			enrichedPlayer["nhrl_last_fight_date"] = recentFights[0].Date
		}

		// Get streak stats
		if streakStats, err := getNHRLStreakStats(botName); err == nil && streakStats != nil {
			enrichedPlayer["nhrl_current_streak"] = map[string]interface{}{
				"length": streakStats.CurrentStreak,
				"type":   streakStats.CurrentStreakType,
			}
		}

		// Add bot picture URLs
		formattedBotName := strings.ReplaceAll(botName, " ", "_")
		enrichedPlayer["bot_picture"] = map[string]interface{}{
			"thumbnail_url": fmt.Sprintf("https://brettzone.nhrl.io/brettZone/getBotPic.php?bot=%s&thumb", formattedBotName),
			"full_size_url": fmt.Sprintf("https://brettzone.nhrl.io/brettZone/getBotPic.php?bot=%s", formattedBotName),
		}
	}

	return enrichedPlayer
}

// Get NHRL qualification system explanation
func getNHRLQualificationSystemTool(args map[string]interface{}) (string, error) {
	explanation := getQualificationPathExplanation()

	// Add more detailed information if specific round is asked
	if roundCode, ok := args["round_code"].(string); ok && roundCode != "" {
		roundInfo := getRoundInfo(roundCode)

		result := map[string]interface{}{
			"qualification_system": explanation,
			"specific_round": map[string]interface{}{
				"code":        roundInfo.Code,
				"name":        roundInfo.Name,
				"description": roundInfo.Description,
				"win_result":  roundInfo.WinResult,
				"lose_result": roundInfo.LoseResult,
			},
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal result: %w", err)
		}
		return string(jsonData), nil
	}

	// Return general explanation
	result := map[string]interface{}{
		"qualification_system": explanation,
		"round_codes": map[string]interface{}{
			"Q1":  "Opening - First qualifying match",
			"Q2W": "The Cusp - For Opening winners",
			"Q2L": "Redemption - For Opening losers",
			"Q3":  "Bubble - Final qualifying round",
		},
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}
	return string(jsonData), nil
}

// Get live fight stats between two bots
func getNHRLLiveFightStatsTool(args map[string]interface{}) (string, error) {
	bot1, ok := args["bot1"].(string)
	if !ok {
		return "", fmt.Errorf("bot1 is required for get_live_fight_stats operation")
	}

	bot2, ok := args["bot2"].(string)
	if !ok {
		return "", fmt.Errorf("bot2 is required for get_live_fight_stats operation")
	}

	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required for get_live_fight_stats operation")
	}

	stats, err := getNHRLLiveFightStats(bot1, bot2, tournamentID)
	if err != nil {
		return "", fmt.Errorf("failed to get live fight stats: %w", err)
	}

	// The API returns stats for bot2 with head-to-head data against bot1
	result := map[string]interface{}{
		"bot1":          bot1,
		"bot2":          bot2,
		"tournament_id": tournamentID,
		"stats_count":   len(stats),
	}

	// Add bot picture URLs for both bots
	bot1Formatted := strings.ReplaceAll(bot1, " ", "_")
	bot2Formatted := strings.ReplaceAll(bot2, " ", "_")

	result["bot_pictures"] = map[string]interface{}{
		"bot1": map[string]interface{}{
			"name":          bot1,
			"thumbnail_url": fmt.Sprintf("https://brettzone.nhrl.io/brettZone/getBotPic.php?bot=%s&thumb", bot1Formatted),
			"full_size_url": fmt.Sprintf("https://brettzone.nhrl.io/brettZone/getBotPic.php?bot=%s", bot1Formatted),
		},
		"bot2": map[string]interface{}{
			"name":          bot2,
			"thumbnail_url": fmt.Sprintf("https://brettzone.nhrl.io/brettZone/getBotPic.php?bot=%s&thumb", bot2Formatted),
			"full_size_url": fmt.Sprintf("https://brettzone.nhrl.io/brettZone/getBotPic.php?bot=%s", bot2Formatted),
		},
	}

	if len(stats) > 0 {
		// Extract the main bot stats (for bot2)
		botStats := stats[0]
		result["bot_stats"] = map[string]interface{}{
			"bot_name":             botStats.BotName,
			"driver_name":          botStats.DriverName,
			"driver_pronunciation": botStats.DriverPronunciation,
			"city":                 botStats.City,
			"state_province":       botStats.StateProvince,
			"country":              botStats.Country,
			"pronouns":             botStats.Pronouns,
			"team_name":            botStats.TeamName,
			"bot_pronunciation":    botStats.BotPronunciation,
			"ranking":              botStats.Ranking,
			"bot_type":             botStats.BotType,
			"builder_background":   botStats.BuilderBackground,
			"interesting_fact":     botStats.InterestingFact,
			"interesting_fact_2":   botStats.InterestingFact2,
		}

		// Overall stats
		result["overall_stats"] = map[string]interface{}{
			"events":                botStats.Events,
			"fights":                botStats.Fights,
			"wins":                  botStats.W,
			"losses":                botStats.L,
			"win_pct":               botStats.WinPct,
			"knockouts":             botStats.WKO,
			"knocked_out":           botStats.LKO,
			"judge_decision_wins":   botStats.WJD,
			"judge_decision_losses": botStats.LJD,
		}

		// Head-to-head stats against bot1
		result["head_to_head"] = map[string]interface{}{
			"total_wins":          botStats.HthW,
			"knockout_wins":       botStats.HthWKO,
			"judge_decision_wins": botStats.HthWJD,
			"last_meeting":        botStats.LastMeeting,
		}

		// Add note about head-to-head record
		if botStats.HthW > 0 || botStats.LastMeeting != nil {
			// If they've fought, add a note about the wins
			result["head_to_head_note"] = fmt.Sprintf("%s has %d wins against %s", botStats.BotName, botStats.HthW, bot1)
		}
	} else {
		result["message"] = "No stats found for these bots in the specified tournament"
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get bot picture URL
func getNHRLBotPictureURLTool(args map[string]interface{}) (string, error) {
	botName, ok := args["bot_name"].(string)
	if !ok {
		return "", fmt.Errorf("bot_name is required for get_bot_picture_url operation")
	}

	// Convert spaces to underscores in bot name
	formattedBotName := strings.ReplaceAll(botName, " ", "_")

	// Generate the picture URLs
	thumbnailURL := fmt.Sprintf("https://brettzone.nhrl.io/brettZone/getBotPic.php?bot=%s&thumb", formattedBotName)
	fullSizeURL := fmt.Sprintf("https://brettzone.nhrl.io/brettZone/getBotPic.php?bot=%s", formattedBotName)

	result := map[string]interface{}{
		"bot_name":      botName,
		"thumbnail_url": thumbnailURL,
		"full_size_url": fullSizeURL,
		"note":          "These URLs return PNG images. The thumbnail is smaller and loads faster.",
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Helper function to enrich tournament data with NHRL weight class context
func enrichTournamentWithNHRLContext(tournament map[string]interface{}) map[string]interface{} {
	enrichedTournament := make(map[string]interface{})
	for k, v := range tournament {
		enrichedTournament[k] = v
	}

	// Try to determine weight class from tournament title or other fields
	title := ""
	if t, ok := tournament["title"].(string); ok {
		title = t
	}

	weightClass := ""
	if strings.Contains(strings.ToLower(title), "3lb") || strings.Contains(strings.ToLower(title), "beetle") {
		weightClass = "3lb"
	} else if strings.Contains(strings.ToLower(title), "12lb") || strings.Contains(strings.ToLower(title), "ant") {
		weightClass = "12lb"
	} else if strings.Contains(strings.ToLower(title), "30lb") || strings.Contains(strings.ToLower(title), "hobby") {
		weightClass = "30lb"
	}

	if weightClass != "" {
		enrichedTournament["detected_weight_class"] = weightClass

		// Add recent champions for context
		if eventWinners, err := getNHRLEventWinners(weightClass); err == nil && len(eventWinners) > 0 {
			recentWinners := eventWinners
			if len(eventWinners) > 3 {
				recentWinners = eventWinners[:3]
			}
			enrichedTournament["nhrl_recent_champions"] = recentWinners
		}
	}

	// Enrich players with NHRL stats if available
	if players, ok := tournament["players"].([]interface{}); ok {
		enrichedPlayers := make([]interface{}, len(players))
		for i, p := range players {
			if player, ok := p.(map[string]interface{}); ok {
				enrichedPlayers[i] = enrichPlayerWithNHRLStats(player)
			} else {
				enrichedPlayers[i] = p
			}
		}
		enrichedTournament["players"] = enrichedPlayers
	}

	return enrichedTournament
}
