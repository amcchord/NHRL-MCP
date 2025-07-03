# NHRL Wiki Tool Documentation

## Overview

The `nhrl_wiki` tool provides access to the NHRL (National Havoc Robot League) wiki at https://wiki.nhrl.io, allowing LLMs to search for and retrieve information about robot combat rules, regulations, specifications, and other NHRL-related content.

## Features

- **Search**: Find wiki pages by keywords
- **Get Page**: Retrieve the full content of a specific wiki page
- **Get Page Extract**: Get a plain text summary/extract of a wiki page

## Tool Configuration

The wiki tool is enabled by default and is considered a read-only tool, making it safe to use in all tool modes:
- ✅ Available in `reporting` mode
- ✅ Available in `full-safe` mode  
- ✅ Available in `full` mode

## Operations

### 1. Search (`search`)

Search for wiki pages containing specific keywords.

**Parameters:**
- `operation`: "search" (required)
- `query`: Search terms (required)
- `limit`: Maximum number of results (optional, default: 10, max: 50)

**Example:**
```json
{
  "name": "nhrl_wiki",
  "arguments": {
    "operation": "search",
    "query": "weight class specifications",
    "limit": 5
  }
}
```

### 2. Get Page (`get_page`)

Retrieve the full wiki markup content of a specific page.

**Parameters:**
- `operation`: "get_page" (required)
- `title`: Exact page title (required, case-sensitive)

**Example:**
```json
{
  "name": "nhrl_wiki",
  "arguments": {
    "operation": "get_page",
    "title": "Main Page"
  }
}
```

### 3. Get Page Extract (`get_page_extract`)

Get a plain text extract/summary of a wiki page.

**Parameters:**
- `operation`: "get_page_extract" (required)
- `title`: Exact page title (required, case-sensitive)

**Example:**
```json
{
  "name": "nhrl_wiki",
  "arguments": {
    "operation": "get_page_extract",
    "title": "Safety Requirements"
  }
}
```

## Common Use Cases

1. **Looking up rules and regulations**
   - Search for "robot rules" or "competition rules"
   - Get specific rule pages

2. **Finding weight class information**
   - Search for "weight class" or specific classes like "beetleweight"
   - Get detailed specifications

3. **Safety requirements**
   - Search for "safety" or "safety requirements"
   - Get safety checklists and requirements

4. **Tournament information**
   - Search for "tournament format" or "qualification"
   - Get tournament procedures and formats

5. **Technical specifications**
   - Search for specific technical terms
   - Get detailed technical requirements

## Error Handling

The tool handles common errors:
- Missing required parameters
- Invalid operations
- Pages not found
- Network errors

## Integration Notes

- The tool uses the MediaWiki API
- Results are returned in JSON format
- Wiki markup is cleaned for better readability when using `get_page`
- Extracts are plain text for easy consumption

## Testing

A test script is provided at `scripts/test-wiki-tool.sh` to verify the wiki tool functionality.

Run it with:
```bash
./scripts/test-wiki-tool.sh
``` 