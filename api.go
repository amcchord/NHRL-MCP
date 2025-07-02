package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Global configuration
var (
	APIBaseURL = "https://truefinals.com/api" // Default base URL
	apiKey     string                         // API key for authentication
	apiUserID  string                         // API user ID for authentication
)

// HTTP client with timeout
var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

// TrueFinals API structures based on OpenAPI spec
type Tournament struct {
	ID                   string           `json:"id"`
	CreatorID            string           `json:"creatorID"`
	CreatorProfileID     string           `json:"creatorProfileID"`
	Title                string           `json:"title"`
	GameTitleInfo        *GameTitleInfo   `json:"gameTitleInfo"`
	EventLocation        string           `json:"eventLocation"`
	Privacy              string           `json:"privacy"`
	DisplayCheckInStatus bool             `json:"displayCheckInStatus"`
	LogoURL              *string          `json:"logoUrl"`
	ThumbnailURL         *string          `json:"thumbnailUrl"`
	CreateTime           int64            `json:"createTime"`
	ScheduledStartTime   *int64           `json:"scheduledStartTime"`
	StartTime            *int64           `json:"startTime"`
	EndTime              *int64           `json:"endTime"`
	UpdateTime           int64            `json:"updateTime"`
	Description          *string          `json:"description"`
	OverlayParams        *OverlayParams   `json:"overlayParams"`
	Players              []Player         `json:"players"`
	Locations            []Location       `json:"locations"`
	Games                []Game           `json:"games"`
	Format               TournamentFormat `json:"format"`
}

type GameTitleInfo struct {
	Name         string  `json:"name"`
	ID           *string `json:"id"`
	ThumbnailURL *string `json:"thumbnailUrl"`
}

type OverlayParams struct {
	GameID      *string     `json:"gameID"`
	OverlayData OverlayData `json:"overlayData"`
	Theme       Theme       `json:"theme"`
}

type OverlayData struct {
	TournamentName string          `json:"tournamentName"`
	LogoURL        *string         `json:"logoUrl"`
	BracketName    string          `json:"bracketName"`
	RoundName      string          `json:"roundName"`
	ShortRoundName string          `json:"shortRoundName"`
	ScoreToWin     int             `json:"scoreToWin"`
	Players        []OverlayPlayer `json:"players"`
	Hidden         bool            `json:"hidden"`
	Swapped        bool            `json:"swapped"`
}

type OverlayPlayer struct {
	Name          string  `json:"name"`
	ScoreText     string  `json:"scoreText"`
	Tag           *string `json:"tag"`
	PhotoURL      *string `json:"photoUrl"`
	Pronouns      *string `json:"pronouns"`
	TwitterHandle *string `json:"twitterHandle"`
	Wins          int     `json:"wins"`
	Losses        int     `json:"losses"`
	Ties          int     `json:"ties"`
	Seed          *int    `json:"seed"`
}

type Theme struct {
	Shape               string     `json:"shape"`
	BgParams            BgParams   `json:"bgParams"`
	PrimaryTextParams   TextParams `json:"primaryTextParams"`
	SecondaryTextParams TextParams `json:"secondaryTextParams"`
	BodyTextParams      TextParams `json:"bodyTextParams"`
	ScoreTextParams     TextParams `json:"scoreTextParams"`
}

type BgParams struct {
	GradientDir      string  `json:"gradientDir"`
	PrimaryBgColor1  string  `json:"primaryBgColor1"`
	PrimaryBgColor2  string  `json:"primaryBgColor2"`
	SecondaryBgColor string  `json:"secondaryBgColor"`
	BackdropColor    string  `json:"backdropColor"`
	ScoreColor       string  `json:"scoreColor"`
	AccentColor      string  `json:"accentColor"`
	AccentWidthPx    float64 `json:"accentWidthPx"`
}

type TextParams struct {
	FontFamily string  `json:"fontFamily"`
	FontSizePx float64 `json:"fontSizePx"`
	FontColor  string  `json:"fontColor"`
	Transform  string  `json:"transform"`
	Bold       bool    `json:"bold"`
	Italic     bool    `json:"italic"`
}

type Player struct {
	ID                string       `json:"id"`
	Name              string       `json:"name"`
	PhotoURL          *string      `json:"photoUrl"`
	Seed              *int         `json:"seed"`
	Wins              int          `json:"wins"`
	Losses            int          `json:"losses"`
	Ties              int          `json:"ties"`
	IsBye             bool         `json:"isBye"`
	IsDisqualified    bool         `json:"isDisqualified"`
	LastPlayTime      *int64       `json:"lastPlayTime"`
	LastBracketGameID *string      `json:"lastBracketGameID"`
	Placement         *int         `json:"placement"`
	ProfileInfo       *ProfileInfo `json:"profileInfo"`
}

type ProfileInfo struct {
	ID              string  `json:"id"`
	Tag             string  `json:"tag"`
	Name            string  `json:"name"`
	PhotoURL        *string `json:"photoUrl"`
	Pronouns        string  `json:"pronouns"`
	TwitchHandle    *string `json:"twitchHandle"`
	TwitterHandle   *string `json:"twitterHandle"`
	StartggPlayerID *int    `json:"startggPlayerID"`
}

type Location struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	ActiveGameID        *string  `json:"activeGameID"`
	LastCompletedGameID *string  `json:"lastCompletedGameID"`
	Queue               []string `json:"queue"`
	UnavailableQueue    []string `json:"unavailableQueue"`
	BlockActive         bool     `json:"blockActive"`
}

type Game struct {
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	BracketID        string     `json:"bracketID"`
	Round            int        `json:"round"`
	ScoreToWin       int        `json:"scoreToWin"`
	Slots            []GameSlot `json:"slots"`
	State            string     `json:"state"`
	ActiveSince      *int64     `json:"activeSince"`
	AvailableSince   *int64     `json:"availableSince"`
	CalledSince      *int64     `json:"calledSince"`
	HeldSince        *int64     `json:"heldSince"`
	EndTime          *int64     `json:"endTime"`
	ScheduledTime    *int64     `json:"scheduledTime"`
	NextGameSlotIDs  []string   `json:"nextGameSlotIDs"`
	LocationID       *string    `json:"locationID"`
	ResultAnnotation *string    `json:"resultAnnotation"`
	WinnerPlacement  *int       `json:"winnerPlacement"`
	LoserPlacement   *int       `json:"loserPlacement"`
}

type GameSlot struct {
	GameID      string  `json:"gameID"`
	SlotIdx     int     `json:"slotIdx"`
	PlayerID    *string `json:"playerID"`
	CheckInTime *int64  `json:"checkInTime"`
	WaitingTime *int64  `json:"waitingTime"`
	PrevGameID  *string `json:"prevGameID"`
	Score       float64 `json:"score"` // Can be integer or -1
	SlotState   string  `json:"slotState"`
}

type TournamentFormat struct {
	Type    string      `json:"type"`
	Scoring interface{} `json:"scoring"`
	// Additional fields depend on the tournament type
}

// API Response structures
type APIErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Issues  []struct {
		Message string `json:"message"`
	} `json:"issues,omitempty"`
}

type TournamentListResponse []TournamentListItem

type TournamentListItem struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Privacy    string `json:"privacy"`
	CreateTime int64  `json:"createTime"`
	EndTime    *int64 `json:"endTime"`
}

// makeAPIRequest performs HTTP requests to the TrueFinals API
func makeAPIRequest(method, endpoint string, body interface{}) ([]byte, error) {
	var reqBody io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	// Build full URL
	fullURL := APIBaseURL + endpoint

	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers for TrueFinals API
	req.Header.Set("x-api-user-id", apiUserID)
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP error status codes
	if resp.StatusCode >= 400 {
		var apiError APIErrorResponse
		if err := json.Unmarshal(responseBody, &apiError); err == nil {
			return nil, fmt.Errorf("API error (%d): %s - %s", resp.StatusCode, apiError.Code, apiError.Message)
		}
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

// buildQueryParams builds query parameters for GET requests
func buildQueryParams(params map[string]interface{}) string {
	if len(params) == 0 {
		return ""
	}

	values := url.Values{}
	for key, value := range params {
		if value != nil {
			values.Add(key, fmt.Sprintf("%v", value))
		}
	}

	return "?" + values.Encode()
}

// Helper function to build endpoint with query parameters
func buildEndpoint(base string, params map[string]interface{}) string {
	return base + buildQueryParams(params)
}

// Helper functions to enrich data with human-readable information

// EnrichGameWithPlayerAndLocationInfo adds player names and location name to game data
func enrichGameWithPlayerAndLocationInfo(game map[string]interface{}, tournamentID string) map[string]interface{} {
	// Get players and locations data
	playersData, _ := getPlayersData(tournamentID)
	locationsData, _ := getLocationsData(tournamentID)

	// Create player ID to enriched player map
	playerMap := make(map[string]map[string]interface{})
	playerNameMap := make(map[string]string)
	if players, ok := playersData.([]interface{}); ok {
		for _, p := range players {
			if player, ok := p.(map[string]interface{}); ok {
				if id, ok := player["id"].(string); ok {
					// Enrich player with NHRL stats
					enrichedPlayer := enrichPlayerWithNHRLStats(player)
					playerMap[id] = enrichedPlayer
					if name, ok := player["name"].(string); ok {
						playerNameMap[id] = name
					}
				}
			}
		}
	}

	// Create location ID to name map
	locationMap := make(map[string]string)
	if locations, ok := locationsData.([]interface{}); ok {
		for _, l := range locations {
			if location, ok := l.(map[string]interface{}); ok {
				if id, ok := location["id"].(string); ok {
					if name, ok := location["name"].(string); ok {
						locationMap[id] = name
					}
				}
			}
		}
	}

	// Enrich game data
	if slots, ok := game["slots"].([]interface{}); ok {
		enrichedSlots := make([]interface{}, len(slots))
		for i, s := range slots {
			if slot, ok := s.(map[string]interface{}); ok {
				enrichedSlot := make(map[string]interface{})
				for k, v := range slot {
					enrichedSlot[k] = v
				}
				// Add player name and NHRL stats if player ID exists
				if playerID, ok := slot["playerID"].(string); ok && playerID != "" {
					enrichedSlot["playerName"] = playerNameMap[playerID]
					// Add NHRL stats if available
					if enrichedPlayer, ok := playerMap[playerID]; ok {
						if nhrlRank, ok := enrichedPlayer["nhrl_rank"]; ok {
							enrichedSlot["nhrl_rank"] = nhrlRank
						}
						if nhrlStreak, ok := enrichedPlayer["nhrl_current_streak"]; ok {
							enrichedSlot["nhrl_current_streak"] = nhrlStreak
						}
					}
				}
				enrichedSlots[i] = enrichedSlot
			}
		}
		game["slots"] = enrichedSlots
	}

	// Add location name if location ID exists
	if locationID, ok := game["locationID"].(string); ok && locationID != "" {
		game["locationName"] = locationMap[locationID]
	}

	// Add round qualification information if game name matches qualification rounds
	if gameName, ok := game["name"].(string); ok && gameName != "" {
		roundInfo := getRoundInfo(gameName)
		// Only add qualification info if it's actually a qualification round
		if roundInfo.Code == gameName && (gameName == "Q1" || gameName == "Q2W" || gameName == "Q2L" || gameName == "Q3") {
			game["roundName"] = roundInfo.Name
			game["roundDescription"] = roundInfo.Description
			game["winImplication"] = roundInfo.WinResult
			game["loseImplication"] = roundInfo.LoseResult
			game["isQualificationRound"] = true
		}
	}

	return game
}

// Helper to get players data
func getPlayersData(tournamentID string) (interface{}, error) {
	endpoint := fmt.Sprintf("/v1/tournaments/%s/players", tournamentID)
	data, err := makeAPIRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var players []interface{}
	if err := json.Unmarshal(data, &players); err != nil {
		return nil, err
	}

	return players, nil
}

// Helper to get locations data
func getLocationsData(tournamentID string) (interface{}, error) {
	endpoint := fmt.Sprintf("/v1/tournaments/%s/locations", tournamentID)
	data, err := makeAPIRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var locations []interface{}
	if err := json.Unmarshal(data, &locations); err != nil {
		return nil, err
	}

	return locations, nil
}

// EnrichLocationWithGameInfo adds active game information with player names
func enrichLocationWithGameInfo(location map[string]interface{}, tournamentID string) map[string]interface{} {
	// If there's an active game ID, fetch the game info
	if activeGameID, ok := location["activeGameID"].(string); ok && activeGameID != "" {
		endpoint := fmt.Sprintf("/v1/tournaments/%s/games/%s", tournamentID, activeGameID)
		gameData, err := makeAPIRequest("GET", endpoint, nil)
		if err == nil {
			var game map[string]interface{}
			if err := json.Unmarshal(gameData, &game); err == nil {
				// Enrich the game with player info
				enrichedGame := enrichGameWithPlayerAndLocationInfo(game, tournamentID)

				// Extract player names for easy display
				var playerNames []string
				if slots, ok := enrichedGame["slots"].([]interface{}); ok {
					for _, s := range slots {
						if slot, ok := s.(map[string]interface{}); ok {
							if name, ok := slot["playerName"].(string); ok && name != "" {
								playerNames = append(playerNames, name)
							}
						}
					}
				}

				activeGameInfo := map[string]interface{}{
					"gameID":      activeGameID,
					"gameName":    enrichedGame["name"],
					"playerNames": playerNames,
					"state":       enrichedGame["state"],
				}

				// Add qualification round info if present
				if roundName, ok := enrichedGame["roundName"].(string); ok {
					activeGameInfo["roundName"] = roundName
				}
				if isQual, ok := enrichedGame["isQualificationRound"].(bool); ok && isQual {
					activeGameInfo["isQualificationRound"] = true
				}

				location["activeGameInfo"] = activeGameInfo
			}
		}
	}

	return location
}

// EnrichPlayerDataWithProfileDetails adds more details from profile info and NHRL stats
func enrichPlayerDataWithProfileDetails(player map[string]interface{}) map[string]interface{} {
	// Extract useful info from profileInfo if it exists
	if profileInfo, ok := player["profileInfo"].(map[string]interface{}); ok {
		// Add tag if available
		if tag, ok := profileInfo["tag"].(string); ok && tag != "" {
			player["tag"] = tag
		}
		// Add pronouns if available
		if pronouns, ok := profileInfo["pronouns"].(string); ok && pronouns != "" {
			player["pronouns"] = pronouns
		}
		// Add social handles if available
		if twitchHandle, ok := profileInfo["twitchHandle"].(string); ok && twitchHandle != "" {
			player["twitchHandle"] = twitchHandle
		}
		if twitterHandle, ok := profileInfo["twitterHandle"].(string); ok && twitterHandle != "" {
			player["twitterHandle"] = twitterHandle
		}
	}

	// Create a display name that includes tag if available
	if tag, ok := player["tag"].(string); ok && tag != "" {
		player["displayName"] = tag
	} else {
		player["displayName"] = player["name"]
	}

	// Enrich with NHRL stats
	player = enrichPlayerWithNHRLStats(player)

	return player
}

// EnrichTournamentData adds human-readable information to tournament data
func enrichTournamentData(tournament map[string]interface{}) map[string]interface{} {
	// First apply existing TrueFinals enrichment
	// Enrich players
	if players, ok := tournament["players"].([]interface{}); ok {
		enrichedPlayers := make([]interface{}, len(players))
		for i, p := range players {
			if player, ok := p.(map[string]interface{}); ok {
				enrichedPlayers[i] = enrichPlayerDataWithProfileDetails(player)
			} else {
				enrichedPlayers[i] = p
			}
		}
		tournament["players"] = enrichedPlayers
		tournament["playersCount"] = len(enrichedPlayers)
	}

	// Enrich games with player names and location names
	if games, ok := tournament["games"].([]interface{}); ok {
		// Create player ID to name map
		playerMap := make(map[string]string)
		if players, ok := tournament["players"].([]interface{}); ok {
			for _, p := range players {
				if player, ok := p.(map[string]interface{}); ok {
					if id, ok := player["id"].(string); ok {
						if displayName, ok := player["displayName"].(string); ok {
							playerMap[id] = displayName
						} else if name, ok := player["name"].(string); ok {
							playerMap[id] = name
						}
					}
				}
			}
		}

		// Create location ID to name map
		locationMap := make(map[string]string)
		if locations, ok := tournament["locations"].([]interface{}); ok {
			for _, l := range locations {
				if location, ok := l.(map[string]interface{}); ok {
					if id, ok := location["id"].(string); ok {
						if name, ok := location["name"].(string); ok {
							locationMap[id] = name
						}
					}
				}
			}
		}

		// Enrich each game
		enrichedGames := make([]interface{}, len(games))
		for i, g := range games {
			if game, ok := g.(map[string]interface{}); ok {
				// Add location name
				if locationID, ok := game["locationID"].(string); ok && locationID != "" {
					game["locationName"] = locationMap[locationID]
				}

				// Enrich slots with player names
				if slots, ok := game["slots"].([]interface{}); ok {
					enrichedSlots := make([]interface{}, len(slots))
					for j, s := range slots {
						if slot, ok := s.(map[string]interface{}); ok {
							enrichedSlot := make(map[string]interface{})
							for k, v := range slot {
								enrichedSlot[k] = v
							}
							// Add player name if player ID exists
							if playerID, ok := slot["playerID"].(string); ok && playerID != "" {
								enrichedSlot["playerName"] = playerMap[playerID]
							}
							enrichedSlots[j] = enrichedSlot
						}
					}
					game["slots"] = enrichedSlots
				}

				// Add round qualification information if game name matches qualification rounds
				if gameName, ok := game["name"].(string); ok && gameName != "" {
					roundInfo := getRoundInfo(gameName)
					// Only add qualification info if it's actually a qualification round
					if roundInfo.Code == gameName && (gameName == "Q1" || gameName == "Q2W" || gameName == "Q2L" || gameName == "Q3") {
						game["roundName"] = roundInfo.Name
						game["roundDescription"] = roundInfo.Description
						game["winImplication"] = roundInfo.WinResult
						game["loseImplication"] = roundInfo.LoseResult
						game["isQualificationRound"] = true
					}
				}

				enrichedGames[i] = game
			} else {
				enrichedGames[i] = g
			}
		}
		tournament["games"] = enrichedGames
		tournament["gamesCount"] = len(enrichedGames)
	}

	// Enrich locations with active game info
	if locations, ok := tournament["locations"].([]interface{}); ok {
		tournament["locationsCount"] = len(locations)
		// Note: We could enrich locations with active game player names but that would be redundant
		// since the games are already enriched above
	}

	tournament["enrichmentNote"] = "Player names, display names, and location names are included throughout for better readability"

	// Then apply NHRL enrichment
	enrichedWithNHRL := enrichTournamentWithNHRLContext(tournament)

	return enrichedWithNHRL
}
