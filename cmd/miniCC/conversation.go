package main

import (
	"fmt"
	"os"

	"miniCC/llm"
)

const maxBytes = 200 * 1024 * 1024 // 200MB

// Conversation holds the message history and manages automatic compression.
type Conversation struct {
	provider llm.Provider
	history  []llm.Message
}

func NewConversation(p llm.Provider) *Conversation {
	return &Conversation{provider: p}
}

// totalBytes returns the sum of message content lengths.
func (c *Conversation) totalBytes() int {
	n := 0
	for _, m := range c.history {
		n += len(m.Content)
	}
	return n
}

// Say sends a user message and returns the assistant reply.
func (c *Conversation) Say(msg string) (string, error) {
	return c.SayRole(llm.RoleUser, msg)
}

// SayRole sends a message with an arbitrary role and returns the assistant reply.
func (c *Conversation) SayRole(role, msg string) (string, error) {
	c.history = append(c.history, llm.Message{Role: role, Content: msg})

	reply, err := c.sayStream(c.history, nil)
	if err != nil {
		return "", err
	}

	c.history = append(c.history, llm.Message{Role: llm.RoleAssistant, Content: reply})

	if c.totalBytes() >= maxBytes {
		if err := c.compress(); err != nil {
			fmt.Fprintf(os.Stderr, "压缩对话失败: %v\n", err)
		}
	}

	return reply, nil
}

// SayStream sends a user message and streams the reply chunk by chunk.
func (c *Conversation) SayStream(msg string, onChunk func(string)) (string, error) {
	return c.SayRoleStream(llm.RoleUser, msg, onChunk)
}

// SayRoleStream sends a message with an arbitrary role and streams the reply.
func (c *Conversation) SayRoleStream(role, msg string, onChunk func(string)) (string, error) {
	c.history = append(c.history, llm.Message{Role: role, Content: msg})

	reply, err := c.sayStream(c.history, onChunk)
	if err != nil {
		return "", err
	}

	c.history = append(c.history, llm.Message{Role: llm.RoleAssistant, Content: reply})

	if c.totalBytes() >= maxBytes {
		if err := c.compress(); err != nil {
			fmt.Fprintf(os.Stderr, "压缩对话失败: %v\n", err)
		}
	}

	return reply, nil
}

// sayStream calls the provider, using streaming if available.
func (c *Conversation) sayStream(messages []llm.Message, onChunk func(string)) (string, error) {
	if sp, ok := c.provider.(llm.StreamingProvider); ok && onChunk != nil {
		return sp.ChatStream(messages, onChunk)
	}
	return c.provider.Chat(messages)
}

// Reset clears conversation history.
func (c *Conversation) Reset() {
	c.history = nil
}

// compress summarizes older messages into a system prompt, keeping recent
// exchanges intact so the LLM doesn't lose immediate context.
func (c *Conversation) compress() error {
	const keepRecent = 4 // keep last 2 exchanges (user + assistant × 2)

	old := c.history[:len(c.history)-keepRecent]
	recent := c.history[len(c.history)-keepRecent:]

	prompt := make([]llm.Message, 0, len(old)+3)
	prompt = append(prompt, llm.Message{
		Role:    llm.RoleSystem,
		Content: "你是对话摘要助手。请用一段中文总结以下对话的核心内容，保留关键上下文信息。",
	})
	prompt = append(prompt, old...)
	prompt = append(prompt, llm.Message{
		Role:    llm.RoleUser,
		Content: "请用一段话总结以上对话中已经讨论过的内容，包括关键决定、代码细节、用户偏好等。",
	})

	summary, err := c.provider.Chat(prompt)
	if err != nil {
		return fmt.Errorf("对话压缩失败: %w", err)
	}

	c.history = append(
		[]llm.Message{{Role: llm.RoleSystem, Content: "对话摘要：" + summary}},
		recent...,
	)
	return nil
}
