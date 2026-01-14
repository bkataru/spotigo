// Package jsonutil provides shared JSON loading and parsing utilities
// for working with Spotify backup data.
package jsonutil

import (
	"encoding/json"
	"os"
)

// LoadJSONFile reads and unmarshals a JSON file into the target interface.
func LoadJSONFile(path string, target interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

// GetString safely extracts a string from a map[string]interface{}.
// Returns empty string if key doesn't exist or value is not a string.
func GetString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

// GetNestedString safely extracts a nested string value from a map.
// Keys are traversed in order, returning empty string if path doesn't exist.
func GetNestedString(m map[string]interface{}, keys ...string) string {
	current := m
	for i, key := range keys {
		if i == len(keys)-1 {
			return GetString(current, key)
		}
		if next, ok := current[key].(map[string]interface{}); ok {
			current = next
		} else {
			return ""
		}
	}
	return ""
}

// GetStringSlice safely extracts a []string from a map[string]interface{}.
// Returns nil if key doesn't exist or value is not the expected type.
func GetStringSlice(m map[string]interface{}, key string) []string {
	if v, ok := m[key].([]interface{}); ok {
		var result []string
		for _, item := range v {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}
	return nil
}

// GetArtistNames extracts artist names from a track's artists array.
// Handles both direct track format and spotify.SavedTrack format.
func GetArtistNames(m map[string]interface{}) []string {
	if artists, ok := m["artists"].([]interface{}); ok {
		var names []string
		for _, a := range artists {
			if artist, ok := a.(map[string]interface{}); ok {
				if name := GetString(artist, "name"); name != "" {
					names = append(names, name)
				}
			}
		}
		return names
	}
	return nil
}

// GetTrackArtists extracts artist names from a track map.
// Handles both spotify.SavedTrack format (with nested "track") and plain track format.
func GetTrackArtists(track map[string]interface{}) []string {
	// Handle spotify.SavedTrack format
	trackData := track
	if t, ok := track["track"].(map[string]interface{}); ok {
		trackData = t
	}
	return GetArtistNames(trackData)
}

// GetTrackAlbum extracts the album name from a track map.
// Handles both spotify.SavedTrack format (with nested "track") and plain track format.
func GetTrackAlbum(track map[string]interface{}) string {
	// Handle spotify.SavedTrack format
	trackData := track
	if t, ok := track["track"].(map[string]interface{}); ok {
		trackData = t
	}

	if album, ok := trackData["album"].(map[string]interface{}); ok {
		if name, ok := album["name"].(string); ok {
			return name
		}
	}
	return ""
}

// GetArtistGenres extracts genres from an artist map.
func GetArtistGenres(artist map[string]interface{}) []string {
	var genres []string
	if g, ok := artist["genres"].([]interface{}); ok {
		for _, genre := range g {
			if s, ok := genre.(string); ok {
				genres = append(genres, s)
			}
		}
	}
	return genres
}

// GetPlaylistName extracts the name from a playlist map.
func GetPlaylistName(playlist map[string]interface{}) string {
	if name, ok := playlist["name"].(string); ok {
		return name
	}
	return ""
}

// GetPlaylistTrackCount returns the number of tracks in a playlist.
func GetPlaylistTrackCount(playlist map[string]interface{}) int {
	if tracks, ok := playlist["tracks"].([]interface{}); ok {
		return len(tracks)
	}
	return 0
}

// GetPlaylistOwner extracts the owner from a playlist map.
func GetPlaylistOwner(playlist map[string]interface{}) string {
	if owner, ok := playlist["owner"].(string); ok {
		return owner
	}
	return ""
}

// Truncate shortens a string to maxLen characters, adding "..." if truncated.
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// Min returns the smaller of two integers.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
