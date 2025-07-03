package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// Wiki API configuration
const (
	WikiBaseURL = "https://wiki.nhrl.io/wiki/api.php"
)

// HTTP client for Wiki API with timeout
var wikiHttpClient = &http.Client{
	Timeout: 30 * time.Second,
}

// Wiki API response structures
type WikiSearchResult struct {
	Pageid    int    `json:"pageid"`
	Ns        int    `json:"ns"`
	Title     string `json:"title"`
	Size      int    `json:"size"`
	Wordcount int    `json:"wordcount"`
	Snippet   string `json:"snippet"`
	Timestamp string `json:"timestamp"`
}

type WikiSearchResponse struct {
	Query struct {
		SearchInfo struct {
			Totalhits int `json:"totalhits"`
		} `json:"searchinfo"`
		Search []WikiSearchResult `json:"search"`
	} `json:"query"`
}

type WikiPageContent struct {
	Parse struct {
		Title    string `json:"title"`
		Pageid   int    `json:"pageid"`
		Wikitext struct {
			Text string `json:"*"`
		} `json:"wikitext"`
		Text struct {
			Text string `json:"*"`
		} `json:"text"`
	} `json:"parse"`
}

type WikiPageExtract struct {
	Query struct {
		Pages map[string]struct {
			Pageid  int    `json:"pageid"`
			Ns      int    `json:"ns"`
			Title   string `json:"title"`
			Extract string `json:"extract"`
		} `json:"pages"`
	} `json:"query"`
}

// handleNHRLWikiTool handles all NHRL wiki operations
func handleNHRLWikiTool(args map[string]interface{}) (string, error) {
	operation, ok := args["operation"].(string)
	if !ok {
		return "", fmt.Errorf("operation parameter is required")
	}

	// Check if operation is allowed in current tools mode
	if !isOperationAllowed("nhrl_wiki", operation) {
		return "", fmt.Errorf(getOperationNotAllowedError(operation))
	}

	switch operation {
	case "search":
		return searchNHRLWiki(args)
	case "get_page":
		return getNHRLWikiPage(args)
	case "get_page_extract":
		return getNHRLWikiPageExtract(args)
	default:
		return "", fmt.Errorf("unknown operation: %s", operation)
	}
}

// getNHRLWikiToolInfo returns the tool definition for NHRL wiki operations
func getNHRLWikiToolInfo() ToolInfo {
	return ToolInfo{
		Name: "nhrl_wiki",
		Description: `Browse and search the NHRL (National Havoc Robot League) wiki at https://wiki.nhrl.io.

This tool provides access to:
- Search functionality to find relevant wiki pages
- Page content retrieval to read specific wiki articles
- Page extracts for quick summaries

The NHRL wiki contains detailed information about:
- Robot combat rules and regulations
- Weight class specifications
- Tournament formats and procedures
- Safety requirements
- Technical specifications
- Historical information
- Builder resources`,
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"operation": map[string]interface{}{
					"type": "string",
					"description": `The wiki operation to perform:

- search: Search for wiki pages by keywords
- get_page: Get the full content of a specific wiki page
- get_page_extract: Get a plain text extract/summary of a wiki page`,
					"enum": []string{"search", "get_page", "get_page_extract"},
				},
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Search query for finding wiki pages (required for search operation)",
				},
				"title": map[string]interface{}{
					"type":        "string",
					"description": "The exact title of the wiki page to retrieve (required for get_page and get_page_extract operations). Case-sensitive.",
				},
				"limit": map[string]interface{}{
					"type":        "number",
					"description": "Maximum number of search results to return. Defaults to 10, max 50.",
				},
			},
			"required": []string{"operation"},
		},
	}
}

// searchNHRLWiki searches the wiki for pages matching the query
func searchNHRLWiki(args map[string]interface{}) (string, error) {
	query, ok := args["query"].(string)
	if !ok || query == "" {
		return "", fmt.Errorf("query is required for search operation")
	}

	limit := 10
	if l, ok := args["limit"].(float64); ok && l > 0 && l <= 50 {
		limit = int(l)
	}

	// Build query parameters
	params := url.Values{}
	params.Set("action", "query")
	params.Set("list", "search")
	params.Set("srsearch", query)
	params.Set("srlimit", fmt.Sprintf("%d", limit))
	params.Set("format", "json")
	params.Set("srinfo", "totalhits")
	params.Set("srprop", "size|wordcount|timestamp|snippet")

	// Make API request
	resp, err := wikiHttpClient.Get(WikiBaseURL + "?" + params.Encode())
	if err != nil {
		return "", fmt.Errorf("failed to search wiki: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var searchResp WikiSearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return "", fmt.Errorf("failed to parse search response: %w", err)
	}

	// Format results
	results := make([]map[string]interface{}, 0, len(searchResp.Query.Search))
	for _, result := range searchResp.Query.Search {
		// Clean up snippet HTML
		snippet := regexp.MustCompile(`<[^>]+>`).ReplaceAllString(result.Snippet, "")
		snippet = strings.ReplaceAll(snippet, "&quot;", "\"")
		snippet = strings.ReplaceAll(snippet, "&amp;", "&")
		snippet = strings.ReplaceAll(snippet, "&lt;", "<")
		snippet = strings.ReplaceAll(snippet, "&gt;", ">")

		results = append(results, map[string]interface{}{
			"title":     result.Title,
			"pageid":    result.Pageid,
			"size":      result.Size,
			"wordcount": result.Wordcount,
			"snippet":   snippet,
			"timestamp": result.Timestamp,
		})
	}

	output := map[string]interface{}{
		"query":        query,
		"total_hits":   searchResp.Query.SearchInfo.Totalhits,
		"result_count": len(results),
		"results":      results,
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal results: %w", err)
	}

	return string(jsonData), nil
}

// getNHRLWikiPage retrieves the full content of a wiki page
func getNHRLWikiPage(args map[string]interface{}) (string, error) {
	title, ok := args["title"].(string)
	if !ok || title == "" {
		return "", fmt.Errorf("title is required for get_page operation")
	}

	// Build query parameters for wikitext
	params := url.Values{}
	params.Set("action", "parse")
	params.Set("page", title)
	params.Set("prop", "wikitext")
	params.Set("format", "json")

	// Make API request
	resp, err := wikiHttpClient.Get(WikiBaseURL + "?" + params.Encode())
	if err != nil {
		return "", fmt.Errorf("failed to get wiki page: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Check for error in response
	var errorResp map[string]interface{}
	if err := json.Unmarshal(body, &errorResp); err == nil {
		if errInfo, ok := errorResp["error"].(map[string]interface{}); ok {
			return "", fmt.Errorf("wiki API error: %v", errInfo["info"])
		}
	}

	var pageContent WikiPageContent
	if err := json.Unmarshal(body, &pageContent); err != nil {
		return "", fmt.Errorf("failed to parse page response: %w", err)
	}

	// Format the wikitext for better readability
	wikitext := pageContent.Parse.Wikitext.Text

	// Convert some common wiki markup to more readable format
	// Remove template calls that are too complex
	wikitext = regexp.MustCompile(`\{\{[^}]+\}\}`).ReplaceAllStringFunc(wikitext, func(match string) string {
		// Keep simple templates like {{NHRL}} but remove complex ones
		if strings.Count(match, "|") <= 1 && len(match) < 50 {
			return match
		}
		return ""
	})

	// Clean up multiple blank lines
	wikitext = regexp.MustCompile(`\n{3,}`).ReplaceAllString(wikitext, "\n\n")

	output := map[string]interface{}{
		"title":   pageContent.Parse.Title,
		"pageid":  pageContent.Parse.Pageid,
		"content": wikitext,
		"url":     fmt.Sprintf("https://wiki.nhrl.io/wiki/index.php/%s", url.QueryEscape(strings.ReplaceAll(title, " ", "_"))),
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal results: %w", err)
	}

	return string(jsonData), nil
}

// getNHRLWikiPageExtract retrieves a plain text extract of a wiki page
func getNHRLWikiPageExtract(args map[string]interface{}) (string, error) {
	title, ok := args["title"].(string)
	if !ok || title == "" {
		return "", fmt.Errorf("title is required for get_page_extract operation")
	}

	// Build query parameters
	params := url.Values{}
	params.Set("action", "query")
	params.Set("titles", title)
	params.Set("prop", "extracts")
	params.Set("exintro", "true")
	params.Set("explaintext", "true")
	params.Set("exsectionformat", "plain")
	params.Set("format", "json")

	// Make API request
	resp, err := wikiHttpClient.Get(WikiBaseURL + "?" + params.Encode())
	if err != nil {
		return "", fmt.Errorf("failed to get wiki page extract: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var extractResp WikiPageExtract
	if err := json.Unmarshal(body, &extractResp); err != nil {
		return "", fmt.Errorf("failed to parse extract response: %w", err)
	}

	// Extract the page data (MediaWiki returns pages as a map with page ID as key)
	var pageData struct {
		Pageid  int    `json:"pageid"`
		Title   string `json:"title"`
		Extract string `json:"extract"`
	}

	for _, page := range extractResp.Query.Pages {
		pageData.Pageid = page.Pageid
		pageData.Title = page.Title
		pageData.Extract = page.Extract
		break // We only expect one page
	}

	if pageData.Extract == "" {
		return "", fmt.Errorf("page not found or has no content: %s", title)
	}

	output := map[string]interface{}{
		"title":   pageData.Title,
		"pageid":  pageData.Pageid,
		"extract": pageData.Extract,
		"url":     fmt.Sprintf("https://wiki.nhrl.io/wiki/index.php/%s", url.QueryEscape(strings.ReplaceAll(title, " ", "_"))),
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal results: %w", err)
	}

	return string(jsonData), nil
}
