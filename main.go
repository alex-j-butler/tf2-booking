package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"alex-j-butler.com/tf2-booking/commands"
	"alex-j-butler.com/tf2-booking/commands/ingame"
	"alex-j-butler.com/tf2-booking/commands/ingame/loghandler"
	"alex-j-butler.com/tf2-booking/config"
	"alex-j-butler.com/tf2-booking/database"
	"alex-j-butler.com/tf2-booking/servers"
	"alex-j-butler.com/tf2-booking/steamauth"
	"alex-j-butler.com/tf2-booking/util"
	"alex-j-butler.com/tf2-booking/wait"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
)

var c *cron.Cron

// Session is an instance of the Discord client.
var Session *discordgo.Session

// BotID represents the ID of the current user.
var BotID string

// Users maps user IDs to booking state (true = booked, false = unbooked)
var Users map[string]bool

// UserServers maps user IDs to server pointers.
var UserServers map[string]*servers.Server

// UserReportTimeouts maps user's steamids to the time after which they can report again.
var UserReportTimeouts map[string]time.Time

// Command system
var Command *commands.Command
var IngameCommand *ingame.Command

var httpServer *steamauth.HTTPServer

func main() {
	config.InitialiseConfiguration()
	servers.InitialiseServers()

	database.Initialise()
	defer database.DB.Close()

	httpServer := steamauth.New(config.Conf.SteamAuthServer.Address, config.Conf.SteamAuthServer.Port, config.Conf.SteamAuthServer.RootURL)
	httpServer.AddHandler(AccountLinked)
	go httpServer.Run()

	SetupCron()

	logs, err := loghandler.Dial(config.Conf.LogServer.LogAddress, config.Conf.LogServer.LogPort)
	if err != nil {
		log.Println("LogHandler failed to connect:", err)
	} else {
		log.Println(fmt.Sprintf("LogHandler listening on %s:%d", logs.Address, logs.Port))
	}

	logs.AddHandler(IngameMessageCreate)

	// Register the commands and their command handlers.
	Command = commands.New("")
	Command.Add(
		commands.NewCommand(BookServer),
		"book",
	)
	Command.Add(
		commands.NewCommand(UnbookServer),
		"return",
		"unbook",
	)
	Command.Add(
		commands.NewCommand(ExtendServer),
		"extend",
	)
	Command.Add(
		commands.NewCommand(PrintStats).
			Permissions(discordgo.PermissionManageServer).
			RespondToDM(true),
		"stats",
	)
	Command.Add(
		commands.NewCommand(Update).
			Permissions(discordgo.PermissionManageServer).
			RespondToDM(true),
		"update",
	)
	Command.Add(
		commands.NewCommand(Exit).
			Permissions(discordgo.PermissionManageServer).
			RespondToDM(true),
		"exit",
	)

	Command.Add(
		commands.NewCommand(Link),
		"link",
	)

	// Register the ingame commands and their command handlers.
	IngameCommand = ingame.New("!")
	IngameCommand.Add(
		ingame.NewCommand(ReportServer),
		"report",
	)
	IngameCommand.Add(
		ingame.NewCommand(TimeLeft),
		"time",
	)

	// Create maps.
	Users = make(map[string]bool)
	UserServers = make(map[string]*servers.Server)
	UserReportTimeouts = make(map[string]time.Time)

	// Create the Discord client from the bot token in the configuration.
	dg, err := discordgo.New(fmt.Sprintf("Bot %s", config.Conf.Discord.Token))
	if err != nil {
		log.Println("Failed to create Discord session:", err)
		return
	}

	if config.Conf.Discord.Debug {
		dg.LogLevel = discordgo.LogDebug
	}

	// Get user information of the Discord user that is currently logged in (the bot).
	u, err := dg.User("@me")
	if err != nil {
		log.Println("Failed to obtain Discord bot information:", err)
		return
	}

	// Save variables for later use.
	Session = dg
	BotID = u.ID

	// Register the OnReady handler.
	dg.AddHandler(OnReady)
	dg.AddHandler(OnGuildReady)

	// Open the Discord websocket.
	err = dg.Open()
	if err != nil {
		log.Println("Failed to open Discord websocket:", err)
		return
	}

	// Keep running until Control-C pressed.
	// <-make(chan struct{})
	wait.Wait()

	Session.Close()

	// Stop cron.
	c.Stop()
}

// OnReady handler for Discord.
// Called when the connection has been completely setup.
func OnReady(s *discordgo.Session, r *discordgo.Ready) {
	// Restore state from the state file, if it exists.
	if HasState(".state.json") {
		err, servers_, users, userServers := LoadState(".state.json")

		if err != nil {
			log.Println("Found state file, failed to restore:", err)
		} else {
			log.Println("Found state file, restoring from previous state.")

			if err = DeleteState(".state.json"); err != nil {
				log.Println("Failed to delete state file:", err)
			}

			servers.Servers = servers_
			Users = users
			UserServers = userServers
		}
	}

	log.Println("Updated game string.")
	UpdateGameString()

	log.Println("Discord bot successfully started.")
}

// OnGuildReady handler for Discord.
// Called when all the guilds have been lazy loaded.
func OnGuildReady(s *discordgo.Session, r *discordgo.GuildReady) {
	// Register a message create handler.
	// This must be done in the OnGuildReady event, otherwise guild lookups would fail because of
	// it not having the list of guilds yet.
	s.AddHandler(MessageCreate)

	log.Println("Discord guilds successfully loaded.")
}

// MessageCreate handler for Discord.
// Called when a message is received by the Discord client.
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Do not process messages that were created by the current user.
	if m.Author.ID == BotID {
		return
	}

	permissionsChannelID := m.ChannelID

	// Lookup Discord channel.
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Println("Failed to lookup channels.", err)
	}

	if channel.IsPrivate {
		permissionsChannelID = config.Conf.Discord.DefaultChannel
	}

	// Configuration has a string slice containing channels the bot should operate in.
	// If the channel of the newly received message is not in the slice, stop now.
	if !util.Contains(config.Conf.Discord.AcceptableChannels, m.ChannelID) && !channel.IsPrivate {
		return
	}

	Permissions, err := Session.State.UserChannelPermissions(m.Author.ID, permissionsChannelID)
	if err != nil {
		log.Println("Failed to lookup permissions.", err)
	}

	// Send the message content to the command handler to be dispatched appropriately.
	Command.Handle(Session, m, strings.ToLower(m.Content), Permissions)
}

func IngameMessageCreate(lh *loghandler.LogHandler, server *servers.Server, event *loghandler.SayEvent) {
	log.Println(fmt.Sprintf("Received command from '%s' on server '%s': %s", event.Username, server.Name, event.Message))
	IngameCommand.Handle(ingame.CommandInfo{SayEvent: *event, Server: server}, event.Message, 0)
}

// AccountLinked is the event handler for the LinkSuccessEvent event.
func AccountLinked(httpServer *steamauth.HTTPServer, event *steamauth.LinkSuccessEvent) {
	log.Println("Account linked:", event)
}

// SetupCron creates the cron scheduler and adds the functions and their respective schedules.
// and finally starts the cron scheduler.
func SetupCron() {
	c = cron.New()
	c.AddFunc("*/1 * * * *", CheckUnbookServers)
	c.AddFunc("0 * * * *", CheckIdleMinutes)
	c.AddFunc("*/10 * * * *", CheckStats)
	c.Start()
}

func DeleteMessage(channelID string, messageID string, duration time.Duration) error {
	time.Sleep(duration)
	return Session.ChannelMessageDelete(channelID, messageID)
}
