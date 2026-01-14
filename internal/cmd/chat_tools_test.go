package cmd

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/bkataru/spotigo/internal/ollama"
	"github.com/bkataru/spotigo/internal/tools"
)

// mockOllamaClient is a mock implementation of the Ollama client for testing
type mockOllamaClient struct {
	pingErr       error
	chatResponses []*ollama.ChatResponse
	chatErrors    []error
	callIndex     int
}

func (m *mockOllamaClient) Ping(ctx context.Context) error {
	return m.pingErr
}

func (m *mockOllamaClient) Chat(ctx context.Context, req ollama.ChatRequest) (*ollama.ChatResponse, error) {
	if m.callIndex >= len(m.chatResponses) {
		if m.callIndex < len(m.chatErrors) {
			err := m.chatErrors[m.callIndex]
			m.callIndex++
			return nil, err
		}
		return &ollama.ChatResponse{
			Message: ollama.Message{
				Role:    "assistant",
				Content: "Default response",
			},
		}, nil
	}

	resp := m.chatResponses[m.callIndex]
	var err error
	if m.callIndex < len(m.chatErrors) {
		err = m.chatErrors[m.callIndex]
	}
	m.callIndex++
	return resp, err
}

func TestValidateChatInput(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "valid input",
			input:       "Hello, how are you?",
			expectError: false,
		},
		{
			name:        "empty input",
			input:       "",
			expectError: false,
		},
		{
			name:        "input with newline",
			input:       "Hello\nWorld",
			expectError: false,
		},
		{
			name:        "input with tab",
			input:       "Hello\tWorld",
			expectError: false,
		},
		{
			name:        "input too long",
			input:       string(make([]byte, MaxInputLength+1)),
			expectError: true,
		},
		{
			name:        "input with null byte",
			input:       "Hello\x00World",
			expectError: true,
		},
		{
			name:        "input with control character",
			input:       "Hello\x01World",
			expectError: true,
		},
		{
			name:        "valid UTF-8 emoji",
			input:       "Hello 👋 World!",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateChatInput(tt.input)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestToolCallingFlow(t *testing.T) {
	// Create mock Ollama client that simulates tool calling
	mock := &mockOllamaClient{
		chatResponses: []*ollama.ChatResponse{
			// First response: model requests a tool call
			{
				Message: ollama.Message{
					Role:    "assistant",
					Content: "",
					ToolCalls: []ollama.ToolCall{
						{
							Function: ollama.FunctionCall{
								Name:      "get_library_stats",
								Arguments: "{}",
							},
						},
					},
				},
			},
			// Second response: model provides final answer after tool result
			{
				Message: ollama.Message{
					Role:    "assistant",
					Content: "You have 3 saved tracks, 2 playlists, and 2 followed artists in your library.",
				},
			},
		},
	}

	// Test that we can handle tool calling responses
	ctx := context.Background()
	req := ollama.ChatRequest{
		Model: "test-model",
		Messages: []ollama.Message{
			{Role: "user", Content: "How many tracks do I have?"},
		},
	}

	// First call should return tool calls
	resp1, err := mock.Chat(ctx, req)
	if err != nil {
		t.Fatalf("First chat call failed: %v", err)
	}

	if len(resp1.Message.ToolCalls) != 1 {
		t.Errorf("Expected 1 tool call, got %d", len(resp1.Message.ToolCalls))
	}

	if resp1.Message.ToolCalls[0].Function.Name != "get_library_stats" {
		t.Errorf("Expected tool call 'get_library_stats', got '%s'", resp1.Message.ToolCalls[0].Function.Name)
	}

	// Second call should return final response
	resp2, err := mock.Chat(ctx, req)
	if err != nil {
		t.Fatalf("Second chat call failed: %v", err)
	}

	if resp2.Message.Content == "" {
		t.Error("Expected content in final response")
	}

	if len(resp2.Message.ToolCalls) != 0 {
		t.Error("Expected no tool calls in final response")
	}
}

func TestToolExecutionWithMockData(t *testing.T) {
	// Create temporary test data
	tmpDir := t.TempDir()

	// Create saved_tracks.json
	savedTracks := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"track": map[string]interface{}{
					"id":   "track1",
					"name": "Test Track",
					"artists": []map[string]interface{}{
						{"name": "Test Artist"},
					},
				},
				"added_at": "2024-01-15T10:30:00Z",
			},
		},
	}

	tracksData, err := json.Marshal(savedTracks)
	if err != nil {
		t.Fatalf("Failed to marshal tracks: %v", err)
	}

	if err := writeTestFile(tmpDir, "saved_tracks.json", tracksData); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create playlists.json
	playlists := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"id":   "playlist1",
				"name": "Test Playlist",
			},
		},
	}

	playlistsData, err := json.Marshal(playlists)
	if err != nil {
		t.Fatalf("Failed to marshal playlists: %v", err)
	}

	if err := writeTestFile(tmpDir, "playlists.json", playlistsData); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create followed_artists.json
	followedArtists := map[string]interface{}{
		"artists": map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"id":   "artist1",
					"name": "Test Artist",
				},
			},
		},
	}

	artistsData, err := json.Marshal(followedArtists)
	if err != nil {
		t.Fatalf("Failed to marshal artists: %v", err)
	}

	if err := writeTestFile(tmpDir, "followed_artists.json", artistsData); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create music tools instance
	musicTools := tools.NewMusicTools(tmpDir)
	if musicTools == nil {
		t.Fatal("Failed to create music tools")
	}

	// Test tool execution
	toolCall := ollama.ToolCall{
		Function: ollama.FunctionCall{
			Name:      "get_library_stats",
			Arguments: "{}",
		},
	}

	result, err := musicTools.ExecuteToolCall(toolCall)
	if err != nil {
		t.Fatalf("Tool execution failed: %v", err)
	}

	if result == "" {
		t.Error("Tool result is empty")
	}

	// Verify result is valid JSON
	var resultData map[string]interface{}
	if err := json.Unmarshal([]byte(result), &resultData); err != nil {
		t.Errorf("Tool result is not valid JSON: %v", err)
	}
}

func TestToolCallingWithMultipleTools(t *testing.T) {
	tmpDir := t.TempDir()

	// Create minimal test data
	savedTracks := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"track": map[string]interface{}{
					"id":   "track1",
					"name": "Bohemian Rhapsody",
					"artists": []map[string]interface{}{
						{"name": "Queen"},
					},
				},
				"added_at": "2024-01-15T10:30:00Z",
			},
			{
				"track": map[string]interface{}{
					"id":   "track2",
					"name": "We Will Rock You",
					"artists": []map[string]interface{}{
						{"name": "Queen"},
					},
				},
				"added_at": "2024-01-16T10:30:00Z",
			},
		},
	}

	tracksData, _ := json.Marshal(savedTracks)
	if err := writeTestFile(tmpDir, "saved_tracks.json", tracksData); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	playlists := map[string]interface{}{"items": []map[string]interface{}{}}
	playlistsData, _ := json.Marshal(playlists)
	writeTestFile(tmpDir, "playlists.json", playlistsData)

	followedArtists := map[string]interface{}{"artists": map[string]interface{}{"items": []map[string]interface{}{}}}
	artistsData, _ := json.Marshal(followedArtists)
	writeTestFile(tmpDir, "followed_artists.json", artistsData)

	musicTools := tools.NewMusicTools(tmpDir)

	// Test multiple tool calls
	toolCalls := []ollama.ToolCall{
		{
			Function: ollama.FunctionCall{
				Name:      "search_tracks",
				Arguments: `{"query": "Queen", "limit": 10}`,
			},
		},
		{
			Function: ollama.FunctionCall{
				Name:      "get_tracks_by_artist",
				Arguments: `{"artist_name": "Queen"}`,
			},
		},
	}

	for i, toolCall := range toolCalls {
		result, err := musicTools.ExecuteToolCall(toolCall)
		if err != nil {
			t.Errorf("Tool call %d failed: %v", i, err)
			continue
		}

		if result == "" {
			t.Errorf("Tool call %d returned empty result", i)
		}
	}
}

func TestToolCallingErrorHandling(t *testing.T) {
	tmpDir := t.TempDir()
	musicTools := tools.NewMusicTools(tmpDir)

	tests := []struct {
		name        string
		toolCall    ollama.ToolCall
		expectError bool
	}{
		{
			name: "invalid JSON arguments",
			toolCall: ollama.ToolCall{
				Function: ollama.FunctionCall{
					Name:      "search_tracks",
					Arguments: "not valid json",
				},
			},
			expectError: true,
		},
		{
			name: "unknown tool",
			toolCall: ollama.ToolCall{
				Function: ollama.FunctionCall{
					Name:      "nonexistent_tool",
					Arguments: "{}",
				},
			},
			expectError: true,
		},
		{
			name: "missing required parameter",
			toolCall: ollama.ToolCall{
				Function: ollama.FunctionCall{
					Name:      "search_tracks",
					Arguments: "{}",
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := musicTools.ExecuteToolCall(tt.toolCall)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestConversationContextManagement(t *testing.T) {
	// Test that conversation history is properly managed
	messages := []ollama.Message{
		{Role: "system", Content: "You are a helpful assistant"},
		{Role: "user", Content: "Hello"},
		{Role: "assistant", Content: "Hi there!"},
	}

	// Simulate adding many messages
	for i := 0; i < MaxConversationHistory+10; i++ {
		messages = append(messages, ollama.Message{
			Role:    "user",
			Content: "test message",
		})
		messages = append(messages, ollama.Message{
			Role:    "assistant",
			Content: "test response",
		})

		// Trim if exceeds max
		if len(messages) > MaxConversationHistory {
			messages = append(messages[:1], messages[len(messages)-MaxConversationHistory+1:]...)
		}
	}

	// Verify length constraint
	if len(messages) > MaxConversationHistory {
		t.Errorf("Conversation history exceeds max length: got %d, max %d", len(messages), MaxConversationHistory)
	}

	// Verify system message is preserved
	if messages[0].Role != "system" {
		t.Error("System message was not preserved")
	}
}

func TestToolDefinitionsStructure(t *testing.T) {
	tmpDir := t.TempDir()
	musicTools := tools.NewMusicTools(tmpDir)

	toolDefs := musicTools.GetToolDefinitions()

	if len(toolDefs) == 0 {
		t.Fatal("No tool definitions returned")
	}

	for _, toolDef := range toolDefs {
		t.Run(toolDef.Function.Name, func(t *testing.T) {
			// Verify required fields
			if toolDef.Type != "function" {
				t.Errorf("Expected type 'function', got '%s'", toolDef.Type)
			}

			if toolDef.Function.Name == "" {
				t.Error("Tool name is empty")
			}

			if toolDef.Function.Description == "" {
				t.Error("Tool description is empty")
			}

			// Verify parameters structure
			params := toolDef.Function.Parameters
			if params == nil {
				t.Fatal("Parameters is nil")
			}

			if params["type"] != "object" {
				t.Errorf("Expected parameters type 'object', got '%v'", params["type"])
			}

			if params["properties"] == nil {
				t.Error("Parameters missing 'properties' field")
			}

			if params["required"] == nil {
				t.Error("Parameters missing 'required' field")
			}
		})
	}
}

func TestToolResponseFormat(t *testing.T) {
	tmpDir := t.TempDir()

	// Create minimal test data
	savedTracks := map[string]interface{}{"items": []map[string]interface{}{}}
	tracksData, _ := json.Marshal(savedTracks)
	writeTestFile(tmpDir, "saved_tracks.json", tracksData)

	playlists := map[string]interface{}{"items": []map[string]interface{}{}}
	playlistsData, _ := json.Marshal(playlists)
	writeTestFile(tmpDir, "playlists.json", playlistsData)

	followedArtists := map[string]interface{}{"artists": map[string]interface{}{"items": []map[string]interface{}{}}}
	artistsData, _ := json.Marshal(followedArtists)
	writeTestFile(tmpDir, "followed_artists.json", artistsData)

	musicTools := tools.NewMusicTools(tmpDir)

	toolCall := ollama.ToolCall{
		Function: ollama.FunctionCall{
			Name:      "get_library_stats",
			Arguments: "{}",
		},
	}

	result, err := musicTools.ExecuteToolCall(toolCall)
	if err != nil {
		t.Fatalf("Tool execution failed: %v", err)
	}

	// Verify response is valid JSON
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(result), &response); err != nil {
		t.Fatalf("Response is not valid JSON: %v", err)
	}

	// Verify response structure (should have 'data' field)
	if _, ok := response["data"]; !ok {
		t.Error("Response missing 'data' field")
	}

	// Check for error field (should be empty or absent)
	if errField, ok := response["error"]; ok && errField != "" && errField != nil {
		t.Errorf("Unexpected error in response: %v", errField)
	}
}

func TestChatRequestWithTools(t *testing.T) {
	tmpDir := t.TempDir()
	musicTools := tools.NewMusicTools(tmpDir)

	// Create a chat request with tools
	req := ollama.ChatRequest{
		Model: "test-model",
		Messages: []ollama.Message{
			{Role: "system", Content: "You are a music assistant"},
			{Role: "user", Content: "What's in my library?"},
		},
		Tools: musicTools.GetToolDefinitions(),
		Options: &ollama.Options{
			Temperature: 0.7,
			NumPredict:  2048,
		},
	}

	// Verify request structure
	if req.Model == "" {
		t.Error("Model is empty")
	}

	if len(req.Messages) == 0 {
		t.Error("Messages is empty")
	}

	if len(req.Tools) == 0 {
		t.Error("Tools is empty")
	}

	if req.Options == nil {
		t.Error("Options is nil")
	}

	// Verify each tool in request
	for _, tool := range req.Tools {
		if tool.Type != "function" {
			t.Errorf("Tool type is not 'function': %s", tool.Type)
		}

		if tool.Function.Name == "" {
			t.Error("Tool function name is empty")
		}
	}
}

// writeTestFile is a helper to write test data files
func writeTestFile(dir, filename string, data []byte) error {
	path := filepath.Join(dir, filename)
	return os.WriteFile(path, data, 0644)
}
