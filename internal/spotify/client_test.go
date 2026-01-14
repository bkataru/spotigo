// Package spotify provides Spotify API client functionality
package spotify

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid config without token file",
			cfg: Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				RedirectURI:  "http://localhost:8080/callback",
			},
			wantErr: false,
		},
		{
			name: "valid config with non-existent token file",
			cfg: Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				RedirectURI:  "http://localhost:8080/callback",
				TokenFile:    "/non/existent/path/token.json",
			},
			wantErr: false,
		},
		{
			name: "empty config",
			cfg:  Config{},
			// Should not error, just creates unauthenticated client
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client without error")
			}
		})
	}
}

func TestClient_IsAuthenticated(t *testing.T) {
	tests := []struct {
		name     string
		client   *Client
		expected bool
	}{
		{
			name:     "unauthenticated client",
			client:   &Client{},
			expected: false,
		},
		{
			name: "client with token but no spotify client",
			client: &Client{
				token: &oauth2.Token{
					AccessToken: "test-token",
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.client.IsAuthenticated(); got != tt.expected {
				t.Errorf("IsAuthenticated() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestClient_GetAuthURL(t *testing.T) {
	cfg := Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURI:  "http://localhost:8080/callback",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	state := "test-state-123"
	url := client.GetAuthURL(state)

	// Verify URL contains expected components
	if url == "" {
		t.Error("GetAuthURL() returned empty string")
	}

	// URL should contain the state parameter
	if !contains(url, "state="+state) {
		t.Errorf("GetAuthURL() URL missing state parameter, got: %s", url)
	}

	// URL should contain client_id
	if !contains(url, "client_id="+cfg.ClientID) {
		t.Errorf("GetAuthURL() URL missing client_id, got: %s", url)
	}

	// URL should contain redirect_uri (URL encoded)
	if !contains(url, "redirect_uri=") {
		t.Errorf("GetAuthURL() URL missing redirect_uri, got: %s", url)
	}
}

func TestClient_SaveToken(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name     string
		client   *Client
		filename string
		wantErr  bool
	}{
		{
			name:     "no token to save",
			client:   &Client{},
			filename: filepath.Join(tmpDir, "token1.json"),
			wantErr:  true,
		},
		{
			name: "valid token",
			client: &Client{
				token: &oauth2.Token{
					AccessToken:  "test-access-token",
					TokenType:    "Bearer",
					RefreshToken: "test-refresh-token",
					Expiry:       time.Now().Add(time.Hour),
				},
			},
			filename: filepath.Join(tmpDir, "token2.json"),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.client.SaveToken(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify file was created
				if _, err := os.Stat(tt.filename); os.IsNotExist(err) {
					t.Error("SaveToken() did not create file")
				}
			}
		})
	}
}

func TestClient_loadToken(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a valid token file
	validToken := &oauth2.Token{
		AccessToken:  "test-access-token",
		TokenType:    "Bearer",
		RefreshToken: "test-refresh-token",
		Expiry:       time.Now().Add(time.Hour),
	}
	validTokenData, _ := json.Marshal(validToken)
	validTokenPath := filepath.Join(tmpDir, "valid_token.json")
	if err := os.WriteFile(validTokenPath, validTokenData, 0600); err != nil {
		t.Fatalf("Failed to create test token file: %v", err)
	}

	// Create an invalid token file
	invalidTokenPath := filepath.Join(tmpDir, "invalid_token.json")
	if err := os.WriteFile(invalidTokenPath, []byte("not valid json"), 0600); err != nil {
		t.Fatalf("Failed to create invalid token file: %v", err)
	}

	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "valid token file",
			filename: validTokenPath,
			wantErr:  false,
		},
		{
			name:     "non-existent file",
			filename: filepath.Join(tmpDir, "nonexistent.json"),
			wantErr:  true,
		},
		{
			name:     "invalid json",
			filename: invalidTokenPath,
			wantErr:  true,
		},
	}

	client := &Client{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := client.loadToken(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && token == nil {
				t.Error("loadToken() returned nil token without error")
			}

			if !tt.wantErr && token != nil {
				if token.AccessToken != validToken.AccessToken {
					t.Errorf("loadToken() AccessToken = %v, want %v", token.AccessToken, validToken.AccessToken)
				}
			}
		})
	}
}

func TestClient_UnauthenticatedMethods(t *testing.T) {
	client := &Client{}
	ctx := context.Background()

	// Test GetCurrentUser without authentication
	_, err := client.GetCurrentUser(ctx)
	if err == nil {
		t.Error("GetCurrentUser() should return error when not authenticated")
	}

	// Test GetSavedTracks without authentication
	_, err = client.GetSavedTracks(ctx)
	if err == nil {
		t.Error("GetSavedTracks() should return error when not authenticated")
	}

	// Test GetPlaylists without authentication
	_, err = client.GetPlaylists(ctx)
	if err == nil {
		t.Error("GetPlaylists() should return error when not authenticated")
	}

	// Test GetPlaylistTracks without authentication
	_, err = client.GetPlaylistTracks(ctx, "test-playlist-id")
	if err == nil {
		t.Error("GetPlaylistTracks() should return error when not authenticated")
	}

	// Test GetFollowedArtists without authentication
	_, err = client.GetFollowedArtists(ctx)
	if err == nil {
		t.Error("GetFollowedArtists() should return error when not authenticated")
	}

	// Test GetTopTracks without authentication
	_, err = client.GetTopTracks(ctx, "medium")
	if err == nil {
		t.Error("GetTopTracks() should return error when not authenticated")
	}

	// Test GetTopArtists without authentication
	_, err = client.GetTopArtists(ctx, "medium")
	if err == nil {
		t.Error("GetTopArtists() should return error when not authenticated")
	}

	// Test GetRecentlyPlayed without authentication
	_, err = client.GetRecentlyPlayed(ctx)
	if err == nil {
		t.Error("GetRecentlyPlayed() should return error when not authenticated")
	}
}

func TestParseTimeRange(t *testing.T) {
	tests := []struct {
		name      string
		timeRange string
		expected  string
	}{
		{
			name:      "short term",
			timeRange: "short",
			expected:  "short_term",
		},
		{
			name:      "medium term",
			timeRange: "medium",
			expected:  "medium_term",
		},
		{
			name:      "long term",
			timeRange: "long",
			expected:  "long_term",
		},
		{
			name:      "invalid defaults to medium",
			timeRange: "invalid",
			expected:  "medium_term",
		},
		{
			name:      "empty defaults to medium",
			timeRange: "",
			expected:  "medium_term",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseTimeRange(tt.timeRange)
			if string(result) != tt.expected {
				t.Errorf("parseTimeRange(%q) = %v, want %v", tt.timeRange, result, tt.expected)
			}
		})
	}
}

func TestScopes(t *testing.T) {
	// Verify all required scopes are defined
	expectedScopes := []string{
		"user-library-read",
		"playlist-read-private",
		"playlist-read-collaborative",
		"user-top-read",
		"user-read-recently-played",
		"user-follow-read",
	}

	if len(Scopes) != len(expectedScopes) {
		t.Errorf("Expected %d scopes, got %d", len(expectedScopes), len(Scopes))
	}

	// Check that scopes are non-empty
	for i, scope := range Scopes {
		if scope == "" {
			t.Errorf("Scope at index %d is empty", i)
		}
	}
}

func TestClient_HandleCallback(t *testing.T) {
	// Create a mock OAuth server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate token endpoint
		if r.URL.Path == "/api/token" {
			token := map[string]interface{}{
				"access_token":  "mock-access-token",
				"token_type":    "Bearer",
				"refresh_token": "mock-refresh-token",
				"expires_in":    3600,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(token)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer mockServer.Close()

	// Create client
	cfg := Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURI:  "http://localhost:8080/callback",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test with invalid request (no code parameter)
	req := httptest.NewRequest("GET", "/callback?state=test-state", nil)
	ctx := context.Background()

	// This should fail because there's no authorization code
	err = client.HandleCallback(ctx, "test-state", req, cfg.RedirectURI)
	if err == nil {
		t.Error("HandleCallback() should return error with invalid request")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Benchmark tests
func BenchmarkNewClient(b *testing.B) {
	cfg := Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURI:  "http://localhost:8080/callback",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewClient(cfg)
	}
}

func BenchmarkGetAuthURL(b *testing.B) {
	cfg := Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURI:  "http://localhost:8080/callback",
	}

	client, _ := NewClient(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = client.GetAuthURL("test-state")
	}
}

func BenchmarkParseTimeRange(b *testing.B) {
	ranges := []string{"short", "medium", "long", "invalid", ""}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = parseTimeRange(ranges[i%len(ranges)])
	}
}
