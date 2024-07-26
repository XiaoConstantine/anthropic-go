package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const messagesEndpoint = "/messages"

// CreateMessage sends a request to create a new message.
func (s *Client) Create(ctx context.Context, params *MessageParams) (*Message, error) {
	url := s.baseURL + messagesEndpoint

	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", s.APIKey)
	req.Header.Set("anthropic-version", s.APIVersion)

	// Set Accep header based on whether streaming is requested
	if params.IsStreaming() {
		req.Header.Set("Accept", "text/event-stream")
	} else {
		req.Header.Set("Accept", "application/json")
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if params.IsStreaming() {
		return parseStreamingMessageResponse(ctx, resp.Body, params)
	}

	var message Message
	err = json.NewDecoder(resp.Body).Decode(&message)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &message, nil
}
