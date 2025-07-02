package main

import (
	"encoding/json"
	"fmt"
	"sort"
)

// handleBracketTool handles bracket visualization operations
func handleBracketTool(args map[string]interface{}) (string, error) {
	operation, ok := args["operation"].(string)
	if !ok {
		return "", fmt.Errorf("operation parameter is required")
	}

	// Check if operation is allowed in current tools mode
	if !isOperationAllowed("truefinals_bracket", operation) {
		return "", fmt.Errorf("operation '%s' not allowed in '%s' mode", operation, toolsMode)
	}

	switch operation {
	case "get":
		return getBracket(args)
	case "get_round":
		return getBracketRound(args)
	case "get_standings":
		return getBracketStandings(args)
	default:
		return "", fmt.Errorf("unknown operation: %s", operation)
	}
}

// getBracketToolInfo returns the tool definition for bracket operations
func getBracketToolInfo() ToolInfo {
	return ToolInfo{
		Name:        "truefinals_bracket",
		Description: "Get tournament bracket information formatted for display. Returns structured data showing tournament progression, current matches, and results.",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"operation": map[string]interface{}{
					"type":        "string",
					"description": "The operation to perform",
					"enum":        []string{"get", "get_round", "get_standings"},
				},
				"tournament_id": map[string]interface{}{
					"type":        "string",
					"description": "Tournament ID - required for all operations",
				},
				"round": map[string]interface{}{
					"type":        "integer",
					"description": "Specific round number to get (for get_round operation)",
				},
				"bracket_type": map[string]interface{}{
					"type":        "string",
					"description": "Filter by bracket type",
					"enum":        []string{"winners", "losers", "all"},
				},
			},
			"required": []string{"operation", "tournament_id"},
		},
	}
}

// BracketRound represents a round in the bracket
type BracketRound struct {
	Round       int                      `json:"round"`
	RoundName   string                   `json:"roundName"`
	BracketType string                   `json:"bracketType"`
	Games       []map[string]interface{} `json:"games"`
}

// Get complete bracket information
func getBracket(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	// Get full tournament data
	endpoint := fmt.Sprintf("/v1/tournaments/%s", tournamentID)
	data, err := makeAPIRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get tournament: %w", err)
	}

	var tournament map[string]interface{}
	if err := json.Unmarshal(data, &tournament); err != nil {
		return "", fmt.Errorf("failed to parse tournament response: %w", err)
	}

	// Get format info
	format, _ := tournament["format"].(map[string]interface{})
	formatType, _ := format["type"].(string)

	// Build player map for easy lookup
	playerMap := make(map[string]map[string]interface{})
	if players, ok := tournament["players"].([]interface{}); ok {
		for _, p := range players {
			if player, ok := p.(map[string]interface{}); ok {
				if id, ok := player["id"].(string); ok {
					// Enrich player data
					enrichedPlayer := enrichPlayerDataWithProfileDetails(player)
					playerMap[id] = enrichedPlayer
				}
			}
		}
	}

	// Process games and organize by rounds
	var allGames []map[string]interface{}
	if games, ok := tournament["games"].([]interface{}); ok {
		for _, g := range games {
			if game, ok := g.(map[string]interface{}); ok {
				enrichedGame := enrichGameForBracket(game, playerMap)
				allGames = append(allGames, enrichedGame)
			}
		}
	}

	// Organize games by round and bracket type
	roundsMap := make(map[string][]map[string]interface{})
	for _, game := range allGames {
		round := int(game["round"].(float64))

		// Determine bracket type
		bracketType := "main"
		if formatType == "double_elimination" {
			if round < 0 {
				bracketType = "losers"
			} else {
				bracketType = "winners"
			}
		}

		key := fmt.Sprintf("%s-%d", bracketType, abs(round))
		roundsMap[key] = append(roundsMap[key], game)
	}

	// Convert to structured rounds
	var rounds []BracketRound
	for key, games := range roundsMap {
		var bracketType string
		var round int
		fmt.Sscanf(key, "%[^-]-%d", &bracketType, &round)

		// Sort games within round by name
		sort.Slice(games, func(i, j int) bool {
			return games[i]["name"].(string) < games[j]["name"].(string)
		})

		rounds = append(rounds, BracketRound{
			Round:       round,
			RoundName:   getRoundName(round, bracketType, formatType),
			BracketType: bracketType,
			Games:       games,
		})
	}

	// Sort rounds
	sort.Slice(rounds, func(i, j int) bool {
		if rounds[i].BracketType != rounds[j].BracketType {
			return rounds[i].BracketType == "winners"
		}
		return rounds[i].Round < rounds[j].Round
	})

	// Build bracket summary
	bracketSummary := map[string]interface{}{
		"tournamentID":   tournamentID,
		"tournamentName": tournament["title"],
		"format":         formatType,
		"status":         getTournamentStatus(tournament),
		"rounds":         rounds,
		"playerCount":    len(playerMap),
		"totalGames":     len(allGames),
		"completedGames": countCompletedGames(allGames),
		"activeGames":    countActiveGames(allGames),
		"currentRound":   getCurrentRound(allGames),
		"champions":      getChampions(allGames, playerMap),
		"displayTips": map[string]string{
			"rounds":     "Games are organized by round and bracket type (winners/losers for double elimination)",
			"gameStates": "Game states: 'unavailable' = waiting for previous games, 'available' = ready to play, 'active' = in progress, 'done' = completed",
			"scores":     "Scores of -1 indicate a player hasn't competed yet or was eliminated",
			"navigation": "Use round numbers to focus on specific rounds, negative rounds are losers bracket",
		},
	}

	jsonData, err := json.MarshalIndent(bracketSummary, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get specific round information
func getBracketRound(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	roundNum, ok := args["round"].(float64)
	if !ok {
		return "", fmt.Errorf("round is required")
	}
	round := int(roundNum)

	bracketType := "all"
	if bt, ok := args["bracket_type"].(string); ok {
		bracketType = bt
	}

	// Get full tournament data
	endpoint := fmt.Sprintf("/v1/tournaments/%s", tournamentID)
	data, err := makeAPIRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get tournament: %w", err)
	}

	var tournament map[string]interface{}
	if err := json.Unmarshal(data, &tournament); err != nil {
		return "", fmt.Errorf("failed to parse tournament response: %w", err)
	}

	// Build player map
	playerMap := make(map[string]map[string]interface{})
	if players, ok := tournament["players"].([]interface{}); ok {
		for _, p := range players {
			if player, ok := p.(map[string]interface{}); ok {
				if id, ok := player["id"].(string); ok {
					enrichedPlayer := enrichPlayerDataWithProfileDetails(player)
					playerMap[id] = enrichedPlayer
				}
			}
		}
	}

	// Filter games for the specific round
	var roundGames []map[string]interface{}
	if games, ok := tournament["games"].([]interface{}); ok {
		for _, g := range games {
			if game, ok := g.(map[string]interface{}); ok {
				gameRound := int(game["round"].(float64))
				if abs(gameRound) == abs(round) {
					// Check bracket type filter
					if bracketType != "all" {
						isLosers := gameRound < 0
						if (bracketType == "losers" && !isLosers) || (bracketType == "winners" && isLosers) {
							continue
						}
					}

					enrichedGame := enrichGameForBracket(game, playerMap)
					roundGames = append(roundGames, enrichedGame)
				}
			}
		}
	}

	// Sort games by name
	sort.Slice(roundGames, func(i, j int) bool {
		return roundGames[i]["name"].(string) < roundGames[j]["name"].(string)
	})

	format, _ := tournament["format"].(map[string]interface{})
	formatType, _ := format["type"].(string)

	result := map[string]interface{}{
		"tournamentID":   tournamentID,
		"tournamentName": tournament["title"],
		"round":          round,
		"roundName":      getRoundName(abs(round), getBracketTypeFromRound(round), formatType),
		"bracketType":    getBracketTypeFromRound(round),
		"games":          roundGames,
		"gameCount":      len(roundGames),
		"completedCount": countCompletedGames(roundGames),
		"activeCount":    countActiveGames(roundGames),
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Get current standings
func getBracketStandings(args map[string]interface{}) (string, error) {
	tournamentID, ok := args["tournament_id"].(string)
	if !ok {
		return "", fmt.Errorf("tournament_id is required")
	}

	// Get full tournament data
	endpoint := fmt.Sprintf("/v1/tournaments/%s", tournamentID)
	data, err := makeAPIRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get tournament: %w", err)
	}

	var tournament map[string]interface{}
	if err := json.Unmarshal(data, &tournament); err != nil {
		return "", fmt.Errorf("failed to parse tournament response: %w", err)
	}

	// Get and enrich players with their placement
	var standings []map[string]interface{}
	if players, ok := tournament["players"].([]interface{}); ok {
		for _, p := range players {
			if player, ok := p.(map[string]interface{}); ok {
				// Skip bye players
				if isBye, ok := player["isBye"].(bool); ok && isBye {
					continue
				}

				enrichedPlayer := enrichPlayerDataWithProfileDetails(player)

				// Create standing entry
				standing := map[string]interface{}{
					"playerID":       player["id"],
					"displayName":    enrichedPlayer["displayName"],
					"name":           player["name"],
					"placement":      player["placement"],
					"wins":           player["wins"],
					"losses":         player["losses"],
					"ties":           player["ties"],
					"isDisqualified": player["isDisqualified"],
					"seed":           player["seed"],
				}

				// Add profile info if available
				if tag, ok := enrichedPlayer["tag"]; ok {
					standing["tag"] = tag
				}

				standings = append(standings, standing)
			}
		}
	}

	// Sort by placement (nil placement goes to end)
	sort.Slice(standings, func(i, j int) bool {
		placeI := standings[i]["placement"]
		placeJ := standings[j]["placement"]

		if placeI == nil && placeJ == nil {
			// Both nil, sort by wins then losses
			winsI := standings[i]["wins"].(float64)
			winsJ := standings[j]["wins"].(float64)
			if winsI != winsJ {
				return winsI > winsJ
			}
			lossesI := standings[i]["losses"].(float64)
			lossesJ := standings[j]["losses"].(float64)
			return lossesI < lossesJ
		}
		if placeI == nil {
			return false
		}
		if placeJ == nil {
			return true
		}
		return placeI.(float64) < placeJ.(float64)
	})

	result := map[string]interface{}{
		"tournamentID":   tournamentID,
		"tournamentName": tournament["title"],
		"status":         getTournamentStatus(tournament),
		"standings":      standings,
		"playerCount":    len(standings),
		"note":           "Players are sorted by placement. Players without placement are sorted by record.",
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonData), nil
}

// Helper function to enrich game data for bracket display
func enrichGameForBracket(game map[string]interface{}, playerMap map[string]map[string]interface{}) map[string]interface{} {
	enrichedGame := make(map[string]interface{})

	// Copy basic game info
	enrichedGame["id"] = game["id"]
	enrichedGame["name"] = game["name"]
	enrichedGame["round"] = game["round"]
	enrichedGame["bracketID"] = game["bracketID"]
	enrichedGame["state"] = game["state"]
	enrichedGame["scoreToWin"] = game["scoreToWin"]

	// Add timing info
	if endTime, ok := game["endTime"]; ok {
		enrichedGame["endTime"] = endTime
	}
	if scheduledTime, ok := game["scheduledTime"]; ok {
		enrichedGame["scheduledTime"] = scheduledTime
	}

	// Process slots with player info
	if slots, ok := game["slots"].([]interface{}); ok {
		enrichedSlots := make([]map[string]interface{}, len(slots))
		for i, s := range slots {
			if slot, ok := s.(map[string]interface{}); ok {
				enrichedSlot := map[string]interface{}{
					"slotIdx":   slot["slotIdx"],
					"score":     slot["score"],
					"slotState": slot["slotState"],
				}

				// Add player info if available
				if playerID, ok := slot["playerID"].(string); ok {
					enrichedSlot["playerID"] = playerID
					if player, found := playerMap[playerID]; found {
						enrichedSlot["playerName"] = player["name"]
						enrichedSlot["displayName"] = player["displayName"]
						if seed, ok := player["seed"]; ok {
							enrichedSlot["seed"] = seed
						}
					}
				}

				// Add previous game info
				if prevGameID, ok := slot["prevGameID"]; ok {
					enrichedSlot["prevGameID"] = prevGameID
				}

				enrichedSlots[i] = enrichedSlot
			}
		}
		enrichedGame["slots"] = enrichedSlots
	}

	// Add next game info
	if nextGameSlotIDs, ok := game["nextGameSlotIDs"]; ok {
		enrichedGame["nextGameSlotIDs"] = nextGameSlotIDs
	}

	// Add placement info
	if winnerPlacement, ok := game["winnerPlacement"]; ok {
		enrichedGame["winnerPlacement"] = winnerPlacement
	}
	if loserPlacement, ok := game["loserPlacement"]; ok {
		enrichedGame["loserPlacement"] = loserPlacement
	}

	// Determine winner if game is done
	if state, ok := game["state"].(string); ok && state == "done" {
		if slots, ok := enrichedGame["slots"].([]map[string]interface{}); ok && len(slots) >= 2 {
			score1, _ := slots[0]["score"].(float64)
			score2, _ := slots[1]["score"].(float64)
			if score1 > score2 {
				enrichedGame["winnerSlotIdx"] = 0
			} else if score2 > score1 {
				enrichedGame["winnerSlotIdx"] = 1
			}
		}
	}

	return enrichedGame
}

// Helper functions

func getRoundName(round int, bracketType string, formatType string) string {
	if formatType == "single_elimination" || bracketType == "winners" {
		switch round {
		case 1:
			return "Finals"
		case 2:
			return "Semifinals"
		case 3:
			return "Quarterfinals"
		default:
			return fmt.Sprintf("Round of %d", 1<<uint(round))
		}
	} else if bracketType == "losers" {
		// Losers bracket naming is more complex
		return fmt.Sprintf("Losers Round %d", round)
	}
	return fmt.Sprintf("Round %d", round)
}

func getTournamentStatus(tournament map[string]interface{}) string {
	if endTime, ok := tournament["endTime"]; ok && endTime != nil {
		return "completed"
	}
	if startTime, ok := tournament["startTime"]; ok && startTime != nil {
		return "in_progress"
	}
	return "pending"
}

func countCompletedGames(games []map[string]interface{}) int {
	count := 0
	for _, game := range games {
		if state, ok := game["state"].(string); ok && state == "done" {
			count++
		}
	}
	return count
}

func countActiveGames(games []map[string]interface{}) int {
	count := 0
	for _, game := range games {
		if state, ok := game["state"].(string); ok && (state == "active" || state == "called") {
			count++
		}
	}
	return count
}

func getCurrentRound(games []map[string]interface{}) map[string]interface{} {
	// Find the highest round with any non-completed games
	maxWinnersRound := 0
	maxLosersRound := 0

	for _, game := range games {
		if state, ok := game["state"].(string); ok && state != "done" {
			if round, ok := game["round"].(float64); ok {
				r := int(round)
				if r > 0 && r > maxWinnersRound {
					maxWinnersRound = r
				} else if r < 0 && abs(r) > maxLosersRound {
					maxLosersRound = abs(r)
				}
			}
		}
	}

	result := make(map[string]interface{})
	if maxWinnersRound > 0 {
		result["winnersRound"] = maxWinnersRound
	}
	if maxLosersRound > 0 {
		result["losersRound"] = maxLosersRound
	}

	return result
}

func getChampions(games []map[string]interface{}, playerMap map[string]map[string]interface{}) []map[string]interface{} {
	var champions []map[string]interface{}

	// Find games with winner placement = 1
	for _, game := range games {
		if wp, ok := game["winnerPlacement"].(float64); ok && wp == 1 {
			if state, ok := game["state"].(string); ok && state == "done" {
				// Find the winner
				if winnerSlotIdx, ok := game["winnerSlotIdx"].(int); ok {
					if slots, ok := game["slots"].([]map[string]interface{}); ok && len(slots) > winnerSlotIdx {
						if playerID, ok := slots[winnerSlotIdx]["playerID"].(string); ok {
							if player, found := playerMap[playerID]; found {
								champion := map[string]interface{}{
									"playerID":    playerID,
									"displayName": player["displayName"],
									"name":        player["name"],
									"gameID":      game["id"],
									"gameName":    game["name"],
								}
								champions = append(champions, champion)
							}
						}
					}
				}
			}
		}
	}

	return champions
}

func getBracketTypeFromRound(round int) string {
	if round < 0 {
		return "losers"
	}
	return "winners"
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
