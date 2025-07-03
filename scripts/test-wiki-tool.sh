#!/bin/bash
# Test script for NHRL MCP Wiki Tool

# Color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to run a test
run_test() {
    local description="$1"
    local method="$2"
    local params="$3"
    local test_num="$4"
    
    echo -e "\n${BLUE}Test $test_num: $description${NC}"
    echo "Request: $method with params: $params"
    
    # Create request JSON
    request_json="{\"jsonrpc\":\"2.0\",\"id\":$test_num,\"method\":\"$method\",\"params\":$params}"
    
    # Send request and capture response
    response=$(echo "$request_json" | ./nhrl-mcp-server --api-key="$TRUEFINALS_API_KEY" --api-user-id="$TRUEFINALS_API_USER_ID" --exit-after-first 2>/dev/null | tail -1)
    
    # Check if response contains error
    if echo "$response" | grep -q '"error"'; then
        echo -e "${RED}FAILED${NC}"
        echo "Response: $response"
    else
        echo -e "${GREEN}SUCCESS${NC}"
        # Extract and format the result
        if echo "$response" | grep -q '"content"'; then
            echo "Response received (showing first 200 chars):"
            echo "$response" | jq -r '.result.content[0].text' 2>/dev/null | head -c 200
            echo "..."
        fi
    fi
}

echo "Starting NHRL MCP Wiki Tool Tests"
echo "================================="

# Check for required environment variables
if [ -z "$TRUEFINALS_API_KEY" ] || [ -z "$TRUEFINALS_API_USER_ID" ]; then
    echo -e "${RED}Error: TRUEFINALS_API_KEY and TRUEFINALS_API_USER_ID environment variables must be set${NC}"
    exit 1
fi

# Initialize the server first
echo -e "\n${BLUE}Initializing server...${NC}"
init_request='{"jsonrpc":"2.0","id":0,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{}}}'
echo "$init_request" | ./nhrl-mcp-server --api-key="$TRUEFINALS_API_KEY" --api-user-id="$TRUEFINALS_API_USER_ID" --exit-after-first >/dev/null 2>&1

# Test 1: Search for robot rules
run_test "Search for robot rules" \
    "tools/call" \
    '{"name":"nhrl_wiki","arguments":{"operation":"search","query":"robot rules","limit":5}}' \
    1

# Test 2: Search for weight classes
run_test "Search for weight classes" \
    "tools/call" \
    '{"name":"nhrl_wiki","arguments":{"operation":"search","query":"weight class specifications"}}' \
    2

# Test 3: Get a specific page extract
run_test "Get Main Page extract" \
    "tools/call" \
    '{"name":"nhrl_wiki","arguments":{"operation":"get_page_extract","title":"Main Page"}}' \
    3

# Test 4: Get full page content (if a page exists)
run_test "Get full page content for a known page" \
    "tools/call" \
    '{"name":"nhrl_wiki","arguments":{"operation":"get_page","title":"Main Page"}}' \
    4

# Test 5: Search for safety requirements
run_test "Search for safety requirements" \
    "tools/call" \
    '{"name":"nhrl_wiki","arguments":{"operation":"search","query":"safety requirements"}}' \
    5

# Test 6: Test error handling - missing query
run_test "Error handling - missing query" \
    "tools/call" \
    '{"name":"nhrl_wiki","arguments":{"operation":"search"}}' \
    6

# Test 7: Test error handling - invalid operation
run_test "Error handling - invalid operation" \
    "tools/call" \
    '{"name":"nhrl_wiki","arguments":{"operation":"invalid_op"}}' \
    7

echo -e "\n${GREEN}Wiki tool tests completed!${NC}" 