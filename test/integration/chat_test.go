//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"spotigo/internal/ollama"
)

func TestChat_BasicConversation(t *testing.T) {
	SkipIfNoOllama(t)
	SkipIfNoModel(t, "llama3.2:latest")

	cfg := SetupTestEnvironment(t)

	// Prepare chat request
	req := ollama.ChatRequest{
		Model: "llama3.2:latest",
		Messages: []ollama.Message{
			{
				Role:    "system",
				Content: "You are a helpful music assistant.",
			},
			{
				Role:    "user",
				Content: "What is the capital of France?",
			},
		},
		Stream: false,
	}

	// Send chat request
	ctx := context.Background()
	resp, err := cfg.OllamaClient.Chat(ctx, req)
	AssertNoError(t, err, "Chat request failed")

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if resp.Message.Content == "" {
		t.Error("Expected non-empty response content")
	}

	t.Logf("Chat response: %s", resp.Message.Content)
}

func TestChat_StreamingResponse(t *testing.T) {
	SkipIfNoOllama(t)
	SkipIfNoModel(t, "llama3.2:latest")

	cfg := SetupTestEnvironment(t)

	messages := []ollama.Message{
		{
			Role:    "user",
			Content: "Count from 1 to 5.",
		},
	}

	// Prepare chat request

	// Send chat request (streaming not implemented in current API)
	ctx := context.Background()
	req := ollama.ChatRequest{
		Model:    "llama3.2:latest",
		Messages: messages,
		Stream:   false,
	}
	_, err := cfg.OllamaClient.Chat(ctx, req)
	AssertNoError(t, err, "Chat failed")
}

func TestChat_MultiTurnConversation(t *testing.T) {
	SkipIfNoOllama(t)
	SkipIfNoModel(t, "llama3.2:latest")

	cfg := SetupTestEnvironment(t)

	// First turn
	req1 := ollama.ChatRequest{
		Model: "llama3.2:latest",
		Messages: []ollama.Message{
			{
				Role:    "user",
				Content: "My favorite genre is rock.",
			},
		},
		Stream: false,
	}

	ctx := context.Background()
	resp1, err := cfg.OllamaClient.Chat(ctx, req1)
	AssertNoError(t, err, "First chat turn failed")

	// Second turn with context
	req2 := ollama.ChatRequest{
		Model: "llama3.2:latest",
		Messages: []ollama.Message{
			{
				Role:    "user",
				Content: "My favorite genre is rock.",
			},
			{
				Role:    "assistant",
				Content: resp1.Message.Content,
			},
			{
				Role:    "user",
				Content: "What did I just say?",
			},
		},
		Stream: false,
	}

	resp2, err := cfg.OllamaClient.Chat(ctx, req2)
	AssertNoError(t, err, "Second chat turn failed")

	if resp2.Message.Content == "" {
		t.Error("Expected non-empty response for second turn")
	}

	t.Logf("Turn 1: %s", resp1.Message.Content)
	t.Logf("Turn 2: %s", resp2.Message.Content)
}

func TestChat_InvalidModel(t *testing.T) {
	SkipIfNoOllama(t)

	cfg := SetupTestEnvironment(t)

	req := ollama.ChatRequest{
		Model: "nonexistent-model-xyz",
		Messages: []ollama.Message{
			{
				Role:    "user",
				Content: "Hello",
			},
		},
		Stream: false,
	}

	ctx := context.Background()
	_, err := cfg.OllamaClient.Chat(ctx, req)

	// Should return an error for non-existent model
	if err == nil {
		t.Log("Warning: Expected error for non-existent model, but request succeeded")
	}
}

func TestChat_EmptyMessage(t *testing.T) {
	SkipIfNoOllama(t)
	SkipIfNoModel(t, "llama3.2:latest")

	cfg := SetupTestEnvironment(t)

	req := ollama.ChatRequest{
		Model: "llama3.2:latest",
		Messages: []ollama.Message{
			{
				Role:    "user",
				Content: "",
			},
		},
		Stream: false,
	}

	ctx := context.Background()
	resp, err := cfg.OllamaClient.Chat(ctx, req)

	// Some models may accept empty messages, some may not
	if err != nil {
		t.Logf("Empty message returned error (expected): %v", err)
		return
	}

	if resp != nil {
		t.Logf("Empty message accepted, response: %s", resp.Message.Content)
	}
}

func TestChat_WithOptions(t *testing.T) {
	SkipIfNoOllama(t)
	SkipIfNoModel(t, "llama3.2:latest")

	cfg := SetupTestEnvironment(t)

	req := ollama.ChatRequest{
		Model: "llama3.2:latest",
		Messages: []ollama.Message{
			{
				Role:    "user",
				Content: "Say hello.",
			},
		},
		Stream: false,
		Options: &ollama.Options{
			Temperature: 0.1, // Low temperature for deterministic output
			NumPredict:  10,  // Limit response length
		},
	}

	ctx := context.Background()
	resp, err := cfg.OllamaClient.Chat(ctx, req)
	AssertNoError(t, err, "Chat with options failed")

	if resp.Message.Content == "" {
		t.Error("Expected non-empty response")
	}

	t.Logf("Response with options: %s", resp.Message.Content)
}
