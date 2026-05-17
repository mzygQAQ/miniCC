package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type DeepSeek struct {
	apiKey string
	model  string
	client *http.Client
}

func NewDeepSeek(apiKey string) *DeepSeek {
	return &DeepSeek{
		apiKey: apiKey,
		model:  "deepseek-chat",
		client: http.DefaultClient,
	}
}

func (d *DeepSeek) Chat(messages []Message) (string, error) {
	body := map[string]any{
		"model":    d.model,
		"messages": messages,
	}
	return d.postAndRead(body)
}

func (d *DeepSeek) ChatWithTools(messages []Message, tools []map[string]any) (string, []ToolCall, error) {
	body := map[string]any{
		"model":    d.model,
		"messages": messages,
	}
	if len(tools) > 0 {
		body["tools"] = tools
	}
	return d.postAndReadTools(body)
}

func (d *DeepSeek) postAndRead(body map[string]any) (string, error) {
	data, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+d.apiKey)

	resp, err := d.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("llm: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("llm: API error %d: %s", resp.StatusCode, respBody)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("llm: decode response: %w", err)
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("llm: no choices in response")
	}
	return result.Choices[0].Message.Content, nil
}

func (d *DeepSeek) postAndReadTools(body map[string]any) (string, []ToolCall, error) {
	data, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+d.apiKey)

	resp, err := d.client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("llm: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("llm: API error %d: %s", resp.StatusCode, respBody)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content   string     `json:"content"`
				ToolCalls []ToolCall `json:"tool_calls"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", nil, fmt.Errorf("llm: decode response: %w", err)
	}
	if len(result.Choices) == 0 {
		return "", nil, fmt.Errorf("llm: no choices in response")
	}
	msg := result.Choices[0].Message
	return msg.Content, msg.ToolCalls, nil
}

func (d *DeepSeek) ChatStream(messages []Message, onChunk func(string)) (string, error) {
	body := map[string]any{
		"model":    d.model,
		"messages": messages,
		"stream":   true,
	}
	return d.streamAndRead(body, onChunk)
}

func (d *DeepSeek) ChatStreamWithTools(messages []Message, tools []map[string]any, onChunk func(string)) (string, []ToolCall, error) {
	body := map[string]any{
		"model":    d.model,
		"messages": messages,
		"stream":   true,
	}
	if len(tools) > 0 {
		body["tools"] = tools
	}
	return d.streamAndReadTools(body, onChunk)
}

// streamEvent represents a single SSE event from the streaming API.
type streamEvent struct {
	Choices []struct {
		Delta struct {
			Content   string `json:"content"`
			ToolCalls []struct {
				Index    int    `json:"index"`
				ID       string `json:"id"`
				Type     string `json:"type"`
				Function struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				} `json:"function"`
			} `json:"tool_calls"`
		} `json:"delta"`
	} `json:"choices"`
}

func (d *DeepSeek) streamAndRead(body map[string]any, onChunk func(string)) (string, error) {
	data, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+d.apiKey)

	resp, err := d.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("llm: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("llm: API error %d: %s", resp.StatusCode, respBody)
	}

	var full strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		payload := strings.TrimPrefix(line, "data: ")
		if payload == "[DONE]" {
			break
		}
		var evt streamEvent
		if err := json.Unmarshal([]byte(payload), &evt); err != nil {
			continue
		}
		if len(evt.Choices) > 0 {
			chunk := evt.Choices[0].Delta.Content
			full.WriteString(chunk)
			if onChunk != nil {
				onChunk(chunk)
			}
		}
	}
	return full.String(), scanner.Err()
}

func (d *DeepSeek) streamAndReadTools(body map[string]any, onChunk func(string)) (string, []ToolCall, error) {
	data, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+d.apiKey)

	resp, err := d.client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("llm: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", nil, fmt.Errorf("llm: API error %d: %s", resp.StatusCode, respBody)
	}

	type partialToolCall struct {
		index int
		id    string
		name  string
		args  strings.Builder
	}
	toolCallMap := make(map[int]*partialToolCall)
	var hasToolCalls bool
	var full strings.Builder

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		payload := strings.TrimPrefix(line, "data: ")
		if payload == "[DONE]" {
			break
		}
		var evt streamEvent
		if err := json.Unmarshal([]byte(payload), &evt); err != nil {
			continue
		}
		if len(evt.Choices) == 0 {
			continue
		}
		delta := evt.Choices[0].Delta

		if len(delta.ToolCalls) > 0 {
			hasToolCalls = true
			for _, tc := range delta.ToolCalls {
				existing, ok := toolCallMap[tc.Index]
				if !ok {
					existing = &partialToolCall{index: tc.Index}
					toolCallMap[tc.Index] = existing
				}
				if tc.ID != "" {
					existing.id = tc.ID
				}
				if tc.Function.Name != "" {
					existing.name = tc.Function.Name
				}
				if tc.Function.Arguments != "" {
					existing.args.WriteString(tc.Function.Arguments)
				}
			}
		}
		if delta.Content != "" {
			full.WriteString(delta.Content)
			if onChunk != nil {
				onChunk(delta.Content)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", nil, err
	}

	if hasToolCalls {
		toolCalls := make([]ToolCall, 0, len(toolCallMap))
		for i := 0; i < len(toolCallMap); i++ {
			tc := toolCallMap[i]
			toolCalls = append(toolCalls, ToolCall{
				ID:   tc.id,
				Type: "function",
				Function: struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				}{
					Name:      tc.name,
					Arguments: tc.args.String(),
				},
			})
		}
		return "", toolCalls, nil
	}

	return full.String(), nil, nil
}
