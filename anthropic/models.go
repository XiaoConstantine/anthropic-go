package anthropic

import (
	"context"
	"encoding/json"
	"time"
)

// ContentBlock represents a block of content in a message.
type ContentBlock struct {
	Type       string      `json:"type"`
	Text       string      `json:"text,omitempty"`
	Source     *Image      `json:"source,omitempty"`
	ToolCall   *ToolCall   `json:"tool_call,omitempty"`
	ToolOutput *ToolOutput `json:"tool_output,omitempty"`
	Thinking   string      `json:"thinking,omitempty"`
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
	Beta         *BetaMetadata  `json:"beta,omitempty"`
}

// Usage represents the token usage information.
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// ThinkingConfig represents configuration for the thinking mode.
type ThinkingConfig struct {
	Type     string `json:"type"`
	Enabled  bool   `json:"enabled,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
}

// ThinkingBlock represents a thinking block in a message.
type ThinkingBlock struct {
	Type     string `json:"type"`
	Thinking string `json:"thinking"`
}

// RedactedThinkingBlock represents a redacted thinking block.
type RedactedThinkingBlock struct {
	Type string `json:"type"`
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
	Tools         []Tool                              `json:"tools,omitempty"`
	ToolChoice    *ToolChoice                         `json:"tool_choice,omitempty"`
	Thinking      *ThinkingConfig                     `json:"thinking,omitempty"`
}

type BetaMetadata struct {
	CacheControl CacheControl `json:"cache_control,omitempty"`

	Thinking interface{} `json:"thinking,omitempty"`
}

type CacheControl string

const (
	CacheControlEphemeral CacheControl = "ephemeral"
)

type ToolChoice struct {
	Type  string `json:"type"`
	Tool  *Tool  `json:"tool,omitempty"`
	Tools []Tool `json:"tools,omitempty"`
}

const (
	ToolChoiceTypeAuto = "auto"
	ToolChoiceTypeNone = "none"
	ToolChoiceTypeTool = "tool"
)

// ComputerUseTool represents a computer use tool.
type ComputerUseTool struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// BashTool represents a bash tool.
type BashTool struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// TextEditorTool represents a text editor tool.
type TextEditorTool struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// IsStreaming returns true if the MessageParams is configured for streaming.
func (p *MessageParams) IsStreaming() bool {
	return p.StreamFunc != nil
}

// MarshalJSON implements custom JSON marshaling for MessageParams.
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
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"input_schema"`
}

type InputSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
}

// ToolCall represents a call to a tool made by the model.
type ToolCall struct {
	ID    string          `json:"id"`
	Type  string          `json:"type"`
	Name  string          `json:"name"`
	Input json.RawMessage `json:"input"`
}

// ToolOutput represents the output of a tool call.
type ToolOutput struct {
	ToolCallID string `json:"tool_call_id"`
	Output     string `json:"output"`
}

// Error represents an error returned by the API.
type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// MessageEvent represents an event in the message streaming process.
type MessageEvent struct {
	Response *Message
	Err      error
}

// ModelID represents the available model IDs.
type ModelID string

// Model represents an Anthropic AI model.
type Model struct {
	ID   ModelID `json:"id"`
	Name string  `json:"name"`
}

// Constants for available model IDs.
const (
	ModelHaiku            ModelID = "claude-3-haiku-20240307"
	ModelSonnet           ModelID = "claude-3-7-sonnet-20250219"
	ModelOpus             ModelID = "claude-3-opus-20240229"
	ModelClaude3Embedding ModelID = "claude-3-embedding-20240229"
	ModelSonnetOld        ModelID = "claude-3-5-sonnet-20240620"
)
