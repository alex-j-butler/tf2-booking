package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

	"github.com/Qixalite/booking-api/client"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"

	"github.com/codegangsta/cli"
)

var version = "unknown"

var c *cron.Cron

// Session is an instance of the Discord client.
var Session *discordgo.Session

// BotID represents the ID of the current user.
var BotID string

// UserReportTimeouts maps user's steamids to the time after which they can report again.
var UserReportTimeouts map[string]time.Time

// Redis script to retrieve a key, and if that key does not exist, then set a default value.
var GetDefaultValue *redis.Script

// Command system
var Command *commands.CommandSystem
var IngameCommand *ingame.Command

var pool servers.ServerPool

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
			Name:    "graphtest",
			Aliases: []string{"gt"},
			Usage:   "test graphs",
			Action:  GraphTest,
		},
	}

	app.Run(os.Args)
}

func GraphTest(ctx *cli.Context) {
	now := time.Now()
	to := now.Add(-24 * (time.Hour * 7))

	graphClient := NewGraphClient("dd5930a2b34093f052aea1eeb290f11b", "a138ebaed3072d4c04e063b8ee66f686980aa794")
	// graphClient.Graph("avg:trace.rpc.request.duration{*} by {resource_name}", now, to)
	graph, err := graphClient.Graph("Chart", "sum:booking_api.servers_running{*}.rollup(max)", now, to)
	if err != nil {
		log.Fatalln(err)
	}
	ioutil.WriteFile("graph.png", graph.Bytes(), 0644)
}

// RunServer is the subcommand handler that starts the TF2 Booking server.
func RunServer(ctx *cli.Context) {
	// Create the Booking client.
	bookingClient := client.New(
		config.Conf.Booking.APIAddress,
		config.Conf.Booking.APIPort,
		config.Conf.Booking.APIKey,
	)

	// Initialise the server pool.
	pool = &servers.APIServerPool{Tag: config.Conf.Booking.Tag, APIClient: bookingClient}
	err := pool.Initialise()
	if err != nil {
		log.Println(err)
		return
	}

	SetupCron()

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
		log.Fatalln("Redis ping failed:", err)
	}
	globals.RedisClient = client

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

	// Attempt to update all our servers (that we just got from the server pool) with the information from Redis.
	// If no Redis entry exists, update Redis with the default server information.
	for _, server := range pool.GetServers() {
		err := server.Synchronise(globals.RedisClient)
		if err != nil {
			server.Update(globals.RedisClient)
		}
	}

	// Create the loghandler server
	// and bind it to the appropriate address & port.
	logs, err := loghandler.Dial(config.Conf.LogServer.LogAddress, config.Conf.LogServer.LogPort)
	if err != nil {
		// Loghandler server couldn't bind properly.
		// Not a problem, results in ingame commands not being received by the
		// booking bot.
		log.Println("LogHandler bind failed:", err)
		log.Println("NOTE: This will disable ingame commands from functioning correctly.")
	} else {
		log.Println(fmt.Sprintf("LogHandler listening on %s:%d", logs.Address, logs.Port))
	}

	logs.AddHandler(IngameMessageCreate)

	// Register the commands and their command handlers.
	Command = commands.New()
	Command.Add(
		commands.NewCommand(SyncServers),
		"sync",
	)
	Command.Add(
		commands.NewCommand(Help),
		"help",
		"/help",
	)
	Command.Add(
		commands.NewCommand(DemoLink),
		"demo",
		"demos",
		"/demo",
		"/demos",
	)
	Command.Add(
		commands.NewCommand(BookServer),
		"book",
		"/book",
	)
	Command.Add(
		commands.NewCommand(UnbookServer),
		"return",
		"unbook",
		"/return",
		"/unbook",
	)
	Command.Add(
		commands.NewCommand(ExtendServer),
		"extend",
		"/extend",
	)
	Command.Add(
		commands.NewCommand(SendPassword),
		"password",
		"string",
		"/string",
	)
	Command.Add(
		commands.NewCommand(Console),
		"console",
	)
	Command.Add(
		commands.NewCommand(Graph),
		"graph",
	)
	Command.Add(
		commands.NewCommand(PrintStats).
			Permissions(discordgo.PermissionManageServer).
			RespondToDM(true),
		"stats",
	)
	Command.Add(
		commands.NewCommand(Version).
			Permissions(discordgo.PermissionManageServer).
			RespondToDM(true),
		"version",
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
		log.Println("Discord session creation failed:", err)
		return
	}

	if config.Conf.Discord.Debug {
		dg.LogLevel = discordgo.LogDebug
	}

	// Get user information of the Discord user that is currently logged in (the bot).
	u, err := dg.User("@me")
	if err != nil {
		log.Println("Discord bot information obtain failed:", err)
		return
	}

	// Save variables for later use.
	Session = dg
	BotID = u.ID

	// Register the OnReady handler.
	dg.AddHandler(OnReady)

	// Open the Discord websocket.
	err = dg.Open()
	if err != nil {
		log.Println("Discord websocket opening failed:", err)
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
		log.Println("Game string update failed:", err)
	} else {
		log.Println("Successfully updated game string.")
	}

	// Register a message create handler.
	// This must be done in the OnGuildReady event, otherwise guild lookups would fail because of
	// it not having the list of guilds yet.
	if MessageCreateFunc != nil {
		MessageCreateFunc()
	}
	MessageCreateFunc = s.AddHandler(MessageCreate)

	log.Println("Discord bot successfully started.")
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
		log.Println("Channel lookup failed:", err)
	}

	if channel.Type == discordgo.ChannelTypeDM {
		permissionsChannelID = config.Conf.Discord.DefaultChannel
	}

	// Configuration has a string slice containing channels the bot should operate in.
	// If the channel of the newly received message is not in the slice, stop now.
	if !util.Contains(config.Conf.Discord.AcceptableChannels, m.ChannelID) && channel.Type != discordgo.ChannelTypeDM {
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
	Command.Handle(Session, m, m.Content, Permissions)
}

// IngameMessageCreate handler for the ingame TF2 log handler.
// Called when a message is sent in any TF2 server that is logging to the remote logging server.
func IngameMessageCreate(lh *loghandler.LogHandler, server *servers.Server, event *loghandler.SayEvent) {
	log.Println(fmt.Sprintf("Received command from '%s' on server '%s': %s", event.Username, server.Name, event.Message))
	IngameCommand.Handle(ingame.CommandInfo{SayEvent: *event, Server: server}, event.Message, 0)
}

// SetupCron creates the cron scheduler and adds the functions and their respective schedules.
// and finally starts the cron scheduler.
func SetupCron() {
	c = cron.New()
	c.AddFunc("0 * * * *", CheckIdleMinutes)

	c.AddFunc("@every 1m", Cron1Minute)

	c.Start()
}

// DeleteMessage deletes the specified Discord message after a certain duration has passed.
func DeleteMessage(channelID string, messageID string, duration time.Duration) error {
	time.Sleep(duration)
	return Session.ChannelMessageDelete(channelID, messageID)
}
