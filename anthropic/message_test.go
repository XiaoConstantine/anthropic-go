package anthropic

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMessagesService_Create(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != "POST" {
			t.Errorf("Expected 'POST' request, got '%s'", r.Method)
		}

		// Check request headers
		if r.Header.Get("X-API-Key") != "test-key" {
			t.Errorf("Expected API Key header 'test-key', got '%s'", r.Header.Get("X-API-Key"))
		}

		// Prepare response
		response := Message{
			ID:   "msg_123",
			Role: "assistant",
			Content: []ContentBlock{
				{Type: "text", Text: "Hello, how can I help you?"},
			},
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			return
		}
	}))
	defer server.Close()

	// Create client with mock server URL
	client, _ := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	// Create message
	params := &MessageParams{
		Model: string(ModelSonnet),
		Messages: []MessageParam{
			{
				Role: "user",
				Content: []ContentBlock{
					{Type: "text", Text: "Hello"},
				},
			},
		},
	}

	message, err := client.Messages().Create(context.Background(), params)
	if err != nil {
		t.Fatalf("Failed to create message: %v", err)
	}

	// Check response
	if message.ID != "msg_123" {
		t.Errorf("Expected message ID 'msg_123', got '%s'", message.ID)
	}
	if message.Role != "assistant" {
		t.Errorf("Expected role 'assistant', got '%s'", message.Role)
	}
	if len(message.Content) != 1 || message.Content[0].Text != "Hello, how can I help you?" {
		t.Errorf("Unexpected message content: %+v", message.Content)
	}
}

func TestMessagesService_CreateStreaming(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and headers
		if r.Method != "POST" {
			t.Errorf("Expected 'POST' request, got '%s'", r.Method)
		}
		if r.Header.Get("Accept") != "text/event-stream" {
			t.Errorf("Expected Accept header 'text/event-stream', got '%s'", r.Header.Get("Accept"))
		}

		// Send streaming response
		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Fatal("Expected http.ResponseWriter to be an http.Flusher")
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		events := []string{
			`{"type":"message_start","message":{"id":"msg_123","role":"assistant","model":"claude-3-sonnet-20240229","usage":{"input_tokens":10}}}`,
			`{"type":"content_block_start","index":0,"content_block":{"type":"text"}}`,
			`{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hello"}}`,
			`{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":", how can I help you?"}}`,
			`{"type":"content_block_stop","index":0}`,
			`{"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":20}}`,
			`{"type":"message_stop"}`,
		}

		for _, event := range events {
			_, err := w.Write([]byte("data: " + event + "\n\n"))
			if err != nil {
				t.Fatalf("Failed to write event: %v", err)
			}
			flusher.Flush()
		}
	}))
	defer server.Close()

	// Create client with mock server URL
	client, _ := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	// Create streaming message
	params := &MessageParams{
		Model: string(ModelSonnet),
		Messages: []MessageParam{
			{
				Role: "user",
				Content: []ContentBlock{
					{Type: "text", Text: "Hello"},
				},
			},
		},
		StreamFunc: func(ctx context.Context, chunk []byte) error {
			// In a real test, you might want to accumulate these chunks
			// and check the final result
			return nil
		},
	}

	message, err := client.Messages().Create(context.Background(), params)
	if err != nil {
		t.Fatalf("Failed to create streaming message: %v", err)
	}

	// Check response
	if message.ID != "msg_123" {
		t.Errorf("Expected message ID 'msg_123', got '%s'", message.ID)
	}
	if message.Role != "assistant" {
		t.Errorf("Expected role 'assistant', got '%s'", message.Role)
	}
	if len(message.Content) != 1 || message.Content[0].Text != "Hello, how can I help you?" {
		t.Errorf("Unexpected message content: %+v", message.Content)
	}
}

func TestMessagesService_CreateWithTools(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected 'POST' request, got '%s'", r.Method)
		}

		if r.Header.Get("X-API-Key") != "test-key" {
			t.Errorf("Expected API Key header 'test-key', got '%s'", r.Header.Get("X-API-Key"))
		}

		var requestBody MessageParams
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if len(requestBody.Tools) != 1 || requestBody.Tools[0].Name != "get_stock_price" {
			t.Errorf("Expected one tool named 'get_stock_price', got: %+v", requestBody.Tools)
		}

		response := Message{
			ID:   "msg_123",
			Role: "assistant",
			Content: []ContentBlock{
				{
					Type: "text",
					Text: "Certainly! I'll check the current S&P 500 price for you.",
				},
				{
					Type: "tool_call",
					ToolCall: &ToolCall{
						ID:   "call_123",
						Type: "function",
						Name: "get_stock_price",
						Input: json.RawMessage(`{
                            "ticker": "^GSPC"
                        }`),
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			return
		}
	}))
	defer server.Close()

	client, _ := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	params := &MessageParams{
		Model: string(ModelSonnet),
		Messages: []MessageParam{
			{
				Role: "user",
				Content: []ContentBlock{
					{Type: "text", Text: "What's the S&P 500 at today?"},
				},
			},
		},
		Tools: []Tool{
			{
				Name:        "get_stock_price",
				Description: "Get the current stock price for a given ticker symbol.",
				InputSchema: InputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"ticker": map[string]interface{}{
							"type":        "string",
							"description": "The stock ticker symbol, e.g. AAPL for Apple Inc.",
						},
					},
				},
			},
		},
	}

	message, err := client.Messages().Create(context.Background(), params)
	if err != nil {
		t.Fatalf("Failed to create message: %v", err)
	}

	if message.ID != "msg_123" {
		t.Errorf("Expected message ID 'msg_123', got '%s'", message.ID)
	}
	if message.Role != "assistant" {
		t.Errorf("Expected role 'assistant', got '%s'", message.Role)
	}
	if len(message.Content) != 2 {
		t.Fatalf("Expected 2 content blocks, got %d", len(message.Content))
	}
	if message.Content[0].Type != "text" || message.Content[0].Text != "Certainly! I'll check the current S&P 500 price for you." {
		t.Errorf("Unexpected text content: %+v", message.Content[0])
	}
	if message.Content[1].Type != "tool_call" || message.Content[1].ToolCall == nil {
		t.Errorf("Expected tool_call content, got: %+v", message.Content[1])
	}
	if message.Content[1].ToolCall.Name != "get_stock_price" {
		t.Errorf("Expected tool call name 'get_stock_price', got '%s'", message.Content[1].ToolCall.Name)
	}
	var input map[string]string
	err = json.Unmarshal(message.Content[1].ToolCall.Input, &input)
	if err != nil {
		t.Fatalf("Failed to unmarshal tool call input: %v", err)
	}
	if input["ticker"] != "^GSPC" {
		t.Errorf("Expected ticker '^GSPC', got '%s'", input["ticker"])
	}
}
