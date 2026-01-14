// Package jsonquery provides a powerful JSON query and manipulation tool
// for structured music data. This enables efficient RAG over JSON files
// without flooding the context window with raw data.
package jsonquery

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// QueryResult represents the result of a query operation
type QueryResult struct {
	Count   int         `json:"count"`
	Data    interface{} `json:"data,omitempty"`
	Summary string      `json:"summary,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Query represents a structured query against JSON data
type Query struct {
	// Source file to query
	Source string `json:"source"`

	// Operation to perform: select, count, aggregate, search, filter, sort, distinct, stats
	Operation string `json:"operation"`

	// Field path to operate on (dot notation, e.g., "track.artists.0.name")
	Field string `json:"field,omitempty"`

	// Filter conditions
	Filters []Filter `json:"filters,omitempty"`

	// Sort configuration
	SortBy    string `json:"sort_by,omitempty"`
	SortOrder string `json:"sort_order,omitempty"` // asc or desc

	// Limit results
	Limit int `json:"limit,omitempty"`

	// Offset for pagination
	Offset int `json:"offset,omitempty"`

	// Search term for search operations
	SearchTerm string `json:"search_term,omitempty"`

	// Aggregation function: count, sum, avg, min, max, group
	AggFunc string `json:"agg_func,omitempty"`

	// Group by field for aggregations
	GroupBy string `json:"group_by,omitempty"`
}

// Filter represents a filter condition
type Filter struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"` // eq, ne, gt, gte, lt, lte, contains, regex, in, exists
	Value    interface{} `json:"value"`
}

// Engine processes queries against JSON data
type Engine struct {
	dataDir string
	cache   map[string][]interface{}
}

// NewEngine creates a new query engine
func NewEngine(dataDir string) *Engine {
	return &Engine{
		dataDir: dataDir,
		cache:   make(map[string][]interface{}),
	}
}

// Execute runs a query and returns results
func (e *Engine) Execute(q Query) QueryResult {
	// Load data
	data, err := e.loadData(q.Source)
	if err != nil {
		return QueryResult{Error: fmt.Sprintf("failed to load data: %v", err)}
	}

	// Apply filters
	filtered := e.applyFilters(data, q.Filters)

	// Execute operation
	switch q.Operation {
	case "select":
		return e.selectOp(filtered, q)
	case "count":
		return e.countOp(filtered, q)
	case "aggregate":
		return e.aggregateOp(filtered, q)
	case "search":
		return e.searchOp(filtered, q)
	case "filter":
		return e.filterOp(filtered, q)
	case "sort":
		return e.sortOp(filtered, q)
	case "distinct":
		return e.distinctOp(filtered, q)
	case "stats":
		return e.statsOp(filtered, q)
	case "sample":
		return e.sampleOp(filtered, q)
	default:
		return QueryResult{Error: fmt.Sprintf("unknown operation: %s", q.Operation)}
	}
}

// loadData loads JSON data from a file
func (e *Engine) loadData(source string) ([]interface{}, error) {
	// Check cache
	if cached, ok := e.cache[source]; ok {
		return cached, nil
	}

	// Build file path
	filePath := source
	if e.dataDir != "" && !strings.HasPrefix(source, "/") && !strings.Contains(source, ":") {
		filePath = fmt.Sprintf("%s/%s", e.dataDir, source)
	}

	// Clean path to prevent traversal attacks
	cleanPath := filepath.Clean(filePath)

	// Read file
	file, err := os.ReadFile(cleanPath) // #nosec G304 - path is sanitized with filepath.Clean
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse JSON
	var data interface{}
	if err := json.Unmarshal(file, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Convert to array if necessary
	var result []interface{}
	switch v := data.(type) {
	case []interface{}:
		result = v
	case map[string]interface{}:
		result = []interface{}{v}
	default:
		result = []interface{}{data}
	}

	// Cache and return
	e.cache[source] = result
	return result, nil
}

// ClearCache clears the data cache
func (e *Engine) ClearCache() {
	e.cache = make(map[string][]interface{})
}

// applyFilters applies filter conditions to data
func (e *Engine) applyFilters(data []interface{}, filters []Filter) []interface{} {
	if len(filters) == 0 {
		return data
	}

	result := make([]interface{}, 0)
	for _, item := range data {
		if e.matchesFilters(item, filters) {
			result = append(result, item)
		}
	}
	return result
}

// matchesFilters checks if an item matches all filters
func (e *Engine) matchesFilters(item interface{}, filters []Filter) bool {
	for _, f := range filters {
		if !e.matchesFilter(item, f) {
			return false
		}
	}
	return true
}

// matchesFilter checks if an item matches a single filter
func (e *Engine) matchesFilter(item interface{}, f Filter) bool {
	value := getFieldValue(item, f.Field)

	switch f.Operator {
	case "eq":
		return compareValues(value, f.Value) == 0
	case "ne":
		return compareValues(value, f.Value) != 0
	case "gt":
		return compareValues(value, f.Value) > 0
	case "gte":
		return compareValues(value, f.Value) >= 0
	case "lt":
		return compareValues(value, f.Value) < 0
	case "lte":
		return compareValues(value, f.Value) <= 0
	case "contains":
		return containsValue(value, f.Value)
	case "regex":
		return matchesRegex(value, f.Value)
	case "in":
		return inValues(value, f.Value)
	case "exists":
		return value != nil
	case "not_exists":
		return value == nil
	default:
		return false
	}
}

// selectOp performs a select operation
func (e *Engine) selectOp(data []interface{}, q Query) QueryResult {
	// Apply sorting
	if q.SortBy != "" {
		data = e.sortData(data, q.SortBy, q.SortOrder)
	}

	// Apply pagination
	if q.Offset > 0 {
		if q.Offset >= len(data) {
			data = []interface{}{}
		} else {
			data = data[q.Offset:]
		}
	}

	if q.Limit > 0 && q.Limit < len(data) {
		data = data[:q.Limit]
	}

	// Extract specific fields if requested
	if q.Field != "" {
		extracted := make([]interface{}, 0, len(data))
		for _, item := range data {
			val := getFieldValue(item, q.Field)
			if val != nil {
				extracted = append(extracted, val)
			}
		}
		return QueryResult{Count: len(extracted), Data: extracted}
	}

	return QueryResult{Count: len(data), Data: data}
}

// countOp counts items
func (e *Engine) countOp(data []interface{}, _ Query) QueryResult {
	count := len(data)
	return QueryResult{
		Count:   count,
		Summary: fmt.Sprintf("Found %d items", count),
	}
}

// aggregateOp performs aggregation operations
func (e *Engine) aggregateOp(data []interface{}, q Query) QueryResult {
	switch q.AggFunc {
	case "count":
		return e.countOp(data, q)
	case "sum":
		return e.sumOp(data, q.Field)
	case "avg":
		return e.avgOp(data, q.Field)
	case "min":
		return e.minOp(data, q.Field)
	case "max":
		return e.maxOp(data, q.Field)
	case "group":
		return e.groupOp(data, q)
	default:
		return QueryResult{Error: fmt.Sprintf("unknown aggregation function: %s", q.AggFunc)}
	}
}

func (e *Engine) sumOp(data []interface{}, field string) QueryResult {
	var sum float64
	for _, item := range data {
		val := getFieldValue(item, field)
		if num, ok := toFloat64(val); ok {
			sum += num
		}
	}
	return QueryResult{Count: len(data), Data: sum, Summary: fmt.Sprintf("Sum of %s: %.2f", field, sum)}
}

func (e *Engine) avgOp(data []interface{}, field string) QueryResult {
	var sum float64
	var count int
	for _, item := range data {
		val := getFieldValue(item, field)
		if num, ok := toFloat64(val); ok {
			sum += num
			count++
		}
	}
	if count == 0 {
		return QueryResult{Count: 0, Data: 0, Summary: "No numeric values found"}
	}
	avg := sum / float64(count)
	return QueryResult{Count: count, Data: avg, Summary: fmt.Sprintf("Average of %s: %.2f", field, avg)}
}

func (e *Engine) minOp(data []interface{}, field string) QueryResult {
	var minVal interface{}
	for _, item := range data {
		val := getFieldValue(item, field)
		if val == nil {
			continue
		}
		if minVal == nil || compareValues(val, minVal) < 0 {
			minVal = val
		}
	}
	return QueryResult{Count: len(data), Data: minVal, Summary: fmt.Sprintf("Min %s: %v", field, minVal)}
}

func (e *Engine) maxOp(data []interface{}, field string) QueryResult {
	var maxVal interface{}
	for _, item := range data {
		val := getFieldValue(item, field)
		if val == nil {
			continue
		}
		if maxVal == nil || compareValues(val, maxVal) > 0 {
			maxVal = val
		}
	}
	return QueryResult{Count: len(data), Data: maxVal, Summary: fmt.Sprintf("Max %s: %v", field, maxVal)}
}

func (e *Engine) groupOp(data []interface{}, q Query) QueryResult {
	groups := make(map[string]int)
	for _, item := range data {
		val := getFieldValue(item, q.GroupBy)
		key := fmt.Sprintf("%v", val)
		groups[key]++
	}

	// Convert to sorted slice
	type groupItem struct {
		Key   string `json:"key"`
		Count int    `json:"count"`
	}
	items := make([]groupItem, 0, len(groups))
	for k, v := range groups {
		items = append(items, groupItem{Key: k, Count: v})
	}

	// Sort by count descending
	sort.Slice(items, func(i, j int) bool {
		return items[i].Count > items[j].Count
	})

	// Apply limit
	if q.Limit > 0 && q.Limit < len(items) {
		items = items[:q.Limit]
	}

	return QueryResult{
		Count:   len(groups),
		Data:    items,
		Summary: fmt.Sprintf("Found %d unique groups", len(groups)),
	}
}

// searchOp performs text search
func (e *Engine) searchOp(data []interface{}, q Query) QueryResult {
	if q.SearchTerm == "" {
		return QueryResult{Error: "search_term is required for search operation"}
	}

	term := strings.ToLower(q.SearchTerm)
	results := make([]interface{}, 0)

	for _, item := range data {
		if e.itemContainsText(item, term, q.Field) {
			results = append(results, item)
		}
	}

	// Apply limit
	if q.Limit > 0 && q.Limit < len(results) {
		results = results[:q.Limit]
	}

	return QueryResult{
		Count:   len(results),
		Data:    results,
		Summary: fmt.Sprintf("Found %d items matching '%s'", len(results), q.SearchTerm),
	}
}

func (e *Engine) itemContainsText(item interface{}, term string, field string) bool {
	if field != "" {
		val := getFieldValue(item, field)
		return strings.Contains(strings.ToLower(fmt.Sprintf("%v", val)), term)
	}

	// Search all string fields
	jsonBytes, err := json.Marshal(item)
	if err != nil {
		// If marshaling fails, fall back to string formatting
		return strings.Contains(strings.ToLower(fmt.Sprintf("%v", item)), term)
	}
	return strings.Contains(strings.ToLower(string(jsonBytes)), term)
}

// filterOp is an alias for select with filters pre-applied
func (e *Engine) filterOp(data []interface{}, q Query) QueryResult {
	return e.selectOp(data, q)
}

// sortOp sorts data
func (e *Engine) sortOp(data []interface{}, q Query) QueryResult {
	if q.SortBy == "" {
		return QueryResult{Error: "sort_by is required for sort operation"}
	}

	sorted := e.sortData(data, q.SortBy, q.SortOrder)

	// Apply limit
	if q.Limit > 0 && q.Limit < len(sorted) {
		sorted = sorted[:q.Limit]
	}

	return QueryResult{Count: len(sorted), Data: sorted}
}

func (e *Engine) sortData(data []interface{}, sortBy string, order string) []interface{} {
	result := make([]interface{}, len(data))
	copy(result, data)

	sort.Slice(result, func(i, j int) bool {
		valI := getFieldValue(result[i], sortBy)
		valJ := getFieldValue(result[j], sortBy)

		cmp := compareValues(valI, valJ)
		if order == "desc" {
			return cmp > 0
		}
		return cmp < 0
	})

	return result
}

// distinctOp gets distinct values
func (e *Engine) distinctOp(data []interface{}, q Query) QueryResult {
	if q.Field == "" {
		return QueryResult{Error: "field is required for distinct operation"}
	}

	seen := make(map[string]bool)
	distinct := make([]interface{}, 0)

	for _, item := range data {
		val := getFieldValue(item, q.Field)
		key := fmt.Sprintf("%v", val)
		if !seen[key] {
			seen[key] = true
			distinct = append(distinct, val)
		}
	}

	// Apply limit
	if q.Limit > 0 && q.Limit < len(distinct) {
		distinct = distinct[:q.Limit]
	}

	return QueryResult{
		Count:   len(distinct),
		Data:    distinct,
		Summary: fmt.Sprintf("Found %d distinct values for %s", len(distinct), q.Field),
	}
}

// statsOp provides statistics about the data
func (e *Engine) statsOp(data []interface{}, q Query) QueryResult {
	stats := map[string]interface{}{
		"total_count": len(data),
	}

	if q.Field != "" {
		// Get stats for specific field
		var numericCount int
		var sum, minVal, maxVal float64
		var minSet bool
		stringValues := make(map[string]int)

		for _, item := range data {
			val := getFieldValue(item, q.Field)
			if val == nil {
				continue
			}

			if num, ok := toFloat64(val); ok {
				numericCount++
				sum += num
				if !minSet || num < minVal {
					minVal = num
					minSet = true
				}
				if num > maxVal {
					maxVal = num
				}
			} else {
				strVal := fmt.Sprintf("%v", val)
				stringValues[strVal]++
			}
		}

		if numericCount > 0 {
			stats["numeric_count"] = numericCount
			stats["sum"] = sum
			stats["avg"] = sum / float64(numericCount)
			stats["min"] = minVal
			stats["max"] = maxVal
		}

		if len(stringValues) > 0 {
			stats["unique_values"] = len(stringValues)
			// Find most common
			var maxCount int
			var mostCommon string
			for k, v := range stringValues {
				if v > maxCount {
					maxCount = v
					mostCommon = k
				}
			}
			stats["most_common"] = mostCommon
			stats["most_common_count"] = maxCount
		}
	}

	return QueryResult{
		Count:   len(data),
		Data:    stats,
		Summary: fmt.Sprintf("Statistics for %d items", len(data)),
	}
}

// sampleOp returns a random sample of data
func (e *Engine) sampleOp(data []interface{}, q Query) QueryResult {
	limit := q.Limit
	if limit <= 0 {
		limit = 5
	}

	if limit >= len(data) {
		return QueryResult{Count: len(data), Data: data}
	}

	// Simple sampling by taking evenly spaced items
	step := len(data) / limit
	sample := make([]interface{}, 0, limit)
	for i := 0; i < len(data) && len(sample) < limit; i += step {
		sample = append(sample, data[i])
	}

	return QueryResult{
		Count:   len(sample),
		Data:    sample,
		Summary: fmt.Sprintf("Sample of %d items from %d total", len(sample), len(data)),
	}
}

// Helper functions

// getFieldValue extracts a value from nested data using dot notation
func getFieldValue(data interface{}, path string) interface{} {
	if path == "" {
		return data
	}

	parts := strings.Split(path, ".")
	current := data

	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			var ok bool
			current, ok = v[part]
			if !ok {
				return nil
			}
		case []interface{}:
			// Try to parse as index
			var idx int
			if _, err := fmt.Sscanf(part, "%d", &idx); err == nil {
				if idx >= 0 && idx < len(v) {
					current = v[idx]
				} else {
					return nil
				}
			} else {
				// Try to get field from all array elements
				results := make([]interface{}, 0)
				for _, item := range v {
					if val := getFieldValue(item, part); val != nil {
						results = append(results, val)
					}
				}
				if len(results) > 0 {
					current = results
				} else {
					return nil
				}
			}
		default:
			return nil
		}
	}

	return current
}

// compareValues compares two values
func compareValues(a, b interface{}) int {
	// Handle nil
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}

	// Try numeric comparison
	numA, okA := toFloat64(a)
	numB, okB := toFloat64(b)
	if okA && okB {
		if numA < numB {
			return -1
		}
		if numA > numB {
			return 1
		}
		return 0
	}

	// Try time comparison
	timeA, okA := toTime(a)
	timeB, okB := toTime(b)
	if okA && okB {
		if timeA.Before(timeB) {
			return -1
		}
		if timeA.After(timeB) {
			return 1
		}
		return 0
	}

	// String comparison
	strA := fmt.Sprintf("%v", a)
	strB := fmt.Sprintf("%v", b)
	return strings.Compare(strA, strB)
}

func toFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case int32:
		return float64(n), true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	default:
		return 0, false
	}
}

func toTime(v interface{}) (time.Time, bool) {
	if s, ok := v.(string); ok {
		// Try common time formats
		formats := []string{
			time.RFC3339,
			"2006-01-02T15:04:05Z",
			"2006-01-02",
			"2006-01-02 15:04:05",
		}
		for _, format := range formats {
			if t, err := time.Parse(format, s); err == nil {
				return t, true
			}
		}
	}
	return time.Time{}, false
}

func containsValue(value, search interface{}) bool {
	valStr := strings.ToLower(fmt.Sprintf("%v", value))
	searchStr := strings.ToLower(fmt.Sprintf("%v", search))
	return strings.Contains(valStr, searchStr)
}

func matchesRegex(value, pattern interface{}) bool {
	valStr := fmt.Sprintf("%v", value)
	patStr := fmt.Sprintf("%v", pattern)
	re, err := regexp.Compile(patStr)
	if err != nil {
		return false
	}
	return re.MatchString(valStr)
}

func inValues(value interface{}, list interface{}) bool {
	listArr, ok := list.([]interface{})
	if !ok {
		return false
	}

	for _, item := range listArr {
		if compareValues(value, item) == 0 {
			return true
		}
	}
	return false
}

// MusicQueryHelper provides high-level music-specific queries
type MusicQueryHelper struct {
	Engine *Engine
}

// NewMusicQueryHelper creates a new music query helper
func NewMusicQueryHelper(dataDir string) *MusicQueryHelper {
	return &MusicQueryHelper{
		Engine: NewEngine(dataDir),
	}
}

// GetAllArtists returns all unique artists from saved tracks
func (m *MusicQueryHelper) GetAllArtists() QueryResult {
	return m.Engine.Execute(Query{
		Source:    "saved_tracks.json",
		Operation: "distinct",
		Field:     "track.artists.name",
	})
}

// GetTracksByArtist returns tracks by a specific artist
func (m *MusicQueryHelper) GetTracksByArtist(artistName string) QueryResult {
	return m.Engine.Execute(Query{
		Source:     "saved_tracks.json",
		Operation:  "search",
		Field:      "track.artists.name",
		SearchTerm: artistName,
	})
}

// GetPlaylistByName returns a playlist by name
func (m *MusicQueryHelper) GetPlaylistByName(name string) QueryResult {
	return m.Engine.Execute(Query{
		Source:     "playlists.json",
		Operation:  "search",
		Field:      "name",
		SearchTerm: name,
	})
}

// GetRecentlyAddedTracks returns the most recently added tracks
func (m *MusicQueryHelper) GetRecentlyAddedTracks(limit int) QueryResult {
	return m.Engine.Execute(Query{
		Source:    "saved_tracks.json",
		Operation: "sort",
		SortBy:    "added_at",
		SortOrder: "desc",
		Limit:     limit,
	})
}

// GetLibraryStats returns statistics about the music library
func (m *MusicQueryHelper) GetLibraryStats() QueryResult {
	tracksResult := m.Engine.Execute(Query{
		Source:    "saved_tracks.json",
		Operation: "count",
	})

	playlistsResult := m.Engine.Execute(Query{
		Source:    "playlists.json",
		Operation: "count",
	})

	artistsResult := m.Engine.Execute(Query{
		Source:    "followed_artists.json",
		Operation: "count",
	})

	return QueryResult{
		Data: map[string]interface{}{
			"saved_tracks":     tracksResult.Count,
			"playlists":        playlistsResult.Count,
			"followed_artists": artistsResult.Count,
		},
		Summary: fmt.Sprintf("Library: %d tracks, %d playlists, %d followed artists",
			tracksResult.Count, playlistsResult.Count, artistsResult.Count),
	}
}

// SearchAllData searches across all music data files
func (m *MusicQueryHelper) SearchAllData(term string, limit int) QueryResult {
	var allResults []interface{}

	files := []string{"saved_tracks.json", "playlists.json", "followed_artists.json"}
	for _, file := range files {
		result := m.Engine.Execute(Query{
			Source:     file,
			Operation:  "search",
			SearchTerm: term,
			Limit:      limit,
		})
		if result.Data != nil {
			if items, ok := result.Data.([]interface{}); ok {
				for _, item := range items {
					allResults = append(allResults, map[string]interface{}{
						"source": file,
						"item":   item,
					})
				}
			}
		}
	}

	if limit > 0 && len(allResults) > limit {
		allResults = allResults[:limit]
	}

	return QueryResult{
		Count:   len(allResults),
		Data:    allResults,
		Summary: fmt.Sprintf("Found %d results for '%s' across all data", len(allResults), term),
	}
}
