package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Anthropic struct {
	apiKey string
	model  string
	client *http.Client
}

func NewAnthropic(apiKey string) *Anthropic {
	return &Anthropic{
		apiKey: apiKey,
		model:  "claude-sonnet-4-6",
		client: http.DefaultClient,
	}
}

func (a *Anthropic) Chat(messages []Message) (string, error) {
	// Anthropic separates system from the messages array.
	var system string
	apiMessages := make([]map[string]string, 0, len(messages))
	for _, m := range messages {
		text := ""
		if m.Content != nil {
			text = *m.Content
		}
		switch m.Role {
		case RoleSystem:
			system = text
		default:
			apiMessages = append(apiMessages, map[string]string{"role": m.Role, "content": text})
		}
	}

	body := map[string]any{
		"model":      a.model,
		"max_tokens": 1024,
		"messages":   apiMessages,
	}
	if system != "" {
		body["system"] = system
	}
	data, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", a.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("llm: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("llm: API error %d: %s", resp.StatusCode, respBody)
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("llm: decode response: %w", err)
	}
	if len(result.Content) == 0 {
		return "", fmt.Errorf("llm: no content in response")
	}
	return result.Content[0].Text, nil
}
