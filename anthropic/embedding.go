package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// EmbeddingParams represents the parameters for creating embeddings
type EmbeddingParams struct {
	Model    string                 `json:"model"`
	Input    []string               `json:"input"`
	Encoding string                 `json:"encoding,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// EmbeddingResponse represents the response from creating embeddings
type EmbeddingResponse struct {
	ID         string      `json:"id"`
	Model      string      `json:"model"`
	Type       string      `json:"type"`
	Embeddings [][]float64 `json:"embeddings"`
	Usage      Usage       `json:"usage"`
}

// Create generates embeddings for the provided input texts
func (s *EmbeddingsService) Create(ctx context.Context, params *EmbeddingParams) (*EmbeddingResponse, error) {
	url := s.client.baseURL + embeddingsEndpoint

	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", s.client.APIKey)
	req.Header.Set("anthropic-version", s.client.APIVersion)

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var response EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}
