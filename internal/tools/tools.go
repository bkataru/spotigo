// Package tools provides function calling tools for AI chat
package tools

import (
	"encoding/json"
	"fmt"

	"github.com/bkataru/spotigo/internal/jsonquery"
	"github.com/bkataru/spotigo/internal/ollama"
)

// MusicTools provides tools for querying music data
type MusicTools struct {
	queryHelper *jsonquery.MusicQueryHelper
}

// NewMusicTools creates a new music tools instance
func NewMusicTools(dataDir string) *MusicTools {
	return &MusicTools{
		queryHelper: jsonquery.NewMusicQueryHelper(dataDir),
	}
}

// GetToolDefinitions returns all available tool definitions
func (m *MusicTools) GetToolDefinitions() []ollama.Tool {
	return []ollama.Tool{
		{
			Type: "function",
			Function: ollama.FunctionDef{
				Name:        "get_library_stats",
				Description: "Get statistics about the user's music library including total tracks, playlists, and followed artists",
				Parameters: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
					"required":   []string{},
				},
			},
		},
		{
			Type: "function",
			Function: ollama.FunctionDef{
				Name:        "search_tracks",
				Description: "Search for tracks by artist name, song title, or any text",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type":        "string",
							"description": "Search query (artist name, song title, or any text)",
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum number of results to return (default: 10)",
						},
					},
					"required": []string{"query"},
				},
			},
		},
		{
			Type: "function",
			Function: ollama.FunctionDef{
				Name:        "get_tracks_by_artist",
				Description: "Get all tracks by a specific artist",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"artist_name": map[string]interface{}{
							"type":        "string",
							"description": "Name of the artist",
						},
					},
					"required": []string{"artist_name"},
				},
			},
		},
		{
			Type: "function",
			Function: ollama.FunctionDef{
				Name:        "get_recently_added_tracks",
				Description: "Get the most recently added tracks to the library",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Number of tracks to return (default: 10)",
						},
					},
					"required": []string{},
				},
			},
		},
		{
			Type: "function",
			Function: ollama.FunctionDef{
				Name:        "get_all_artists",
				Description: "Get all unique artists in the library",
				Parameters: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
					"required":   []string{},
				},
			},
		},
		{
			Type: "function",
			Function: ollama.FunctionDef{
				Name:        "get_playlist_by_name",
				Description: "Find a playlist by name",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"playlist_name": map[string]interface{}{
							"type":        "string",
							"description": "Name of the playlist to find",
						},
					},
					"required": []string{"playlist_name"},
				},
			},
		},
		{
			Type: "function",
			Function: ollama.FunctionDef{
				Name:        "query_music_data",
				Description: "Execute a custom JSON query on music data. Use this for complex queries like filtering, sorting, aggregating data.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"source": map[string]interface{}{
							"type":        "string",
							"description": "Data source: 'saved_tracks.json', 'playlists.json', or 'followed_artists.json'",
						},
						"operation": map[string]interface{}{
							"type":        "string",
							"description": "Operation: 'select', 'count', 'filter', 'search', 'sort', 'distinct', 'aggregate', 'stats'",
						},
						"filters": map[string]interface{}{
							"type":        "array",
							"description": "Filter conditions",
							"items": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"field": map[string]interface{}{
										"type":        "string",
										"description": "Field to filter on (supports dot notation)",
									},
									"operator": map[string]interface{}{
										"type":        "string",
										"description": "Operator: 'eq', 'ne', 'gt', 'gte', 'lt', 'lte', 'contains', 'regex'",
									},
									"value": map[string]interface{}{
										"description": "Value to compare against",
									},
								},
							},
						},
						"sort_by": map[string]interface{}{
							"type":        "string",
							"description": "Field to sort by (dot notation supported)",
						},
						"sort_order": map[string]interface{}{
							"type":        "string",
							"description": "Sort order: 'asc' or 'desc'",
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Limit number of results",
						},
						"field": map[string]interface{}{
							"type":        "string",
							"description": "Specific field to extract",
						},
					},
					"required": []string{"source", "operation"},
				},
			},
		},
	}
}

// ExecuteToolCall executes a tool call and returns the result
func (m *MusicTools) ExecuteToolCall(toolCall ollama.ToolCall) (string, error) {
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	switch toolCall.Function.Name {
	case "get_library_stats":
		return m.executeGetLibraryStats()
	case "search_tracks":
		return m.executeSearchTracks(args)
	case "get_tracks_by_artist":
		return m.executeGetTracksByArtist(args)
	case "get_recently_added_tracks":
		return m.executeGetRecentlyAddedTracks(args)
	case "get_all_artists":
		return m.executeGetAllArtists()
	case "get_playlist_by_name":
		return m.executeGetPlaylistByName(args)
	case "query_music_data":
		return m.executeQueryMusicData(args)
	default:
		return "", fmt.Errorf("unknown tool: %s", toolCall.Function.Name)
	}
}

func (m *MusicTools) executeGetLibraryStats() (string, error) {
	result := m.queryHelper.GetLibraryStats()
	if result.Error != "" {
		return "", fmt.Errorf("query error: %s", result.Error)
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (m *MusicTools) executeSearchTracks(args map[string]interface{}) (string, error) {
	query, ok := args["query"].(string)
	if !ok {
		return "", fmt.Errorf("query parameter required")
	}

	limit := 10
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}

	result := m.queryHelper.SearchAllData(query, limit)
	if result.Error != "" {
		return "", fmt.Errorf("query error: %s", result.Error)
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (m *MusicTools) executeGetTracksByArtist(args map[string]interface{}) (string, error) {
	artistName, ok := args["artist_name"].(string)
	if !ok {
		return "", fmt.Errorf("artist_name parameter required")
	}

	result := m.queryHelper.GetTracksByArtist(artistName)
	if result.Error != "" {
		return "", fmt.Errorf("query error: %s", result.Error)
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (m *MusicTools) executeGetRecentlyAddedTracks(args map[string]interface{}) (string, error) {
	limit := 10
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}

	result := m.queryHelper.GetRecentlyAddedTracks(limit)
	if result.Error != "" {
		return "", fmt.Errorf("query error: %s", result.Error)
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (m *MusicTools) executeGetAllArtists() (string, error) {
	result := m.queryHelper.GetAllArtists()
	if result.Error != "" {
		return "", fmt.Errorf("query error: %s", result.Error)
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (m *MusicTools) executeGetPlaylistByName(args map[string]interface{}) (string, error) {
	playlistName, ok := args["playlist_name"].(string)
	if !ok {
		return "", fmt.Errorf("playlist_name parameter required")
	}

	result := m.queryHelper.GetPlaylistByName(playlistName)
	if result.Error != "" {
		return "", fmt.Errorf("query error: %s", result.Error)
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (m *MusicTools) executeQueryMusicData(args map[string]interface{}) (string, error) {
	// Build query from args
	query := jsonquery.Query{}

	if source, ok := args["source"].(string); ok {
		query.Source = source
	} else {
		return "", fmt.Errorf("source parameter required")
	}

	if operation, ok := args["operation"].(string); ok {
		query.Operation = operation
	} else {
		return "", fmt.Errorf("operation parameter required")
	}

	if field, ok := args["field"].(string); ok {
		query.Field = field
	}

	if sortBy, ok := args["sort_by"].(string); ok {
		query.SortBy = sortBy
	}

	if sortOrder, ok := args["sort_order"].(string); ok {
		query.SortOrder = sortOrder
	}

	if limit, ok := args["limit"].(float64); ok {
		query.Limit = int(limit)
	}

	if filters, ok := args["filters"].([]interface{}); ok {
		for _, f := range filters {
			filterMap, ok := f.(map[string]interface{})
			if !ok {
				continue
			}

			filter := jsonquery.Filter{}
			if field, ok := filterMap["field"].(string); ok {
				filter.Field = field
			}
			if operator, ok := filterMap["operator"].(string); ok {
				filter.Operator = operator
			}
			if value, ok := filterMap["value"]; ok {
				filter.Value = value
			}

			query.Filters = append(query.Filters, filter)
		}
	}

	// Execute query
	result := m.queryHelper.Engine.Execute(query)
	if result.Error != "" {
		return "", fmt.Errorf("query error: %s", result.Error)
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
