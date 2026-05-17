package tool

import (
	"encoding/json"
	"fmt"
)

type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
	Execute     func(args json.RawMessage) (string, error)
}

type Registry struct {
	tools map[string]*Tool
}

func NewRegistry() *Registry {
	return &Registry{tools: make(map[string]*Tool)}
}

func (r *Registry) Register(t *Tool) {
	r.tools[t.Name] = t
}

func (r *Registry) Get(name string) (*Tool, bool) {
	t, ok := r.tools[name]
	return t, ok
}

func (r *Registry) List() []*Tool {
	list := make([]*Tool, 0, len(r.tools))
	for _, t := range r.tools {
		list = append(list, t)
	}
	return list
}

func (r *Registry) Schema() []map[string]any {
	schemas := make([]map[string]any, 0, len(r.tools))
	for _, t := range r.tools {
		schemas = append(schemas, map[string]any{
			"type": "function",
			"function": map[string]any{
				"name":        t.Name,
				"description": t.Description,
				"parameters":  t.Parameters,
			},
		})
	}
	return schemas
}

// Default is the global tool registry, auto-populated via init().
var Default = NewRegistry()

func init() {
	Default.Register(BashTool())
	Default.Register(ReadTool())
	Default.Register(WriteTool())
	Default.Register(EditTool())
}

func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + fmt.Sprintf("\n...(truncated %d bytes)", len(s))
}
