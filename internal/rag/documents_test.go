package rag

import (
	"strings"
	"testing"
)

func TestTrackToDocument(t *testing.T) {
	track := TrackData{
		ID:      "track123",
		Name:    "Bohemian Rhapsody",
		Artists: []string{"Queen"},
		Album:   "A Night at the Opera",
		Genres:  []string{"rock", "progressive rock"},
	}

	doc := TrackToDocument(track)

	if doc.ID != "track:track123" {
		t.Errorf("expected ID 'track:track123', got '%s'", doc.ID)
	}

	if doc.Type != "track" {
		t.Errorf("expected Type 'track', got '%s'", doc.Type)
	}

	if !strings.Contains(doc.Content, "Bohemian Rhapsody") {
		t.Error("Content should contain track name")
	}

	if !strings.Contains(doc.Content, "Queen") {
		t.Error("Content should contain artist name")
	}

	if !strings.Contains(doc.Content, "A Night at the Opera") {
		t.Error("Content should contain album name")
	}

	if !strings.Contains(doc.Content, "rock") {
		t.Error("Content should contain genre")
	}

	if doc.Metadata["name"] != "Bohemian Rhapsody" {
		t.Errorf("expected metadata name 'Bohemian Rhapsody', got '%s'", doc.Metadata["name"])
	}

	if doc.Metadata["artists"] != "Queen" {
		t.Errorf("expected metadata artists 'Queen', got '%s'", doc.Metadata["artists"])
	}
}

func TestTrackToDocument_MultipleArtists(t *testing.T) {
	track := TrackData{
		ID:      "track456",
		Name:    "Collaboration Song",
		Artists: []string{"Artist A", "Artist B", "Artist C"},
		Album:   "Collab Album",
	}

	doc := TrackToDocument(track)

	expectedArtists := "Artist A, Artist B, Artist C"
	if doc.Metadata["artists"] != expectedArtists {
		t.Errorf("expected artists '%s', got '%s'", expectedArtists, doc.Metadata["artists"])
	}

	if !strings.Contains(doc.Content, "Artist A, Artist B, Artist C") {
		t.Error("Content should contain all artist names")
	}
}

func TestTrackToDocument_NoAlbumOrGenres(t *testing.T) {
	track := TrackData{
		ID:      "track789",
		Name:    "Single Track",
		Artists: []string{"Solo Artist"},
	}

	doc := TrackToDocument(track)

	// Content should still be valid without album/genres
	if !strings.Contains(doc.Content, "Single Track by Solo Artist") {
		t.Error("Content should contain basic track info")
	}

	// Should not contain "from album" if no album
	if track.Album == "" && strings.Contains(doc.Content, "from album") {
		t.Error("Content should not mention album when empty")
	}
}

func TestArtistToDocument(t *testing.T) {
	artist := ArtistData{
		ID:     "artist123",
		Name:   "The Beatles",
		Genres: []string{"rock", "pop", "british invasion"},
	}

	doc := ArtistToDocument(artist)

	if doc.ID != "artist:artist123" {
		t.Errorf("expected ID 'artist:artist123', got '%s'", doc.ID)
	}

	if doc.Type != "artist" {
		t.Errorf("expected Type 'artist', got '%s'", doc.Type)
	}

	if !strings.Contains(doc.Content, "The Beatles") {
		t.Error("Content should contain artist name")
	}

	if !strings.Contains(doc.Content, "rock") {
		t.Error("Content should contain genre")
	}

	if doc.Metadata["name"] != "The Beatles" {
		t.Errorf("expected metadata name 'The Beatles', got '%s'", doc.Metadata["name"])
	}

	expectedGenres := "rock, pop, british invasion"
	if doc.Metadata["genres"] != expectedGenres {
		t.Errorf("expected genres '%s', got '%s'", expectedGenres, doc.Metadata["genres"])
	}
}

func TestArtistToDocument_NoGenres(t *testing.T) {
	artist := ArtistData{
		ID:   "artist456",
		Name: "New Artist",
	}

	doc := ArtistToDocument(artist)

	if doc.Content != "New Artist" {
		t.Errorf("expected content 'New Artist', got '%s'", doc.Content)
	}

	if doc.Metadata["genres"] != "" {
		t.Error("genres should be empty when not provided")
	}
}

func TestPlaylistToDocument(t *testing.T) {
	playlist := PlaylistData{
		ID:          "playlist123",
		Name:        "Workout Mix",
		Description: "High energy songs for the gym",
		Owner:       "testuser",
		TrackCount:  50,
		TrackNames:  []string{"Eye of the Tiger", "Lose Yourself", "Stronger"},
	}

	doc := PlaylistToDocument(playlist)

	if doc.ID != "playlist:playlist123" {
		t.Errorf("expected ID 'playlist:playlist123', got '%s'", doc.ID)
	}

	if doc.Type != "playlist" {
		t.Errorf("expected Type 'playlist', got '%s'", doc.Type)
	}

	if !strings.Contains(doc.Content, "Workout Mix") {
		t.Error("Content should contain playlist name")
	}

	if !strings.Contains(doc.Content, "High energy songs") {
		t.Error("Content should contain description")
	}

	if !strings.Contains(doc.Content, "Eye of the Tiger") {
		t.Error("Content should contain track names")
	}

	if doc.Metadata["name"] != "Workout Mix" {
		t.Errorf("expected metadata name 'Workout Mix', got '%s'", doc.Metadata["name"])
	}

	if doc.Metadata["track_count"] != "50" {
		t.Errorf("expected track_count '50', got '%s'", doc.Metadata["track_count"])
	}
}

func TestPlaylistToDocument_ManyTracks(t *testing.T) {
	// Create playlist with more than 10 tracks
	trackNames := make([]string, 20)
	for i := 0; i < 20; i++ {
		trackNames[i] = "Track " + string(rune('A'+i))
	}

	playlist := PlaylistData{
		ID:         "playlist456",
		Name:       "Big Playlist",
		TrackCount: 20,
		TrackNames: trackNames,
	}

	doc := PlaylistToDocument(playlist)

	// Should only include first 10 tracks
	if strings.Contains(doc.Content, "Track K") {
		t.Error("Content should only include first 10 track names")
	}

	if !strings.Contains(doc.Content, "Track A") {
		t.Error("Content should include first track name")
	}

	if !strings.Contains(doc.Content, "Track J") {
		t.Error("Content should include 10th track name")
	}
}

func TestPlaylistToDocument_NoDescription(t *testing.T) {
	playlist := PlaylistData{
		ID:         "playlist789",
		Name:       "My Playlist",
		TrackCount: 10,
	}

	doc := PlaylistToDocument(playlist)

	// Should have valid content without description
	if !strings.HasPrefix(doc.Content, "Playlist: My Playlist") {
		t.Errorf("unexpected content format: %s", doc.Content)
	}

	if doc.Metadata["description"] != "" {
		t.Error("description should be empty when not provided")
	}
}

// TestDocumentMetadataIntegrity ensures metadata is properly set
func TestDocumentMetadataIntegrity(t *testing.T) {
	tests := []struct {
		name     string
		doc      Document
		expected map[string]string
	}{
		{
			name: "track metadata",
			doc: TrackToDocument(TrackData{
				ID:      "t1",
				Name:    "Test",
				Artists: []string{"Artist"},
				Album:   "Album",
				Genres:  []string{"genre"},
			}),
			expected: map[string]string{
				"id":      "t1",
				"name":    "Test",
				"artists": "Artist",
				"album":   "Album",
				"genres":  "genre",
			},
		},
		{
			name: "artist metadata",
			doc: ArtistToDocument(ArtistData{
				ID:     "a1",
				Name:   "Artist Name",
				Genres: []string{"rock", "pop"},
			}),
			expected: map[string]string{
				"id":     "a1",
				"name":   "Artist Name",
				"genres": "rock, pop",
			},
		},
		{
			name: "playlist metadata",
			doc: PlaylistToDocument(PlaylistData{
				ID:          "p1",
				Name:        "Playlist Name",
				Description: "A description",
				Owner:       "owner",
				TrackCount:  10,
			}),
			expected: map[string]string{
				"id":          "p1",
				"name":        "Playlist Name",
				"description": "A description",
				"owner":       "owner",
				"track_count": "10",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, expectedValue := range tt.expected {
				if tt.doc.Metadata[key] != expectedValue {
					t.Errorf("metadata[%s]: expected '%s', got '%s'",
						key, expectedValue, tt.doc.Metadata[key])
				}
			}
		})
	}
}
