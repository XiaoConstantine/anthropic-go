Anthropic Go Client
-------------------
[![Release](https://github.com/XiaoConstantine/anthropic-go/actions/workflows/release.yml/badge.svg)](https://github.com/XiaoConstantine/anthropic-go/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/XiaoConstantine/anthropic-go)](https://goreportcard.com/report/github.com/XiaoConstantine/anthropic-go)
[![codecov](https://codecov.io/gh/XiaoConstantine/anthropic-go/graph/badge.svg?token=DZCEY7IFBG)](https://codecov.io/gh/XiaoConstantine/anthropic-go)
[![Documentation](https://github.com/XiaoConstantine/anthropic-go/actions/workflows/doc.yml/badge.svg)](https://github.com/XiaoConstantine/anthropic-go/actions/workflows/doc.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/XiaoConstantine/anthropic-go.svg)](https://pkg.go.dev/github.com/XiaoConstantine/anthropic-go)

Go sdk for interacting with anthropic API


Installation
------------

```bash
go get github.com/XiaoConstantine/anthropic-go

```

Usage
------

### Creating a Client
```go
client, err := anthropic.NewClient(
    anthropic.WithAPIKey(""), // Uses the ANTHROPIC_API_KEY environment variable if empty
    anthropic.WithTimeout(30*time.Second),
)
if err != nil {
    log.Fatalf("Failed to create client: %v", err)
}
```


### Interacting with Models

#### Listing Available Models

```go
	// List available models - Since there's no public api for this, currently the result is
	// hard coded
	models, _ := client.Models().List()
	fmt.Println("Available models:")
	for _, model := range models {
		fmt.Printf("- %s (%s)\n", model.Name, model.ID)
	}
	fmt.Println()
```

### Working with Messages

#### Sending a basic Message
-------

```go
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
		MaxTokens: 2048,
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

#### Streaming Response

```go
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
		MaxTokens: 2048,
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

### Handling Images

#### Processing an Image
```go
	imageData, err := ioutil.ReadFile("/<path_to_image>/image.png")
	if err != nil {
		log.Fatalf("Failed to read image file: %v", err)
	}

	// Encode the image data to base64
	base64Image := base64.StdEncoding.EncodeToString(imageData)

	// Create the message params
	params := &anthropic.MessageParams{
		Model:     string(anthropic.ModelSonnet), // model has vision capability
		MaxTokens: 4096,
		Messages: []anthropic.MessageParam{
			{
				Role: "user",
				Content: []anthropic.ContentBlock{
					{
						Type: "text",
						Text: "Here's an image. Can you describe it?",
					},
					{
						Type: "image",
						Source: &anthropic.Image{
							Type:      "base64",
							MediaType: "image/png", // Adjust based on your image type
							Data:      base64Image,
						},
					},
				},
			},
		},
	}

	// Send the request
	message, err := client.Messages().Create(context.Background(), params)
	if err != nil {
		log.Fatalf("Failed to create message: %v", err)
	}

	// Print the response
	fmt.Printf("Response: %+v\n", message)

	// If you want to print just the text content:
	for _, content := range message.Content {
		if content.Type == "text" {
			fmt.Println(content.Text)
		}
	}
```

License
-------
MIT
