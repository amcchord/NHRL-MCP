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
	NHRLBaseURL = "https://stats.nhrl.io/statsbook"
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
