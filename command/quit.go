package command

import "os"

type quitCmd struct{}

func (quitCmd) Name() string { return "/exit" }
func (quitCmd) Desc() string { return "退出程序" }

func (quitCmd) Execute(_ Bot, _ string) error {
	os.Exit(0)
	return nil
}
