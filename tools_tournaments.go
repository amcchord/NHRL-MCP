package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// handleTournamentsTool handles all tournament operations
func handleTournamentsTool(args map[string]interface{}) (string, error) {
	operation, ok := args["operation"].(string)
	if !ok {
		return "", fmt.Errorf("operation parameter is required")
	}

	// Check if operation is allowed in current tools mode
	if !isOperationAllowed("truefinals_tournaments", operation) {
		return "", fmt.Errorf(getOperationNotAllowedError(operation))
	}

	switch operation {
	case "list":
		return listTournaments(args)
	case "get":
		return getTournament(args)
	case "details":
		return getTournamentDetails(args)
	case "format":
		return getTournamentFormat(args)
	case "overlay_params":
		return getTournamentOverlayParams(args)
	case "description":
		return getTournamentDescription(args)
	case "private":
		return getTournamentPrivate(args)
	case "webhooks":
		return getTournamentWebhooks(args)
	case "create":
		return createTournament(args)
	case "update":
		return updateTournament(args)
	case "update_description":
		return updateTournamentDescription(args)
	case "update_overlay_params":
		return updateTournamentOverlayParams(args)
	case "update_webhooks":
		return updateTournamentWebhooks(args)
	case "start":
		return startTournament(args)
	case "reset":
		return resetTournament(args)
	case "push_schedule":
		return pushTournamentSchedule(args)
	case "delete":
		return deleteTournament(args)
	default:
		return "", fmt.Errorf("unknown operation: %s", operation)
	}
}

// getTournamentsToolInfo returns the tool definition for tournament operations
func getTournamentsToolInfo() ToolInfo {
	return ToolInfo{
		Name: "truefinals_tournaments",
		Description: `Manage NHRL tournaments in the TrueFinals system. This is the primary tool for tournament administration and bracket management.

TrueFinals is NHRL's tournament management platform that handles:
- Tournament creation and configuration
- Bracket generation and management  
- Match scheduling and progression
- Player/bot registration and seeding
- Live scoring and results tracking

Use this tool when you need to:
- Create or modify tournaments
- View tournament brackets and standings
- Check tournament settings and participants
- Start tournaments or reset brackets
- Manage tournament progression

Note: Test tournaments (containing 'TEST' or 'test' in the name) are filtered out by default when listing. Use include_test_tournaments=true to show all tournaments.`,
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"operation": map[string]interface{}{
					"type": "string",
					"description": `The tournament operation to perform:

QUERY OPERATIONS (read-only):
- list: Get all tournaments you have access to (filters test tournaments by default)
- get: Get complete tournament data including bracket, games, and participants
- details: Get lightweight tournament info without full bracket data
- format: Get tournament format settings (single elim, double elim, round robin)
- overlay_params: Get streaming overlay configuration
- description: Get tournament description text
- private: Get private tournament data (webhooks, etc.)
- webhooks: Get configured webhooks for tournament events

MODIFICATION OPERATIONS (require write access):
- create: Create a new tournament with specified settings
- update: Update tournament basic settings (title, location, privacy)
- update_description: Update tournament description text
- update_overlay_params: Update streaming overlay settings
- update_webhooks: Update webhook configurations

TOURNAMENT CONTROL:
- start: Start the tournament (locks bracket and begins matches)
- reset: Reset tournament bracket (bracket_only or all)
- push_schedule: Delay all scheduled matches by specified minutes
- delete: Delete the tournament completely`,
					"enum": []string{
						"list", "get", "details", "format", "overlay_params", "description", "private", "webhooks",
						"create", "update", "update_description", "update_overlay_params", "update_webhooks",
						"start", "reset", "push_schedule", "delete",
					},
				},
				"tournament_id": map[string]interface{}{
					"type":        "string",
					"description": "Unique tournament identifier. Required for all operations except 'list' and 'create'. Format is typically lowercase with underscores (e.g., 'nhrl_dec24_3lb')",
				},
				// Tournament creation/update fields
				"title": map[string]interface{}{
					"type":        "string",
					"description": "Tournament display title (max 64 characters). Example: 'NHRL December 2024 - 3lb Beetleweight'",
				},
				"creator_profile_id": map[string]interface{}{
					"type":        "string",
					"description": "Profile ID of tournament creator. Required for create operation.",
				},
				"game_title_info": map[string]interface{}{
					"type":        "object",
					"description": "Game/sport information for the tournament. For NHRL, typically includes weight class info.",
				},
				"event_location": map[string]interface{}{
					"type":        "string",
					"description": "Physical venue location (max 128 characters). Example: 'House of Havoc - Norwalk, CT'",
				},
				"scheduled_start_time": map[string]interface{}{
					"type":        "integer",
					"description": "Tournament start time as Unix timestamp (seconds since epoch). Used for scheduling.",
				},
				"privacy": map[string]interface{}{
					"type":        "string",
					"description": "Tournament visibility: 'public' (visible to all), 'unlisted' (accessible via link), 'private' (invite only)",
					"enum":        []string{"public", "unlisted", "private"},
				},
				"display_check_in_status": map[string]interface{}{
					"type":        "boolean",
					"description": "Whether to show player check-in status on bracket. Useful for tracking attendance.",
				},
				"logo_url": map[string]interface{}{
					"type":        "string",
					"description": "URL for tournament logo image. Displayed on bracket and overlays.",
				},
				"thumbnail_url": map[string]interface{}{
					"type":        "string",
					"description": "URL for tournament thumbnail image. Used in tournament listings.",
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "Detailed tournament description (max 50000 characters). Can include rules, schedule, prizes, etc. Supports markdown formatting.",
				},
				"format_options": map[string]interface{}{
					"type":        "object",
					"description": "Tournament format configuration. Specify type (single_elimination, double_elimination, round_robin) and related settings like third place match, group stage, etc.",
				},
				"participants": map[string]interface{}{
					"type":        "array",
					"description": "List of tournament participants with seeding information. Each entry includes player/team data and optional seed number.",
				},
				"non_participants": map[string]interface{}{
					"type":        "array",
					"description": "List of registered players not participating in bracket (e.g., alternates, staff).",
				},
				"locations": map[string]interface{}{
					"type":        "array",
					"description": "Tournament venue locations. For NHRL, typically includes cage/arena assignments.",
				},
				"webhooks": map[string]interface{}{
					"type":        "array",
					"description": "Webhook configurations for tournament events (match start, match end, tournament complete, etc.)",
				},
				"overlay_params": map[string]interface{}{
					"type":        "object",
					"description": "Streaming overlay parameters including colors, fonts, sponsor logos, and layout settings.",
				},
				"score_updates": map[string]interface{}{
					"type":        "object",
					"description": "Configuration for how scores are updated (manual, API integration, etc.)",
				},
				"reset_mode": map[string]interface{}{
					"type":        "string",
					"description": "Reset scope for tournament reset: 'bracket_only' (keep players, reset matches) or 'all' (full reset)",
					"enum":        []string{"bracket_only", "all"},
				},
				"delay_minutes": map[string]interface{}{
					"type":        "integer",
					"description": "Number of minutes to delay all scheduled matches. Used with push_schedule operation.",
				},
				"include_test_tournaments": map[string]interface{}{
					"type":        "boolean",
					"description": "Include test tournaments in list results. Default: false (test tournaments are hidden)",
				},
			},
			"required": []string{"operation"},
		},
	}
}

// List tournaments owned by the API user
func listTournaments(args map[string]interface{}) (string, error) {
	endpoint := "/v1/user/tournaments"

	data, err := makeAPIRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to list tournaments: %w", err)
	}

	var tournaments TournamentListResponse
	if err := json.Unmarshal(data, &tournaments); err != nil {
		return "", fmt.Errorf("failed to parse tournaments response: %w", err)
	}

	// Check if we should include test tournaments
	includeTestTournaments := false
	if include, ok := args["include_test_tournaments"].(bool); ok {
		includeTestTournaments = include
	}

	// Filter out test tournaments unless specifically requested
	var filteredTournaments TournamentListResponse
	for _, tournament := range tournaments {
		// Check if title contains "TEST" or "test"
		if !includeTestTournaments && (strings.Contains(tournament.Title, "TEST") || strings.Contains(tournament.Title, "test")) {
			continue
		}
		filteredTournaments = append(filteredTournaments, tournament)
	}

	result := map[string]interface{}{
		"tournaments": filteredTournaments,
		"count":       len(filteredTournaments),
	}

	// Add a note if any tournaments were filtered
	if len(tournaments) > len(filteredTournaments) {
		result["note"] = fmt.Sprintf("%d test tournament(s) filtered out. Use include_test_tournaments=true to show all.", len(tournaments)-len(filteredTournaments))
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get tournament by ID
func getTournament(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s", tournamentID)

	data, err := makeAPIRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get tournament: %w", err)
	}

	// Pretty print the response
	var tournament map[string]interface{}
	if err := json.Unmarshal(data, &tournament); err != nil {
		return "", fmt.Errorf("failed to parse tournament response: %w", err)
	}

	// Enrich tournament data with human-readable information
	enrichedTournament := enrichTournamentData(tournament)

	jsonData, err := json.MarshalIndent(enrichedTournament, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get tournament details (lighter version)
func getTournamentDetails(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/details", tournamentID)

	data, err := makeAPIRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get tournament details: %w", err)
	}

	var details map[string]interface{}
	if err := json.Unmarshal(data, &details); err != nil {
		return "", fmt.Errorf("failed to parse tournament details response: %w", err)
	}

	jsonData, err := json.MarshalIndent(details, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get tournament format
func getTournamentFormat(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/format", tournamentID)

	data, err := makeAPIRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get tournament format: %w", err)
	}

	var format map[string]interface{}
	if err := json.Unmarshal(data, &format); err != nil {
		return "", fmt.Errorf("failed to parse tournament format response: %w", err)
	}

	jsonData, err := json.MarshalIndent(format, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get tournament overlay parameters
func getTournamentOverlayParams(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/overlayParams", tournamentID)

	data, err := makeAPIRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get tournament overlay params: %w", err)
	}

	var overlayParams map[string]interface{}
	if err := json.Unmarshal(data, &overlayParams); err != nil {
		return "", fmt.Errorf("failed to parse overlay params response: %w", err)
	}

	jsonData, err := json.MarshalIndent(overlayParams, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get tournament description
func getTournamentDescription(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/description", tournamentID)

	data, err := makeAPIRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get tournament description: %w", err)
	}

	// This endpoint returns a JSON string, so parse it appropriately
	var description interface{}
	if err := json.Unmarshal(data, &description); err != nil {
		return "", fmt.Errorf("failed to parse description response: %w", err)
	}

	result := map[string]interface{}{
		"description": description,
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get tournament private data (webhooks)
func getTournamentPrivate(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/private", tournamentID)

	data, err := makeAPIRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get tournament private data: %w", err)
	}

	var privateData map[string]interface{}
	if err := json.Unmarshal(data, &privateData); err != nil {
		return "", fmt.Errorf("failed to parse private data response: %w", err)
	}

	jsonData, err := json.MarshalIndent(privateData, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get tournament webhooks
func getTournamentWebhooks(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/private/webhooks", tournamentID)

	data, err := makeAPIRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get tournament webhooks: %w", err)
	}

	var webhooks []interface{}
	if err := json.Unmarshal(data, &webhooks); err != nil {
		return "", fmt.Errorf("failed to parse webhooks response: %w", err)
	}

	result := map[string]interface{}{
		"webhooks": webhooks,
		"count":    len(webhooks),
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Create a new tournament
func createTournament(args map[string]interface{}) (string, error) {
	// Build the request body from args
	requestBody := map[string]interface{}{}

	// Required fields
	if tournamentID, ok := args["tournament_id"].(string); ok {
		requestBody["tournamentID"] = tournamentID
	}
	if creatorProfileID, ok := args["creator_profile_id"].(string); ok {
		requestBody["creatorProfileID"] = creatorProfileID
	}
	if title, ok := args["title"].(string); ok {
		requestBody["title"] = title
	}
	if eventLocation, ok := args["event_location"].(string); ok {
		requestBody["eventLocation"] = eventLocation
	}
	if privacy, ok := args["privacy"].(string); ok {
		requestBody["privacy"] = privacy
	}

	// Optional fields
	if gameTitleInfo, ok := args["game_title_info"]; ok {
		requestBody["gameTitleInfo"] = gameTitleInfo
	}
	if scheduledStartTime, ok := args["scheduled_start_time"]; ok {
		requestBody["scheduledStartTime"] = scheduledStartTime
	}
	if displayCheckInStatus, ok := args["display_check_in_status"]; ok {
		requestBody["displayCheckInStatus"] = displayCheckInStatus
	}
	if logoURL, ok := args["logo_url"]; ok {
		requestBody["logoUrl"] = logoURL
	}
	if thumbnailURL, ok := args["thumbnail_url"]; ok {
		requestBody["thumbnailUrl"] = thumbnailURL
	}
	if description, ok := args["description"]; ok {
		requestBody["description"] = description
	}
	if locations, ok := args["locations"]; ok {
		requestBody["locations"] = locations
	}
	if webhooks, ok := args["webhooks"]; ok {
		requestBody["webhooks"] = webhooks
	}
	if formatOptions, ok := args["format_options"]; ok {
		requestBody["formatOptions"] = formatOptions
	}
	if participants, ok := args["participants"]; ok {
		requestBody["participants"] = participants
	}
	if nonParticipants, ok := args["non_participants"]; ok {
		requestBody["nonParticipants"] = nonParticipants
	}

	endpoint := "/v1/tournaments"

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create tournament: %w", err)
	}

	var tournament map[string]interface{}
	if err := json.Unmarshal(data, &tournament); err != nil {
		return "", fmt.Errorf("failed to parse tournament response: %w", err)
	}

	// Enrich tournament data with human-readable information
	enrichedTournament := enrichTournamentData(tournament)

	jsonData, err := json.MarshalIndent(enrichedTournament, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Update tournament settings
func updateTournament(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	// Build the request body from args (similar to create but without tournament_id in body)
	requestBody := map[string]interface{}{}

	if creatorProfileID, ok := args["creator_profile_id"].(string); ok {
		requestBody["creatorProfileID"] = creatorProfileID
	}
	if title, ok := args["title"].(string); ok {
		requestBody["title"] = title
	}
	if eventLocation, ok := args["event_location"].(string); ok {
		requestBody["eventLocation"] = eventLocation
	}
	if privacy, ok := args["privacy"].(string); ok {
		requestBody["privacy"] = privacy
	}

	// Add other optional fields as needed...

	endpoint := fmt.Sprintf("/v1/tournaments/%s", tournamentID)

	data, err := makeAPIRequest("PUT", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to update tournament: %w", err)
	}

	var tournament map[string]interface{}
	if err := json.Unmarshal(data, &tournament); err != nil {
		return "", fmt.Errorf("failed to parse tournament response: %w", err)
	}

	// Enrich tournament data with human-readable information
	enrichedTournament := enrichTournamentData(tournament)

	jsonData, err := json.MarshalIndent(enrichedTournament, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Update tournament description
func updateTournamentDescription(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	description, ok := args["description"].(string)
	if !ok {
		return "", fmt.Errorf("description is required")
	}

	requestBody := map[string]interface{}{
		"description": description,
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/description", tournamentID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to update tournament description: %w", err)
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

// Update tournament overlay parameters
func updateTournamentOverlayParams(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	overlayParams, ok := args["overlay_params"]
	if !ok {
		return "", fmt.Errorf("overlay_params is required")
	}

	requestBody := map[string]interface{}{
		"overlayParams": overlayParams,
	}

	if scoreUpdates, ok := args["score_updates"]; ok {
		requestBody["scoreUpdates"] = scoreUpdates
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/overlayParams", tournamentID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to update overlay params: %w", err)
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

// Update tournament webhooks
func updateTournamentWebhooks(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	webhooks, ok := args["webhooks"]
	if !ok {
		return "", fmt.Errorf("webhooks is required")
	}

	requestBody := map[string]interface{}{
		"webhooks": webhooks,
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/private/webhooks", tournamentID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to update webhooks: %w", err)
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

// Start a tournament
func startTournament(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	requestBody := map[string]interface{}{}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/start", tournamentID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to start tournament: %w", err)
	}

	var tournament map[string]interface{}
	if err := json.Unmarshal(data, &tournament); err != nil {
		return "", fmt.Errorf("failed to parse tournament response: %w", err)
	}

	// Enrich tournament data with human-readable information
	enrichedTournament := enrichTournamentData(tournament)

	jsonData, err := json.MarshalIndent(enrichedTournament, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Reset a tournament
func resetTournament(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	mode, ok := args["reset_mode"].(string)
	if !ok {
		mode = "bracket_only" // Default mode
	}

	requestBody := map[string]interface{}{
		"mode": mode,
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/reset", tournamentID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to reset tournament: %w", err)
	}

	var tournament map[string]interface{}
	if err := json.Unmarshal(data, &tournament); err != nil {
		return "", fmt.Errorf("failed to parse tournament response: %w", err)
	}

	// Enrich tournament data with human-readable information
	enrichedTournament := enrichTournamentData(tournament)

	jsonData, err := json.MarshalIndent(enrichedTournament, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Push tournament game schedule
func pushTournamentSchedule(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	delayMinutes, ok := args["delay_minutes"].(float64)
	if !ok {
		return "", fmt.Errorf("delay_minutes is required")
	}

	requestBody := map[string]interface{}{
		"delayMinutes": int(delayMinutes),
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/pushGameSchedule", tournamentID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to push tournament schedule: %w", err)
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

// Delete a tournament
func deleteTournament(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s", tournamentID)

	data, err := makeAPIRequest("DELETE", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to delete tournament: %w", err)
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
