# AI Chat Tool Calling Documentation

This document describes how Spotigo's AI chat uses function calling (tool calling) to query your music library efficiently.

## Overview

When you ask the AI chat questions about your music library, it doesn't load all your data into context. Instead, it uses **function calling** to execute precise queries against your structured JSON data. This approach:

- **Reduces context window usage** - Only relevant data is retrieved
- **Improves accuracy** - Structured queries are more reliable than text embeddings
- **Speeds up responses** - Direct JSON queries are fast
- **Preserves structure** - No loss of information from chunking

## How It Works

1. **User asks a question** - e.g., "How many tracks do I have?"
2. **Model decides to call a tool** - Selects `get_library_stats` function
3. **Tool executes** - Queries JSON files and returns structured data
4. **Model receives results** - Gets actual counts from your library
5. **Model provides answer** - "You have 1,234 tracks in your library"

## Available Tools

### 1. `get_library_stats`

Get overall statistics about your music library.

**Parameters:** None

**Example Query:**
- "How many tracks do I have?"
- "What's in my library?"
- "Give me my library stats"

**Returns:**
```json
{
  "data": {
    "saved_tracks": 1234,
    "playlists": 25,
    "followed_artists": 42
  },
  "summary": "Library: 1234 tracks, 25 playlists, 42 followed artists"
}
```

### 2. `search_tracks`

Search for tracks by text (matches artist, song title, album).

**Parameters:**
- `query` (required): Search term
- `limit` (optional): Max results to return (default: 10)

**Example Queries:**
- "Find songs by Queen"
- "Search for Bohemian Rhapsody"
- "Do I have any Beatles tracks?"

**Returns:**
```json
{
  "count": 2,
  "data": [
    {
      "track": {
        "name": "Bohemian Rhapsody",
        "artists": [{"name": "Queen"}],
        "album": {"name": "A Night at the Opera"}
      }
    }
  ]
}
```

### 3. `get_tracks_by_artist`

Get all tracks by a specific artist.

**Parameters:**
- `artist_name` (required): Name of the artist

**Example Queries:**
- "Show me all Queen songs"
- "What tracks do I have by Led Zeppelin?"
- "List all songs from The Beatles"

**Returns:**
```json
{
  "count": 15,
  "data": [/* array of tracks */]
}
```

### 4. `get_recently_added_tracks`

Get the most recently added tracks.

**Parameters:**
- `limit` (optional): Number of tracks (default: 10)

**Example Queries:**
- "What did I add recently?"
- "Show my last 20 saved tracks"
- "What's new in my library?"

**Returns:**
```json
{
  "count": 10,
  "data": [/* tracks sorted by added_at desc */]
}
```

### 5. `get_all_artists`

Get all unique artists in your library.

**Parameters:** None

**Example Queries:**
- "List all my artists"
- "Who do I listen to?"
- "Show me all artists I follow"

**Returns:**
```json
{
  "count": 42,
  "data": ["Queen", "Led Zeppelin", "The Beatles", ...]
}
```

### 6. `get_playlist_by_name`

Find a playlist by name.

**Parameters:**
- `playlist_name` (required): Name of the playlist

**Example Queries:**
- "Find my 'Workout Mix' playlist"
- "Show me the 'Road Trip' playlist"
- "Do I have a playlist called 'Chill Vibes'?"

**Returns:**
```json
{
  "data": {
    "name": "Workout Mix",
    "tracks": {"total": 30},
    "description": "High energy tracks"
  }
}
```

### 7. `query_music_data`

Execute custom queries with filtering, sorting, aggregation.

**Parameters:**
- `source` (required): `saved_tracks.json`, `playlists.json`, or `followed_artists.json`
- `operation` (required): `select`, `count`, `filter`, `search`, `sort`, `distinct`, `aggregate`, `stats`
- `filters` (optional): Array of filter conditions
- `sort_by` (optional): Field to sort by
- `sort_order` (optional): `asc` or `desc`
- `limit` (optional): Limit results
- `field` (optional): Specific field to extract

**Example Queries:**
- "Find tracks with popularity over 90"
- "Show me my oldest playlists"
- "What's my most popular saved track?"

**Filter Operators:**
- `eq` - Equal
- `ne` - Not equal
- `gt` - Greater than
- `gte` - Greater than or equal
- `lt` - Less than
- `lte` - Less than or equal
- `contains` - Contains substring
- `regex` - Regular expression match

**Example:**
```json
{
  "source": "saved_tracks.json",
  "operation": "filter",
  "filters": [
    {
      "field": "track.popularity",
      "operator": "gte",
      "value": 90
    }
  ],
  "sort_by": "track.popularity",
  "sort_order": "desc",
  "limit": 10
}
```

## Usage Examples

### Starting the Chat with Tools

```bash
# Enable tools (default)
spotigo chat

# Disable tools
spotigo chat --tools=false

# Specify custom data directory
spotigo chat --data-dir ./my-music-data
```

### Example Conversation

```
You: How many tracks do I have?

🔧 Calling tool: get_library_stats
   Arguments: {}

Spotigo: You have 1,234 saved tracks, 25 playlists, and 42 followed artists in your library.
```

```
You: Find my most popular tracks

🔧 Calling tool: query_music_data
   Arguments: {
     "source": "saved_tracks.json",
     "operation": "select",
     "sort_by": "track.popularity",
     "sort_order": "desc",
     "limit": 5
   }

Spotigo: Here are your top 5 most popular tracks:
1. "Bohemian Rhapsody" by Queen (popularity: 95)
2. "Stairway to Heaven" by Led Zeppelin (popularity: 92)
3. "Hey Jude" by The Beatles (popularity: 91)
...
```

```
You: What did I add this month?

🔧 Calling tool: get_recently_added_tracks
   Arguments: {"limit": 30}

Spotigo: You've added 23 tracks this month, including "New Song" by Artist X and "Another Track" by Artist Y...
```

## Tool Execution Flow

```
┌─────────────────┐
│   User Query    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Model Decides  │
│  Which Tool(s)  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Execute Tool   │
│  (Query JSON)   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Return Results  │
│  to Model       │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Model Formats  │
│ Natural Answer  │
└─────────────────┘
```

## Advanced Features

### Multiple Tool Calls

The model can call multiple tools in sequence to answer complex questions:

```
You: Compare my Queen and Beatles tracks

Tool 1: get_tracks_by_artist("Queen")
Tool 2: get_tracks_by_artist("The Beatles")

Spotigo: You have 15 Queen tracks and 12 Beatles tracks...
```

### Chained Queries

The model can use results from one tool to inform the next:

```
You: What's my most followed artist's most popular song?

Tool 1: query_music_data(sort by followers)
Tool 2: get_tracks_by_artist(top artist from Tool 1)
Tool 3: query_music_data(filter by that artist, sort by popularity)

Spotigo: Your most followed artist is Queen with 35M followers, and their most popular song in your library is "Bohemian Rhapsody".
```

## Configuration

### Model Configuration

Edit `config/models.yaml` to configure the chat model:

```yaml
agents:
  chat_agent:
    role: "chat"
    model: "granite4:1b"
    fallback: "qwen3:0.6b"
    system_prompt: "You are Spotigo, a friendly music intelligence assistant..."
```

### CLI Flags

- `--model`: Override the default chat model
- `--context`: Set context window size (default: 4096)
- `--tools`: Enable/disable tool calling (default: true)
- `--data-dir`: Set music data directory (default: ./data)

## Best Practices

1. **Be specific in queries** - "Find tracks by Queen" is better than "Queen songs"
2. **Use natural language** - The model understands conversational queries
3. **Ask follow-up questions** - Context is preserved in the conversation
4. **Use 'clear' to reset** - Start fresh if context gets confused

## Troubleshooting

### Tools Not Working

1. **Check data directory**: Ensure JSON files exist in the data directory
   ```bash
   ls ./data
   # Should show: saved_tracks.json, playlists.json, followed_artists.json
   ```

2. **Verify tools are enabled**:
   ```bash
   spotigo chat --tools=true
   ```

3. **Check model compatibility**: Some models may not support function calling
   - Use granite4:1b or qwen3:0.6b for best results

### Empty Results

- Run `spotigo backup` to ensure you have recent data
- Check that JSON files are not empty
- Verify file permissions

### Tool Errors

If you see tool execution errors:
- Check the arguments in the debug output
- Ensure JSON files are valid
- Look for error messages in the tool result

## Architecture Notes

### Why Function Calling?

Traditional RAG (Retrieval Augmented Generation) for JSON data has issues:

- **Text chunking loses structure** - JSON arrays and objects get split awkwardly
- **Embeddings are imprecise** - Semantic search doesn't work well for exact queries
- **Context window bloat** - Loading entire JSON files wastes tokens
- **No aggregations** - Can't count, sum, or group data efficiently

Function calling solves these by:

- **Preserving structure** - Queries work on actual JSON objects
- **Exact matching** - Filter by precise field values
- **Minimal context** - Only return what's needed
- **Rich operations** - Count, sort, aggregate, filter

### JSON Query Engine

Under the hood, tools use the `internal/jsonquery` package:

- Loads JSON files on demand
- Caches parsed data for performance
- Supports complex filtering and sorting
- Returns structured results

See `internal/jsonquery/query.go` for implementation details.

## Examples in Code

### Programmatic Tool Usage

```go
import (
    "github.com/bkataru/spotigo/internal/tools"
    "github.com/bkataru/spotigo/internal/ollama"
)

// Create music tools
musicTools := tools.NewMusicTools("./data")

// Execute a tool
toolCall := ollama.ToolCall{
    Function: ollama.FunctionCall{
        Name:      "search_tracks",
        Arguments: `{"query": "Queen", "limit": 10}`,
    },
}

result, err := musicTools.ExecuteToolCall(toolCall)
// result contains JSON string with search results
```

### Testing Tools

See `internal/tools/tools_test.go` for comprehensive test examples.

## Further Reading

- [README.md](../README.md) - General Spotigo documentation
- [JSON Query Documentation](../internal/jsonquery/query.go) - Query engine details
- [Ollama Function Calling](https://ollama.com/blog/function-calling) - Ollama docs

## Contributing

To add new tools:

1. Add tool definition in `internal/tools/tools.go` → `GetToolDefinitions()`
2. Implement execution logic in `executeXXX()` method
3. Add comprehensive tests in `internal/tools/tools_test.go`
4. Update this documentation with examples

See existing tools as templates.