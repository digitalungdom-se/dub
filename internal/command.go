package internal

import "github.com/digitalungdom-se/dub/pkg"

type (
	Command struct {
		Name        string
		Description string
		Aliases     []string
		Group       string
		Usage       string
		Example     string
		ServerOnly  bool
		AdminOnly   bool
		Execute     func(*pkg.Context, *Server) error
	}

	Commands map[string]Command

	CommandHandler struct {
		commands Commands
	}
)

func NewCommandHandler() CommandHandler {
	var commandHandler CommandHandler
	commandHandler.commands = make(Commands)

	return commandHandler
}

func (handler *CommandHandler) GetCommands(group string) Commands {
	if group == "" {
		return handler.commands
	} else {
		commands := make(Commands)

		for name, command := range handler.commands {
			if command.Group == group {
				commands[name] = command
			}
		}
		return commands
	}
}

func (handler *CommandHandler) GetCommand(name string) (Command, bool) {
	command, found := handler.commands[name]

	if found {
		return command, found
	}

	commands := handler.GetCommands("")

	for _, command := range commands {
		if pkg.StringInSlice(name, command.Aliases) {
			return command, true
		}
	}

	return Command{}, false
}

func (handler *CommandHandler) Register(command Command) {
	handler.commands[command.Name] = command
}
