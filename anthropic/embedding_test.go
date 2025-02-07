package anthropic

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEmbeddingsService_Create(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != "POST" {
			t.Errorf("Expected 'POST' request, got '%s'", r.Method)
		}

		// Check request path
		if r.URL.Path != "/embeddings" {
			t.Errorf("Expected path '/embeddings', got '%s'", r.URL.Path)
		}

		// Check request headers
		if r.Header.Get("X-API-Key") != "test-key" {
			t.Errorf("Expected API Key header 'test-key', got '%s'", r.Header.Get("X-API-Key"))
		}

		// Check content type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header 'application/json', got '%s'", r.Header.Get("Content-Type"))
		}

		// Prepare mock response
		response := EmbeddingResponse{
			ID:    "emb_123",
			Model: string(ModelClaude3Embedding),
			Type:  "text_embedding",
			Embeddings: [][]float64{
				{0.1, 0.2, 0.3},
				{0.4, 0.5, 0.6},
			},
			Usage: Usage{
				InputTokens:  10,
				OutputTokens: 0,
			},
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Create client with mock server URL
	client, _ := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	// Create embeddings
	params := &EmbeddingParams{
		Model: string(ModelClaude3Embedding),
		Input: []string{
			"Hello, world!",
			"How are you?",
		},
	}

	response, err := client.Embeddings().Create(context.Background(), params)
	if err != nil {
		t.Fatalf("Failed to create embeddings: %v", err)
	}

	// Check response
	if response.ID != "emb_123" {
		t.Errorf("Expected embedding ID 'emb_123', got '%s'", response.ID)
	}
	if response.Model != string(ModelClaude3Embedding) {
		t.Errorf("Expected model '%s', got '%s'", ModelClaude3Embedding, response.Model)
	}
	if len(response.Embeddings) != 2 {
		t.Errorf("Expected 2 embeddings, got %d", len(response.Embeddings))
	}
	if response.Usage.InputTokens != 10 {
		t.Errorf("Expected 10 input tokens, got %d", response.Usage.InputTokens)
	}
}
