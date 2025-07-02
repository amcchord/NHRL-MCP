#!/bin/bash

# Test script to verify get_bot_picture_url works in read-only mode

echo "Testing get_bot_picture_url in read-only mode..."
echo ""

# Build the server first
echo "Building NHRL MCP server..."
go build -o nhrl-mcp-server

# Test with read-only flag
echo "Testing with --read-only flag..."
echo '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"nhrl_stats","arguments":{"operation":"get_bot_picture_url","bot_name":"Lynx"}}}' | ./nhrl-mcp-server --read-only --api-key "$TRUEFINALS_API_KEY" --api-user-id "$TRUEFINALS_API_USER_ID" --exit-after-first | jq -r '.result.content[0].text' | jq .

echo ""
echo "Testing with TRUEFINALS_READ_ONLY=true environment variable..."
TRUEFINALS_READ_ONLY=true ./nhrl-mcp-server --api-key "$TRUEFINALS_API_KEY" --api-user-id "$TRUEFINALS_API_USER_ID" --exit-after-first <<EOF | jq -r '.result.content[0].text' | jq .
{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"nhrl_stats","arguments":{"operation":"get_bot_picture_url","bot_name":"Huge"}}}
EOF

echo ""
echo "Testing that write operations are blocked in read-only mode..."
echo '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"truefinals_tournaments","arguments":{"operation":"create","title":"Test Tournament"}}}' | ./nhrl-mcp-server --read-only --api-key "$TRUEFINALS_API_KEY" --api-user-id "$TRUEFINALS_API_USER_ID" --exit-after-first | jq .

echo ""
echo "✅ Test complete! get_bot_picture_url should work in read-only mode." 