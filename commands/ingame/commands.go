package ingame

import (
	"reflect"
	"strings"

	"alex-j-butler.com/tf2-booking/commands/ingame/loghandler"
	"alex-j-butler.com/tf2-booking/servers"
)

type CommandInfo struct {
	loghandler.SayEvent
	Server *servers.Server
}

type CommandFunction func(CommandInfo, string, []string)

type CommandHandler struct {
	function CommandFunction
}

type Command struct {
	Prefix   string
	Handlers map[string]*CommandHandler
}

func NewCommand(function CommandFunction) *CommandHandler {
	return &CommandHandler{
		function: function,
	}
}

// New creates a new instance of the Command system
// with the specified prefix.
func New(prefix string) *Command {
	return &Command{
		Prefix:   prefix,
		Handlers: make(map[string]*CommandHandler),
	}
}

// Add creates a new entry in the command handlers map.
// First argument accepts a command handler implementing the type `CommandHandler`,
// Second argument accepts a variable amount of strings specifying the commands to register.
func (c *Command) Add(handler *CommandHandler, commands ...string) {
	for _, command := range commands {
		c.Handlers[command] = handler
	}
}

// Remove deletes entries from the command handlers map.
func (c *Command) Remove(commands ...string) {
	for _, command := range commands {
		delete(c.Handlers, command)
	}
}

// Handle the incoming commands and dispatches them to the appropriate
// command handler, after parsing them.
func (c *Command) Handle(info CommandInfo, command string, permissions int) {
	for str, handler := range c.Handlers {
		if !strings.HasPrefix(command, c.Prefix) && c.Prefix != "" {
			continue
		}

		handlerSplit := strings.Split(str, " ")
		commandSplit := strings.Split(command[len(c.Prefix):], " ")

		if len(commandSplit) < len(handlerSplit) {
			continue
		}

		if reflect.DeepEqual(handlerSplit, commandSplit[:len(handlerSplit)]) {
			handler.function(info, strings.Join(handlerSplit, " "), commandSplit[len(handlerSplit):])

			break
		}
	}
}
