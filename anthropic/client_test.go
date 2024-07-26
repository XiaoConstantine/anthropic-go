package anthropic

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
    client, err := NewClient(WithAPIKey("test-key"))
    if err != nil {
        t.Fatalf("Failed to create client: %v", err)
    }

    if client.APIKey != "test-key" {
        t.Errorf("Expected API key to be 'test-key', got '%s'", client.APIKey)
    }

    if client.baseURL != defaultBaseURL {
        t.Errorf("Expected base URL to be '%s', got '%s'", defaultBaseURL, client.baseURL)
    }
}

func TestClientOptions(t *testing.T) {
    customURL := "https://custom.anthropic.com"
    customTimeout := 30 * time.Second

    client, err := NewClient(
        WithAPIKey("test-key"),
        WithBaseURL(customURL),
        WithTimeout(customTimeout),
    )
    if err != nil {
        t.Fatalf("Failed to create client with options: %v", err)
    }

    if client.baseURL != customURL {
        t.Errorf("Expected base URL to be '%s', got '%s'", customURL, client.baseURL)
    }

    if client.httpClient.Timeout != customTimeout {
        t.Errorf("Expected timeout to be %v, got %v", customTimeout, client.httpClient.Timeout)
    }
}

func TestModelsService_List(t *testing.T) {
    client, _ := NewClient(WithAPIKey("test-key"))
    models, err := client.Models().List()
    if err != nil {
        t.Fatalf("Failed to list models: %v", err)
    }

    expectedModels := []Model{
        {ID: ModelHaiku, Name: "Claude 3 Haiku"},
        {ID: ModelSonnet, Name: "Claude 3 Sonnet"},
        {ID: ModelOpus, Name: "Claude 3 Opus"},
    }

    if len(models) != len(expectedModels) {
        t.Fatalf("Expected %d models, got %d", len(expectedModels), len(models))
    }

    for i, model := range models {
        if model != expectedModels[i] {
            t.Errorf("Expected model %+v, got %+v", expectedModels[i], model)
        }
    }
}
