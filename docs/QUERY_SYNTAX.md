# JSON Query Syntax Documentation

This document describes the powerful JSON query syntax used by Spotigo's query engine for structured music data queries.

## Table of Contents

- [Overview](#overview)
- [Query Structure](#query-structure)
- [Operations](#operations)
- [Filter Operators](#filter-operators)
- [Field Paths](#field-paths)
- [Examples](#examples)
- [Best Practices](#best-practices)

## Overview

The JSON query engine provides a structured way to query, filter, sort, and aggregate your Spotify music data without loading entire JSON files into memory or relying on imprecise text embeddings.

**Key Features:**
- 🎯 **Precise filtering** - Exact field matching with multiple operators
- 📊 **Aggregation support** - Count, sum, average, min, max, group by
- 🔍 **Text search** - Full-text search across fields
- 📈 **Sorting & pagination** - Order results and limit output
- ⚡ **Performance** - In-memory caching for repeated queries
- 🎵 **Music-specific helpers** - Pre-built queries for common operations

## Query Structure

A query is a JSON object with the following structure:

```json
{
  "source": "saved_tracks.json",
  "operation": "filter",
  "filters": [
    {
      "field": "track.popularity",
      "operator": "gte",
      "value": 80
    }
  ],
  "sort_by": "track.name",
  "sort_order": "asc",
  "limit": 10,
  "offset": 0
}
```

### Required Fields

- **`source`** (string): The JSON file to query
  - `saved_tracks.json` - Your saved tracks
  - `playlists.json` - Your playlists
  - `followed_artists.json` - Artists you follow

- **`operation`** (string): The operation to perform
  - `select` - Retrieve items
  - `count` - Count matching items
  - `filter` - Filter items by conditions
  - `search` - Full-text search
  - `sort` - Sort items
  - `distinct` - Get unique values
  - `aggregate` - Perform aggregations
  - `stats` - Generate statistics

### Optional Fields

- **`field`** (string): Specific field to extract or operate on
- **`filters`** (array): Filter conditions (see [Filter Operators](#filter-operators))
- **`sort_by`** (string): Field to sort by
- **`sort_order`** (string): Sort direction (`asc` or `desc`)
- **`limit`** (integer): Maximum number of results
- **`offset`** (integer): Skip this many results (pagination)
- **`search_term`** (string): Text to search for
- **`agg_func`** (string): Aggregation function (`count`, `sum`, `avg`, `min`, `max`, `group`)
- **`group_by`** (string): Field to group by for aggregations

## Operations

### 1. Select

Retrieve items from the data source.

```json
{
  "source": "saved_tracks.json",
  "operation": "select",
  "limit": 10
}
```

**Result:**
```json
{
  "count": 10,
  "data": [/* array of track objects */],
  "summary": "Retrieved 10 items"
}
```

### 2. Count

Count items matching criteria.

```json
{
  "source": "saved_tracks.json",
  "operation": "count"
}
```

**Result:**
```json
{
  "count": 1234,
  "summary": "Found 1234 items"
}
```

### 3. Filter

Filter items by conditions.

```json
{
  "source": "saved_tracks.json",
  "operation": "filter",
  "filters": [
    {
      "field": "track.artists.0.name",
      "operator": "eq",
      "value": "Queen"
    }
  ]
}
```

### 4. Search

Full-text search across all fields.

```json
{
  "source": "saved_tracks.json",
  "operation": "search",
  "search_term": "bohemian rhapsody",
  "limit": 5
}
```

### 5. Sort

Sort items by a field.

```json
{
  "source": "saved_tracks.json",
  "operation": "sort",
  "sort_by": "track.popularity",
  "sort_order": "desc",
  "limit": 10
}
```

### 6. Distinct

Get unique values for a field.

```json
{
  "source": "saved_tracks.json",
  "operation": "distinct",
  "field": "track.artists.0.name"
}
```

**Result:**
```json
{
  "count": 42,
  "data": ["Queen", "The Beatles", "Led Zeppelin", ...],
  "summary": "Found 42 unique values"
}
```

### 7. Aggregate

Perform aggregations on numeric fields.

```json
{
  "source": "saved_tracks.json",
  "operation": "aggregate",
  "field": "track.popularity",
  "agg_func": "avg"
}
```

**Available aggregation functions:**
- `count` - Count items
- `sum` - Sum numeric values
- `avg` - Average of values
- `min` - Minimum value
- `max` - Maximum value
- `group` - Group by field

### 8. Stats

Generate comprehensive statistics.

```json
{
  "source": "saved_tracks.json",
  "operation": "stats",
  "field": "track.duration_ms"
}
```

**Result:**
```json
{
  "data": {
    "count": 1234,
    "min": 120000,
    "max": 480000,
    "avg": 240000,
    "sum": 296160000
  },
  "summary": "Statistics for track.duration_ms"
}
```

## Filter Operators

### Comparison Operators

#### `eq` - Equal

```json
{
  "field": "track.name",
  "operator": "eq",
  "value": "Bohemian Rhapsody"
}
```

#### `ne` - Not Equal

```json
{
  "field": "track.explicit",
  "operator": "ne",
  "value": true
}
```

#### `gt` - Greater Than

```json
{
  "field": "track.popularity",
  "operator": "gt",
  "value": 80
}
```

#### `gte` - Greater Than or Equal

```json
{
  "field": "track.duration_ms",
  "operator": "gte",
  "value": 300000
}
```

#### `lt` - Less Than

```json
{
  "field": "track.popularity",
  "operator": "lt",
  "value": 50
}
```

#### `lte` - Less Than or Equal

```json
{
  "field": "track.duration_ms",
  "operator": "lte",
  "value": 180000
}
```

### String Operators

#### `contains` - Contains Substring (case-insensitive)

```json
{
  "field": "track.name",
  "operator": "contains",
  "value": "love"
}
```

#### `regex` - Regular Expression Match

```json
{
  "field": "track.name",
  "operator": "regex",
  "value": "^The\\s"
}
```

### Existence Operators

#### `exists` - Field Exists

```json
{
  "field": "track.preview_url",
  "operator": "exists",
  "value": true
}
```

#### `in` - Value in List

```json
{
  "field": "track.artists.0.name",
  "operator": "in",
  "value": ["Queen", "The Beatles", "Led Zeppelin"]
}
```

## Field Paths

Field paths use **dot notation** to access nested fields.

### Saved Tracks Structure

```json
{
  "items": [
    {
      "track": {
        "id": "track_id",
        "name": "Song Name",
        "artists": [
          {
            "name": "Artist Name",
            "id": "artist_id"
          }
        ],
        "album": {
          "name": "Album Name",
          "release_date": "2024-01-15"
        },
        "popularity": 95,
        "duration_ms": 354000,
        "explicit": false
      },
      "added_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

### Common Field Paths

| Field Path | Description |
|------------|-------------|
| `track.name` | Song name |
| `track.artists.0.name` | First artist name |
| `track.album.name` | Album name |
| `track.album.release_date` | Release date |
| `track.popularity` | Popularity score (0-100) |
| `track.duration_ms` | Duration in milliseconds |
| `track.explicit` | Explicit content flag |
| `added_at` | Date added to library |

### Playlist Structure

```json
{
  "items": [
    {
      "id": "playlist_id",
      "name": "Playlist Name",
      "owner": {
        "display_name": "User Name"
      },
      "tracks": {
        "total": 50
      },
      "description": "Playlist description",
      "public": true
    }
  ]
}
```

### Followed Artists Structure

```json
{
  "artists": {
    "items": [
      {
        "id": "artist_id",
        "name": "Artist Name",
        "genres": ["rock", "classic rock"],
        "followers": {
          "total": 35000000
        },
        "popularity": 88
      }
    ]
  }
}
```

## Examples

### Example 1: Find Popular Rock Songs

```json
{
  "source": "saved_tracks.json",
  "operation": "filter",
  "filters": [
    {
      "field": "track.popularity",
      "operator": "gte",
      "value": 80
    },
    {
      "field": "track.album.genres",
      "operator": "contains",
      "value": "rock"
    }
  ],
  "sort_by": "track.popularity",
  "sort_order": "desc",
  "limit": 20
}
```

### Example 2: Find Recently Added Songs

```json
{
  "source": "saved_tracks.json",
  "operation": "select",
  "sort_by": "added_at",
  "sort_order": "desc",
  "limit": 10
}
```

### Example 3: Count Songs by Artist

```json
{
  "source": "saved_tracks.json",
  "operation": "filter",
  "filters": [
    {
      "field": "track.artists.0.name",
      "operator": "eq",
      "value": "Queen"
    }
  ]
}
```

### Example 4: Find Long Songs (>5 minutes)

```json
{
  "source": "saved_tracks.json",
  "operation": "filter",
  "filters": [
    {
      "field": "track.duration_ms",
      "operator": "gt",
      "value": 300000
    }
  ],
  "sort_by": "track.duration_ms",
  "sort_order": "desc"
}
```

### Example 5: Search for Songs with "Night" in Title

```json
{
  "source": "saved_tracks.json",
  "operation": "search",
  "search_term": "night",
  "limit": 20
}
```

### Example 6: Get All Unique Artists

```json
{
  "source": "saved_tracks.json",
  "operation": "distinct",
  "field": "track.artists.0.name"
}
```

### Example 7: Average Song Popularity

```json
{
  "source": "saved_tracks.json",
  "operation": "aggregate",
  "field": "track.popularity",
  "agg_func": "avg"
}
```

### Example 8: Find Playlists with Many Tracks

```json
{
  "source": "playlists.json",
  "operation": "filter",
  "filters": [
    {
      "field": "tracks.total",
      "operator": "gte",
      "value": 50
    }
  ],
  "sort_by": "tracks.total",
  "sort_order": "desc"
}
```

### Example 9: Find Artists by Genre

```json
{
  "source": "followed_artists.json",
  "operation": "filter",
  "filters": [
    {
      "field": "genres",
      "operator": "contains",
      "value": "jazz"
    }
  ],
  "sort_by": "followers.total",
  "sort_order": "desc"
}
```

### Example 10: Multiple Filter Conditions (AND)

```json
{
  "source": "saved_tracks.json",
  "operation": "filter",
  "filters": [
    {
      "field": "track.popularity",
      "operator": "gte",
      "value": 70
    },
    {
      "field": "track.duration_ms",
      "operator": "lte",
      "value": 240000
    },
    {
      "field": "track.explicit",
      "operator": "eq",
      "value": false
    }
  ],
  "limit": 50
}
```

## Best Practices

### 1. Use Specific Fields

✅ **Good:**
```json
{
  "field": "track.popularity",
  "operation": "select"
}
```

❌ **Bad:**
```json
{
  "operation": "select"  // Returns entire documents
}
```

### 2. Limit Results

Always use `limit` for large datasets:

```json
{
  "operation": "select",
  "limit": 100
}
```

### 3. Use Efficient Filters

✅ **Efficient:**
```json
{
  "field": "track.popularity",
  "operator": "gte",
  "value": 80
}
```

❌ **Less Efficient:**
```json
{
  "field": "track.name",
  "operator": "regex",
  "value": ".*"  // Matches everything
}
```

### 4. Combine Operations

Use filters with sorting for better results:

```json
{
  "operation": "filter",
  "filters": [...],
  "sort_by": "track.popularity",
  "sort_order": "desc",
  "limit": 10
}
```

### 5. Use Pagination for Large Results

```json
{
  "operation": "select",
  "limit": 50,
  "offset": 0  // First page
}

// Next page:
{
  "operation": "select",
  "limit": 50,
  "offset": 50  // Second page
}
```

### 6. Cache Queries

The query engine automatically caches loaded JSON files. Repeated queries on the same source are fast.

### 7. Use Music-Specific Helpers

Instead of writing complex queries, use helper functions:

```go
// Instead of complex query JSON:
helper := jsonquery.NewMusicQueryHelper(dataDir)

// Use simple helpers:
stats := helper.GetLibraryStats()
tracks := helper.GetTracksByArtist("Queen")
recent := helper.GetRecentlyAddedTracks(20)
```

## Query Result Format

All queries return a `QueryResult` object:

```json
{
  "count": 10,
  "data": [/* array or object */],
  "summary": "Human-readable summary",
  "error": ""  // Empty if no error
}
```

### Fields

- **`count`** (integer): Number of items returned or counted
- **`data`** (any): The actual result data (array, object, or value)
- **`summary`** (string): Human-readable description of the result
- **`error`** (string): Error message if query failed (empty on success)

## Performance Tips

1. **Use filters first** - Reduce dataset before sorting or aggregating
2. **Limit early** - Don't retrieve more data than needed
3. **Cache results** - Store query results if you'll use them repeatedly
4. **Use specific fields** - Extract only the fields you need
5. **Avoid complex regex** - Use simple string operators when possible

## Integration with AI Chat

The query engine is designed to work seamlessly with AI chat through function calling:

```
User: "Find my most popular Queen songs"

AI decides to call: query_music_data
Arguments:
{
  "source": "saved_tracks.json",
  "operation": "filter",
  "filters": [
    {
      "field": "track.artists.0.name",
      "operator": "eq",
      "value": "Queen"
    }
  ],
  "sort_by": "track.popularity",
  "sort_order": "desc",
  "limit": 10
}

AI formats response: "Here are your top 10 Queen songs..."
```

## Advanced Features

### Nested Field Access

Access deeply nested fields:

```json
{
  "field": "track.album.artists.0.name"
}
```

### Array Access

Access array elements by index:

```json
{
  "field": "track.artists.0.name"  // First artist
}
```

### Boolean Operations

All filters in the `filters` array are combined with **AND** logic:

```json
{
  "filters": [
    {"field": "A", "operator": "eq", "value": 1},
    {"field": "B", "operator": "gt", "value": 5}
  ]
  // Returns items where A=1 AND B>5
}
```

## Error Handling

If a query fails, the `error` field will contain a description:

```json
{
  "count": 0,
  "data": null,
  "error": "field 'track.invalid_field' not found"
}
```

Common errors:
- `field not found` - Invalid field path
- `failed to load data` - JSON file not found or invalid
- `invalid operator` - Unknown filter operator
- `type mismatch` - Comparing incompatible types

## Further Reading

- [Tool Calling Documentation](TOOLS.md) - How AI chat uses queries
- [JSON Query Engine Code](../internal/jsonquery/query.go) - Implementation
- [Music Query Helpers](../internal/jsonquery/query.go#L700) - Pre-built queries

## Contributing

To add new query operations:

1. Add operation name to switch statement in `Execute()`
2. Implement operation handler (e.g., `statsOp()`)
3. Add tests in `query_test.go`
4. Update this documentation

See existing operations in `internal/jsonquery/query.go` as examples.