package command

import (
	"fmt"
	"os/exec"
	"strings"
)

type bashCmd struct{}

func (bashCmd) Name() string { return "/bash" }
func (bashCmd) Desc() string { return "执行 shell 命令，如 /bash ls -l" }

func (bashCmd) Execute(_ Bot, args string) error {
	if strings.TrimSpace(args) == "" {
		return fmt.Errorf("请输入要执行的命令")
	}

	cmd := exec.Command("sh", "-c", args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("退出码: %v\n", err)
	}
	fmt.Print(string(out))
	return nil
}
