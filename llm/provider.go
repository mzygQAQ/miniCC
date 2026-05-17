package llm

import "encoding/json"

// Standard roles.
const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleTool      = "tool"
)

// Message represents a single turn in a conversation.
type Message struct {
	Role       string          `json:"role"`
	Content    *string         `json:"content"`
	ToolCallID string          `json:"tool_call_id,omitempty"`
	ToolCalls  json.RawMessage `json:"tool_calls,omitempty"`
}

// Str returns a pointer to s, for use in Message.Content.
func Str(s string) *string { return &s }

// ToolCall represents a tool invocation from the LLM.
type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

// Provider is the abstraction over different LLM backends.
type Provider interface {
	Chat(messages []Message) (string, error)
}

// StreamingProvider providers can optionally implement for streaming output.
type StreamingProvider interface {
	ChatStream(messages []Message, onChunk func(string)) (string, error)
}

// ToolProvider extends Provider with tool calling support.
type ToolProvider interface {
	ChatWithTools(messages []Message, tools []map[string]any) (string, []ToolCall, error)
}

// StreamingToolProvider extends StreamingProvider with tool calling support.
type StreamingToolProvider interface {
	ChatStreamWithTools(messages []Message, tools []map[string]any, onChunk func(string)) (string, []ToolCall, error)
}
