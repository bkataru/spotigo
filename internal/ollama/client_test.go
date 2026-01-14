package ollama

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("http://localhost:11434", 30*time.Second)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.baseURL != "http://localhost:11434" {
		t.Errorf("expected baseURL 'http://localhost:11434', got '%s'", client.baseURL)
	}

	if client.httpClient == nil {
		t.Error("httpClient should be initialized")
	}

	if client.httpClient.Timeout != 30*time.Second {
		t.Errorf("expected timeout 30s, got %v", client.httpClient.Timeout)
	}
}

func TestClient_Chat(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		// Verify request body
		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}

		if req.Model != "test-model" {
			t.Errorf("expected model 'test-model', got '%s'", req.Model)
		}

		// Send response
		resp := ChatResponse{
			Model: "test-model",
			Message: Message{
				Role:    "assistant",
				Content: "Hello! I'm a test response.",
			},
			Done: true,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, 10*time.Second)

	req := ChatRequest{
		Model: "test-model",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
	}

	resp, err := client.Chat(context.Background(), req)
	if err != nil {
		t.Fatalf("Chat failed: %v", err)
	}

	if resp.Model != "test-model" {
		t.Errorf("expected model 'test-model', got '%s'", resp.Model)
	}
	if resp.Message.Role != "assistant" {
		t.Errorf("expected role 'assistant', got '%s'", resp.Message.Role)
	}
	if resp.Message.Content != "Hello! I'm a test response." {
		t.Errorf("unexpected content: %s", resp.Message.Content)
	}
}

func TestClient_Chat_WithOptions(t *testing.T) {
	var receivedReq ChatRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedReq)

		resp := ChatResponse{
			Model:   "test-model",
			Message: Message{Role: "assistant", Content: "Response"},
			Done:    true,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, 10*time.Second)

	req := ChatRequest{
		Model: "test-model",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
		Options: &Options{
			Temperature: 0.7,
			NumPredict:  100,
			TopK:        40,
			TopP:        0.9,
		},
	}

	_, err := client.Chat(context.Background(), req)
	if err != nil {
		t.Fatalf("Chat failed: %v", err)
	}

	if receivedReq.Options == nil {
		t.Fatal("Options should be sent")
	}
	if receivedReq.Options.Temperature != 0.7 {
		t.Errorf("expected Temperature 0.7, got %f", receivedReq.Options.Temperature)
	}
	if receivedReq.Options.NumPredict != 100 {
		t.Errorf("expected NumPredict 100, got %d", receivedReq.Options.NumPredict)
	}
}

func TestClient_Chat_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := NewClient(server.URL, 10*time.Second)

	req := ChatRequest{
		Model:    "test-model",
		Messages: []Message{{Role: "user", Content: "Hello"}},
	}

	_, err := client.Chat(context.Background(), req)
	if err == nil {
		t.Error("Chat should fail on 500 error")
	}
}

func TestClient_Embed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/embed" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		// Verify request
		var req EmbedRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Model != "nomic-embed-text" {
			t.Errorf("expected model 'nomic-embed-text', got '%s'", req.Model)
		}
		if req.Input != "Test text for embedding" {
			t.Errorf("unexpected input: %s", req.Input)
		}

		// Send response
		resp := EmbedResponse{
			Model:      "nomic-embed-text",
			Embeddings: [][]float64{{0.1, 0.2, 0.3, 0.4, 0.5}},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, 10*time.Second)

	embedding, err := client.Embed(context.Background(), "nomic-embed-text", "Test text for embedding")
	if err != nil {
		t.Fatalf("Embed failed: %v", err)
	}

	if len(embedding) != 5 {
		t.Errorf("expected 5 dimensions, got %d", len(embedding))
	}
	if embedding[0] != 0.1 {
		t.Errorf("expected first value 0.1, got %f", embedding[0])
	}
}

func TestClient_Embed_NoEmbeddings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := EmbedResponse{
			Model:      "test-model",
			Embeddings: [][]float64{}, // Empty embeddings
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, 10*time.Second)

	_, err := client.Embed(context.Background(), "test-model", "test")
	if err == nil {
		t.Error("Embed should fail when no embeddings returned")
	}
}

func TestClient_ListModels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/tags" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}

		resp := TagsResponse{
			Models: []ModelInfo{
				{
					Name:       "granite4:1b",
					ModifiedAt: time.Now(),
					Size:       1000000000,
				},
				{
					Name:       "nomic-embed-text",
					ModifiedAt: time.Now(),
					Size:       500000000,
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, 10*time.Second)

	models, err := client.ListModels(context.Background())
	if err != nil {
		t.Fatalf("ListModels failed: %v", err)
	}

	if len(models) != 2 {
		t.Errorf("expected 2 models, got %d", len(models))
	}
	if models[0].Name != "granite4:1b" {
		t.Errorf("expected first model 'granite4:1b', got '%s'", models[0].Name)
	}
}

func TestClient_Ping(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/tags" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(TagsResponse{})
	}))
	defer server.Close()

	client := NewClient(server.URL, 10*time.Second)

	err := client.Ping(context.Background())
	if err != nil {
		t.Errorf("Ping should succeed, got error: %v", err)
	}
}

func TestClient_Ping_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewClient(server.URL, 10*time.Second)

	err := client.Ping(context.Background())
	if err == nil {
		t.Error("Ping should fail on 503 error")
	}
}

func TestClient_Ping_Unreachable(t *testing.T) {
	// Use an unreachable address
	client := NewClient("http://localhost:99999", 1*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := client.Ping(ctx)
	if err == nil {
		t.Error("Ping should fail for unreachable host")
	}
}

func TestMessage_Structure(t *testing.T) {
	msg := Message{
		Role:    "user",
		Content: "Hello, world!",
	}

	if msg.Role != "user" {
		t.Errorf("expected role 'user', got '%s'", msg.Role)
	}
	if msg.Content != "Hello, world!" {
		t.Errorf("unexpected content: %s", msg.Content)
	}
}

func TestOptions_Structure(t *testing.T) {
	opts := Options{
		Temperature: 0.8,
		NumPredict:  2048,
	}

	if opts.Temperature != 0.8 {
		t.Errorf("expected Temperature 0.8, got %f", opts.Temperature)
	}
	if opts.NumPredict != 2048 {
		t.Errorf("expected NumPredict 2048, got %d", opts.NumPredict)
	}
}

func TestChatRequest_JSON(t *testing.T) {
	req := ChatRequest{
		Model: "test-model",
		Messages: []Message{
			{Role: "system", Content: "You are helpful."},
			{Role: "user", Content: "Hello"},
		},
		Stream: false,
		Options: &Options{
			Temperature: 0.7,
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded ChatRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Model != req.Model {
		t.Errorf("model mismatch")
	}
	if len(decoded.Messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(decoded.Messages))
	}
}
