package main

import (
	"encoding/json"
	"fmt"
)

// handlePlayersTool handles all player operations
func handlePlayersTool(args map[string]interface{}) (string, error) {
	operation, ok := args["operation"].(string)
	if !ok {
		return "", fmt.Errorf("operation parameter is required")
	}

	// Check if operation is allowed in current tools mode
	if !isOperationAllowed("truefinals_players", operation) {
		return "", fmt.Errorf("operation '%s' not allowed in '%s' mode", operation, toolsMode)
	}

	switch operation {
	case "list":
		return listPlayers(args)
	case "get":
		return getPlayer(args)
	case "add":
		return addPlayer(args)
	case "update":
		return updatePlayer(args)
	case "delete":
		return deletePlayer(args)
	case "reseed":
		return reseedPlayer(args)
	case "randomize":
		return randomizeSeeding(args)
	case "bulk_update":
		return bulkUpdatePlayers(args)
	case "checkin":
		return checkinPlayer(args)
	case "disqualify":
		return disqualifyPlayer(args)
	default:
		return "", fmt.Errorf("unknown operation: %s", operation)
	}
}

// getPlayersToolInfo returns the tool definition for player operations
func getPlayersToolInfo() ToolInfo {
	return ToolInfo{
		Name:        "truefinals_players",
		Description: "Manage tournament players in TrueFinals. Add, update, delete players, manage seeding, check-ins, and disqualifications.",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"operation": map[string]interface{}{
					"type":        "string",
					"description": "The operation to perform",
					"enum": []string{
						"list", "get", "add", "update", "delete", "reseed",
						"randomize", "bulk_update", "checkin", "disqualify",
					},
				},
				"tournament_id": map[string]interface{}{
					"type":        "string",
					"description": "Tournament ID - required for all operations",
				},
				"player_id": map[string]interface{}{
					"type":        "string",
					"description": "Player ID - required for single player operations",
				},
				"include_bye_players": map[string]interface{}{
					"type":        "boolean",
					"description": "Whether to include bye players in results",
				},
				"seed_idx": map[string]interface{}{
					"type":        "integer",
					"description": "Seed index for player position",
				},
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Player name (1-32 characters)",
					"minLength":   1,
					"maxLength":   32,
				},
				"photo_url": map[string]interface{}{
					"type":        "string",
					"description": "Player photo URL",
					"maxLength":   1000,
				},
				"profile_info": map[string]interface{}{
					"type":        "object",
					"description": "Player profile information",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type": "string",
						},
						"tag": map[string]interface{}{
							"type":    "string",
							"pattern": "^.*#[0-9]{4}",
						},
						"name": map[string]interface{}{
							"type":      "string",
							"minLength": 1,
							"maxLength": 32,
						},
						"photo_url": map[string]interface{}{
							"type":      "string",
							"maxLength": 1000,
						},
						"pronouns": map[string]interface{}{
							"type":      "string",
							"maxLength": 32,
						},
						"twitch_handle": map[string]interface{}{
							"type": "string",
						},
						"twitter_handle": map[string]interface{}{
							"type": "string",
						},
						"startgg_player_id": map[string]interface{}{
							"type":    "integer",
							"minimum": 1,
						},
					},
				},
				"participants": map[string]interface{}{
					"type":        "array",
					"description": "Array of participant player information for bulk operations",
				},
				"non_participants": map[string]interface{}{
					"type":        "array",
					"description": "Array of non-participant player information for bulk operations",
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

// List all players in a tournament
func listPlayers(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/players", tournamentID)

	// Add query parameter if specified
	if includeByePlayers, ok := args["include_bye_players"].(bool); ok {
		if includeByePlayers {
			endpoint += "?includeByePlayers=true"
		} else {
			endpoint += "?includeByePlayers=false"
		}
	}

	data, err := makeAPIRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to list players: %w", err)
	}

	var players []interface{}
	if err := json.Unmarshal(data, &players); err != nil {
		return "", fmt.Errorf("failed to parse players response: %w", err)
	}

	// Enrich each player with profile details
	enrichedPlayers := make([]interface{}, len(players))
	for i, p := range players {
		if player, ok := p.(map[string]interface{}); ok {
			enrichedPlayers[i] = enrichPlayerDataWithProfileDetails(player)
		} else {
			enrichedPlayers[i] = p
		}
	}

	result := map[string]interface{}{
		"players": enrichedPlayers,
		"count":   len(enrichedPlayers),
		"note":    "Player display names and profile details are included for better readability",
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get a specific player by ID
func getPlayer(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	playerID, ok := args["player_id"].(string)
	if !ok {
		return "", fmt.Errorf("player_id is required")
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/players/%s", tournamentID, playerID)

	data, err := makeAPIRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get player: %w", err)
	}

	var player map[string]interface{}
	if err := json.Unmarshal(data, &player); err != nil {
		return "", fmt.Errorf("failed to parse player response: %w", err)
	}

	// Enrich player with profile details
	enrichedPlayer := enrichPlayerDataWithProfileDetails(player)

	jsonData, err := json.MarshalIndent(enrichedPlayer, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Add a new player to a tournament
func addPlayer(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	playerID, ok := args["player_id"].(string)
	if !ok {
		return "", fmt.Errorf("player_id is required")
	}

	name, ok := args["name"].(string)
	if !ok {
		return "", fmt.Errorf("name is required")
	}

	requestBody := map[string]interface{}{
		"playerID": playerID,
		"name":     name,
	}

	// Optional fields
	if seedIdx, ok := args["seed_idx"]; ok {
		requestBody["seedIdx"] = seedIdx
	}
	if photoUrl, ok := args["photo_url"]; ok {
		requestBody["photoUrl"] = photoUrl
	}
	if profileInfo, ok := args["profile_info"]; ok {
		requestBody["profileInfo"] = profileInfo
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/players", tournamentID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to add player: %w", err)
	}

	var player map[string]interface{}
	if err := json.Unmarshal(data, &player); err != nil {
		return "", fmt.Errorf("failed to parse player response: %w", err)
	}

	// Enrich player with profile details
	enrichedPlayer := enrichPlayerDataWithProfileDetails(player)

	jsonData, err := json.MarshalIndent(enrichedPlayer, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Update an existing player
func updatePlayer(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	playerID, ok := args["player_id"].(string)
	if !ok {
		return "", fmt.Errorf("player_id is required")
	}

	name, ok := args["name"].(string)
	if !ok {
		return "", fmt.Errorf("name is required")
	}

	requestBody := map[string]interface{}{
		"name": name,
	}

	// Optional fields
	if photoUrl, ok := args["photo_url"]; ok {
		requestBody["photoUrl"] = photoUrl
	}
	if profileInfo, ok := args["profile_info"]; ok {
		requestBody["profileInfo"] = profileInfo
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/players/%s", tournamentID, playerID)

	data, err := makeAPIRequest("PUT", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to update player: %w", err)
	}

	var player map[string]interface{}
	if err := json.Unmarshal(data, &player); err != nil {
		return "", fmt.Errorf("failed to parse player response: %w", err)
	}

	// Enrich player with profile details
	enrichedPlayer := enrichPlayerDataWithProfileDetails(player)

	jsonData, err := json.MarshalIndent(enrichedPlayer, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Delete a player from a tournament
func deletePlayer(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	playerID, ok := args["player_id"].(string)
	if !ok {
		return "", fmt.Errorf("player_id is required")
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/players/%s", tournamentID, playerID)

	data, err := makeAPIRequest("DELETE", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to delete player: %w", err)
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

// Reseed a player in a tournament
func reseedPlayer(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	playerID, ok := args["player_id"].(string)
	if !ok {
		return "", fmt.Errorf("player_id is required")
	}

	requestBody := map[string]interface{}{}

	if seedIdx, ok := args["seed_idx"]; ok {
		requestBody["seedIdx"] = seedIdx
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/players/%s/reseed", tournamentID, playerID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to reseed player: %w", err)
	}

	var player map[string]interface{}
	if err := json.Unmarshal(data, &player); err != nil {
		return "", fmt.Errorf("failed to parse player response: %w", err)
	}

	// Enrich player with profile details
	enrichedPlayer := enrichPlayerDataWithProfileDetails(player)

	jsonData, err := json.MarshalIndent(enrichedPlayer, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Randomize the seeding of a tournament
func randomizeSeeding(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	requestBody := map[string]interface{}{}

	if includeByePlayers, ok := args["include_bye_players"].(bool); ok {
		requestBody["includeByePlayers"] = includeByePlayers
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/bulkPlayers/randomize", tournamentID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to randomize seeding: %w", err)
	}

	var players []interface{}
	if err := json.Unmarshal(data, &players); err != nil {
		return "", fmt.Errorf("failed to parse players response: %w", err)
	}

	// Enrich each player with profile details
	enrichedPlayers := make([]interface{}, len(players))
	for i, p := range players {
		if player, ok := p.(map[string]interface{}); ok {
			enrichedPlayers[i] = enrichPlayerDataWithProfileDetails(player)
		} else {
			enrichedPlayers[i] = p
		}
	}

	result := map[string]interface{}{
		"players": enrichedPlayers,
		"count":   len(enrichedPlayers),
		"note":    "Player display names and profile details are included for better readability",
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Bulk update players (complete list replacement)
func bulkUpdatePlayers(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	participants, ok := args["participants"].([]interface{})
	if !ok {
		return "", fmt.Errorf("participants is required")
	}

	nonParticipants, ok := args["non_participants"].([]interface{})
	if !ok {
		return "", fmt.Errorf("non_participants is required")
	}

	requestBody := map[string]interface{}{
		"participants":    participants,
		"nonParticipants": nonParticipants,
	}

	if includeByePlayers, ok := args["include_bye_players"].(bool); ok {
		requestBody["includeByePlayers"] = includeByePlayers
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/bulkPlayers/update", tournamentID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to bulk update players: %w", err)
	}

	var players []interface{}
	if err := json.Unmarshal(data, &players); err != nil {
		return "", fmt.Errorf("failed to parse players response: %w", err)
	}

	// Enrich each player with profile details
	enrichedPlayers := make([]interface{}, len(players))
	for i, p := range players {
		if player, ok := p.(map[string]interface{}); ok {
			enrichedPlayers[i] = enrichPlayerDataWithProfileDetails(player)
		} else {
			enrichedPlayers[i] = p
		}
	}

	result := map[string]interface{}{
		"players": enrichedPlayers,
		"count":   len(enrichedPlayers),
		"note":    "Player display names and profile details are included for better readability",
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Check a player into their upcoming match
func checkinPlayer(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	playerID, ok := args["player_id"].(string)
	if !ok {
		return "", fmt.Errorf("player_id is required")
	}

	checkInStatus, ok := args["check_in_status"].(string)
	if !ok {
		return "", fmt.Errorf("check_in_status is required")
	}

	requestBody := map[string]interface{}{
		"checkInStatus": checkInStatus,
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/players/%s/checkin", tournamentID, playerID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to check in player: %w", err)
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

// Disqualify a player from a tournament
func disqualifyPlayer(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	playerID, ok := args["player_id"].(string)
	if !ok {
		return "", fmt.Errorf("player_id is required")
	}

	requestBody := map[string]interface{}{}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/players/%s/disqualify", tournamentID, playerID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to disqualify player: %w", err)
	}

	var player map[string]interface{}
	if err := json.Unmarshal(data, &player); err != nil {
		return "", fmt.Errorf("failed to parse player response: %w", err)
	}

	// Enrich player with profile details
	enrichedPlayer := enrichPlayerDataWithProfileDetails(player)

	jsonData, err := json.MarshalIndent(enrichedPlayer, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}
