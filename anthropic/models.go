package anthropic

import (
	"context"
	"encoding/json"
	"time"
)

// ContentBlock represents a block of content in a message.
type ContentBlock struct {
	Type  string `json:"type"`
	Text  string `json:"text,omitempty"`
	Image *Image `json:"image,omitempty"`
}

// Image represents an image in a content block.
type Image struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

// Message represents a complete message from the API.
type Message struct {
	ID           string         `json:"id"`
	Type         string         `json:"type"`
	Role         string         `json:"role"`
	Content      []ContentBlock `json:"content"`
	Model        string         `json:"model"`
	StopReason   string         `json:"stop_reason"`
	StopSequence string         `json:"stop_sequence"`
	Usage        Usage          `json:"usage"`
	CreatedAt    time.Time      `json:"created_at"`
}

// Usage represents the token usage information.
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// MessageParams represents the parameters for creating a message.
type MessageParams struct {
	Model         string                              `json:"model"`
	Messages      []MessageParam                      `json:"messages"`
	MaxTokens     int                                 `json:"max_tokens,omitempty"`
	Temperature   float64                             `json:"temperature,omitempty"`
	TopP          float64                             `json:"top_p,omitempty"`
	TopK          int                                 `json:"top_k,omitempty"`
	StopSequences []string                            `json:"stop_sequences,omitempty"`
	Metadata      map[string]interface{}              `json:"metadata,omitempty"`
	StreamFunc    func(context.Context, []byte) error `json:"-"`
}

func (p *MessageParams) IsStreaming() bool {
	return p.StreamFunc != nil
}

func (p *MessageParams) MarshalJSON() ([]byte, error) {
	type Alias MessageParams
	return json.Marshal(&struct {
		*Alias
		Stream bool `json:"stream"`
	}{
		Alias:  (*Alias)(p),
		Stream: p.IsStreaming(),
	})
}

// MessageParam represents a single message in the conversation history.
type MessageParam struct {
	Role    string         `json:"role"`
	Content []ContentBlock `json:"content"`
}

// TextBlock is a convenience type for creating text content blocks.
type TextBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ImageBlock is a convenience type for creating image content blocks.
type ImageBlock struct {
	Type   string `json:"type"`
	Source struct {
		Type      string `json:"type"`
		MediaType string `json:"media_type"`
		Data      string `json:"data"`
	} `json:"source"`
}

// Tool represents a tool that can be used by the model.
type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

// ToolFunction represents the function definition of a tool.
type ToolFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// ToolCall represents a call to a tool made by the model.
type ToolCall struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"`
	Function ToolCallFunction `json:"function"`
}

// ToolCallFunction represents the function call made by the model.
type ToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ToolOutput represents the output of a tool call.
type ToolOutput struct {
	ToolCallID string `json:"tool_call_id"`
	Output     string `json:"output"`
}

// MessageStreamEvent represents an event in the message stream.
type MessageStreamEvent struct {
	Type    string `json:"type"`
	Message struct {
		ID           string         `json:"id"`
		Type         string         `json:"type"`
		Role         string         `json:"role"`
		Model        string         `json:"model"`
		Content      []ContentBlock `json:"content"`
		StopReason   *string        `json:"stop_reason"`
		StopSequence *string        `json:"stop_sequence"`
		Usage        struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	} `json:"message,omitempty"`
	Index        int `json:"index,omitempty"`
	ContentBlock struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content_block,omitempty"`
	Delta struct {
		Type         string  `json:"type,omitempty"`
		Text         string  `json:"text,omitempty"`
		StopReason   *string `json:"stop_reason,omitempty"`
		StopSequence *string `json:"stop_sequence,omitempty"`
	} `json:"delta,omitempty"`
	Usage struct {
		OutputTokens int `json:"output_tokens,omitempty"`
	} `json:"usage,omitempty"`
}

// Error represents an error returned by the API.
type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// APIResponse represents a generic API response.
type APIResponse struct {
	Object string `json:"object"`
	Data   []struct {
		ID        string    `json:"id"`
		Object    string    `json:"object"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"data"`
}

type MessageResponsePayload struct {
	Content      []ContentBlock `json:"content"`
	ID           string         `json:"id"`
	Model        string         `json:"model"`
	Role         string         `json:"role"`
	StopReason   string         `json:"stop_reason"`
	StopSequence string         `json:"stop_sequence"`
	Type         string         `json:"type"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

// Model represents an AI model available in the Anthropic API.
type Model struct {
	ID        ModelID `json:"id"`
	Name      string  `json:"name"`
	Created   int64   `json:"created"`
	Owned     bool    `json:"owned"`
	Supported bool    `json:"supported"`
}

// ModelsResponse represents the response from the models endpoint.
type ModelsResponse struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

// ModelID represents the available model IDs.
type ModelID string

const (
	ModelHaiku  ModelID = "claude-3-haiku-20240307"
	ModelSonnet ModelID = "claude-3-sonnet-20240229"
	ModelOpus   ModelID = "claude-3-opus-20240229"
)
