package anthropic

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestParseStreamingMessageResponse(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected *Message
		hasError bool
	}{
		{
			name: "Valid stream",
			input: `data: {"type":"message_start","message":{"id":"msg_123","role":"assistant","model":"claude-3-sonnet-20240229","usage":{"input_tokens":10}}}

data: {"type":"content_block_start","index":0,"content_block":{"type":"text"}}

data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hello"}}

data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":", world!"}}

data: {"type":"content_block_stop","index":0}

data: {"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":20}}

data: {"type":"message_stop"}
`,
			expected: &Message{
				ID:    "msg_123",
				Role:  "assistant",
				Model: "claude-3-sonnet-20240229",
				Content: []ContentBlock{
					{Type: "text", Text: "Hello, world!"},
				},
				StopReason: "end_turn",
				Usage: Usage{
					InputTokens:  10,
					OutputTokens: 20,
				},
			},
			hasError: false,
		},
		{
			name:     "Invalid JSON",
			input:    "data: {invalid json}",
			expected: nil,
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.input)
			params := &MessageParams{
				StreamFunc: func(ctx context.Context, chunk []byte) error {
					return nil
				},
			}
			result, err := parseStreamingMessageResponse(context.Background(), reader, params)

			if tc.hasError && err == nil {
				t.Errorf("Expected an error, but got none")
			}
			if !tc.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, but got %+v", tc.expected, result)
			}
		})
	}
}

func TestParseStreamEvent(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected map[string]interface{}
		hasError bool
	}{
		{
			name:     "Valid JSON",
			input:    `{"type":"message_start","message":{"id":"msg_123"}}`,
			expected: map[string]interface{}{"type": "message_start", "message": map[string]interface{}{"id": "msg_123"}},
			hasError: false,
		},
		{
			name:     "Invalid JSON",
			input:    `{invalid json}`,
			expected: nil,
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parseStreamEvent(tc.input)

			if tc.hasError && err == nil {
				t.Errorf("Expected an error, but got none")
			}
			if !tc.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, but got %+v", tc.expected, result)
			}
		})
	}
}

func TestProcessStreamEvent(t *testing.T) {
	testCases := []struct {
		name     string
		event    map[string]interface{}
		response Message
		expected Message
		hasError bool
	}{
		{
			name: "Message Start Event",
			event: map[string]interface{}{
				"type": "message_start",
				"message": map[string]interface{}{
					"id":    "msg_123",
					"role":  "assistant",
					"model": "claude-3-sonnet-20240229",
					"usage": map[string]interface{}{
						"input_tokens": float64(10),
					},
				},
			},
			response: Message{},
			expected: Message{
				ID:    "msg_123",
				Role:  "assistant",
				Model: "claude-3-sonnet-20240229",
				Usage: Usage{InputTokens: 10},
			},
			hasError: false,
		},
		{
			name: "Content Block Start Event",
			event: map[string]interface{}{
				"type":  "content_block_start",
				"index": float64(0),
				"content_block": map[string]interface{}{
					"type": "text",
				},
			},
			response: Message{},
			expected: Message{
				Content: []ContentBlock{{Type: "text"}},
			},
			hasError: false,
		},
		{
			name: "Content Block Delta Event",
			event: map[string]interface{}{
				"type":  "content_block_delta",
				"index": float64(0),
				"delta": map[string]interface{}{
					"type": "text_delta",
					"text": "Hello",
				},
			},
			response: Message{Content: []ContentBlock{{Type: "text"}}},
			expected: Message{Content: []ContentBlock{{Type: "text", Text: "Hello"}}},
			hasError: false,
		},
		{
			name: "Message Delta Event",
			event: map[string]interface{}{
				"type": "message_delta",
				"delta": map[string]interface{}{
					"stop_reason": "end_turn",
				},
				"usage": map[string]interface{}{
					"output_tokens": float64(20),
				},
			},
			response: Message{},
			expected: Message{
				StopReason: "end_turn",
				Usage:      Usage{OutputTokens: 20},
			},
			hasError: false,
		},
		{
			name: "Unknown Event Type",
			event: map[string]interface{}{
				"type": "unknown_event",
			},
			response: Message{},
			expected: Message{},
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			payload := &MessageParams{
				StreamFunc: func(ctx context.Context, chunk []byte) error {
					return nil
				},
			}
			eventChan := make(chan MessageEvent, 1)

			result, err := processStreamEvent(ctx, tc.event, payload, tc.response, eventChan)

			if tc.hasError && err == nil {
				t.Errorf("Expected an error, but got none")
			}
			if !tc.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, but got %+v", tc.expected, result)
			}
		})
	}
}

func TestHandleMessageStartEvent(t *testing.T) {
	testCases := []struct {
		name     string
		event    map[string]interface{}
		response Message
		expected Message
		hasError bool
	}{
		{
			name: "Valid Message Start Event",
			event: map[string]interface{}{
				"message": map[string]interface{}{
					"id":    "msg_123",
					"model": "claude-3-sonnet-20240229",
					"role":  "assistant",
					"type":  "message",
					"usage": map[string]interface{}{
						"input_tokens": float64(10),
					},
				},
			},
			response: Message{},
			expected: Message{
				ID:    "msg_123",
				Model: "claude-3-sonnet-20240229",
				Role:  "assistant",
				Type:  "message",
				Usage: Usage{InputTokens: 10},
			},
			hasError: false,
		},
		{
			name: "Invalid Message Field",
			event: map[string]interface{}{
				"message": "invalid",
			},
			response: Message{},
			expected: Message{},
			hasError: true,
		},
		{
			name: "Invalid Usage Field",
			event: map[string]interface{}{
				"message": map[string]interface{}{
					"usage": "invalid",
				},
			},
			response: Message{},
			expected: Message{},
			hasError: true,
		},
		{
			name: "Invalid Input Tokens Field",
			event: map[string]interface{}{
				"message": map[string]interface{}{
					"usage": map[string]interface{}{
						"input_tokens": "invalid",
					},
				},
			},
			response: Message{},
			expected: Message{},
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := handleMessageStartEvent(tc.event, tc.response)

			if tc.hasError && err == nil {
				t.Errorf("Expected an error, but got none")
			}
			if !tc.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, but got %+v", tc.expected, result)
			}
		})
	}
}

func TestHandleContentBlockStartEvent(t *testing.T) {
	testCases := []struct {
		name     string
		event    map[string]interface{}
		response Message
		expected Message
		hasError bool
	}{
		{
			name: "Valid Content Block Start Event",
			event: map[string]interface{}{
				"index": float64(0),
				"content_block": map[string]interface{}{
					"type": "text",
				},
			},
			response: Message{},
			expected: Message{
				Content: []ContentBlock{{Type: "text"}},
			},
			hasError: false,
		},
		{
			name: "Invalid Index Field",
			event: map[string]interface{}{
				"index": "invalid",
			},
			response: Message{},
			expected: Message{},
			hasError: true,
		},
		{
			name: "Invalid Content Block Field",
			event: map[string]interface{}{
				"index":         float64(0),
				"content_block": "invalid",
			},
			response: Message{},
			expected: Message{},
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := handleContentBlockStartEvent(tc.event, tc.response)

			if tc.hasError && err == nil {
				t.Errorf("Expected an error, but got none")
			}
			if !tc.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, but got %+v", tc.expected, result)
			}
		})
	}
}

func TestHandleContentBlockDeltaEvent(t *testing.T) {
	testCases := []struct {
		name     string
		event    map[string]interface{}
		payload  *MessageParams
		response Message
		expected Message
		hasError bool
	}{
		{
			name: "Valid Content Block Delta Event",
			event: map[string]interface{}{
				"index": float64(0),
				"delta": map[string]interface{}{
					"type": "text_delta",
					"text": "Hello",
				},
			},
			payload: &MessageParams{
				StreamFunc: func(ctx context.Context, chunk []byte) error {
					return nil
				},
			},
			response: Message{Content: []ContentBlock{{Type: "text"}}},
			expected: Message{Content: []ContentBlock{{Type: "text", Text: "Hello"}}},
			hasError: false,
		},
		{
			name: "Invalid Index Field",
			event: map[string]interface{}{
				"index": "invalid",
			},
			payload:  &MessageParams{},
			response: Message{},
			expected: Message{},
			hasError: true,
		},
		{
			name: "Invalid Delta Field",
			event: map[string]interface{}{
				"index": float64(0),
				"delta": "invalid",
			},
			payload:  &MessageParams{},
			response: Message{},
			expected: Message{},
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := handleContentBlockDeltaEvent(context.Background(), tc.event, tc.payload, tc.response)

			if tc.hasError && err == nil {
				t.Errorf("Expected an error, but got none")
			}
			if !tc.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, but got %+v", tc.expected, result)
			}
		})
	}
}

func TestHandleMessageDeltaEvent(t *testing.T) {
	testCases := []struct {
		name     string
		event    map[string]interface{}
		response Message
		expected Message
		hasError bool
	}{
		{
			name: "Valid Message Delta Event",
			event: map[string]interface{}{
				"delta": map[string]interface{}{
					"stop_reason": "end_turn",
				},
				"usage": map[string]interface{}{
					"output_tokens": float64(20),
				},
			},
			response: Message{},
			expected: Message{
				StopReason: "end_turn",
				Usage:      Usage{OutputTokens: 20},
			},
			hasError: false,
		},
		{
			name: "Invalid Delta Field",
			event: map[string]interface{}{
				"delta": "invalid",
			},
			response: Message{},
			expected: Message{},
			hasError: true,
		},
		{
			name: "Invalid Usage Field",
			event: map[string]interface{}{
				"delta": map[string]interface{}{},
				"usage": "invalid",
			},
			response: Message{},
			expected: Message{},
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := handleMessageDeltaEvent(tc.event, tc.response)

			if tc.hasError && err == nil {
				t.Errorf("Expected an error, but got none")
			}
			if !tc.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, but got %+v", tc.expected, result)
			}
		})
	}
}


func TestGetString(t *testing.T) {
	testCases := []struct {
		name     string
		input    map[string]interface{}
		key      string
		expected string
	}{
		{
			name:     "Valid string value",
			input:    map[string]interface{}{"key": "value"},
			key:      "key",
			expected: "value",
		},
		{
			name:     "Non-existent key",
			input:    map[string]interface{}{"key": "value"},
			key:      "nonexistent",
			expected: "",
		},
		{
			name:     "Non-string value",
			input:    map[string]interface{}{"key": 123},
			key:      "key",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getString(tc.input, tc.key)
			if result != tc.expected {
				t.Errorf("Expected %s, but got %s", tc.expected, result)
			}
		})
	}
}

func TestParseStreamingMessageResponseWithInvalidScanner(t *testing.T) {
	invalidReader := &errorReader{}
	params := &MessageParams{
		StreamFunc: func(ctx context.Context, chunk []byte) error {
			return nil
		},
	}
	_, err := parseStreamingMessageResponse(context.Background(), invalidReader, params)
	if err == nil {
		t.Errorf("Expected an error, but got none")
	}
}

// errorReader is a custom io.Reader that always returns an error.
type errorReader struct{}

func (er *errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("error reading")
}

func TestHandleContentBlockDeltaEventWithStreamFuncError(t *testing.T) {
	event := map[string]interface{}{
		"index": float64(0),
		"delta": map[string]interface{}{
			"type": "text_delta",
			"text": "Hello",
		},
	}
	payload := &MessageParams{
		StreamFunc: func(ctx context.Context, chunk []byte) error {
			return fmt.Errorf("streaming error")
		},
	}
	response := Message{Content: []ContentBlock{{Type: "text"}}}

	_, err := handleContentBlockDeltaEvent(context.Background(), event, payload, response)
	if err == nil {
		t.Errorf("Expected an error, but got none")
	}
}

func TestProcessStreamEventWithUnknownType(t *testing.T) {
	event := map[string]interface{}{
		"type": "unknown_event",
	}
	payload := &MessageParams{}
	response := Message{}
	eventChan := make(chan MessageEvent, 1)

	result, err := processStreamEvent(context.Background(), event, payload, response, eventChan)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !reflect.DeepEqual(result, response) {
		t.Errorf("Expected %+v, but got %+v", response, result)
	}
}

func TestParseStreamingMessageResponseWithMessageStopEvent(t *testing.T) {
	input := `data: {"type":"message_stop"}
`
	reader := strings.NewReader(input)
	params := &MessageParams{
		StreamFunc: func(ctx context.Context, chunk []byte) error {
			return nil
		},
	}
	result, err := parseStreamingMessageResponse(context.Background(), reader, params)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result == nil {
		t.Errorf("Expected a non-nil result, but got nil")
	}
}

func TestParseStreamingMessageResponseWithPingEvent(t *testing.T) {
	input := `data: {"type":"ping"}
`
	reader := strings.NewReader(input)
	params := &MessageParams{
		StreamFunc: func(ctx context.Context, chunk []byte) error {
			return nil
		},
	}
	result, err := parseStreamingMessageResponse(context.Background(), reader, params)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != nil {
		t.Errorf("Expected a nil result, but got %+v", result)
	}
}

func TestHandleContentBlockDeltaEventWithNewContentBlock(t *testing.T) {
	event := map[string]interface{}{
		"index": float64(1),
		"delta": map[string]interface{}{
			"type": "text_delta",
			"text": "New content",
		},
	}
	payload := &MessageParams{
		StreamFunc: func(ctx context.Context, chunk []byte) error {
			return nil
		},
	}
	response := Message{Content: []ContentBlock{{Type: "text", Text: "Existing content"}}}

	result, err := handleContentBlockDeltaEvent(context.Background(), event, payload, response)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := Message{Content: []ContentBlock{
		{Type: "text", Text: "Existing content"},
		{Type: "text", Text: "New content"},
	}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, but got %+v", expected, result)
	}
}