/*
Package anthropic provides a Go client for interacting with the Anthropic API.

This SDK allows Go developers to easily integrate Anthropic's AI models into their
applications. It provides a simple and idiomatic way to send requests to the Anthropic API,
create messages, and handle responses, including support for streaming.

Getting Started:

To use this package, you need to have an API key from Anthropic. You can create a new client like this:

	client, err := anthropic.NewClient(anthropic.WithAPIKey("your-api-key"))
	if err != nil {
	    // Handle error
	}

Creating a Message:

To create a new message:

	params := &anthropic.MessageParams{
	    Model: string(anthropic.ModelSonnet),
	    Messages: []anthropic.MessageParam{
	        {
	            Role: "user",
	            Content: []anthropic.ContentBlock{
	                {Type: "text", Text: "Hello, Claude!"},
	            },
	        },
	    },
	}

	message, err := client.Messages().Create(context.Background(), params)
	if err != nil {
	    // Handle error
	}

	fmt.Println(message.Content)

Streaming:

For streaming responses, use the StreamFunc in MessageParams:

	params.StreamFunc = func(ctx context.Context, chunk []byte) error {
	    fmt.Print(string(chunk))
	    return nil
	}

	message, err := client.Messages().Create(context.Background(), params)
	if err != nil {
	    // Handle error
	}

Available Models:

The SDK supports the following Anthropic models:
- Claude 3 Haiku (anthropic.ModelHaiku)
- Claude 3 Sonnet (anthropic.ModelSonnet)
- Claude 3 Opus (anthropic.ModelOpus)

For more detailed information about specific components, please refer to the
documentation of individual files and functions.
*/
package anthropic
