package tool

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func EditTool() *Tool {
	return &Tool{
		Name:        "edit",
		Description: "修改文件中的指定内容（精确字符串替换），适合修改代码",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"file_path": {
					"type": "string",
					"description": "文件路径"
				},
				"old_string": {
					"type": "string",
					"description": "被替换的旧内容（必须是文件中存在的唯一匹配）"
				},
				"new_string": {
					"type": "string",
					"description": "替换后的新内容"
				}
			},
			"required": ["file_path", "old_string", "new_string"]
		}`),
		Execute: func(args json.RawMessage) (string, error) {
			var input struct {
				Path     string `json:"file_path"`
				Old      string `json:"old_string"`
				New      string `json:"new_string"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", fmt.Errorf("解析参数失败: %w", err)
			}

			data, err := os.ReadFile(input.Path)
			if err != nil {
				return "", fmt.Errorf("读取文件失败: %w", err)
			}

			content := string(data)
			count := strings.Count(content, input.Old)
			if count == 0 {
				return "", fmt.Errorf("未在文件中找到匹配的内容")
			}
			if count > 1 {
				return "", fmt.Errorf("找到 %d 处匹配，请提供更精确的匹配内容（需要唯一匹配）", count)
			}

			newContent := strings.Replace(content, input.Old, input.New, 1)
			if err := os.WriteFile(input.Path, []byte(newContent), 0644); err != nil {
				return "", fmt.Errorf("写入文件失败: %w", err)
			}

			return fmt.Sprintf("已修改 %s（替换了 %d 个字符）", input.Path, len(input.Old)), nil
		},
	}
}
