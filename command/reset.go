package command

import "fmt"

type resetCmd struct{}

func (resetCmd) Name() string { return "/reset" }
func (resetCmd) Desc() string { return "清空对话历史" }

func (resetCmd) Execute(bot Bot, _ string) error {
	bot.Reset()
	fmt.Println("记忆已清空")
	return nil
}
