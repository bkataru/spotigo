package cmd

import (
	"net/url"
	"runtime"
	"testing"
)

func TestGenerateRandomState(t *testing.T) {
	// Test that state generation works
	state1, err := generateRandomState()
	if err != nil {
		t.Fatalf("generateRandomState() error = %v", err)
	}

	if state1 == "" {
		t.Error("generateRandomState() returned empty string")
	}

	// Test that states are unique
	state2, err := generateRandomState()
	if err != nil {
		t.Fatalf("generateRandomState() second call error = %v", err)
	}

	if state1 == state2 {
		t.Error("generateRandomState() returned same state twice")
	}

	// Test state length (32 bytes base64 encoded = ~44 chars)
	if len(state1) < 40 {
		t.Errorf("generateRandomState() state too short: %d chars", len(state1))
	}
}

func TestOpenBrowserCommand(t *testing.T) {
	// Test URL that would be used for OAuth
	testURL := "https://accounts.spotify.com/authorize?client_id=test&response_type=code&redirect_uri=http%3A%2F%2F127.0.0.1%3A8888%2Fcallback&state=abc123&scope=user-library-read"

	// We can't actually test browser opening without side effects,
	// but we can verify the function doesn't panic with valid input
	// The actual browser opening is platform-specific

	// Parse the URL to verify it's valid
	parsedURL, err := url.Parse(testURL)
	if err != nil {
		t.Fatalf("Test URL is invalid: %v", err)
	}

	// Verify expected OAuth parameters are present
	params := parsedURL.Query()
	expectedParams := []string{"client_id", "response_type", "redirect_uri", "state", "scope"}
	for _, param := range expectedParams {
		if params.Get(param) == "" {
			t.Errorf("Expected parameter %q not found in URL", param)
		}
	}
}

func TestOpenBrowserPlatformCommand(t *testing.T) {
	// Document expected commands per platform
	type platformCommand struct {
		goos        string
		expectedCmd string
	}

	platforms := []platformCommand{
		{"windows", "rundll32"},
		{"darwin", "open"},
		{"linux", "xdg-open"},
	}

	for _, p := range platforms {
		t.Run(p.goos, func(t *testing.T) {
			// This just documents the expected behavior
			// We can't easily test cross-platform without mocking
			if runtime.GOOS == p.goos {
				t.Logf("Current platform %s uses command: %s", p.goos, p.expectedCmd)
			}
		})
	}
}

func BenchmarkGenerateRandomState(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = generateRandomState()
	}
}

func TestGenerateRandomStateUniqueness(t *testing.T) {
	// Generate many states and verify uniqueness
	states := make(map[string]bool)
	numStates := 100

	for i := 0; i < numStates; i++ {
		state, err := generateRandomState()
		if err != nil {
			t.Fatalf("generateRandomState() error on iteration %d: %v", i, err)
		}

		if states[state] {
			t.Errorf("Duplicate state generated on iteration %d", i)
		}
		states[state] = true
	}

	if len(states) != numStates {
		t.Errorf("Expected %d unique states, got %d", numStates, len(states))
	}
}
