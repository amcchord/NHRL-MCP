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
		return "", fmt.Errorf(getOperationNotAllowedError(operation))
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
		Name: "truefinals_players",
		Description: `Manage tournament participants (bots and their operators) in TrueFinals. This tool handles registration, seeding, and participant management for NHRL tournaments.

In NHRL context, "players" represent combat robots and their human operators/teams. Each tournament has participants who compete in brackets based on weight class.

Use this tool when you need to:
- Add or remove tournament participants
- Set or update seeding for bracket generation
- Manage participant check-in status
- Handle disqualifications or withdrawals
- Update bot/team information
- Swap participants or correct entries

Note: Seeding is crucial for bracket fairness - higher seeds face lower seeds in early rounds.`,
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"operation": map[string]interface{}{
					"type": "string",
					"description": `The player/participant operation to perform:

QUERY OPERATIONS:
- list: Get all participants in a tournament
- get: Get detailed information about a specific participant

PARTICIPANT MANAGEMENT (require write access):
- add: Register a new bot/team to the tournament
- update: Update participant details (name, team info, etc.)
- delete: Remove a participant from the tournament
- set_seed: Assign or update seeding number (1 = top seed)
- swap: Exchange positions of two participants in bracket
- check_in: Mark participant as checked in and ready
- undo_check_in: Clear check-in status
- disqualify: Mark participant as disqualified
- undisqualify: Remove disqualification status`,
					"enum": []string{
						"list", "get", "add", "update", "delete",
						"set_seed", "swap", "check_in", "undo_check_in",
						"disqualify", "undisqualify",
					},
				},
				"tournament_id": map[string]interface{}{
					"type":        "string",
					"description": "Tournament identifier. Required for all operations. Format: 'nhrl_month##_weightclass'",
				},
				"player_id": map[string]interface{}{
					"type":        "string",
					"description": "Participant/player profile ID. Required for single participant operations.",
				},
				"player_id_1": map[string]interface{}{
					"type":        "string",
					"description": "First participant ID for swap operation. This participant will take the position of player_id_2.",
				},
				"player_id_2": map[string]interface{}{
					"type":        "string",
					"description": "Second participant ID for swap operation. This participant will take the position of player_id_1.",
				},
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Bot/robot name. Examples: 'Ripperoni', 'Bloodsport', 'Malice'",
				},
				"display_name": map[string]interface{}{
					"type":        "string",
					"description": "Display name for bracket/overlay. Can include team name or sponsors.",
				},
				"team_name": map[string]interface{}{
					"type":        "string",
					"description": "Team or operator name. Examples: 'Team Velocity', 'Chaos Corps'",
				},
				"seed": map[string]interface{}{
					"type":        "integer",
					"description": "Seeding number for bracket placement. 1 = highest seed (best ranking). Lower seeds face higher seeds in early rounds.",
				},
				"checked_in": map[string]interface{}{
					"type":        "boolean",
					"description": "Whether the participant has checked in for the tournament. Required before matches can begin.",
				},
				"disqualified": map[string]interface{}{
					"type":        "boolean",
					"description": "Whether the participant is disqualified. Disqualified bots cannot compete but remain in bracket.",
				},
				"profile_data": map[string]interface{}{
					"type":        "object",
					"description": "Additional participant data including contact info, bot specifications, sponsors, etc.",
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
