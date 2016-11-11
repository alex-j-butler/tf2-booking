package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"alex-j-butler.com/tf2-booking/commands"
	"alex-j-butler.com/tf2-booking/commands/ingame/loghandler"
	"alex-j-butler.com/tf2-booking/config"
	"alex-j-butler.com/tf2-booking/servers"
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

// Command system
var Command *commands.Command

func main() {
	config.InitialiseConfiguration()
	SetupCron()

	logs, err := loghandler.Dial("", 3001)
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

	// Register the ingame commands and their command handlers.
	// IngameCommand := ingame.New("")
	// IngameCommand.Add(
	// 	ingame.NewCommand(ReportServer),
	// 	"report",
	// )

	// Create maps.
	Users = make(map[string]bool)
	UserServers = make(map[string]*servers.Server)

	// Create the Discord client from the bot token in the configuration.
	dg, err := discordgo.New(fmt.Sprintf("Bot %s", config.Conf.DiscordToken))
	if err != nil {
		log.Println("Failed to create Discord session:", err)
		return
	}

	if config.Conf.DiscordDebug {
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

// BookServer command handler
// Called when a user types the 'book' command into the Discord channel.
// This function checks whether the user has a server booked, if not,
// it books a new server, preventing it from being used by another user,
// sets up the RCON password & Server Password and finally starts the TF2 server.
func BookServer(m *discordgo.MessageCreate, command string, args []string) {
	User := &util.PatchUser{m.Author}

	// Check if the user has already booked a server out.
	if value, _ := Users[m.Author.ID]; value == true {
		// Send a message to let the user know they've already booked a server.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: You've already booked a server. Type `unbook` to return the server.", User.GetMention()))
		return
	}

	// Get the next available server.
	Serv := servers.GetAvailableServer(config.Conf.Servers)

	if Serv != nil {
		// Book the server.
		RCONPassword, ServerPassword, err := Serv.Book(m.Author, config.Conf.BookingDuration.Duration)
		if err != nil {
			Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Something went wrong while trying to book your server, please try again later.", User.GetMention()))
		} else {
			// Start the server.
			go func(Serv *servers.Server, m *discordgo.MessageCreate) {
				err := Serv.Start()

				if err != nil {
					UserChannel, _ := Session.UserChannelCreate(m.Author.ID)
					Session.ChannelMessageSend(
						UserChannel.ID,
						fmt.Sprintf(
							"Uh oh! The server failed to start, contact an admin for further information.",
						),
					)

					Users[m.Author.ID] = false
					UserServers[m.Author.ID] = nil

					UpdateGameString()

					log.Println(fmt.Sprintf("Failed to start server \"%s\" from \"%s\"", Serv.Name, m.Author.ID))
				}
			}(Serv, m)

			// Send message to public channel, without server details.
			Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Server details have been sent via private message.", User.GetMention()))

			// Send message to private DM, with server details.
			UserChannel, _ := Session.UserChannelCreate(m.Author.ID)
			Session.ChannelMessageSend(
				UserChannel.ID,
				fmt.Sprintf(
					"Here is your server:\n\tServer address: %s\n\tRCON Password: %s\n\tPassword: %s\n\tConnect string: `connect %s; password %s`",
					Serv.Address,
					RCONPassword,
					ServerPassword,
					Serv.Address,
					ServerPassword,
				),
			)

			Users[m.Author.ID] = true
			UserServers[m.Author.ID] = Serv

			UpdateGameString()

			log.Println(fmt.Sprintf("Booked server \"%s\" from \"%s\"", Serv.Name, m.Author.ID))
		}
	} else {
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: No servers are currently available.", User.GetMention()))
	}
}

// UnbookServer command handler
// Called when a user types the 'unbook' command into the Discord channel.
// This function checks whether the user has a server booked, if so,
// it unbooks it, allowing it for use by another user, and shutting down
// the TF2 server.
func UnbookServer(m *discordgo.MessageCreate, command string, args []string) {
	User := &util.PatchUser{m.Author}

	// Check if the user has already booked a server out.
	if value, ok := Users[m.Author.ID]; !ok || value == false {
		// Send a message to let the user know they do not have a server booked.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: You haven't booked a server. Type `book` to book a server.", User.GetMention()))

		return
	}

	if Serv, ok := UserServers[m.Author.ID]; ok && Serv != nil {
		// Stop the server.
		go func(Serv *servers.Server, m *discordgo.MessageCreate) {
			err := Serv.Stop()

			if err != nil {
				UserChannel, _ := Session.UserChannelCreate(m.Author.ID)
				Session.ChannelMessageSend(
					UserChannel.ID,
					fmt.Sprintf(
						"Uh oh! The server failed to stop, contact an admin for further information, or leave us to handle it.",
					),
				)
			}
		}(Serv, m)

		// Remove the user's booked state.
		Users[m.Author.ID] = false
		UserServers[m.Author.ID] = nil

		// Unbook the server.
		Serv.Unbook()

		// Upload STV demos
		STVMessage, err := Serv.UploadSTV()

		// Send 'returned' message.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Server returned.", User.GetMention()))

		// Send 'stv' message, if it uploaded successfully.
		if err == nil {
			Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: %s", User.GetMention(), STVMessage))
		}

		UpdateGameString()

		log.Println(fmt.Sprintf("Unbooked server \"%s\" from \"%s\"", Serv.Name, m.Author.ID))
	} else {
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: You haven't booked a server. Type `book` to book a server.", User.GetMention()))

		// We're in an invalid state, reset back to normal.
		Users[m.Author.ID] = false
		UserServers[m.Author.ID] = nil

		return
	}
}

// ExtendServer command handler
// Called when a user types the 'extend' command into the Discord channel.
// This function checks whether the user has a server booked out, if so,
// it will extend the booking by adding time onto the servers return time.
func ExtendServer(m *discordgo.MessageCreate, command string, args []string) {
	User := &util.PatchUser{m.Author}

	// Check if the user has already booked a server out.
	if value, ok := Users[m.Author.ID]; !ok || value == false {
		// Notify Discord channel to let the user know they do not have a server booked.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: You haven't booked a server. Type `book` to book a server.", User.GetMention()))

		return
	}

	if Serv, ok := UserServers[m.Author.ID]; ok && Serv != nil {
		// Extend the booking.
		Serv.ExtendBooking(config.Conf.BookingExtendDuration.Duration)

		// Notify server of successful operation.
		Serv.SendCommand(fmt.Sprintf("say @%s: Your booking has been extended by %s.", m.Author.Username, config.Conf.BookingExtendDurationText))

		// Notify Discord channel of successful operation.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Your booking has been extended by %s.", User.GetMention(), config.Conf.BookingExtendDurationText))
	} else {
		// Notify Discord channel of failed operation.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: You haven't booked a server. Type `book` to book a server.", User.GetMention()))

		// If the program execution reaches here, the state of the users & user-servers map
		// is invalid and should be reset to the 'unbooked' state.
		Users[m.Author.ID] = false
		UserServers[m.Author.ID] = nil

		return
	}
}

func PrintStats(m *discordgo.MessageCreate, command string, args []string) {
	User := &util.PatchUser{m.Author}

	servers := servers.GetBookedServers(config.Conf.Servers)
	message := "Server stats:"
	count := 0

	for i := 0; i < len(servers); i++ {
		server := servers[i]
		if server != nil {
			message = fmt.Sprintf("%s\n\t%s: %f", message, server.Name, server.TickRate)
			count++
		}
	}

	if count == 0 {
		message = "No servers are currently booked."
	}

	Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: %s", User.GetMention(), message))
}

func Update(m *discordgo.MessageCreate, command string, args []string) {
	User := &util.PatchUser{m.Author}

	// Delete the sent message in 10 seconds.
	go DeleteMessage(m.ChannelID, m.ID, time.Second*10)

	if len(args) <= 0 {
		// Send error.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Invalid arguments: `update <url>`", User.GetMention()))
		return
	}

	url := strings.Join(args, " ")

	Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Starting update...", User.GetMention()))

	go func(url string) {
		SaveState(".state.json", config.Conf.Servers, Users, UserServers)
		UpdateExecutable(url)

		m, _ := Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Updated `tf2-booking` & restarting now from URL: %s", User.GetMention(), url))

		// Delete the sent message in 10 seconds.
		go DeleteMessage(m.ChannelID, m.ID, time.Second*10)

		wait.Exit()
	}(url)
}

func Exit(m *discordgo.MessageCreate, command string, args []string) {
	User := &util.PatchUser{m.Author}

	Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Shutting down `tf2-booking`.", User.GetMention()))

	SaveState(".state.json", config.Conf.Servers, Users, UserServers)
	wait.Exit()
}

// OnReady handler for Discord.
// Called when the connection has been completely setup.
func OnReady(s *discordgo.Session, r *discordgo.Ready) {
	// Restore state from the state file, if it exists.
	if HasState(".state.json") {
		err, servers, users, userServers := LoadState(".state.json")

		if err != nil {
			log.Println("Found state file, failed to restore:", err)
		} else {
			log.Println("Found state file, restoring from previous state.")

			if err = DeleteState(".state.json"); err != nil {
				log.Println("Failed to delete state file:", err)
			}

			config.Conf.Servers = servers
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
		permissionsChannelID = config.Conf.DefaultChannel
	}

	// Configuration has a string slice containing channels the bot should operate in.
	// If the channel of the newly received message is not in the slice, stop now.
	if !util.Contains(config.Conf.AcceptableChannels, m.ChannelID) && !channel.IsPrivate {
		return
	}

	Permissions, err := Session.State.UserChannelPermissions(m.Author.ID, permissionsChannelID)
	if err != nil {
		log.Println("Failed to lookup permissions.", err)
	}

	// Send the message content to the command handler to be dispatched appropriately.
	Command.Handle(Session, m, strings.ToLower(m.Content), Permissions)
}

func IngameMessageCreate(lh *loghandler.LogHandler, server *servers.Server, event loghandler.SayEvent) {
	log.Println(fmt.Sprintf("Received command from '%s' on server '%s': %s", event.Username, server.Name, event.Message))
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
