package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type CommandFunction func(*discordgo.MessageCreate, string, []string) bool
type TriggerFunction func(*discordgo.MessageCreate, string)

type Command struct {
	// Function to be called when this function is executed.
	function CommandFunction

	// Specifies the required permissions for the command to run.
	permissions int

	// Specifies whether this command responds to Direct Messages/Private chat messages.
	respondToDM bool

	// Subcommands that this command can call.
	subcommands map[string]*Command
}

type Trigger struct {
	// Function to be called when this trigger is executed.
	function TriggerFunction

	// Specifies the required permissions for the command to run.
	permissions int

	// Specifies whether this command responds to Direct Messages/Private chat messages.
	respondToDM bool
}

func NewCommand(function CommandFunction) *Command {
	return &Command{
		function:    function,
		permissions: -1,
		respondToDM: true,
		subcommands: make(map[string]*Command),
	}
}

func (c *Command) Permissions(permissions int) *Command {
	c.permissions = permissions
	return c
}

func (c *Command) RespondToDM(respondToDM bool) *Command {
	c.respondToDM = respondToDM
	return c
}

func (c *Command) AddSubcommand(commandName string, command *Command) error {
	commandName = strings.ToLower(commandName)
	if _, ok := c.subcommands[commandName]; ok {
		return fmt.Errorf("%s subcommand already exists", commandName)
	}

	c.subcommands[commandName] = command
	return nil
}

func (c *Command) RemoveSubcommand(commandName string) error {
	commandName = strings.ToLower(commandName)
	if _, ok := c.subcommands[commandName]; !ok {
		return fmt.Errorf("%s subcommand does not exist", commandName)
	}

	// Delete the subcommand.
	delete(c.subcommands, commandName)
	return nil
}

func (c *Command) handleCommand(message *discordgo.MessageCreate, input string, args []string) {
	// Call the handler function for our command.
	if c.function != nil {
		// If the handler function for our command returns true, then we shouldn't handle
		// anything else here, otherwise we'll continue processing subcommands.
		if c.function(message, input, args) {
			return
		}
	}

	// If we're here, the handler function was either nil, or returned false informing us we should continue processing
	// the subcommands.
	if len(c.subcommands) > 0 && len(args) > 0 {
		// Check if any subcommands are a match.
		if command, ok := c.subcommands[strings.ToLower(args[0])]; ok {
			input = fmt.Sprintf("%s %s", input, args[0])
			args = args[1:]

			command.handleCommand(message, input, args)
			return
		}
	}
}

func NewTrigger(function TriggerFunction) *Trigger {
	return &Trigger{
		function:    function,
		permissions: -1,
		respondToDM: true,
	}
}

func (t *Trigger) Permissions(permissions int) *Trigger {
	t.permissions = permissions
	return t
}

func (t *Trigger) RespondToDM(respondToDM bool) *Trigger {
	t.respondToDM = respondToDM
	return t
}

type CommandSystem struct {
	// Commands that the command system can handle.
	commands map[string]*Command

	// Triggers that the command system can handle.
	triggers map[string]*Trigger
}

func NewCommandSystem() *CommandSystem {
	return &CommandSystem{
		commands: make(map[string]*Command),
		triggers: make(map[string]*Trigger),
	}
}

func (cs *CommandSystem) AddCommand(commandName string, command *Command) error {
	commandName = strings.ToLower(commandName)
	if _, ok := cs.commands[commandName]; ok {
		return fmt.Errorf("%s command already exists", commandName)
	}

	cs.commands[commandName] = command
	return nil
}

func (cs *CommandSystem) RemoveCommand(commandName string) error {
	commandName = strings.ToLower(commandName)
	if _, ok := cs.commands[commandName]; !ok {
		return fmt.Errorf("%s command does not exist", commandName)
	}

	// Delete the command.
	delete(cs.commands, commandName)
	return nil
}

func (cs *CommandSystem) AddTrigger(triggerName string, trigger *Trigger) error {
	triggerName = strings.ToLower(triggerName)
	if _, ok := cs.triggers[triggerName]; ok {
		return fmt.Errorf("%s trigger already exists", triggerName)
	}

	cs.triggers[triggerName] = trigger
	return nil
}

func (cs *CommandSystem) RemoveTrigger(triggerName string) error {
	triggerName = strings.ToLower(triggerName)
	if _, ok := cs.triggers[triggerName]; !ok {
		return fmt.Errorf("%s trigger does not exist", triggerName)
	}

	// Delete the trigger.
	delete(cs.triggers, triggerName)
	return nil
}

func (cs *CommandSystem) HandleCommand(message *discordgo.MessageCreate, input string, args []string) {
	input = args[0]

	if len(cs.commands) > 0 && len(args) > 0 {
		// Check if any commands are a match.
		if command, ok := cs.commands[strings.ToLower(input)]; ok {
			args = args[1:]

			command.handleCommand(message, input, args)
			return
		}
	}

	if len(cs.triggers) > 0 {
		// Check if any triggers are a match.
		if trigger, ok := cs.triggers[strings.ToLower(message.Content)]; ok {
			// Call the trigger function.
			trigger.function(message, message.Content)
			return
		}
	}
}
