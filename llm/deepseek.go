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

func (d *DeepSeek) ChatStream(messages []Message, onChunk func(string)) (string, error) {
	body := map[string]any{
		"model":    d.model,
		"messages": messages,
		"stream":   true,
	}
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
		var event struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}
		if err := json.Unmarshal([]byte(payload), &event); err != nil {
			continue
		}
		if len(event.Choices) > 0 {
			chunk := event.Choices[0].Delta.Content
			full.WriteString(chunk)
			if onChunk != nil {
				onChunk(chunk)
			}
		}
	}
	return full.String(), scanner.Err()
}
