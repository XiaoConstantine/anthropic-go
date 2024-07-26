Anthropic Go Client
-------------------



Example
------

* List
----

```go
// Create a new client
	client, err := anthropic.NewClient(
		anthropic.WithAPIKey(""), // This will use the ANTHROPIC_API_KEY environment variable
		anthropic.WithTimeout(30*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// List available models - Since there's no public api for this, currently the result is
	// hard coded
	models, _ := client.Models().List()
	fmt.Println("Available models:")
	for _, model := range models {
		fmt.Printf("- %s (%s)\n", model.Name, model.ID)
	}
	fmt.Println()
```

* Message
-------

```go
	// Send a regular message
	message, err := client.Messages().Create(&anthropic.MessageParams{
		Model: string(anthropic.ModelSonnet), // Use the Sonnet model
		Messages: []anthropic.MessageParam{
			{
				Role: "user",
				Content: []anthropic.ContentBlock{
					{Type: "text", Text: "Hello, Claude! What's the capital of France?"},
				},
			},
		},
		MaxTokens: 1000,
	})
	if err != nil {
		log.Fatalf("Failed to create message: %v", err)
	}
	fmt.Println("Regular message response:")
	for _, block := range message.Content {
		if block.Type == "text" {
			fmt.Println(block.Text)
		}
	}
	fmt.Println()
```

* Streaming
--------

```go
  // Streaming response
	fmt.Println("Streaming message response:")
	message, err = client.Messages().Create(context.Background(), &anthropic.MessageParams{
		Model: string(anthropic.ModelSonnet), // Use the Sonnet model
		Messages: []anthropic.MessageParam{
			{
				Role: "user",
				Content: []anthropic.ContentBlock{
					{Type: "text", Text: "Count from 1 to 5 slowly."},
				},
			},
		},
		MaxTokens: 1000,
		StreamFunc: func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		},
	})

	if err != nil {
		fmt.Errorf("got error: %v", err)
	}

	fmt.Printf("\nFinal message: %+v\n", message)
```
