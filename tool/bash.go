package tool

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

func BashTool() *Tool {
	return &Tool{
		Name:        "bash",
		Description: "执行 shell 命令，返回标准输出和错误输出",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"command": {
					"type": "string",
					"description": "要执行的 shell 命令"
				}
			},
			"required": ["command"]
		}`),
		Execute: func(args json.RawMessage) (string, error) {
			var input struct {
				Command string `json:"command"`
			}
			if err := json.Unmarshal(args, &input); err != nil {
				return "", fmt.Errorf("解析参数失败: %w", err)
			}
			cmd := exec.Command("sh", "-c", input.Command)
			out, err := cmd.CombinedOutput()
			result := strings.TrimSpace(string(out))
			if err != nil {
				if result != "" {
					result += "\n"
				}
				result += fmt.Sprintf("(exit: %v)", err)
			}
			return Truncate(result, 10240), nil
		},
	}
}
