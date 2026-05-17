package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/chzyer/readline"

	"miniCC/command"
	"miniCC/llm"
	"miniCC/tool"
)

func main() {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "请设置环境变量 DEEPSEEK_API_KEY")
		os.Exit(1)
	}
	provider := llm.NewDeepSeek(apiKey)
	convo := NewConversation(provider, llm.Message{
		Role:    llm.RoleSystem,
		Content: llm.Str(`你是 miniCC，一个可以直接在用户终端执行命令的 AI 助手。

你有以下工具可用：
- bash：执行 shell 命令，适合文件操作、代码运行、系统管理等
- read：读取文件内容
- write：完整写入文件（覆盖整个文件）
- edit：精确替换文件中的指定内容（修改代码时优先使用，只需指定要替换的片段）

修改代码时优先使用 edit 工具做局部替换，而不是用 write 重写整个文件。
当用户需要操作文件、执行代码或查询系统信息时，请直接调用对应工具。
工具执行结果会返回给你，基于结果继续回答用户问题。`),
	})

	names := make([]string, 0, len(command.Commands))
	for name := range command.Commands {
		names = append(names, name)
	}
	sort.Strings(names)

	completions := make([]readline.PrefixCompleterInterface, len(names))
	for i, name := range names {
		completions[i] = readline.PcItem(name)
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:       ">>> ",
		AutoComplete: readline.NewPrefixCompleter(completions...),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "初始化终端失败: %v\n", err)
		os.Exit(1)
	}
	defer rl.Close()

	fmt.Println("对话开始（↑↓ 浏览历史, Tab 补全命令, 输入 / 查看所有命令）")
	fmt.Println(strings.Repeat("-", 40))

	for {
		input, err := rl.Readline()
		if err != nil {
			break
		}
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		if input == "/" {
			for _, name := range names {
				fmt.Printf("  %s  - %s\n", name, command.Commands[name].Desc())
			}
			continue
		}

		if input[0] == '/' {
			parts := strings.SplitN(input, " ", 2)
			cmd, ok := command.Commands[parts[0]]
			if !ok {
				fmt.Printf("未知命令: %s（输入 / 查看所有命令）\n", parts[0])
				continue
			}
			args := ""
			if len(parts) == 2 {
				args = parts[1]
			}
			if err := cmd.Execute(convo, args); err != nil {
				fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			}
			continue
		}

		if _, err = convo.SayWithTools(input, tool.Default, func(chunk string) {
			fmt.Print(chunk)
		}); err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			continue
		}
		fmt.Println()
		fmt.Println()
	}
}
