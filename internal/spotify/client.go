// Package spotify handles Spotify API interactions
package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"

	"github.com/bkataru/spotigo/internal/crypto"
)

// Scopes required for full library access
var Scopes = []string{
	spotifyauth.ScopeUserLibraryRead,
	spotifyauth.ScopePlaylistReadPrivate,
	spotifyauth.ScopePlaylistReadCollaborative,
	spotifyauth.ScopeUserTopRead,
	spotifyauth.ScopeUserReadRecentlyPlayed,
	spotifyauth.ScopeUserFollowRead,
}

// Client wraps the Spotify client with additional functionality
type Client struct {
	client *spotify.Client
	auth   *spotifyauth.Authenticator
	token  *oauth2.Token
}

// Config holds Spotify client configuration
type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	TokenFile    string
}

// NewClient creates a new Spotify client
func NewClient(cfg Config) (*Client, error) {
	auth := spotifyauth.New(
		spotifyauth.WithClientID(cfg.ClientID),
		spotifyauth.WithClientSecret(cfg.ClientSecret),
		spotifyauth.WithRedirectURL(cfg.RedirectURI),
		spotifyauth.WithScopes(Scopes...),
	)

	c := &Client{
		auth: auth,
	}

	// Try to load existing token
	if cfg.TokenFile != "" {
		if token, err := c.loadToken(cfg.TokenFile); err == nil {
			c.token = token
			httpClient := auth.Client(context.Background(), token)
			c.client = spotify.New(httpClient)
		}
	}

	return c, nil
}

// IsAuthenticated returns true if the client has a valid token
func (c *Client) IsAuthenticated() bool {
	return c.token != nil && c.client != nil
}

// GetAuthURL returns the URL for OAuth authentication
func (c *Client) GetAuthURL(state string) string {
	return c.auth.AuthURL(state)
}

// HandleCallback processes the OAuth callback
func (c *Client) HandleCallback(ctx context.Context, state string, r *http.Request) error {
	token, err := c.auth.Token(ctx, state, r)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	c.token = token
	httpClient := c.auth.Client(ctx, token)
	c.client = spotify.New(httpClient)

	return nil
}

// SaveToken saves the current token to a file (encrypted)
func (c *Client) SaveToken(filename string) error {
	if c.token == nil {
		return fmt.Errorf("no token to save")
	}

	data, err := json.Marshal(c.token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	// Try to encrypt the token
	encryptor, err := crypto.NewTokenEncryptor()
	if err != nil {
		// Fall back to plaintext if encryption fails (with warning)
		fmt.Printf("Warning: Could not encrypt token, saving in plaintext: %v\n", err)
		if err := os.WriteFile(filename, data, 0600); err != nil {
			return fmt.Errorf("failed to write token file: %w", err)
		}
		return nil
	}

	// Save encrypted token
	if err := encryptor.SaveEncryptedFile(filename, data); err != nil {
		return fmt.Errorf("failed to save encrypted token: %w", err)
	}

	return nil
}

func (c *Client) loadToken(filename string) (*oauth2.Token, error) {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, err
	}

	var data []byte
	var err error

	// Check if the file is encrypted
	if crypto.IsEncryptedFile(filename) {
		encryptor, createErr := crypto.NewTokenEncryptor()
		if createErr != nil {
			return nil, fmt.Errorf("failed to create encryptor: %w", createErr)
		}
		data, err = encryptor.LoadEncryptedFile(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt token: %w", err)
		}
	} else {
		// Legacy: load plaintext token
		data, err = os.ReadFile(filename)
		if err != nil {
			return nil, err
		}
	}

	var token oauth2.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, err
	}

	return &token, nil
}

// GetCurrentUser returns the current user's profile
func (c *Client) GetCurrentUser(ctx context.Context) (*spotify.PrivateUser, error) {
	if c.client == nil {
		return nil, fmt.Errorf("client not authenticated")
	}
	return c.client.CurrentUser(ctx)
}

// GetSavedTracks returns all saved tracks
func (c *Client) GetSavedTracks(ctx context.Context) ([]spotify.SavedTrack, error) {
	if c.client == nil {
		return nil, fmt.Errorf("client not authenticated")
	}

	var allTracks []spotify.SavedTrack
	limit := 50
	offset := 0

	for {
		tracks, err := c.client.CurrentUsersTracks(ctx, spotify.Limit(limit), spotify.Offset(offset))
		if err != nil {
			return nil, fmt.Errorf("failed to get tracks: %w", err)
		}

		allTracks = append(allTracks, tracks.Tracks...)

		if len(tracks.Tracks) < limit {
			break
		}
		offset += limit
	}

	return allTracks, nil
}

// GetPlaylists returns all user playlists
func (c *Client) GetPlaylists(ctx context.Context) ([]spotify.SimplePlaylist, error) {
	if c.client == nil {
		return nil, fmt.Errorf("client not authenticated")
	}

	var allPlaylists []spotify.SimplePlaylist
	limit := 50
	offset := 0

	for {
		playlists, err := c.client.CurrentUsersPlaylists(ctx, spotify.Limit(limit), spotify.Offset(offset))
		if err != nil {
			return nil, fmt.Errorf("failed to get playlists: %w", err)
		}

		allPlaylists = append(allPlaylists, playlists.Playlists...)

		if len(playlists.Playlists) < limit {
			break
		}
		offset += limit
	}

	return allPlaylists, nil
}

// GetPlaylistTracks returns all tracks in a playlist
func (c *Client) GetPlaylistTracks(ctx context.Context, playlistID spotify.ID) ([]spotify.PlaylistItem, error) {
	if c.client == nil {
		return nil, fmt.Errorf("client not authenticated")
	}

	var allItems []spotify.PlaylistItem
	limit := 100
	offset := 0

	for {
		items, err := c.client.GetPlaylistItems(ctx, playlistID, spotify.Limit(limit), spotify.Offset(offset))
		if err != nil {
			return nil, fmt.Errorf("failed to get playlist items: %w", err)
		}

		allItems = append(allItems, items.Items...)

		if len(items.Items) < limit {
			break
		}
		offset += limit
	}

	return allItems, nil
}

// GetFollowedArtists returns all followed artists
func (c *Client) GetFollowedArtists(ctx context.Context) ([]spotify.FullArtist, error) {
	if c.client == nil {
		return nil, fmt.Errorf("client not authenticated")
	}

	var allArtists []spotify.FullArtist
	var after string

	for {
		artists, err := c.client.CurrentUsersFollowedArtists(ctx, spotify.Limit(50), spotify.After(after))
		if err != nil {
			return nil, fmt.Errorf("failed to get followed artists: %w", err)
		}

		allArtists = append(allArtists, artists.Artists...)

		if len(artists.Artists) == 0 {
			break
		}
		after = string(artists.Artists[len(artists.Artists)-1].ID)

		// Check if we got fewer than requested (last page)
		if len(artists.Artists) < 50 {
			break
		}
	}

	return allArtists, nil
}

// parseTimeRange converts a time range string to a Spotify Range constant
func parseTimeRange(timeRange string) spotify.Range {
	switch timeRange {
	case "short":
		return spotify.ShortTermRange
	case "medium":
		return spotify.MediumTermRange
	case "long":
		return spotify.LongTermRange
	default:
		return spotify.MediumTermRange
	}
}

// GetTopTracks returns the user's top tracks
func (c *Client) GetTopTracks(ctx context.Context, timeRange string) ([]spotify.FullTrack, error) {
	if c.client == nil {
		return nil, fmt.Errorf("client not authenticated")
	}

	tracks, err := c.client.CurrentUsersTopTracks(ctx, spotify.Limit(50), spotify.Timerange(parseTimeRange(timeRange)))
	if err != nil {
		return nil, fmt.Errorf("failed to get top tracks: %w", err)
	}

	return tracks.Tracks, nil
}

// GetTopArtists returns the user's top artists
func (c *Client) GetTopArtists(ctx context.Context, timeRange string) ([]spotify.FullArtist, error) {
	if c.client == nil {
		return nil, fmt.Errorf("client not authenticated")
	}

	artists, err := c.client.CurrentUsersTopArtists(ctx, spotify.Limit(50), spotify.Timerange(parseTimeRange(timeRange)))
	if err != nil {
		return nil, fmt.Errorf("failed to get top artists: %w", err)
	}

	return artists.Artists, nil
}

// GetRecentlyPlayed returns recently played tracks
func (c *Client) GetRecentlyPlayed(ctx context.Context) ([]spotify.RecentlyPlayedItem, error) {
	if c.client == nil {
		return nil, fmt.Errorf("client not authenticated")
	}

	items, err := c.client.PlayerRecentlyPlayed(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get recently played: %w", err)
	}

	return items, nil
}
