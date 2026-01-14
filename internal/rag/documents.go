package rag

import (
	"fmt"
	"strings"
)

// TrackData represents track information for indexing
type TrackData struct {
	ID       string
	Name     string
	Artists  []string
	Album    string
	Genres   []string
	Duration int // milliseconds
}

// ArtistData represents artist information for indexing
type ArtistData struct {
	ID     string
	Name   string
	Genres []string
}

// PlaylistData represents playlist information for indexing
type PlaylistData struct {
	ID          string
	Name        string
	Description string
	Owner       string
	TrackCount  int
	TrackNames  []string
}

// TrackToDocument converts track data to a searchable document
func TrackToDocument(track TrackData) Document {
	// Create rich text content for embedding
	artists := strings.Join(track.Artists, ", ")
	genres := strings.Join(track.Genres, ", ")

	content := fmt.Sprintf("%s by %s", track.Name, artists)
	if track.Album != "" {
		content += fmt.Sprintf(" from album %s", track.Album)
	}
	if genres != "" {
		content += fmt.Sprintf(". Genres: %s", genres)
	}

	return Document{
		ID:      fmt.Sprintf("track:%s", track.ID),
		Type:    "track",
		Content: content,
		Metadata: map[string]string{
			"id":      track.ID,
			"name":    track.Name,
			"artists": artists,
			"album":   track.Album,
			"genres":  genres,
		},
	}
}

// ArtistToDocument converts artist data to a searchable document
func ArtistToDocument(artist ArtistData) Document {
	genres := strings.Join(artist.Genres, ", ")

	content := artist.Name
	if genres != "" {
		content += fmt.Sprintf(". Genres: %s", genres)
	}

	return Document{
		ID:      fmt.Sprintf("artist:%s", artist.ID),
		Type:    "artist",
		Content: content,
		Metadata: map[string]string{
			"id":     artist.ID,
			"name":   artist.Name,
			"genres": genres,
		},
	}
}

// PlaylistToDocument converts playlist data to a searchable document
func PlaylistToDocument(playlist PlaylistData) Document {
	content := fmt.Sprintf("Playlist: %s", playlist.Name)
	if playlist.Description != "" {
		content += fmt.Sprintf(". %s", playlist.Description)
	}
	if len(playlist.TrackNames) > 0 {
		// Include first few track names for context
		trackSample := playlist.TrackNames
		if len(trackSample) > 10 {
			trackSample = trackSample[:10]
		}
		content += fmt.Sprintf(". Contains tracks like: %s", strings.Join(trackSample, ", "))
	}

	return Document{
		ID:      fmt.Sprintf("playlist:%s", playlist.ID),
		Type:    "playlist",
		Content: content,
		Metadata: map[string]string{
			"id":          playlist.ID,
			"name":        playlist.Name,
			"description": playlist.Description,
			"owner":       playlist.Owner,
			"track_count": fmt.Sprintf("%d", playlist.TrackCount),
		},
	}
}
