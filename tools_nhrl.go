package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// handleNHRLStatsTool handles all NHRL stats operations
func handleNHRLStatsTool(args map[string]interface{}) (string, error) {
	operation, ok := args["operation"].(string)
	if !ok {
		return "", fmt.Errorf("operation parameter is required")
	}

	// Check if operation is allowed in current tools mode
	if !isOperationAllowed("nhrl_stats", operation) {
		return "", fmt.Errorf("operation '%s' not allowed in '%s' mode", operation, toolsMode)
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
	case "get_random_fight":
		return getNHRLRandomFightTool(args)
	case "get_tournament_matches":
		return getBrettZoneTournamentMatchesTool(args)
	case "get_match_review_url":
		return getBrettZoneMatchReviewURLTool(args)
	case "get_qualification_system":
		return getNHRLQualificationSystemTool(args)
	default:
		return "", fmt.Errorf("unknown operation: %s", operation)
	}
}

// getNHRLStatsToolInfo returns the tool definition for NHRL stats operations
func getNHRLStatsToolInfo() ToolInfo {
	return ToolInfo{
		Name:        "nhrl_stats",
		Description: "Query NHRL (National Havoc Robot League) statsbook and BrettZone tournament system for comprehensive combat robot data. Access bot performance data, rankings, fight records, tournament statistics, live match data, and match review videos. Includes both historical statsbook data and real-time tournament match information.",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"operation": map[string]interface{}{
					"type":        "string",
					"description": "The NHRL stats operation to perform",
					"enum": []string{
						"get_bot_rank", "get_bot_fights", "get_bot_head_to_head", "get_bot_stats_by_season",
						"get_bot_streak_stats", "get_bot_event_participants", "get_weight_class_dumpster_count",
						"get_weight_class_event_winners", "get_weight_class_fastest_kos", "get_weight_class_longest_streaks",
						"get_weight_class_stat_summary", "get_random_fight", "get_tournament_matches", "get_match_review_url",
						"get_qualification_system",
					},
				},
				"bot_name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the bot (required for bot-specific operations). Spaces will be automatically converted to underscores.",
				},
				"weight_class": map[string]interface{}{
					"type":        "string",
					"description": "Weight class for the query",
					"enum":        []string{"3lb", "12lb", "30lb", "beetleweight", "antweight", "hobbyweight"},
				},
				"season": map[string]interface{}{
					"type":        "string",
					"description": "Season for stats query",
					"enum":        []string{"current", "all-time", "2018-2019", "2020", "2021", "2022", "2023"},
				},
				"tournament_id": map[string]interface{}{
					"type":        "string",
					"description": "Tournament ID for BrettZone operations (e.g., 'nhrl_june25_30lb')",
				},
				"game_id": map[string]interface{}{
					"type":        "string",
					"description": "Game/Match ID for match review URL generation (e.g., 'W-5')",
				},
				"cage_number": map[string]interface{}{
					"type":        "number",
					"description": "Cage number for match review (1-4, defaults to 1)",
				},
				"time_seconds": map[string]interface{}{
					"type":        "number",
					"description": "Start time in seconds for match review (defaults to 3.0)",
				},
				"round_code": map[string]interface{}{
					"type":        "string",
					"description": "Round code (Q1, Q2W, Q2L, Q3) to get specific information about that qualification round",
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

	fights, err := getNHRLFights(botName)
	if err != nil {
		return "", fmt.Errorf("failed to get bot fights: %w", err)
	}

	result := map[string]interface{}{
		"bot_name":    botName,
		"fight_count": len(fights),
		"fights":      fights,
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

	headToHead, err := getNHRLHeadToHead(botName)
	if err != nil {
		return "", fmt.Errorf("failed to get bot head-to-head: %w", err)
	}

	result := map[string]interface{}{
		"bot_name":           botName,
		"opponent_count":     len(headToHead),
		"head_to_head_stats": headToHead,
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

	participants, err := getNHRLEventParticipants(botName)
	if err != nil {
		return "", fmt.Errorf("failed to get bot event participants: %w", err)
	}

	result := map[string]interface{}{
		"bot_name":    botName,
		"event_count": len(participants),
		"events":      participants,
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

	dumpsterCount, err := getNHRLDumpsterCount(categoryID)
	if err != nil {
		return "", fmt.Errorf("failed to get weight class dumpster count: %w", err)
	}

	result := map[string]interface{}{
		"weight_class":           weightClass,
		"podium_finishers_count": len(dumpsterCount),
		"podium_finishers":       dumpsterCount,
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

	eventWinners, err := getNHRLEventWinners(weightClass)
	if err != nil {
		return "", fmt.Errorf("failed to get weight class event winners: %w", err)
	}

	result := map[string]interface{}{
		"weight_class":  weightClass,
		"event_count":   len(eventWinners),
		"event_winners": eventWinners,
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

	fastestKOs, err := getNHRLFastestKOs(classID)
	if err != nil {
		return "", fmt.Errorf("failed to get weight class fastest KOs: %w", err)
	}

	result := map[string]interface{}{
		"weight_class": weightClass,
		"ko_count":     len(fastestKOs),
		"fastest_kos":  fastestKOs,
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

	longestStreaks, err := getNHRLLongestWinningStreak(categoryID)
	if err != nil {
		return "", fmt.Errorf("failed to get weight class longest streaks: %w", err)
	}

	result := map[string]interface{}{
		"weight_class":    weightClass,
		"streak_count":    len(longestStreaks),
		"longest_streaks": longestStreaks,
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

	statSummary, err := getNHRLStatSummary(categoryID, seasonID)
	if err != nil {
		return "", fmt.Errorf("failed to get weight class stat summary: %w", err)
	}

	result := map[string]interface{}{
		"weight_class": weightClass,
		"season":       season,
		"bot_count":    len(statSummary),
		"stats":        statSummary,
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

	matches, err := getBrettZoneLatestMatches(tournamentID)
	if err != nil {
		return "", fmt.Errorf("failed to get tournament matches: %w", err)
	}

	// Enrich matches with round qualification information
	enrichedBrettZoneMatches := enrichBrettZoneMatches(matches)

	// Convert to response format with additional information
	enrichedMatches := make([]map[string]interface{}, len(enrichedBrettZoneMatches))
	for i, match := range enrichedBrettZoneMatches {
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
		"tournamentName": matches[0].TournamentName,
		"totalMatches":   len(matches),
		"matches":        enrichedMatches,
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
