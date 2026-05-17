package command

import "fmt"

type helpCmd struct{}

func (helpCmd) Name() string { return "/help" }
func (helpCmd) Desc() string { return "显示此帮助" }

func (helpCmd) Execute(_ Bot, _ string) error {
	fmt.Println("可用命令：")
	for _, cmd := range Commands {
		fmt.Printf("  %s  - %s\n", cmd.Name(), cmd.Desc())
	}
	return nil
}
