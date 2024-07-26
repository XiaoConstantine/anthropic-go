package anthropic

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	defaultBaseURL    = "https://api.anthropic.com/v1"
	defaultAPIVersion = "2023-06-01"
	defaultTimeout    = 120 * time.Second
	envAPIKey         = "ANTHROPIC_API_KEY"
)

// Client is the main struct for interacting with the Anthropic API.
type Client struct {
	baseURL    string
	APIKey     string
	APIVersion string
	httpClient *http.Client
}

// ClientOption is a function that modifies a Client.
type ClientOption func(*Client) error

// NewClient creates a new Anthropic client with the given API key and optional configuration options.
func NewClient(opts ...ClientOption) (*Client, error) {
	client := &Client{
		baseURL:    defaultBaseURL,
		APIVersion: defaultAPIVersion,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}

	// Apply any custom options
	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	// Check if API key is set
	if client.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	return client, nil
}

// WithBaseURL sets a custom base URL for the client.
func WithBaseURL(url string) ClientOption {
	return func(c *Client) error {
		c.baseURL = url
		return nil
	}
}

// WithAPIKey sets the API key for the client.
func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) error {
		if apiKey == "" {
			apiKey = os.Getenv(envAPIKey)
		}
		if apiKey == "" {
			return fmt.Errorf("no API key provided and %s environment variable is not set", envAPIKey)
		}
		c.APIKey = apiKey
		return nil
	}
}

// WithAPIVersion sets a custom API version for the client.
func WithAPIVersion(version string) ClientOption {
	return func(c *Client) error {
		c.APIVersion = version
		return nil

	}
}

// WithHTTPClient sets a custom HTTP client for the API client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) error {
		c.httpClient = httpClient
		return nil
	}
}

// WithTimeout sets a custom timeout for the HTTP client.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) error {
		c.httpClient.Timeout = timeout
		return nil
	}
}

// SetAPIKey updates the API key for the client.
func (c *Client) SetAPIKey(apiKey string) {
	c.APIKey = apiKey
}

// Models returns a new ModelsService.
func (c *Client) Models() *ModelsService {
	return &ModelsService{client: c}
}

// Messages returns a new MessagesService.
func (c *Client) Messages() *MessagesService {
	return &MessagesService{client: c}
}

// ModelsService handles operations related to models.
type ModelsService struct {
	client *Client
}

// MessagesService handles operations related to messages.
type MessagesService struct {
	client *Client
}

// List retrieves a list of available models.
func (s *ModelsService) List() ([]Model, error) {
	return []Model{
		{ID: ModelHaiku, Name: "Claude 3 Haiku"},
		{ID: ModelSonnet, Name: "Claude 3 Sonnet"},
		{ID: ModelOpus, Name: "Claude 3 Opus"},
	}, nil
}

// GetModelID returns the ModelID for a given name.
func GetModelID(name string) (ModelID, bool) {
	switch name {
	case "HAIKU":
		return ModelHaiku, true
	case "SONNET":
		return ModelSonnet, true
	case "OPUS":
		return ModelOpus, true
	default:
		return "", false
	}
}

// Create sends a request to create a new message.
func (s *MessagesService) Create(ctx context.Context, params *MessageParams) (*Message, error) {
	return s.client.Create(ctx, params)
}
