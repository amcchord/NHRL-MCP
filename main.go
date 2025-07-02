package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// MCP Protocol structures
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type ToolInfo struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

type ToolResult struct {
	Content []ToolContent `json:"content"`
	IsError bool          `json:"isError"`
}

type ToolContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Configuration
const (
	ServerName = "nhrl-mcp-server"
	Version    = "1.0.0"
)

// Tools filtering modes
const (
	ToolsReporting = "reporting"
	ToolsFullSafe  = "full-safe"
	ToolsFull      = "full"
)

var toolsMode string = ToolsFullSafe // Default to full-safe access
var disabledTools []string           // List of disabled tool names

// Helper functions for tools filtering
func isToolAllowed(toolName string) bool {
	// First check if tool is explicitly disabled
	if isToolDisabled(toolName) {
		return false
	}

	// Then check tools mode permissions
	switch toolsMode {
	case ToolsReporting:
		// Only allow read-only tools
		return isReadOnlyTool(toolName)
	case ToolsFullSafe:
		// Allow everything except dangerous operations
		return !isDangerousTool(toolName)
	case ToolsFull:
		// Allow all tools
		return true
	default:
		return false
	}
}

func isToolDisabled(toolName string) bool {
	for _, disabled := range disabledTools {
		if disabled == toolName {
			return true
		}
	}
	return false
}

func isOperationAllowed(toolName, operation string) bool {
	switch toolsMode {
	case ToolsReporting:
		// Only allow read operations
		return isReadOperation(operation)
	case ToolsFullSafe:
		// Allow everything except dangerous operations
		return !isDangerousOperation(toolName, operation)
	case ToolsFull:
		// Allow all operations
		return true
	default:
		return false
	}
}

func isReadOnlyTool(toolName string) bool {
	readOnlyTools := []string{
		"truefinals_tournaments", "truefinals_games", "truefinals_locations", "truefinals_players", "truefinals_bracket",
		"nhrl_stats",
	}
	for _, tool := range readOnlyTools {
		if tool == toolName {
			return true
		}
	}
	return false
}

func isDangerousTool(toolName string) bool {
	// No tools are completely dangerous in full-safe mode
	// Danger is at the operation level
	return false
}

func isReadOperation(operation string) bool {
	readOps := []string{
		"get", "list", "details", "format", "overlay_params", "description", "private", "webhooks",
		"get_round", "get_standings",
	}
	for _, op := range readOps {
		if op == operation {
			return true
		}
	}
	return false
}

func isDangerousOperation(toolName, operation string) bool {
	// Define dangerous operations that are blocked in full-safe mode
	if operation == "delete" {
		return true
	}
	if operation == "disqualify" {
		return true
	}
	if operation == "reset" {
		return true
	}
	return false
}

func main() {
	// Parse command line flags
	var cliAPIKey = flag.String("api-key", "", "API key for TrueFinals service (overrides TRUEFINALS_API_KEY environment variable)")
	var cliAPIUserID = flag.String("api-user-id", "", "API User ID for TrueFinals service (overrides TRUEFINALS_API_USER_ID environment variable)")
	var cliBaseURL = flag.String("base-url", "", "Base URL for TrueFinals API (overrides TRUEFINALS_BASE_URL environment variable)")
	var cliTools = flag.String("tools", "", "Tools mode: reporting, full-safe, full (overrides TRUEFINALS_TOOLS environment variable)")
	var cliDisabledTools = flag.String("disabled-tools", "", "Comma-separated list of tool names to disable (overrides TRUEFINALS_DISABLED_TOOLS environment variable)")
	var showVersion = flag.Bool("version", false, "Show version information and exit")
	var exitAfterFirst = flag.Bool("exit-after-first", false, "Exit after processing the first request instead of running continuously")
	flag.Parse()

	// Handle version flag
	if *showVersion {
		fmt.Printf("%s version %s\n", ServerName, Version)
		os.Exit(0)
	}

	// Get tools mode from CLI flag or environment variable
	// CLI flag takes precedence over environment variable
	if *cliTools != "" {
		toolsMode = *cliTools
	} else if envTools := os.Getenv("TRUEFINALS_TOOLS"); envTools != "" {
		toolsMode = envTools
	}

	// Get disabled tools from CLI flag or environment variable
	// CLI flag takes precedence over environment variable
	var disabledToolsStr string
	if *cliDisabledTools != "" {
		disabledToolsStr = *cliDisabledTools
	} else if envDisabledTools := os.Getenv("TRUEFINALS_DISABLED_TOOLS"); envDisabledTools != "" {
		disabledToolsStr = envDisabledTools
	}

	// Parse disabled tools list
	if disabledToolsStr != "" {
		toolsList := strings.Split(disabledToolsStr, ",")
		for _, tool := range toolsList {
			tool = strings.TrimSpace(tool)
			if tool != "" {
				disabledTools = append(disabledTools, tool)
			}
		}
		log.Printf("Disabled tools: %v", disabledTools)
	}

	// Validate tools mode
	switch toolsMode {
	case ToolsReporting, ToolsFullSafe, ToolsFull:
		// Valid mode
	default:
		log.Fatalf("Error: Invalid tools mode '%s'. Valid options: reporting, full-safe, full", toolsMode)
	}

	// Get base URL from CLI flag or environment variable
	// CLI flag takes precedence over environment variable
	if *cliBaseURL != "" {
		APIBaseURL = *cliBaseURL
	} else if envBaseURL := os.Getenv("TRUEFINALS_BASE_URL"); envBaseURL != "" {
		APIBaseURL = envBaseURL
	}

	// Get API key from CLI flag or environment variable
	// CLI flag takes precedence over environment variable
	if *cliAPIKey != "" {
		apiKey = *cliAPIKey
	} else {
		apiKey = os.Getenv("TRUEFINALS_API_KEY")
	}

	// Get API user ID from CLI flag or environment variable
	// CLI flag takes precedence over environment variable
	if *cliAPIUserID != "" {
		apiUserID = *cliAPIUserID
	} else {
		apiUserID = os.Getenv("TRUEFINALS_API_USER_ID")
	}

	if apiKey == "" {
		log.Fatal("Error: API key not provided. Use --api-key flag or set TRUEFINALS_API_KEY environment variable")
	}

	if apiUserID == "" {
		log.Fatal("Error: API User ID not provided. Use --api-user-id flag or set TRUEFINALS_API_USER_ID environment variable")
	}

	log.Println("NHRL MCP Server starting...")

	// Start MCP server
	startMCPServer(*exitAfterFirst)
}

func startMCPServer(exitAfterFirst bool) {
	scanner := bufio.NewScanner(os.Stdin)
	requestCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var request MCPRequest
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			// Only send error response if we can determine there was an ID
			var rawMsg map[string]interface{}
			if json.Unmarshal([]byte(line), &rawMsg) == nil {
				if id, exists := rawMsg["id"]; exists {
					response := sendError(id, -32700, "Parse error", nil)
					if responseJSON, err := json.Marshal(response); err == nil {
						fmt.Println(string(responseJSON))
					}
				}
			}
			continue
		}

		// Check if this is a notification (no ID field)
		var rawMsg map[string]interface{}
		json.Unmarshal([]byte(line), &rawMsg)
		_, hasID := rawMsg["id"]

		if !hasID {
			// This is a notification - handle it but don't send a response
			handleNotification(request)
			continue
		}

		// This is a request - handle it and send a response
		response := handleRequest(request)

		responseJSON, err := json.Marshal(response)
		if err != nil {
			errorResponse := sendError(request.ID, -32603, "Internal error", nil)
			if errorJSON, err := json.Marshal(errorResponse); err == nil {
				fmt.Println(string(errorJSON))
			}
			continue
		}

		fmt.Println(string(responseJSON))

		// Increment request count and check if we should exit
		requestCount++
		if exitAfterFirst && requestCount >= 1 {
			log.Println("Exiting after processing first request as requested")
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading input: %v", err)
	}
}

func handleNotification(request MCPRequest) {
	switch request.Method {
	case "notifications/initialized":
		log.Println("Client initialized")
	case "notifications/cancelled":
		log.Println("Request cancelled")
	default:
		log.Printf("Unknown notification: %s", request.Method)
	}
}

func handleRequest(request MCPRequest) MCPResponse {
	switch request.Method {
	case "initialize":
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{},
				},
				"serverInfo": map[string]interface{}{
					"name":    ServerName,
					"version": Version,
				},
			},
		}

	case "tools/list":
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Result: map[string]interface{}{
				"tools": getAllTools(),
			},
		}

	case "tools/call":
		return handleToolCall(request)

	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: map[string]interface{}{
				"code":    -32601,
				"message": "Method not found",
			},
		}
	}
}

func handleToolCall(request MCPRequest) MCPResponse {
	params, ok := request.Params.(map[string]interface{})
	if !ok {
		return sendError(request.ID, -32602, "Invalid params", nil)
	}

	name, ok := params["name"].(string)
	if !ok {
		return sendError(request.ID, -32602, "Tool name required", nil)
	}

	// Check if tool is explicitly disabled
	if isToolDisabled(name) {
		return sendError(request.ID, -32601, fmt.Sprintf("Tool '%s' is disabled", name), nil)
	}

	// Check if tool is allowed in current tools mode
	if !isToolAllowed(name) {
		return sendError(request.ID, -32601, fmt.Sprintf("Tool '%s' not available in '%s' mode", name, toolsMode), nil)
	}

	args, ok := params["arguments"].(map[string]interface{})
	if !ok {
		args = make(map[string]interface{})
	}

	var result ToolResult

	switch name {
	case "truefinals_tournaments":
		data, err := handleTournamentsTool(args)
		if err != nil {
			result = ToolResult{
				Content: []ToolContent{{Type: "text", Text: fmt.Sprintf("Error: %v", err)}},
				IsError: true,
			}
		} else {
			result = ToolResult{
				Content: []ToolContent{{Type: "text", Text: data}},
				IsError: false,
			}
		}

	case "truefinals_games":
		data, err := handleGamesTool(args)
		if err != nil {
			result = ToolResult{
				Content: []ToolContent{{Type: "text", Text: fmt.Sprintf("Error: %v", err)}},
				IsError: true,
			}
		} else {
			result = ToolResult{
				Content: []ToolContent{{Type: "text", Text: data}},
				IsError: false,
			}
		}

	case "truefinals_locations":
		data, err := handleLocationsTool(args)
		if err != nil {
			result = ToolResult{
				Content: []ToolContent{{Type: "text", Text: fmt.Sprintf("Error: %v", err)}},
				IsError: true,
			}
		} else {
			result = ToolResult{
				Content: []ToolContent{{Type: "text", Text: data}},
				IsError: false,
			}
		}

	case "truefinals_players":
		data, err := handlePlayersTool(args)
		if err != nil {
			result = ToolResult{
				Content: []ToolContent{{Type: "text", Text: fmt.Sprintf("Error: %v", err)}},
				IsError: true,
			}
		} else {
			result = ToolResult{
				Content: []ToolContent{{Type: "text", Text: data}},
				IsError: false,
			}
		}

	case "truefinals_bracket":
		data, err := handleBracketTool(args)
		if err != nil {
			result = ToolResult{
				Content: []ToolContent{{Type: "text", Text: fmt.Sprintf("Error: %v", err)}},
				IsError: true,
			}
		} else {
			result = ToolResult{
				Content: []ToolContent{{Type: "text", Text: data}},
				IsError: false,
			}
		}

	case "nhrl_stats":
		data, err := handleNHRLStatsTool(args)
		if err != nil {
			result = ToolResult{
				Content: []ToolContent{{Type: "text", Text: fmt.Sprintf("Error: %v", err)}},
				IsError: true,
			}
		} else {
			result = ToolResult{
				Content: []ToolContent{{Type: "text", Text: data}},
				IsError: false,
			}
		}

	default:
		return sendError(request.ID, -32601, fmt.Sprintf("Unknown tool: %s", name), nil)
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}
}

func sendError(id interface{}, code int, message string, data interface{}) MCPResponse {
	errorObj := map[string]interface{}{
		"code":    code,
		"message": message,
	}
	if data != nil {
		errorObj["data"] = data
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   errorObj,
	}
}

func getAllTools() []ToolInfo {
	var tools []ToolInfo

	// Add tools if they are allowed
	if isToolAllowed("truefinals_tournaments") {
		tools = append(tools, getTournamentsToolInfo())
	}
	if isToolAllowed("truefinals_games") {
		tools = append(tools, getGamesToolInfo())
	}
	if isToolAllowed("truefinals_locations") {
		tools = append(tools, getLocationsToolInfo())
	}
	if isToolAllowed("truefinals_players") {
		tools = append(tools, getPlayersToolInfo())
	}
	if isToolAllowed("truefinals_bracket") {
		tools = append(tools, getBracketToolInfo())
	}
	if isToolAllowed("nhrl_stats") {
		tools = append(tools, getNHRLStatsToolInfo())
	}

	return tools
}
