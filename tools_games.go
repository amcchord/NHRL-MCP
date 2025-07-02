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
		return "", fmt.Errorf("operation '%s' not allowed in '%s' mode", operation, toolsMode)
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
		Name:        "truefinals_games",
		Description: "Manage games within tournaments in TrueFinals. Create, update, delete exhibition games, and manage game states and scores.",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"operation": map[string]interface{}{
					"type":        "string",
					"description": "The operation to perform",
					"enum": []string{
						"list", "get", "add_exhibition", "edit_exhibition", "delete_exhibition",
						"bulk_add_exhibition", "bulk_delete_exhibition", "update", "update_score",
						"update_state", "update_scheduled_time", "update_location", "update_checkin", "undo",
					},
				},
				"tournament_id": map[string]interface{}{
					"type":        "string",
					"description": "Tournament ID - required for all operations",
				},
				"game_id": map[string]interface{}{
					"type":        "string",
					"description": "Game ID - required for single game operations",
				},
				// Exhibition game fields
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Game name (1-6 characters, A-Z0-9_-)",
				},
				"score_to_win": map[string]interface{}{
					"type":        "integer",
					"description": "Score needed to win (1-5)",
					"minimum":     1,
					"maximum":     5,
				},
				"player_ids": map[string]interface{}{
					"type":        "array",
					"description": "Array of player IDs for the game",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
				"games_info": map[string]interface{}{
					"type":        "array",
					"description": "Array of game information for bulk operations",
				},
				"game_ids": map[string]interface{}{
					"type":        "array",
					"description": "Array of game IDs for bulk deletion",
				},
				// Game update fields
				"state": map[string]interface{}{
					"type":        "string",
					"description": "Game state",
					"enum":        []string{"unavailable", "available", "called", "active", "hold", "done"},
				},
				"slots": map[string]interface{}{
					"type":        "array",
					"description": "Game slot information",
				},
				"location_id": map[string]interface{}{
					"type":        "string",
					"description": "Location ID for the game",
				},
				"result_annotation": map[string]interface{}{
					"type":        "string",
					"description": "Result annotation",
					"enum":        []string{"KO", "TO", "JD", "TKO", "HLD", "BY", "DQ", "FF", "T"},
				},
				"scores": map[string]interface{}{
					"type":        "array",
					"description": "Array of scores for game slots",
				},
				"scheduled_time": map[string]interface{}{
					"type":        "number",
					"description": "Scheduled time for the game (Unix timestamp)",
				},
				"location_queue_idx": map[string]interface{}{
					"type":        "integer",
					"description": "Queue index for location assignment",
					"minimum":     0,
				},
				"activate": map[string]interface{}{
					"type":        "boolean",
					"description": "Whether to activate the game at the location",
				},
				"slot_idx": map[string]interface{}{
					"type":        "integer",
					"description": "Slot index for check-in updates",
					"minimum":     0,
					"maximum":     1,
				},
				"check_in_status": map[string]interface{}{
					"type":        "string",
					"description": "Check-in status for player",
					"enum":        []string{"not_ready", "checked_in", "waiting"},
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
