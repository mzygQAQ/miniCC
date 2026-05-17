package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/chzyer/readline"

	"miniCC/command"
	"miniCC/llm"
)

func main() {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "请设置环境变量 DEEPSEEK_API_KEY")
		os.Exit(1)
	}
	provider := llm.NewDeepSeek(apiKey)
	convo := NewConversation(provider)

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

		if _, err = convo.SayStream(input, func(chunk string) {
			fmt.Print(chunk)
		}); err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			continue
		}
		fmt.Println()
		fmt.Println()
	}
}
