package main

import (
	"encoding/json"
	"fmt"
)

// handleLocationsTool handles all location operations
func handleLocationsTool(args map[string]interface{}) (string, error) {
	operation, ok := args["operation"].(string)
	if !ok {
		return "", fmt.Errorf("operation parameter is required")
	}

	// Check if operation is allowed in current tools mode
	if !isOperationAllowed("truefinals_locations", operation) {
		return "", fmt.Errorf(getOperationNotAllowedError(operation))
	}

	switch operation {
	case "list":
		return listLocations(args)
	case "get":
		return getLocation(args)
	case "add":
		return addLocation(args)
	case "update":
		return updateLocation(args)
	case "delete":
		return deleteLocation(args)
	case "start_game":
		return startLocationGame(args)
	case "stop_game":
		return stopLocationGame(args)
	case "update_game_scores":
		return updateLocationGameScores(args)
	default:
		return "", fmt.Errorf("unknown operation: %s", operation)
	}
}

// getLocationsToolInfo returns the tool definition for location operations
func getLocationsToolInfo() ToolInfo {
	return ToolInfo{
		Name: "truefinals_locations",
		Description: `Manage tournament venue locations and combat cages in TrueFinals. This tool handles NHRL's multi-cage setup and match assignment system.

NHRL tournaments typically run with 4 combat cages operating simultaneously. This tool manages:
- Cage/arena assignments for matches
- Match queuing and scheduling at each cage
- Real-time match progression at venues
- Location-specific match controls

Use this tool when you need to:
- Set up tournament venues and cage configurations
- Assign matches to specific cages
- Control match flow at each location
- View cage queues and upcoming matches
- Activate matches when cages are ready

Note: NHRL typically uses Cage 1-4 as location identifiers.`,
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"operation": map[string]interface{}{
					"type": "string",
					"description": `The location operation to perform:

QUERY OPERATIONS:
- list: Get all locations for a tournament
- get: Get specific location details and match queue

LOCATION MANAGEMENT (require write access):
- create: Add a new location/cage to tournament
- update: Update location settings
- delete: Remove a location from tournament

MATCH CONTROL AT LOCATIONS:
- activate_next: Start the next queued match at this location
- update_queue: Reorder matches in location queue
- clear_queue: Remove all matches from location queue`,
					"enum": []string{
						"list", "get", "create", "update", "delete",
						"activate_next", "update_queue", "clear_queue",
					},
				},
				"tournament_id": map[string]interface{}{
					"type":        "string",
					"description": "Tournament identifier. Required for all operations. Format: 'nhrl_month##_weightclass'",
				},
				"location_id": map[string]interface{}{
					"type":        "string",
					"description": "Location/cage identifier. Required for single location operations. For NHRL, typically 'cage1', 'cage2', 'cage3', 'cage4'.",
				},
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Display name for the location. Examples: 'Cage 1', 'Main Arena', 'Blue Cage'",
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "Optional description of the location. Can include notes about equipment, streaming setup, etc.",
				},
				"queue": map[string]interface{}{
					"type":        "array",
					"description": "Ordered list of match IDs queued at this location. First match in array is next to be activated.",
				},
				"active_game_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the currently active match at this location. Only one match can be active per location.",
				},
			},
			"required": []string{"operation", "tournament_id"},
		},
	}
}

// List all locations in a tournament
func listLocations(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/locations", tournamentID)

	data, err := makeAPIRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to list locations: %w", err)
	}

	var locations []interface{}
	if err := json.Unmarshal(data, &locations); err != nil {
		return "", fmt.Errorf("failed to parse locations response: %w", err)
	}

	// Enrich each location with active game info
	enrichedLocations := make([]interface{}, len(locations))
	for i, l := range locations {
		if location, ok := l.(map[string]interface{}); ok {
			enrichedLocations[i] = enrichLocationWithGameInfo(location, tournamentID)
		} else {
			enrichedLocations[i] = l
		}
	}

	result := map[string]interface{}{
		"locations": enrichedLocations,
		"count":     len(enrichedLocations),
		"note":      "Active game information with player names is included for better readability",
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get a specific location by ID
func getLocation(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	locationID, ok := args["location_id"].(string)
	if !ok {
		return "", fmt.Errorf("location_id is required")
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/locations/%s", tournamentID, locationID)

	data, err := makeAPIRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get location: %w", err)
	}

	var location map[string]interface{}
	if err := json.Unmarshal(data, &location); err != nil {
		return "", fmt.Errorf("failed to parse location response: %w", err)
	}

	// Enrich location with active game info
	enrichedLocation := enrichLocationWithGameInfo(location, tournamentID)

	jsonData, err := json.MarshalIndent(enrichedLocation, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Add a new location to a tournament
func addLocation(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	locationID, ok := args["location_id"].(string)
	if !ok {
		return "", fmt.Errorf("location_id is required")
	}

	name, ok := args["name"].(string)
	if !ok {
		return "", fmt.Errorf("name is required")
	}

	blockActive, ok := args["block_active"].(bool)
	if !ok {
		return "", fmt.Errorf("block_active is required")
	}

	requestBody := map[string]interface{}{
		"locationID":  locationID,
		"name":        name,
		"blockActive": blockActive,
	}

	// Optional idx parameter
	if idx, ok := args["idx"].(float64); ok {
		requestBody["idx"] = int(idx)
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/locations", tournamentID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to add location: %w", err)
	}

	var location map[string]interface{}
	if err := json.Unmarshal(data, &location); err != nil {
		return "", fmt.Errorf("failed to parse location response: %w", err)
	}

	jsonData, err := json.MarshalIndent(location, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Update an existing location
func updateLocation(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	locationID, ok := args["location_id"].(string)
	if !ok {
		return "", fmt.Errorf("location_id is required")
	}

	name, ok := args["name"].(string)
	if !ok {
		return "", fmt.Errorf("name is required")
	}

	blockActive, ok := args["block_active"].(bool)
	if !ok {
		return "", fmt.Errorf("block_active is required")
	}

	requestBody := map[string]interface{}{
		"name":        name,
		"blockActive": blockActive,
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/locations/%s", tournamentID, locationID)

	data, err := makeAPIRequest("PUT", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to update location: %w", err)
	}

	var location map[string]interface{}
	if err := json.Unmarshal(data, &location); err != nil {
		return "", fmt.Errorf("failed to parse location response: %w", err)
	}

	jsonData, err := json.MarshalIndent(location, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Delete a location from a tournament
func deleteLocation(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	locationID, ok := args["location_id"].(string)
	if !ok {
		return "", fmt.Errorf("location_id is required")
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/locations/%s", tournamentID, locationID)

	data, err := makeAPIRequest("DELETE", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to delete location: %w", err)
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

// Start the first queued game at a location
func startLocationGame(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	locationID, ok := args["location_id"].(string)
	if !ok {
		return "", fmt.Errorf("location_id is required")
	}

	requestBody := map[string]interface{}{}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/locations/%s/startGame", tournamentID, locationID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to start game at location: %w", err)
	}

	var game map[string]interface{}
	if err := json.Unmarshal(data, &game); err != nil {
		return "", fmt.Errorf("failed to parse game response: %w", err)
	}

	// Enrich game with player and location names
	enrichedGame := enrichGameWithPlayerAndLocationInfo(game, tournamentID)

	// Also get the updated location info
	locationData, _ := makeAPIRequest("GET", fmt.Sprintf("/v1/tournaments/%s/locations/%s", tournamentID, locationID), nil)
	var location map[string]interface{}
	if err := json.Unmarshal(locationData, &location); err == nil {
		enrichedLocation := enrichLocationWithGameInfo(location, tournamentID)
		enrichedGame["locationInfo"] = enrichedLocation
	}

	jsonData, err := json.MarshalIndent(enrichedGame, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Stop the active game at a location
func stopLocationGame(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	locationID, ok := args["location_id"].(string)
	if !ok {
		return "", fmt.Errorf("location_id is required")
	}

	requestBody := map[string]interface{}{}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/locations/%s/stopGame", tournamentID, locationID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to stop game at location: %w", err)
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

// Update the score of the active (or first) game at a location
func updateLocationGameScores(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	locationID, ok := args["location_id"].(string)
	if !ok {
		return "", fmt.Errorf("location_id is required")
	}

	scores, ok := args["scores"].([]interface{})
	if !ok {
		return "", fmt.Errorf("scores is required")
	}

	requestBody := map[string]interface{}{
		"scores": scores,
	}

	// Optional result annotation
	if resultAnnotation, ok := args["result_annotation"]; ok {
		requestBody["resultAnnotation"] = resultAnnotation
	}

	endpoint := fmt.Sprintf("/v1/tournaments/%s/locations/%s/updateGameScores", tournamentID, locationID)

	data, err := makeAPIRequest("POST", endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to update game scores at location: %w", err)
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
