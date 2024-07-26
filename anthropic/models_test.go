package anthropic

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestModelMarshalJSON(t *testing.T) {
	model := Model{
		ID:   ModelHaiku,
		Name: "Claude 3 Haiku",
	}

	jsonData, err := json.Marshal(model)
	if err != nil {
		t.Fatalf("Failed to marshal Model: %v", err)
	}

	expected := `{"id":"claude-3-haiku-20240307","name":"Claude 3 Haiku"}`
	if string(jsonData) != expected {
		t.Errorf("Expected JSON %s, got %s", expected, string(jsonData))
	}
}

func TestMessageParamsMarshalJSON(t *testing.T) {
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
		StreamFunc: func(ctx context.Context, chunk []byte) error { return nil },
	}

	jsonData, err := json.Marshal(params)

	if err != nil {
		t.Fatalf("Failed to marshal MessageParams: %v", err)
	}
	if !params.IsStreaming() {
		t.Errorf("Expected IsStreaming() to return true, got false")
	}
	// Check if "stream":true is in the JSON
	if !strings.Contains(string(jsonData), `"stream":true`) {
		t.Errorf("Expected JSON to contain \"stream\":true, got %s", string(jsonData))
	}
}
