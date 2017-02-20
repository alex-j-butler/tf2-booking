package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	redis "gopkg.in/redis.v5"

	"alex-j-butler.com/tf2-booking/commands"
	"alex-j-butler.com/tf2-booking/commands/ingame"
	"alex-j-butler.com/tf2-booking/commands/ingame/loghandler"
	"alex-j-butler.com/tf2-booking/config"
	"alex-j-butler.com/tf2-booking/globals"
	"alex-j-butler.com/tf2-booking/servers"
	"alex-j-butler.com/tf2-booking/util"
	"alex-j-butler.com/tf2-booking/wait"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"

	"github.com/codegangsta/cli"
	_ "github.com/lib/pq"
)

var c *cron.Cron

// Session is an instance of the Discord client.
var Session *discordgo.Session

// BotID represents the ID of the current user.
var BotID string

// UserReportTimeouts maps user's steamids to the time after which they can report again.
var UserReportTimeouts map[string]time.Time

var GetDefaultValue *redis.Script

// Command system
var Command *commands.Command
var IngameCommand *ingame.Command

// MessageCreateFunc stores the function that deletes the MessageCreate Discord event handler.
// This is a fix for the bot receiving messages twice.
// When the Discord client times out, it reconnects, calling the 'OnGuildReady' event again, which adds a new MessageCreate handler, without removing the old one.
// This means the messages are being received and processed twice.
// By storing the latest MessageCreate delete function, it can delete the previous MessageCreate handler before adding the new one.
var MessageCreateFunc func()

func main() {
	config.InitialiseConfiguration()

	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "run the server",
			Action:  RunServer,
		},
		{
			Name:    "migrate",
			Aliases: []string{"m"},
			Usage:   "migrate the database",
			Action:  Migrate,
		},
	}

	app.Run(os.Args)
}

func Migrate(ctx *cli.Context) {
	log.Println("Error: Migrate is currently unimplemented")
}

func RunServer(ctx *cli.Context) {
	servers.InitialiseServers()
	SetupCron()

	// Connect to the PostgreSQL database.
	db, err := sql.Open("postgres", config.Conf.Database.DSN)
	if err != nil {
		log.Println("Database error:", err)
		os.Exit(1)
	}
	globals.DB = db

	// Ping the database to make sure we're properly connected.
	if err := globals.DB.Ping(); err != nil {
		log.Println("Database error:", err)
		os.Exit(1)
	}

	// Setup the Redis client
	// and PING it to make sure we properly connected
	// and can issue commands to it.
	client := redis.NewClient(&redis.Options{
		Addr:     config.Conf.Redis.Address,
		Password: config.Conf.Redis.Password,
		DB:       config.Conf.Redis.DB,
	})

	_, err = client.Ping().Result()
	if err != nil {
		// Application won't work without a Redis connection.
		log.Println("Redis error:", err)
		os.Exit(1)
	}
	globals.RedisClient = client

	// When the booking bot starts, we need to insert all the servers that we know about
	// that do not currently exist in Redis (which is done through the SETNX Redis command),
	// after which, it will synchronise all of the servers from Redis.
	for i, server := range servers.Servers {
		// Serialise the server as a JSON string.
		serialised, err := json.Marshal(server)
		if err != nil {
			panic(err)
		}

		// Attempt to add the server.
		err = client.SetNX(fmt.Sprintf("server.%s", server.SessionName), serialised, 0).Err()
		if err != nil {
			panic(err)
		}

		// Synchronise the server from Redis, to get information for existing servers.
		err = server.Synchronise(globals.RedisClient)
		if err != nil {
			panic(err)
		}

		// Put the modified server back.
		servers.Servers[i] = server
	}

	// Create Redis scripts.
	// Check if the user has already booked a server out.
	GetDefaultValue = redis.NewScript(`
		local value = redis.call("GET", KEYS[1])
		if (not value) then
			redis.call("SET", KEYS[1], ARGV[1])
			return ARGV[1]
		end
		return value
	`)

	// Create the loghandler server
	// and bind it to the appropriate address & port.
	logs, err := loghandler.Dial(config.Conf.LogServer.LogAddress, config.Conf.LogServer.LogPort)
	if err != nil {
		// Loghandler server couldn't bind properly.
		// Not a problem, results in ingame commands not being received by the
		// booking bot.
		log.Println("LogHandler failed to bind:", err)
		log.Println("NOTE: This will disable ingame commands from functioning correctly.")
	} else {
		log.Println(fmt.Sprintf("LogHandler listening on %s:%d", logs.Address, logs.Port))
	}

	logs.AddHandler(IngameMessageCreate)

	// Register the commands and their command handlers.
	Command = commands.New("")
	Command.Add(
		commands.NewCommand(DebugPrint),
		"debug",
	)
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
		commands.NewCommand(SendPassword),
		"send password",
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
	log.Println("Updating game string with currently booked servers.")
	err := UpdateGameString()
	if err != nil {
		log.Println("Failed updating game string:", err)
	} else {
		log.Println("Successfully updated game string.")
	}
	log.Println("Discord bot successfully started.")
}

// OnGuildReady handler for Discord.
// Called when all the guilds have been lazy loaded.
func OnGuildReady(s *discordgo.Session, r *discordgo.GuildReady) {
	// Register a message create handler.
	// This must be done in the OnGuildReady event, otherwise guild lookups would fail because of
	// it not having the list of guilds yet.
	if MessageCreateFunc != nil {
		MessageCreateFunc()
	}
	MessageCreateFunc = s.AddHandler(MessageCreate)

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
		// Grab the timestamp of this error in GMT+10 time.
		gmt10 := time.FixedZone("GMT+10", 10*60*60)
		timestamp := time.Now().In(gmt10)

		log.Println("discord error: failed to lookup permissions.", err, fmt.Sprintf("(id %s name %s time %s)", m.Author.ID, m.Author.Username, timestamp.String()))

		// Assume permissions = 0
		Permissions = 0

		// Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Sorry, we couldn't look up your Discord permissions, please contact an admin for assistance. (id %s time %s)", fmt.Sprintf("<@%s>", m.Author.ID), m.Author.ID, timestamp.String()))
	}

	// Send the message content to the command handler to be dispatched appropriately.
	Command.Handle(Session, m, strings.ToLower(m.Content), Permissions)
}

func IngameMessageCreate(lh *loghandler.LogHandler, server *servers.Server, event *loghandler.SayEvent) {
	log.Println(fmt.Sprintf("Received command from '%s' on server '%s': %s", event.Username, server.Name, event.Message))
	IngameCommand.Handle(ingame.CommandInfo{SayEvent: *event, Server: server}, event.Message, 0)
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
