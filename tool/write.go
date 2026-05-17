package tool

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func WriteTool() *Tool {
	return &Tool{
		Name:        "write",
		Description: "写入内容到文件（自动创建目录）",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"file_path": {
					"type": "string",
					"description": "文件路径"
				},
				"content": {
					"type": "string",
					"description": "文件内容"
				}
			},
			"required": ["file_path", "content"]
		}`),
		Execute: func(args json.RawMessage) (string, error) {
			var input struct {
				Path    string `json:"file_path"`
				Content string `json:"content"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", fmt.Errorf("解析参数失败: %w", err)
			}

			if err := os.MkdirAll(filepath.Dir(input.Path), 0755); err != nil {
				return "", fmt.Errorf("创建目录失败: %w", err)
			}

			if err := os.WriteFile(input.Path, []byte(input.Content), 0644); err != nil {
				return "", fmt.Errorf("写入文件失败: %w", err)
			}

			return fmt.Sprintf("已写入 %d 字节到 %s", len(input.Content), input.Path), nil
		},
	}
}
