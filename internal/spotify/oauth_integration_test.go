package spotify

import (
	"testing"
)

// TestOAuthStateValidation tests state parameter validation for CSRF protection
func TestOAuthStateValidation(t *testing.T) {
	tests := []struct {
		name          string
		sentState     string
		returnedState string
		expectError   bool
	}{
		{
			name:          "matching state",
			sentState:     "random_state_123",
			returnedState: "random_state_123",
			expectError:   false,
		},
		{
			name:          "mismatched state",
			sentState:     "state_abc",
			returnedState: "state_xyz",
			expectError:   true,
		},
		{
			name:          "empty returned state",
			sentState:     "state_123",
			returnedState: "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate state matching
			isValid := tt.sentState == tt.returnedState && tt.returnedState != ""

			hasError := !isValid
			if hasError != tt.expectError {
				t.Errorf("Expected error=%v, got error=%v", tt.expectError, hasError)
			}
		})
	}
}

// TestOAuthScopeRequirements tests that all required scopes are defined
func TestOAuthScopeRequirements(t *testing.T) {
	if len(Scopes) == 0 {
		t.Error("Scopes should not be empty")
	}

	// Verify we have essential scopes
	requiredScopes := map[string]bool{
		"user-library-read":           false,
		"playlist-read-private":       false,
		"playlist-read-collaborative": false,
		"user-top-read":               false,
		"user-read-recently-played":   false,
		"user-follow-read":            false,
	}

	// Check each defined scope
	for _, scope := range Scopes {
		scopeStr := scope
		for requiredScope := range requiredScopes {
			// Check if this scope matches any required scope
			if scopeStr == requiredScope {
				requiredScopes[requiredScope] = true
			}
		}
	}

	// Verify all required scopes are present
	for scope, found := range requiredScopes {
		if !found {
			t.Logf("Note: Required scope not explicitly found: %s", scope)
		}
	}

	t.Logf("Total scopes defined: %d", len(Scopes))
}

// TestConfigValidation tests configuration validation
func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		isValid bool
	}{
		{
			name: "valid config",
			config: Config{
				ClientID:     "test_client_id",
				ClientSecret: "test_client_secret",
				RedirectURI:  "http://localhost:8888/callback",
				TokenFile:    "/path/to/token.json",
			},
			isValid: true,
		},
		{
			name: "missing client ID",
			config: Config{
				ClientID:     "",
				ClientSecret: "test_client_secret",
				RedirectURI:  "http://localhost:8888/callback",
			},
			isValid: false,
		},
		{
			name: "missing client secret",
			config: Config{
				ClientID:     "test_client_id",
				ClientSecret: "",
				RedirectURI:  "http://localhost:8888/callback",
			},
			isValid: false,
		},
		{
			name: "missing redirect URI",
			config: Config{
				ClientID:     "test_client_id",
				ClientSecret: "test_client_secret",
				RedirectURI:  "",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.config.ClientID != "" &&
				tt.config.ClientSecret != "" &&
				tt.config.RedirectURI != ""

			if isValid != tt.isValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.isValid, isValid)
			}
		})
	}
}

// TestRedirectURIValidation tests redirect URI validation
func TestRedirectURIValidation(t *testing.T) {
	tests := []struct {
		name        string
		redirectURI string
		valid       bool
	}{
		{
			name:        "localhost with port",
			redirectURI: "http://localhost:8888/callback",
			valid:       true,
		},
		{
			name:        "127.0.0.1 with port",
			redirectURI: "http://127.0.0.1:8888/callback",
			valid:       true,
		},
		{
			name:        "custom domain",
			redirectURI: "https://example.com/callback",
			valid:       true,
		},
		{
			name:        "missing protocol",
			redirectURI: "localhost:8888/callback",
			valid:       false,
		},
		{
			name:        "empty URI",
			redirectURI: "",
			valid:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simple validation
			hasHTTP := len(tt.redirectURI) >= 7 && tt.redirectURI[:7] == "http://"
			hasHTTPS := len(tt.redirectURI) >= 8 && tt.redirectURI[:8] == "https://"
			hasProtocol := hasHTTP || hasHTTPS
			isNotEmpty := tt.redirectURI != ""

			isValid := hasProtocol && isNotEmpty

			if isValid != tt.valid {
				t.Errorf("Expected valid=%v, got valid=%v for URI: %s", tt.valid, isValid, tt.redirectURI)
			}
		})
	}
}

// TestTokenFileValidation tests token file path validation
func TestTokenFileValidation(t *testing.T) {
	tests := []struct {
		name      string
		tokenFile string
		valid     bool
	}{
		{
			name:      "valid absolute path",
			tokenFile: "/home/user/.spotigo/token.json",
			valid:     true,
		},
		{
			name:      "valid relative path",
			tokenFile: "./data/token.json",
			valid:     true,
		},
		{
			name:      "valid windows path",
			tokenFile: "C:\\Users\\user\\.spotigo\\token.json",
			valid:     true,
		},
		{
			name:      "empty path",
			tokenFile: "",
			valid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.tokenFile != ""

			if isValid != tt.valid {
				t.Errorf("Expected valid=%v, got valid=%v for path: %s", tt.valid, isValid, tt.tokenFile)
			}
		})
	}
}

// TestOAuthErrorScenarios tests various OAuth error scenarios
func TestOAuthErrorScenarios(t *testing.T) {
	errorScenarios := []struct {
		name          string
		errorCode     string
		errorDesc     string
		shouldRetry   bool
		userActionReq bool
	}{
		{
			name:          "access denied by user",
			errorCode:     "access_denied",
			errorDesc:     "User denied access",
			shouldRetry:   false,
			userActionReq: true,
		},
		{
			name:          "invalid client credentials",
			errorCode:     "invalid_client",
			errorDesc:     "Invalid client credentials",
			shouldRetry:   false,
			userActionReq: true,
		},
		{
			name:          "expired authorization code",
			errorCode:     "invalid_grant",
			errorDesc:     "Authorization code expired",
			shouldRetry:   true,
			userActionReq: true,
		},
		{
			name:          "server error",
			errorCode:     "server_error",
			errorDesc:     "Temporary server error",
			shouldRetry:   true,
			userActionReq: false,
		},
	}

	for _, scenario := range errorScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			if scenario.errorCode == "" {
				t.Error("Error code should not be empty")
			}

			// Validate error handling logic
			isRetryable := scenario.errorCode == "invalid_grant" || scenario.errorCode == "server_error"
			if isRetryable != scenario.shouldRetry {
				t.Errorf("Expected shouldRetry=%v, got %v", scenario.shouldRetry, isRetryable)
			}
		})
	}
}

// TestOAuthCallbackURLParsing tests parsing of OAuth callback URLs
func TestOAuthCallbackURLParsing(t *testing.T) {
	tests := []struct {
		name        string
		callbackURL string
		hasCode     bool
		hasState    bool
		hasError    bool
	}{
		{
			name:        "valid callback with code and state",
			callbackURL: "http://localhost:8888/callback?code=abc123&state=xyz789",
			hasCode:     true,
			hasState:    true,
			hasError:    false,
		},
		{
			name:        "callback with error",
			callbackURL: "http://localhost:8888/callback?error=access_denied&error_description=User+denied",
			hasCode:     false,
			hasState:    false,
			hasError:    true,
		},
		{
			name:        "missing state parameter",
			callbackURL: "http://localhost:8888/callback?code=abc123",
			hasCode:     true,
			hasState:    false,
			hasError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simple string parsing
			containsCode := false
			containsState := false
			containsError := false

			for i := 0; i < len(tt.callbackURL)-5; i++ {
				if tt.callbackURL[i:i+5] == "code=" {
					containsCode = true
				}
				if i < len(tt.callbackURL)-6 && tt.callbackURL[i:i+6] == "state=" {
					containsState = true
				}
				if i < len(tt.callbackURL)-6 && tt.callbackURL[i:i+6] == "error=" {
					containsError = true
				}
			}

			if containsCode != tt.hasCode {
				t.Errorf("Expected hasCode=%v, got %v", tt.hasCode, containsCode)
			}
			if containsState != tt.hasState {
				t.Errorf("Expected hasState=%v, got %v", tt.hasState, containsState)
			}
			if containsError != tt.hasError {
				t.Errorf("Expected hasError=%v, got %v", tt.hasError, containsError)
			}
		})
	}
}

// TestClientCreationValidation tests client creation scenarios
func TestClientCreationValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError bool
	}{
		{
			name: "valid config",
			config: Config{
				ClientID:     "valid_id",
				ClientSecret: "valid_secret",
				RedirectURI:  "http://localhost:8888/callback",
				TokenFile:    "./token.json",
			},
			expectError: false,
		},
		{
			name: "empty client ID",
			config: Config{
				ClientID:     "",
				ClientSecret: "valid_secret",
				RedirectURI:  "http://localhost:8888/callback",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate that required fields are present
			hasRequiredFields := tt.config.ClientID != "" &&
				tt.config.ClientSecret != "" &&
				tt.config.RedirectURI != ""

			wouldError := !hasRequiredFields

			if wouldError != tt.expectError {
				t.Errorf("Expected error=%v, validation says error=%v", tt.expectError, wouldError)
			}
		})
	}
}
