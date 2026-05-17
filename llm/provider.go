package llm

// Standard roles.
const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
)

// Message represents a single turn in a conversation.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Provider is the abstraction over different LLM backends.
type Provider interface {
	Chat(messages []Message) (string, error)
}

// StreamingProvider providers can optionally implement for streaming output.
type StreamingProvider interface {
	ChatStream(messages []Message, onChunk func(string)) (string, error)
}
