package commands

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"alex-j-butler.com/tf2-booking/util"

	"github.com/alecthomas/participle"
	"github.com/bwmarrin/discordgo"
)

type CommandFunction func(*discordgo.MessageCreate, string, CommandPermissions, CommandArgList)

type CommandHandler struct {
	function    CommandFunction
	permissions int
	respondToDM bool
}

type CommandSystem struct {
	Prefix   string
	Handlers map[string]*CommandHandler
	parser   *participle.Parser
}

func NewCommand(function CommandFunction) *CommandHandler {
	return &CommandHandler{
		function:    function,
		permissions: -1,
		respondToDM: false,
	}
}

func (ch *CommandHandler) Permissions(permissions int) *CommandHandler {
	ch.permissions = permissions
	return ch
}

func (ch *CommandHandler) RespondToDM(respondToDM bool) *CommandHandler {
	ch.respondToDM = respondToDM
	return ch
}

// New creates a new instance of the CommandSystem
// with the specified prefix.
func New() *CommandSystem {
	return &CommandSystem{
		Handlers: make(map[string]*CommandHandler),
		parser:   CreateParser(),
	}
}

// Add creates a new entry in the command handlers map.
// First argument accepts a command handler implementing the type `CommandHandler`,
// Second argument accepts a variable amount of strings specifying the commands to register.
func (c *CommandSystem) Add(handler *CommandHandler, commands ...string) {
	for _, command := range commands {
		c.Handlers[command] = handler
	}
}

// Remove deletes entries from the command handlers map.
func (c *CommandSystem) Remove(commands ...string) {
	for _, command := range commands {
		delete(c.Handlers, command)
	}
}

// Handle the incoming commands and dispatches them to the appropriate
// command handler, after parsing them.
func (c *CommandSystem) Handle(session *discordgo.Session, m *discordgo.MessageCreate, command string, permissions int) {
	// Parse the command
	ast := CommandAST{}
	err := c.parser.ParseString(command, &ast)
	if err != nil {
		// Invalid command
		return
	}

	if len(ast.Arguments) > 0 {
		commandName := strings.ToLower(ast.Arguments[0].Command)
		handler, ok := c.Handlers[commandName]
		if !ok {
			// No command found
			return
		}

		if !handler.respondToDM {
			channel, err := session.State.Channel(m.ChannelID)
			if err != nil {
				log.Println("Channel lookup failed", err)
			}

			// Ignore this command if it's sent through direct message.
			if channel.Type == discordgo.ChannelTypeDM {
				return
			}
		}

		if permissions&handler.permissions != 0 || handler.permissions == -1 {
			argumentList := CommandArgList{arguments: ast.Arguments}
			perms := CommandPermissions{userPermissions: permissions}
			handler.function(m, commandName, perms, argumentList)
		} else {
			User := &util.PatchUser{m.Author}
			session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: You don't have permission for that command.", User.GetMention()))
		}
	}
}

type CommandArgList struct {
	arguments []*ArgumentAST
}

func (argList *CommandArgList) Num() int {
	return len(argList.arguments)
}

func (argList *CommandArgList) GetArg(i int) (string, error) {
	if i >= len(argList.arguments) {
		return "", errors.New("out of bounds arg index")
	}

	return argList.arguments[i].Command, nil
}

func (argList *CommandArgList) ToSlice() []string {
	arguments := make([]string, len(argList.arguments))

	for index, argument := range argList.arguments {
		arguments[index] = argument.Command
	}

	return arguments
}

type CommandPermissions struct {
	userPermissions int
}

func (perms *CommandPermissions) Test(permissions int) bool {
	return perms.userPermissions&permissions != 0 || permissions == -1
}
