package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// NHRL Statsbook API client configuration
const (
	NHRLBaseURL        = "https://stats.nhrl.io/statsbook"
	BrettZoneBaseURL   = "https://brettzone.nhrl.io/brettZone/backend"
	BrettZoneReviewURL = "https://brettzone.nhrl.io/brettZone/fightReviewMulti.php"
)

// HTTP client for NHRL API with timeout
var nhrlHttpClient = &http.Client{
	Timeout: 30 * time.Second,
}

// NHRL API response structures
type NHRLDumpsterCount struct {
	BotName string `json:"bot_name"`
	First   int    `json:"first"`
	Second  int    `json:"second"`
	Third   int    `json:"third"`
}

type NHRLEventWinner struct {
	EventDate       string  `json:"event_date"`
	FirstPlaceName  string  `json:"first_place_name"`
	SecondPlaceName string  `json:"second_place_name"`
	ThirdPlaceName  string  `json:"third_place_name"`
	FourthPlaceName *string `json:"fourth_place_name"`
}

type NHRLRanking struct {
	Ranking int `json:"ranking"`
}

type NHRLFight struct {
	Points          string  `json:"points"`
	Date            string  `json:"date"`
	MatchNum        int     `json:"match_num"`
	Round           string  `json:"round"`
	ResultBy        string  `json:"result_by"`
	FightLengthSecs *string `json:"fight_length_secs"`
	VideoLink       *string `json:"video_link"`
}

type NHRLHeadToHead struct {
	OpponentUniqueName string `json:"opponent_unique_name"`
	NumFights          int    `json:"num_fights"`
	Wins               int    `json:"wins"`
	Losses             int    `json:"losses"`
	KOs                int    `json:"kos"`
	KOd                int    `json:"kod"`
	LastMeeting        string `json:"last_meeting"`
}

type NHRLStatSummary struct {
	Bot              string `json:"bot"`
	Ranking          int    `json:"ranking"`
	RankChange       string `json:"rank_change"`
	Points           string `json:"points"`
	Events           int    `json:"events"`
	Fights           int    `json:"fights"`
	W                int    `json:"w"`
	L                int    `json:"l"`
	Pct              string `json:"pct"`
	KOs              int    `json:"kos"`
	KOd              int    `json:"kod"`
	LastAppearance   string `json:"last_appearance"`
	LastEventW       int    `json:"last_event_w"`
	LastEventL       int    `json:"last_event_l"`
	AvgFightTimeSecs string `json:"avg_fight_time_secs"`
}

type NHRLBotStatsBySeason struct {
	Bot              string `json:"bot"`
	Events           int    `json:"events"`
	Fights           int    `json:"fights"`
	W                int    `json:"w"`
	L                int    `json:"l"`
	Pct              string `json:"pct"`
	KOs              int    `json:"kos"`
	KOd              int    `json:"kod"`
	LastAppearance   string `json:"last_appearance"`
	AvgFightTimeSecs string `json:"avg_fight_time_secs"`
}

type NHRLFastestKO struct {
	BotName         string  `json:"bot_name"`
	OpponentName    string  `json:"opponent_name"`
	FightLengthSecs int     `json:"fight_length_secs"`
	Date            string  `json:"date"`
	Round           string  `json:"round"`
	VideoLink       *string `json:"video_link"`
}

type NHRLWinningStreak struct {
	BotName        string `json:"bot_name"`
	StreakLength   int    `json:"streak_length"`
	LastAppearance string `json:"last_appearance"`
}

type NHRLStreakStats struct {
	CurrentStreak     int    `json:"current_streak"`
	CurrentStreakType string `json:"current_streak_type"`
	LongestWinStreak  int    `json:"longest_win_streak"`
	LongestLoseStreak int    `json:"longest_lose_streak"`
}

// NHRLLiveFightStats represents live fight statistics between two bots
type NHRLLiveFightStats struct {
	DriverName          string  `json:"driver_name"`
	DriverPronunciation string  `json:"driver_pronunciation"`
	City                string  `json:"city"`
	StateProvince       string  `json:"state_province"`
	Country             string  `json:"country"`
	Pronouns            string  `json:"pronouns"`
	TeamName            *string `json:"team_name"`
	BotName             string  `json:"bot_name"`
	BotPronunciation    *string `json:"bot_pronunciation"`
	Ranking             int     `json:"ranking"`
	Events              int     `json:"events"`
	Fights              int     `json:"fights"`
	W                   int     `json:"w"`
	L                   int     `json:"l"`
	WKO                 int     `json:"w_ko"`
	LKO                 int     `json:"l_ko"`
	WJD                 int     `json:"w_jd"`
	LJD                 int     `json:"l_jd"`
	WinPct              string  `json:"win_pct"`
	BuilderBackground   string  `json:"builder_background"`
	InterestingFact     string  `json:"interesting_fact"`
	InterestingFact2    *string `json:"interesting_fact_2"`
	BotType             string  `json:"bot_type"`
	HthW                int     `json:"hth_w"`        // Head-to-head wins against opponent
	HthWKO              int     `json:"hth_w_ko"`     // Head-to-head KO wins against opponent
	HthWJD              int     `json:"hth_w_jd"`     // Head-to-head JD wins against opponent
	LastMeeting         *string `json:"last_meeting"` // Last meeting date with opponent
}

// BrettZone Tournament Match data structures
type BrettZoneMatch struct {
	TournamentID   string `json:"tournamentID"`
	ID             string `json:"id"`
	Name           string `json:"name"`
	Round          string `json:"round"`
	Cage           string `json:"cage"`
	Player1        string `json:"player1"`
	Player1Clean   string `json:"player1clean"`
	Player2        string `json:"player2"`
	Player2Clean   string `json:"player2clean"`
	Player1Wins    string `json:"player1wins"`
	Player2Wins    string `json:"player2wins"`
	Cams           string `json:"cams"`
	WinAnnotation  string `json:"winAnnotation"`
	CalledSince    string `json:"calledSince"`
	AvailableSince string `json:"availableSince"`
	EndTime        string `json:"endTime"`
	StartTime      string `json:"startTime"`
	StopTime       string `json:"stopTime"`
	MatchLength    string `json:"matchLength"`
	WeightClass    string `json:"weightClass"`
	TournamentName string `json:"tournamentName"`
	Privacy        string `json:"privacy"`
	IsTest         string `json:"isTest"`
	IsFreestyle    string `json:"isFreestyle"`
}

// Helper function to normalize bot names for API calls (replace spaces with underscores)
func normalizeBotName(botName string) string {
	return strings.ReplaceAll(botName, " ", "_")
}

// Generic function to make NHRL API requests
func makeNHRLAPIRequest(endpoint string, params map[string]string) ([]byte, error) {
	// Build query parameters
	queryParams := url.Values{}
	for key, value := range params {
		queryParams.Add(key, value)
	}

	// Build full URL
	fullURL := fmt.Sprintf("%s/%s", NHRLBaseURL, endpoint)
	if len(queryParams) > 0 {
		fullURL += "?" + queryParams.Encode()
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "NHRL-MCP-Server/1.0.0")

	resp, err := nhrlHttpClient.Do(req)
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
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

// Get dumpster count (podium finishes) for a weight class
func getNHRLDumpsterCount(categoryID string) ([]NHRLDumpsterCount, error) {
	params := map[string]string{
		"category_id": categoryID,
	}

	data, err := makeNHRLAPIRequest("get_dumpster_count.php", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get dumpster count: %w", err)
	}

	var result []NHRLDumpsterCount
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse dumpster count response: %w", err)
	}

	return result, nil
}

// Get event participants for a specific bot
func getNHRLEventParticipants(botName string) ([]map[string]interface{}, error) {
	params := map[string]string{
		"bot_name": normalizeBotName(botName),
	}

	data, err := makeNHRLAPIRequest("get_event_participants.php", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get event participants: %w", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse event participants response: %w", err)
	}

	return result, nil
}

// Get event winners for a weight class
func getNHRLEventWinners(weightClass string) ([]NHRLEventWinner, error) {
	params := map[string]string{
		"weight_class": weightClass,
	}

	data, err := makeNHRLAPIRequest("get_event_winners.php", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get event winners: %w", err)
	}

	var result []NHRLEventWinner
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse event winners response: %w", err)
	}

	return result, nil
}

// Get fastest KOs for a weight class
func getNHRLFastestKOs(classID string) ([]NHRLFastestKO, error) {
	params := map[string]string{
		"class_id": classID,
	}

	data, err := makeNHRLAPIRequest("get_fastest_kos.php", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get fastest KOs: %w", err)
	}

	var result []NHRLFastestKO
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse fastest KOs response: %w", err)
	}

	return result, nil
}

// Get fights for a specific bot
func getNHRLFights(botName string) ([]NHRLFight, error) {
	params := map[string]string{
		"bot_name": normalizeBotName(botName),
	}

	data, err := makeNHRLAPIRequest("get_fights.php", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get fights: %w", err)
	}

	var result []NHRLFight
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse fights response: %w", err)
	}

	return result, nil
}

// Get head-to-head record for a specific bot
func getNHRLHeadToHead(botName string) ([]NHRLHeadToHead, error) {
	params := map[string]string{
		"bot_name": normalizeBotName(botName),
	}

	data, err := makeNHRLAPIRequest("get_head_to_head.php", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get head-to-head: %w", err)
	}

	var result []NHRLHeadToHead
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse head-to-head response: %w", err)
	}

	return result, nil
}

// Get longest winning streaks for a weight class
func getNHRLLongestWinningStreak(categoryID string) ([]NHRLWinningStreak, error) {
	params := map[string]string{
		"category_id": categoryID,
	}

	data, err := makeNHRLAPIRequest("get_longest_winning_streak.php", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get longest winning streak: %w", err)
	}

	var result []NHRLWinningStreak
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse longest winning streak response: %w", err)
	}

	return result, nil
}

// Get random fight
func getNHRLRandomFight() (map[string]interface{}, error) {
	data, err := makeNHRLAPIRequest("get_random_fight.php", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get random fight: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse random fight response: %w", err)
	}

	return result, nil
}

// Get bot rank
func getNHRLBotRank(botName string) (*NHRLRanking, error) {
	params := map[string]string{
		"bot_name": normalizeBotName(botName),
	}

	data, err := makeNHRLAPIRequest("get_rank.php", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get bot rank: %w", err)
	}

	// Handle null response
	if string(data) == "null" {
		return nil, nil
	}

	var result NHRLRanking
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse bot rank response: %w", err)
	}

	return &result, nil
}

// Get stat summary for a weight class and season
func getNHRLStatSummary(categoryID, seasonID string) ([]NHRLStatSummary, error) {
	params := map[string]string{
		"category_id": categoryID,
		"season_id":   seasonID,
	}

	data, err := makeNHRLAPIRequest("get_stat_summary.php", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get stat summary: %w", err)
	}

	var result []NHRLStatSummary
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse stat summary response: %w", err)
	}

	return result, nil
}

// Get stats by season for a specific bot
func getNHRLStatsBySeason(botName, seasonID string) (*NHRLBotStatsBySeason, error) {
	params := map[string]string{
		"bot_name":  normalizeBotName(botName),
		"season_id": seasonID,
	}

	data, err := makeNHRLAPIRequest("get_stats_by_season.php", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats by season: %w", err)
	}

	// Handle null response
	if string(data) == "null" {
		return nil, nil
	}

	var result NHRLBotStatsBySeason
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse stats by season response: %w", err)
	}

	return &result, nil
}

// Get streak stats for a specific bot
func getNHRLStreakStats(botName string) (*NHRLStreakStats, error) {
	params := map[string]string{
		"bot_name": normalizeBotName(botName),
	}

	data, err := makeNHRLAPIRequest("get_streak_stats.php", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get streak stats: %w", err)
	}

	// Handle null response
	if string(data) == "null" {
		return nil, nil
	}

	var result NHRLStreakStats
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse streak stats response: %w", err)
	}

	return &result, nil
}

// Get live fight stats between two bots for a specific tournament
func getNHRLLiveFightStats(bot1, bot2, tournamentID string) ([]NHRLLiveFightStats, error) {
	// Build form data
	formData := url.Values{}
	formData.Set("bot1", bot1)
	formData.Set("bot2", bot2)

	// Build full URL with tournament ID
	fullURL := fmt.Sprintf("https://stats.nhrl.io/live_stats/query/get_fight_stats.php?tournament_id=%s", tournamentID)

	req, err := http.NewRequest("POST", fullURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "NHRL-MCP-Server/1.0.0")

	resp, err := nhrlHttpClient.Do(req)
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
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(responseBody))
	}

	// Handle empty response
	if len(responseBody) == 0 || string(responseBody) == "" {
		return []NHRLLiveFightStats{}, nil
	}

	var result []NHRLLiveFightStats
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse live fight stats response: %w", err)
	}

	return result, nil
}

// BrettZone API Functions

// Generic function to make BrettZone API requests
func makeBrettZoneAPIRequest(endpoint string, params map[string]string) ([]byte, error) {
	// Build query parameters
	queryParams := url.Values{}
	for key, value := range params {
		queryParams.Add(key, value)
	}

	// Build full URL
	fullURL := fmt.Sprintf("%s/%s", BrettZoneBaseURL, endpoint)
	if len(queryParams) > 0 {
		fullURL += "?" + queryParams.Encode()
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "NHRL-MCP-Server/1.0.0")

	resp, err := nhrlHttpClient.Do(req)
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
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

// Get latest matches for a tournament
func getBrettZoneLatestMatches(tournamentID string) ([]BrettZoneMatch, error) {
	params := map[string]string{
		"tournamentID": tournamentID,
	}

	data, err := makeBrettZoneAPIRequest("getLatestMatches.php", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest matches: %w", err)
	}

	var result []BrettZoneMatch
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse latest matches response: %w", err)
	}

	return result, nil
}

// Generate BrettZone fight review URL
func generateBrettZoneReviewURL(gameID, tournamentID string, cageNum int, timeSeconds float64) string {
	cage := fmt.Sprintf("cam-Cage-%d-Overhead-High", cageNum)
	if cageNum == 0 {
		cage = "cam-Cage-1-Overhead-High"
	}

	return fmt.Sprintf("%s?gameID=%s&tournamentID=%s#cams=%s&t=%.2f",
		BrettZoneReviewURL, gameID, tournamentID, cage, timeSeconds)
}

// Helper function to get cage number from cage string
func extractCageNumber(cageStr string) int {
	if strings.Contains(cageStr, "Cage 1") {
		return 1
	} else if strings.Contains(cageStr, "Cage 2") {
		return 2
	} else if strings.Contains(cageStr, "Cage 3") {
		return 3
	} else if strings.Contains(cageStr, "Cage 4") {
		return 4
	}
	return 1 // Default to Cage 1
}

// Helper function to get weight class from tournament name
func extractWeightClass(tournamentName string) string {
	if strings.Contains(strings.ToLower(tournamentName), "3lb") {
		return "3lb"
	} else if strings.Contains(strings.ToLower(tournamentName), "12lb") {
		return "12lb"
	} else if strings.Contains(strings.ToLower(tournamentName), "30lb") {
		return "30lb"
	}
	return "unknown"
}

// Helper function to get weight class category ID from weight class name
func getWeightClassCategoryID(weightClass string) string {
	switch strings.ToLower(weightClass) {
	case "3lb", "3 lb", "beetleweight":
		return "1"
	case "12lb", "12 lb", "antweight":
		return "2"
	case "30lb", "30 lb", "hobbyweight":
		return "4"
	default:
		return "1" // Default to 3lb
	}
}

// Helper function to get season ID from season name/year
func getSeasonID(season string) string {
	switch strings.ToLower(season) {
	case "current", "0":
		return "0"
	case "all-time", "all time", "alltime", "1":
		return "1"
	case "2018-2019", "2018", "2019", "2":
		return "2"
	case "2020", "3":
		return "3"
	case "2021", "4":
		return "4"
	case "2022", "5":
		return "5"
	case "2023", "6":
		return "6"
	default:
		return "1" // Default to all-time
	}
}

// Tournament Round Information
type RoundInfo struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	WinResult   string `json:"win_result"`
	LoseResult  string `json:"lose_result"`
}

// Helper function to get qualification round name from code
func getQualificationRoundName(roundCode string) string {
	switch strings.ToUpper(roundCode) {
	case "Q1":
		return "Opening"
	case "Q2W":
		return "The Cusp"
	case "Q2L":
		return "Redemption"
	case "Q3":
		return "Bubble"
	default:
		// For bracket rounds, return as-is
		return roundCode
	}
}

// Helper function to get detailed round information
func getRoundInfo(roundCode string) RoundInfo {
	switch strings.ToUpper(roundCode) {
	case "Q1":
		return RoundInfo{
			Code:        "Q1",
			Name:        "Opening",
			Description: "First qualifying match - all competitors start here",
			WinResult:   "Advances to The Cusp (Q2W)",
			LoseResult:  "Drops to Redemption (Q2L)",
		}
	case "Q2W":
		return RoundInfo{
			Code:        "Q2W",
			Name:        "The Cusp",
			Description: "Second match for Opening winners - one win away from qualifying",
			WinResult:   "Qualifies for main bracket",
			LoseResult:  "Drops to Bubble (Q3) for last chance",
		}
	case "Q2L":
		return RoundInfo{
			Code:        "Q2L",
			Name:        "Redemption",
			Description: "Second chance for Opening losers",
			WinResult:   "Advances to Bubble (Q3) for last chance",
			LoseResult:  "Eliminated from tournament",
		}
	case "Q3":
		return RoundInfo{
			Code:        "Q3",
			Name:        "Bubble",
			Description: "Final qualifying round - last chance to make the bracket",
			WinResult:   "Qualifies for main bracket",
			LoseResult:  "Eliminated from tournament",
		}
	default:
		// For bracket rounds
		return RoundInfo{
			Code:        roundCode,
			Name:        roundCode,
			Description: "Main bracket round",
			WinResult:   "Advances to next round",
			LoseResult:  "Eliminated from tournament",
		}
	}
}

// Enriched BrettZone Match with round information
type EnrichedBrettZoneMatch struct {
	BrettZoneMatch
	RoundName        string `json:"round_name"`
	RoundDescription string `json:"round_description"`
	WinImplication   string `json:"win_implication"`
	LoseImplication  string `json:"lose_implication"`
}

// Helper function to enrich a BrettZone match with round information
func enrichBrettZoneMatch(match BrettZoneMatch) EnrichedBrettZoneMatch {
	roundInfo := getRoundInfo(match.Round)
	return EnrichedBrettZoneMatch{
		BrettZoneMatch:   match,
		RoundName:        roundInfo.Name,
		RoundDescription: roundInfo.Description,
		WinImplication:   roundInfo.WinResult,
		LoseImplication:  roundInfo.LoseResult,
	}
}

// Helper function to enrich multiple BrettZone matches
func enrichBrettZoneMatches(matches []BrettZoneMatch) []EnrichedBrettZoneMatch {
	enrichedMatches := make([]EnrichedBrettZoneMatch, len(matches))
	for i, match := range matches {
		enrichedMatches[i] = enrichBrettZoneMatch(match)
	}
	return enrichedMatches
}

// Helper function to explain the qualification path
func getQualificationPathExplanation() string {
	return `NHRL Tournament Qualification System:

1. OPENING (Q1): All competitors start here
   - Win → Advance to THE CUSP (Q2W)
   - Lose → Drop to REDEMPTION (Q2L)

2. THE CUSP (Q2W): For Opening winners
   - Win → QUALIFY FOR BRACKET
   - Lose → Drop to BUBBLE (Q3) for last chance

3. REDEMPTION (Q2L): Second chance for Opening losers
   - Win → Advance to BUBBLE (Q3)
   - Lose → ELIMINATED

4. BUBBLE (Q3): Final qualifying round
   - Win → QUALIFY FOR BRACKET
   - Lose → ELIMINATED

Once qualified, competitors enter the main single-elimination bracket.`
}
