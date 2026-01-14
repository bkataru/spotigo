package cmd

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/spf13/cobra"

	"github.com/bkataru/spotigo/internal/spotify"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with Spotify",
	Long: `Authenticate with Spotify to enable library access.

This uses OAuth2 to securely connect to your Spotify account.
Your credentials are stored locally and never sent anywhere else.

Required scopes:
  - user-library-read (saved tracks, albums)
  - playlist-read-private (your playlists)
  - user-top-read (top artists/tracks)
  - user-read-recently-played (recent history)
  - user-follow-read (followed artists)`,
	Run: func(cmd *cobra.Command, args []string) {
		runAuth()
	},
}

func init() {
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authLogoutCmd)
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check authentication status",
	Run: func(cmd *cobra.Command, args []string) {
		checkAuthStatus()
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored credentials",
	Run: func(cmd *cobra.Command, args []string) {
		logout()
	},
}

func runAuth() {
	fmt.Println("Spotify Authentication")
	fmt.Println("======================")
	fmt.Println()

	cfg := GetConfig()
	if cfg == nil {
		fmt.Println("Error: Configuration not loaded")
		return
	}

	// Check for client ID/secret
	if cfg.Spotify.ClientID == "" {
		fmt.Println("Spotify Client ID not configured.")
		fmt.Println()
		fmt.Println("To set up authentication:")
		fmt.Println("1. Go to https://developer.spotify.com/dashboard")
		fmt.Println("2. Create a new application")
		fmt.Println("3. Set the redirect URI to: http://localhost:8888/callback")
		fmt.Println("4. Copy your Client ID and Client Secret")
		fmt.Println("5. Set them in your config file or environment:")
		fmt.Println("   export SPOTIFY_CLIENT_ID=your_client_id")
		fmt.Println("   export SPOTIFY_CLIENT_SECRET=your_client_secret")
		return
	}

	// Create Spotify client
	spotifyCfg := spotify.Config{
		ClientID:     cfg.Spotify.ClientID,
		ClientSecret: cfg.Spotify.ClientSecret,
		RedirectURI:  cfg.Spotify.RedirectURI,
		TokenFile:    cfg.Spotify.TokenFile,
	}

	client, err := spotify.NewClient(spotifyCfg)
	if err != nil {
		fmt.Printf("Error creating Spotify client: %v\n", err)
		return
	}

	// Generate random state for security
	state, err := generateRandomState()
	if err != nil {
		fmt.Printf("Error generating state: %v\n", err)
		return
	}

	// Get auth URL
	authURL := client.GetAuthURL(state)

	fmt.Println("Starting OAuth2 flow...")
	fmt.Printf("Opening browser: %s\n", authURL)
	fmt.Println()

	// Open browser
	if err := openBrowser(authURL); err != nil {
		fmt.Printf("Could not open browser automatically: %v\n", err)
		fmt.Println("Please open this URL manually:")
		fmt.Println(authURL)
	}

	// Start HTTP server to handle callback
	fmt.Println("Waiting for callback...")
	if err := handleCallback(client, state); err != nil {
		fmt.Printf("Error handling callback: %v\n", err)
		return
	}

	fmt.Println("✅ Authentication successful!")

	// Test the authentication
	if err := testAuthentication(client); err != nil {
		fmt.Printf("Warning: Authentication test failed: %v\n", err)
	} else {
		fmt.Println("✅ Authentication verified!")
	}
}

func checkAuthStatus() {
	cfg := GetConfig()
	if cfg == nil {
		fmt.Println("Authentication Status:")
		fmt.Println("  Status: Configuration not loaded")
		return
	}

	spotifyCfg := spotify.Config{
		ClientID:     cfg.Spotify.ClientID,
		ClientSecret: cfg.Spotify.ClientSecret,
		RedirectURI:  cfg.Spotify.RedirectURI,
		TokenFile:    cfg.Spotify.TokenFile,
	}

	client, err := spotify.NewClient(spotifyCfg)
	if err != nil {
		fmt.Printf("Authentication Status:")
		fmt.Println("  Status: Error creating client")
		fmt.Printf("  Error: %v\n", err)
		return
	}

	if client.IsAuthenticated() {
		fmt.Println("Authentication Status:")
		fmt.Println("  Status: ✅ Authenticated")

		// Get user info to verify token is valid
		ctx := context.Background()
		user, err := client.GetCurrentUser(ctx)
		if err != nil {
			fmt.Println("  Token: ⚠️  Invalid or expired")
			fmt.Printf("  Error: %v\n", err)
			fmt.Println("  Run 'spotigo auth' to re-authenticate")
		} else {
			fmt.Printf("  User: %s\n", user.DisplayName)
			fmt.Printf("  ID: %s\n", user.ID)
			fmt.Println("  Token: ✅ Valid")
		}
	} else {
		fmt.Println("Authentication Status:")
		fmt.Println("  Status: ❌ Not authenticated")
		fmt.Println("  Token file: Missing or invalid")
		fmt.Println("  Run 'spotigo auth' to authenticate")
	}
}

func logout() {
	cfg := GetConfig()
	if cfg == nil {
		fmt.Println("Error: Configuration not loaded")
		return
	}

	tokenFile := cfg.Spotify.TokenFile
	if tokenFile == "" {
		tokenFile = ".spotify_token" //nolint:gosec // This is a filename, not a credential
	}

	if err := os.Remove(tokenFile); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No stored credentials found.")
		} else {
			fmt.Printf("Error removing credentials: %v\n", err)
		}
	} else {
		fmt.Println("✅ Credentials removed.")
		fmt.Printf("Deleted token file: %s\n", tokenFile)
	}
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // linux
		cmd = "xdg-open"
		args = []string{url}
	}

	if runtime.GOOS == "windows" {
		args = append(args, url)
	}

	return exec.Command(cmd, args...).Start() //nolint:gosec // This is intentional browser opening
}

func handleCallback(client *spotify.Client, expectedState string) error {
	done := make(chan error, 1)

	server := &http.Server{
		Addr:              ":8888",
		ReadHeaderTimeout: 5 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/callback" {
				http.NotFound(w, r)
				return
			}

			state := r.URL.Query().Get("state")
			if state != expectedState {
				done <- fmt.Errorf("invalid state parameter")
				return
			}

			if err := client.HandleCallback(r.Context(), expectedState, r); err != nil {
				done <- fmt.Errorf("callback failed: %w", err)
				return
			}

			// Save the token
			cfg := GetConfig()
			if err := client.SaveToken(cfg.Spotify.TokenFile); err != nil {
				done <- fmt.Errorf("failed to save token: %w", err)
				return
			}

			// Send success response
			w.Header().Set("Content-Type", "text/html")
			if _, err := fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><title>Spotigo - Authentication Successful</title></head>
<body>
	<h1>✅ Authentication Successful!</h1>
	<p>You can close this window and return to your terminal.</p>
	<script>
		setTimeout(() => window.close(), 3000);
	</script>
</body>
</html>`); err != nil {
				log.Printf("Failed to write auth success response: %v", err)
			}
			done <- nil
		}),
	}

	// Start server in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			done <- err
		}
	}()

	// Wait for callback or timeout
	select {
	case err := <-done:
		// Give server time to respond, then shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if shutdownErr := server.Shutdown(ctx); shutdownErr != nil {
			fmt.Printf("Warning: server shutdown failed: %v\n", shutdownErr)
		}
		return err
	case <-time.After(5 * time.Minute):
		if closeErr := server.Close(); closeErr != nil {
			fmt.Printf("Warning: server close failed: %v\n", closeErr)
		}
		return fmt.Errorf("authentication timeout: no callback received within 5 minutes")
	}
}

func testAuthentication(client *spotify.Client) error {
	ctx := context.Background()
	_, err := client.GetCurrentUser(ctx)
	return err
}
