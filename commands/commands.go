package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type CommandHandler func(*discordgo.MessageCreate, string, []string)

type Command struct {
	Prefix   string
	Handlers map[string]CommandHandler
}

// New creates a new instance of the Command system
// with the specified prefix.
func New(prefix string) *Command {
	return &Command{
		Prefix:   prefix,
		Handlers: make(map[string]CommandHandler),
	}
}

// Add creates a new entry in the command handlers map.
// First argument accepts a command handler implementing the type `CommandHandler`,
// Second argument accepts a variable amount of strings specifying the commands to register.
func (c *Command) Add(handler CommandHandler, commands ...string) {
	for _, command := range commands {
		c.Handlers[fmt.Sprintf("%s%s", c.Prefix, command)] = handler
	}
}

// Remove deletes entries from the command handlers map.
func (c *Command) Remove(commands ...string) {
	for _, command := range commands {
		delete(c.Handlers, command)
	}
}

// Handle
// Handles the incoming commands and dispatches them to the appropriate
// command handler, after parsing them.
func (c *Command) Handle(m *discordgo.MessageCreate, command string) {
	for str, handler := range c.Handlers {
		if command == str {
			var fixedCommand string
			if strings.HasPrefix(command, c.Prefix) {
				fixedCommand = command[len(c.Prefix)-1 : len(command)]
			} else {
				fixedCommand = command
			}

			handler(m, fixedCommand, nil)
		}
	}
}
