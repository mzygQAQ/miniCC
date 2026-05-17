package tool

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func ReadTool() *Tool {
	return &Tool{
		Name:        "read",
		Description: "读取文件内容，可指定行号范围",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"file_path": {
					"type": "string",
					"description": "文件路径"
				},
				"offset": {
					"type": "integer",
					"description": "起始行号，从 0 开始"
				},
				"limit": {
					"type": "integer",
					"description": "最多读取行数"
				}
			},
			"required": ["file_path"]
		}`),
		Execute: func(args json.RawMessage) (string, error) {
			var input struct {
				Path   string `json:"file_path"`
				Offset int    `json:"offset"`
				Limit  int    `json:"limit"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", fmt.Errorf("解析参数失败: %w", err)
			}

			data, err := os.ReadFile(input.Path)
			if err != nil {
				return "", fmt.Errorf("读取文件失败: %w", err)
			}

			lines := strings.Split(string(data), "\n")
			if input.Offset > 0 || input.Limit > 0 {
				start := input.Offset
				if start > len(lines) {
					start = len(lines)
				}
				end := start + input.Limit
				if input.Limit == 0 || end > len(lines) {
					end = len(lines)
				}
				lines = lines[start:end]
			}

			result := strings.Join(lines, "\n")
			return Truncate(result, 10240), nil
		},
	}
}
