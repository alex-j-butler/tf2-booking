package triggers

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"alex-j-butler.com/tf2-booking/util"

	"github.com/bwmarrin/discordgo"
)

type CommandFunction func(*discordgo.MessageCreate, string, []string)

type CommandHandler struct {
	function    CommandFunction
	permissions int
	respondToDM bool
}

type Command struct {
	Prefix   string
	Handlers map[string]*CommandHandler
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
func (c *Command) Handle(session *discordgo.Session, m *discordgo.MessageCreate, command string, permissions int) {
	for str, handler := range c.Handlers {
		if !strings.HasPrefix(command, c.Prefix) && c.Prefix != "" {
			continue
		}

		handlerSplit := strings.Split(str, " ")
		commandSplit := strings.Split(command[len(c.Prefix):], " ")

		if len(commandSplit) < len(handlerSplit) {
			continue
		}

		if !handler.respondToDM {
			channel, err := session.State.Channel(m.ChannelID)
			if err != nil {
				log.Println("Failed to lookup channel.", err)
			}

			if channel.IsPrivate {
				continue
			}
		}

		if reflect.DeepEqual(handlerSplit, commandSplit[:len(handlerSplit)]) {
			log.Println(fmt.Sprintf("Permissions test: %d & %d = %d", permissions, handler.permissions, permissions&handler.permissions))

			if permissions&handler.permissions != 0 || handler.permissions == -1 {
				handler.function(m, strings.Join(handlerSplit, " "), commandSplit[len(handlerSplit):])
			} else {
				User := &util.PatchUser{m.Author}
				session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: You don't have permission for that command.", User.GetMention()))
			}

			break
		}
	}
}
