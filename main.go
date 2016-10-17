package main

import (
	"fmt"
	"log"
	"strings"

	"alex-j-butler.com/tf2-booking/commands"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
)

var c *cron.Cron
var Session *discordgo.Session
var BotID string
var Users map[string]bool
var UserServers map[string]*Server

var Command *commands.Command

func main() {
	InitialiseConfiguration()
	SetupCron()

	// Register the commands and their command handlers.
	Command = commands.New("")
	Command.Add(BookServer, "book a server", "book")
	Command.Add(UnbookServer, "return", "return a server", "unbook", "unbook a server")
	Command.Add(ExtendServer, "extend", "extend a server", "extend my server", "extend booking", "extend my booking")

	Users = make(map[string]bool)
	UserServers = make(map[string]*Server)

	// Create the Discord client from the bot token in the configuration.
	dg, err := discordgo.New(fmt.Sprintf("Bot %s", Conf.DiscordToken))
	if err != nil {
		log.Println("Failed to create Discord session:", err)
		return
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

	// Register a message create handler.
	dg.AddHandler(MessageCreate)

	// Open the Discord websocket.
	err = dg.Open()
	if err != nil {
		log.Println("Failed to open Discord websocket:", err)
		return
	}

	log.Println("Updated game string.")
	UpdateGameString()

	log.Println("Discord bot successfully started.")

	// Keep running until Control-C pressed.
	<-make(chan struct{})

	// Stop cron.
	c.Stop()
}

// BookServer command handler
// Called when a user types the 'book' command into the Discord channel.
// This function checks whether the user has a server booked, if not,
// it books a new server, preventing it from being used by another user,
// sets up the RCON password & Server Password and finally starts the TF2 server.
func BookServer(m *discordgo.MessageCreate, command string, args []string) {
	User := &PatchUser{m.Author}

	// Check if the user has already booked a server out.
	if value, _ := Users[m.Author.ID]; value == true {
		// Send a message to let the user know they've already booked a server.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: You've already booked a server. Type `unbook` to return the server.", User.GetMention()))
		return
	}

	Serv := GetAvailableServer()

	if Serv != nil {
		// Book the server.
		RCONPassword, ServerPassword, err := Serv.Book(m.Author)
		if err != nil {
			Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Something went wrong while trying to book your server, please try again later.", User.GetMention()))
		} else {
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
	User := &PatchUser{m.Author}

	// Check if the user has already booked a server out.
	if value, ok := Users[m.Author.ID]; !ok || value == false {
		// Send a message to let the user know they do not have a server booked.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: You haven't booked a server. Type `book` to book a server.", User.GetMention()))

		return
	}

	if Serv, ok := UserServers[m.Author.ID]; ok && Serv != nil {
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
	User := &PatchUser{m.Author}

	// Check if the user has already booked a server out.
	if value, ok := Users[m.Author.ID]; !ok || value == false {
		// Send a message to let the user know they do not have a server booked.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: You haven't booked a server. Type `book` to book a server.", User.GetMention()))

		return
	}

	if Serv, ok := UserServers[m.Author.ID]; ok && Serv != nil {
		// Extend the booking.
		Serv.ExtendBooking(Conf.BookingExtendDuration.Duration)

		// Send server message.
		Serv.SendCommand(fmt.Sprintf("say @%s: Your booking has been extended by 2 hours.", m.Author.Username))

		// Send Discord message.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Your booking has been extended by 2 hours.", User.GetMention()))
	} else {
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: You haven't booked a server. Type `book` to book a server.", User.GetMention()))

		// We're in an invalid state, reset back to normal.
		Users[m.Author.ID] = false
		UserServers[m.Author.ID] = nil

		return
	}
}

// MessageCreate handler for Discord.
// Called when a message is received by the Discord client.
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}

	// Configuration has a string slice containing channels the bot should operate in.
	// If the channel of the newly received message is not in the slice, stop now.
	if !Contains(Conf.AcceptableChannels, m.ChannelID) {
		return
	}

	// Handle command.
	Command.Handle(m, strings.ToLower(m.Content))
}

// SetupCron creates the cron scheduler and adds the functions and their respective schedules.
// and finally starts the cron scheduler.
func SetupCron() {
	c = cron.New()
	c.AddFunc("*/1 * * * *", CheckUnbookServers)
	c.AddFunc("0 * * * *", CheckIdleMinutes)
	c.Start()
}
