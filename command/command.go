package command

// Bot is the interface commands use to interact with the conversation.
type Bot interface {
	Say(msg string) (string, error)
	SayRole(role, msg string) (string, error)
	Reset()
}

// Command represents a slash command.
type Command interface {
	Name() string
	Desc() string
	Execute(bot Bot, args string) error
}

// Commands is the registry of all available commands.
var Commands map[string]Command

func init() {
	Commands = make(map[string]Command, 4)
	for _, c := range []Command{
		bashCmd{},
		resetCmd{},
		quitCmd{},
		helpCmd{},
	} {
		Commands[c.Name()] = c
	}
}
