package main

import (
	"encoding/json"
	"fmt"
)

// handleGamesTool handles all game operations
func handleGamesTool(args map[string]interface{}) (string, error) {
	operation, ok := args["operation"].(string)
	if !ok {
		return "", fmt.Errorf("operation parameter is required")
	}

	// Check if operation is allowed in current tools mode
	if !isOperationAllowed("truefinals_games", operation) {
		return "", fmt.Errorf(getOperationNotAllowedError(operation))
	}

	switch operation {
	case "list":
		return listGames(args)
	case "get":
		return getGame(args)
	case "add_exhibition":
		return addExhibitionGame(args)
	case "edit_exhibition":
		return editExhibitionGame(args)
	case "delete_exhibition":
		return deleteExhibitionGame(args)
	case "bulk_add_exhibition":
		return bulkAddExhibitionGames(args)
	case "bulk_delete_exhibition":
		return bulkDeleteExhibitionGames(args)
	case "update":
		return updateGame(args)
	case "update_score":
		return updateGameScore(args)
	case "update_state":
		return updateGameState(args)
	case "update_scheduled_time":
		return updateGameScheduledTime(args)
	case "update_location":
		return updateGameLocation(args)
	case "update_checkin":
		return updateGameCheckin(args)
	case "undo":
		return undoGame(args)
	default:
		return "", fmt.Errorf("unknown operation: %s", operation)
	}
}

// getGamesToolInfo returns the tool definition for game operations
func getGamesToolInfo() ToolInfo {
	return ToolInfo{
		Name: "truefinals_games",
		Description: `Manage individual matches/games within NHRL tournaments in TrueFinals. This tool handles match administration, scoring, and results.

In NHRL tournaments, "games" represent individual robot combat matches between two bots. Each match is part of the tournament bracket structure and progresses the tournament when completed.

Use this tool when you need to:
- View match details and current scores
- Update match results and declare winners
- Create exhibition/non-bracket matches
- Check match status and progression
- Record specific win methods (KO, JD, DQ)

Common match identifiers in NHRL:
- Qualifying rounds: Q1-1, Q2W-5, Q2L-3, Q3-2
- Winners bracket: W-1, W-5, WF (Winners Final)
- Losers bracket: L-1, L-8, LF (Losers Final)  
- Finals: GF (Grand Final), GFR (Grand Final Reset)`,
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"operation": map[string]interface{}{
					"type": "string",
					"description": `The game/match operation to perform:

QUERY OPERATIONS:
- list: Get all matches in a tournament with current status
- get: Get detailed information about a specific match

MATCH UPDATES (require write access):
- update: Update match score or result
- create_exhibition: Create a non-bracket exhibition match
- delete_exhibition: Remove an exhibition match
- report_winner: Declare match winner with specific win method
- unreport_winner: Clear match result (reset to pending)
- set_in_progress: Mark match as currently being fought
- set_not_started: Reset match to not started status`,
					"enum": []string{
						"list", "get", "update", "create_exhibition", "delete_exhibition",
						"report_winner", "unreport_winner", "set_in_progress", "set_not_started",
					},
				},
				"tournament_id": map[string]interface{}{
					"type":        "string",
					"description": "Tournament identifier. Required for all operations. Format: 'nhrl_month##_weightclass'",
				},
				"game_id": map[string]interface{}{
					"type":        "string",
					"description": "Match/game identifier within the tournament. Required for single match operations. Examples: 'W-5' (winners bracket), 'Q1-12' (qualifying), 'GF' (grand final)",
				},
				"player1_score": map[string]interface{}{
					"type":        "integer",
					"description": "Score for player 1 (left side of bracket). In NHRL, typically 0 or 1 since matches are single elimination.",
				},
				"player2_score": map[string]interface{}{
					"type":        "integer",
					"description": "Score for player 2 (right side of bracket). In NHRL, typically 0 or 1 since matches are single elimination.",
				},
				"winner_id": map[string]interface{}{
					"type":        "string",
					"description": "Profile ID of the match winner. Used with report_winner operation.",
				},
				"win_annotation": map[string]interface{}{
					"type":        "string",
					"description": "Method of victory for NHRL matches. Common values: 'KO' (knockout), 'JD' (judges decision), 'DQ' (disqualification), 'FF' (forfeit)",
				},
				"match_length": map[string]interface{}{
					"type":        "integer",
					"description": "Match duration in seconds. Used to record how long the fight lasted before conclusion.",
				},
				"status": map[string]interface{}{
					"type":        "string",
					"description": "Match status for filtering or updates. Values: 'pending', 'in_progress', 'complete'",
				},
				"location_id": map[string]interface{}{
					"type":        "string",
					"description": "Venue/cage location ID where the match is scheduled. For NHRL, this typically corresponds to a specific combat cage.",
				},
				"participant_data": map[string]interface{}{
					"type":        "array",
					"description": "Participant information for exhibition matches. Each entry should include bot name and operator details.",
				},
			},
			"required": []string{"operation", "tournament_id"},
		},
	}
}

// List all games in a tournament
func listGames(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/games", tournamentID)

	data, err := makeAPIRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to list games: %w", err)
	}

	var games []interface{}
	if err := json.Unmarshal(data, &games); err != nil {
		return "", fmt.Errorf("failed to parse games response: %w", err)
	}

	// Enrich each game with player and location names
	enrichedGames := make([]interface{}, len(games))
	for i, g := range games {
		if game, ok := g.(map[string]interface{}); ok {
			enrichedGames[i] = enrichGameWithPlayerAndLocationInfo(game, tournamentID)
		} else {
			enrichedGames[i] = g
		}
	}

	result := map[string]interface{}{
		"games": enrichedGames,
		"count": len(enrichedGames),
		"note":  "Player names and location names are included for better readability",
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get a specific game by ID
func getGame(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	gameID, ok := args["game_id"].(string)
	if !ok {
		return "", fmt.Errorf("game_id is required")
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/games/%s", tournamentID, gameID)

	data, err := makeAPIRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get game: %w", err)
	}

	var game map[string]interface{}
	if err := json.Unmarshal(data, &game); err != nil {
		return "", fmt.Errorf("failed to parse game response: %w", err)
	}

	// Enrich game with player and location names
	enrichedGame := enrichGameWithPlayerAndLocationInfo(game, tournamentID)

	jsonData, err := json.MarshalIndent(enrichedGame, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Add an exhibition game
func addExhibitionGame(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	name, ok := args["name"].(string)
	if !ok {
		return "", fmt.Errorf("name is required")
	}

	scoreToWin, ok := args["score_to_win"].(float64)
	if !ok {
		return "", fmt.Errorf("score_to_win is required")
	}

	playerIDs, ok := args["player_ids"].([]interface{})
	if !ok {
		return "", fmt.Errorf("player_ids is required")
	}

	requestBody := map[string]interface{}{
		"name":       name,
		"scoreToWin": int(scoreToWin),
		"playerIDs":  playerIDs,
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/games", tournamentID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to add exhibition game: %w", err)
	}

	var game map[string]interface{}
	if err := json.Unmarshal(data, &game); err != nil {
		return "", fmt.Errorf("failed to parse game response: %w", err)
	}

	// Enrich game with player and location names
	enrichedGame := enrichGameWithPlayerAndLocationInfo(game, tournamentID)

	jsonData, err := json.MarshalIndent(enrichedGame, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Edit an exhibition game
func editExhibitionGame(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	gameID, ok := args["game_id"].(string)
	if !ok {
		return "", fmt.Errorf("game_id is required")
	}

	name, ok := args["name"].(string)
	if !ok {
		return "", fmt.Errorf("name is required")
	}

	scoreToWin, ok := args["score_to_win"].(float64)
	if !ok {
		return "", fmt.Errorf("score_to_win is required")
	}

	playerIDs, ok := args["player_ids"].([]interface{})
	if !ok {
		return "", fmt.Errorf("player_ids is required")
	}

	requestBody := map[string]interface{}{
		"name":       name,
		"scoreToWin": int(scoreToWin),
		"playerIDs":  playerIDs,
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/games/%s", tournamentID, gameID)

	data, err := makeAPIRequest("PUT", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to edit exhibition game: %w", err)
	}

	var game map[string]interface{}
	if err := json.Unmarshal(data, &game); err != nil {
		return "", fmt.Errorf("failed to parse game response: %w", err)
	}

	// Enrich game with player and location names
	enrichedGame := enrichGameWithPlayerAndLocationInfo(game, tournamentID)

	jsonData, err := json.MarshalIndent(enrichedGame, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Delete an exhibition game
func deleteExhibitionGame(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	gameID, ok := args["game_id"].(string)
	if !ok {
		return "", fmt.Errorf("game_id is required")
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/games/%s", tournamentID, gameID)

	data, err := makeAPIRequest("DELETE", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to delete exhibition game: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Bulk add exhibition games
func bulkAddExhibitionGames(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	gamesInfo, ok := args["games_info"].([]interface{})
	if !ok {
		return "", fmt.Errorf("games_info is required")
	}

	requestBody := map[string]interface{}{
		"gamesInfo": gamesInfo,
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/bulkGames/add", tournamentID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to bulk add exhibition games: %w", err)
	}

	var games []interface{}
	if err := json.Unmarshal(data, &games); err != nil {
		return "", fmt.Errorf("failed to parse games response: %w", err)
	}

	// Enrich each game with player and location names
	enrichedGames := make([]interface{}, len(games))
	for i, g := range games {
		if game, ok := g.(map[string]interface{}); ok {
			enrichedGames[i] = enrichGameWithPlayerAndLocationInfo(game, tournamentID)
		} else {
			enrichedGames[i] = g
		}
	}

	result := map[string]interface{}{
		"games": enrichedGames,
		"count": len(enrichedGames),
		"note":  "Player names and location names are included for better readability",
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Bulk delete exhibition games
func bulkDeleteExhibitionGames(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	gameIDs, ok := args["game_ids"]
	if !ok {
		return "", fmt.Errorf("game_ids is required")
	}

	requestBody := map[string]interface{}{
		"gameIDs": gameIDs,
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/bulkGames/delete", tournamentID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to bulk delete exhibition games: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Update a game
func updateGame(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	gameID, ok := args["game_id"].(string)
	if !ok {
		return "", fmt.Errorf("game_id is required")
	}

	requestBody := map[string]interface{}{}

	if state, ok := args["state"].(string); ok {
		requestBody["state"] = state
	}

	if slots, ok := args["slots"]; ok {
		requestBody["slots"] = slots
	}

	if locationID, ok := args["location_id"]; ok {
		requestBody["locationID"] = locationID
	}

	if resultAnnotation, ok := args["result_annotation"]; ok {
		requestBody["resultAnnotation"] = resultAnnotation
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/games/%s", tournamentID, gameID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to update game: %w", err)
	}

	var game map[string]interface{}
	if err := json.Unmarshal(data, &game); err != nil {
		return "", fmt.Errorf("failed to parse game response: %w", err)
	}

	// Enrich game with player and location names
	enrichedGame := enrichGameWithPlayerAndLocationInfo(game, tournamentID)

	jsonData, err := json.MarshalIndent(enrichedGame, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Update game score
func updateGameScore(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	gameID, ok := args["game_id"].(string)
	if !ok {
		return "", fmt.Errorf("game_id is required")
	}

	scores, ok := args["scores"].([]interface{})
	if !ok {
		return "", fmt.Errorf("scores is required")
	}

	requestBody := map[string]interface{}{
		"scores": scores,
	}

	if resultAnnotation, ok := args["result_annotation"]; ok {
		requestBody["resultAnnotation"] = resultAnnotation
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/games/%s/updateScore", tournamentID, gameID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to update game score: %w", err)
	}

	var game map[string]interface{}
	if err := json.Unmarshal(data, &game); err != nil {
		return "", fmt.Errorf("failed to parse game response: %w", err)
	}

	// Enrich game with player and location names
	enrichedGame := enrichGameWithPlayerAndLocationInfo(game, tournamentID)

	jsonData, err := json.MarshalIndent(enrichedGame, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Update game state
func updateGameState(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	gameID, ok := args["game_id"].(string)
	if !ok {
		return "", fmt.Errorf("game_id is required")
	}

	state, ok := args["state"].(string)
	if !ok {
		return "", fmt.Errorf("state is required")
	}

	requestBody := map[string]interface{}{
		"state": state,
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/games/%s/updateState", tournamentID, gameID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to update game state: %w", err)
	}

	var game map[string]interface{}
	if err := json.Unmarshal(data, &game); err != nil {
		return "", fmt.Errorf("failed to parse game response: %w", err)
	}

	// Enrich game with player and location names
	enrichedGame := enrichGameWithPlayerAndLocationInfo(game, tournamentID)

	jsonData, err := json.MarshalIndent(enrichedGame, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Update game scheduled time
func updateGameScheduledTime(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	gameID, ok := args["game_id"].(string)
	if !ok {
		return "", fmt.Errorf("game_id is required")
	}

	scheduledTime, ok := args["scheduled_time"].(float64)
	if !ok {
		return "", fmt.Errorf("scheduled_time is required")
	}

	requestBody := map[string]interface{}{
		"scheduledTime": scheduledTime,
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/games/%s/updateScheduledTime", tournamentID, gameID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to update game scheduled time: %w", err)
	}

	var game map[string]interface{}
	if err := json.Unmarshal(data, &game); err != nil {
		return "", fmt.Errorf("failed to parse game response: %w", err)
	}

	// Enrich game with player and location names
	enrichedGame := enrichGameWithPlayerAndLocationInfo(game, tournamentID)

	jsonData, err := json.MarshalIndent(enrichedGame, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Update game location
func updateGameLocation(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	gameID, ok := args["game_id"].(string)
	if !ok {
		return "", fmt.Errorf("game_id is required")
	}

	requestBody := map[string]interface{}{}

	if locationID, ok := args["location_id"]; ok {
		requestBody["locationID"] = locationID
	}

	if locationQueueIdx, ok := args["location_queue_idx"].(float64); ok {
		requestBody["locationQueueIdx"] = int(locationQueueIdx)
	}

	if activate, ok := args["activate"].(bool); ok {
		requestBody["activate"] = activate
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/games/%s/updateLocation", tournamentID, gameID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to update game location: %w", err)
	}

	var game map[string]interface{}
	if err := json.Unmarshal(data, &game); err != nil {
		return "", fmt.Errorf("failed to parse game response: %w", err)
	}

	// Enrich game with player and location names
	enrichedGame := enrichGameWithPlayerAndLocationInfo(game, tournamentID)

	jsonData, err := json.MarshalIndent(enrichedGame, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Update game check-in status
func updateGameCheckin(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	gameID, ok := args["game_id"].(string)
	if !ok {
		return "", fmt.Errorf("game_id is required")
	}

	slotIdx, ok := args["slot_idx"].(float64)
	if !ok {
		return "", fmt.Errorf("slot_idx is required")
	}

	checkInStatus, ok := args["check_in_status"].(string)
	if !ok {
		return "", fmt.Errorf("check_in_status is required")
	}

	requestBody := map[string]interface{}{
		"slotIdx":       int(slotIdx),
		"checkInStatus": checkInStatus,
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/games/%s/updateCheckin", tournamentID, gameID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to update game check-in: %w", err)
	}

	var game map[string]interface{}
	if err := json.Unmarshal(data, &game); err != nil {
		return "", fmt.Errorf("failed to parse game response: %w", err)
	}

	// Enrich game with player and location names
	enrichedGame := enrichGameWithPlayerAndLocationInfo(game, tournamentID)

	jsonData, err := json.MarshalIndent(enrichedGame, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Undo a game
func undoGame(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	gameID, ok := args["game_id"].(string)
	if !ok {
		return "", fmt.Errorf("game_id is required")
	}

	requestBody := map[string]interface{}{}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/games/%s/undo", tournamentID, gameID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to undo game: %w", err)
	}

	var game map[string]interface{}
	if err := json.Unmarshal(data, &game); err != nil {
		return "", fmt.Errorf("failed to parse game response: %w", err)
	}

	// Enrich game with player and location names
	enrichedGame := enrichGameWithPlayerAndLocationInfo(game, tournamentID)

	jsonData, err := json.MarshalIndent(enrichedGame, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}
