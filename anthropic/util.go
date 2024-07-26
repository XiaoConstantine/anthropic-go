package anthropic

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func parseStreamingMessageResponse(ctx context.Context, r io.Reader, payload *MessageParams) (*Message, error) {
	scanner := bufio.NewScanner(r)
	eventChan := make(chan MessageEvent)

	go func() {
		defer close(eventChan)
		var response Message
		for scanner.Scan() {
			line := scanner.Text()

			if line == "" || !strings.HasPrefix(line, "data:") {
				continue
			}
			data := strings.TrimPrefix(line, "data: ")
			event, err := parseStreamEvent(data)
			if err != nil {
				eventChan <- MessageEvent{Response: nil, Err: fmt.Errorf("failed to parse stream event: %w", err)}
				return
			}
			response, err = processStreamEvent(ctx, event, payload, response, eventChan)
			if err != nil {
				eventChan <- MessageEvent{Response: nil, Err: fmt.Errorf("failed to process stream event: %w", err)}
				return
			}
		}
		if err := scanner.Err(); err != nil {
			eventChan <- MessageEvent{Response: nil, Err: fmt.Errorf("issue scanning response: %w", err)}
		}
	}()

	var lastResponse *Message
	for event := range eventChan {
		if event.Err != nil {
			return nil, event.Err
		}
		lastResponse = event.Response
	}
	return lastResponse, nil
}

func parseStreamEvent(data string) (map[string]interface{}, error) {
	var event map[string]interface{}
	err := json.Unmarshal([]byte(data), &event)
	return event, err
}

func processStreamEvent(ctx context.Context, event map[string]interface{}, payload *MessageParams, response Message, eventChan chan<- MessageEvent) (Message, error) {
	eventType, ok := event["type"].(string)
	if !ok {
		return response, fmt.Errorf("invalid event type")
	}
	switch eventType {
	case "message_start":
		return handleMessageStartEvent(event, response)
	case "content_block_start":
		return handleContentBlockStartEvent(event, response)
	case "content_block_delta":
		return handleContentBlockDeltaEvent(ctx, event, payload, response)
	case "content_block_stop":
		// Nothing to do here
	case "message_delta":
		return handleMessageDeltaEvent(event, response)
	case "message_stop":
		// Nothing to do here
		eventChan <- MessageEvent{Response: &response, Err: nil}
	case "ping":
		// Nothing to do here
	default:
		fmt.Printf("unknown event type: %s\n", eventType)
	}
	return response, nil
}

func handleMessageStartEvent(event map[string]interface{}, response Message) (Message, error) {
	message, ok := event["message"].(map[string]interface{})
	if !ok {
		return response, fmt.Errorf("invalid message field")
	}

	usage, ok := message["usage"].(map[string]interface{})
	if !ok {
		return response, fmt.Errorf("invalid usage field")
	}

	inputTokens, ok := usage["input_tokens"].(float64)
	if !ok {
		return response, fmt.Errorf("invalid input_tokens field")
	}

	response.ID = getString(message, "id")
	response.Model = getString(message, "model")
	response.Role = getString(message, "role")
	response.Type = getString(message, "type")
	response.Usage.InputTokens = int(inputTokens)

	return response, nil
}

func handleContentBlockStartEvent(event map[string]interface{}, response Message) (Message, error) {
	indexValue, ok := event["index"].(float64)
	if !ok {
		return response, fmt.Errorf("invalid index field")
	}
	index := int(indexValue)

	contentBlock, ok := event["content_block"].(map[string]interface{})
	if !ok {
		return response, fmt.Errorf("invalid content_block field")
	}

	contentType := getString(contentBlock, "type")

	if len(response.Content) <= index {
		response.Content = append(response.Content, ContentBlock{
			Type: contentType,
		})
	}
	return response, nil
}

func handleContentBlockDeltaEvent(ctx context.Context, event map[string]interface{}, payload *MessageParams, response Message) (Message, error) {
	indexValue, ok := event["index"].(float64)
	if !ok {
		return response, fmt.Errorf("invalid index field")
	}
	index := int(indexValue)

	delta, ok := event["delta"].(map[string]interface{})
	if !ok {
		return response, fmt.Errorf("invalid delta field")
	}
	deltaType := getString(delta, "type")

	if deltaType == "text_delta" {
		text := getString(delta, "text")
		if len(response.Content) <= index {
			response.Content = append(response.Content, ContentBlock{
				Type: "text",
				Text: text,
			})
		} else {
			response.Content[index].Text += text
		}
	}

	if payload.IsStreaming() {
		text, ok := delta["text"].(string)
		if !ok {
			return response, fmt.Errorf("invalid delta data")
		}
		err := payload.StreamFunc(ctx, []byte(text))
		if err != nil {
			return response, fmt.Errorf("streaming func returned an error: %w", err)
		}
	}

	return response, nil
}

func handleMessageDeltaEvent(event map[string]interface{}, response Message) (Message, error) {
	delta, ok := event["delta"].(map[string]interface{})
	if !ok {
		return response, fmt.Errorf("invalid delta field")
	}
	response.StopReason = getString(delta, "stop_reason")

	usage, ok := event["usage"].(map[string]interface{})
	if !ok {
		return response, fmt.Errorf("invalid usage field")
	}
	outputTokens, ok := usage["output_tokens"].(float64)
	if ok {
		response.Usage.OutputTokens = int(outputTokens)
	}
	return response, nil
}

func getString(m map[string]interface{}, key string) string {
	value, ok := m[key].(string)
	if !ok {
		return ""
	}
	return value
}
