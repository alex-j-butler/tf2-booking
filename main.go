package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"alex-j-butler.com/tf2-booking/commands"
	"alex-j-butler.com/tf2-booking/util"

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
var UserServers map[string]*Server

// Command system
var Command *commands.Command

func main() {
	InitialiseConfiguration()
	SetupCron()

	// Register the commands and their command handlers.
	Command = commands.New("")
	Command.Add(BookServer, "book a server", "book")
	Command.Add(UnbookServer, "return", "return a server", "unbook", "unbook a server")
	Command.Add(ExtendServer, "extend", "extend a server", "extend my server", "extend booking", "extend my booking")
	Command.Add(Update, "update")
	Command.Add(Print, "print", "print state")

	// Create maps.
	Users = make(map[string]bool)
	UserServers = make(map[string]*Server)

	// Restore state.
	if HasState(".state.json") {
		err, servers, users, userServers := LoadState(".state.json")

		if err != nil {
			log.Println("Found state file, failed to restore:", err)
		} else {
			log.Println("Found state file, restoring from previous state.")

			if err = DeleteState(".state.json"); err != nil {
				log.Println("Failed to delete state file:", err)
			}

			Conf.Servers = servers
			Users = users
			UserServers = userServers
		}
	}

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

	// Get the next available server.
	Serv := GetAvailableServer()

	if Serv != nil {
		// Book the server.
		RCONPassword, ServerPassword, err := Serv.Book(m.Author)
		if err != nil {
			Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Something went wrong while trying to book your server, please try again later.", User.GetMention()))
		} else {
			// Start the server.
			go func() {
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
			}()

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
		// Stop the server.
		go func() {
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
		}()

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
		// Notify Discord channel to let the user know they do not have a server booked.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: You haven't booked a server. Type `book` to book a server.", User.GetMention()))

		return
	}

	if Serv, ok := UserServers[m.Author.ID]; ok && Serv != nil {
		// Extend the booking.
		Serv.ExtendBooking(Conf.BookingExtendDuration.Duration)

		// Notify server of successful operation.
		Serv.SendCommand(fmt.Sprintf("say @%s: Your booking has been extended by %s.", m.Author.Username, Conf.BookingExtendDurationText))

		// Notify Discord channel of successful operation.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Your booking has been extended by %s.", User.GetMention(), Conf.BookingExtendDurationText))
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

func Update(m *discordgo.MessageCreate, command string, args []string) {
	User := &PatchUser{m.Author}

	SaveState(".state.json", Conf.Servers, Users, UserServers)
	UpdateExecutable("http://localhost/tf2-booking/latest")

	Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Updated `tf2-booking` & restarting now.", User.GetMention()))

	os.Exit(0)
}

func Print(m *discordgo.MessageCreate, command string, args []string) {
	User := &PatchUser{m.Author}

	Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Update v2", User.GetMention()))
	Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Servers state: %v", User.GetMention(), Conf.Servers))
	Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Users state: %v", User.GetMention(), Users))
	Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: UserServers state: %v", User.GetMention(), UserServers))
}

// MessageCreate handler for Discord.
// Called when a message is received by the Discord client.
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Do not process messages that were created by the current user.
	if m.Author.ID == BotID {
		return
	}

	// Configuration has a string slice containing channels the bot should operate in.
	// If the channel of the newly received message is not in the slice, stop now.
	if !util.Contains(Conf.AcceptableChannels, m.ChannelID) {
		return
	}

	// Send the message content to the command handler to be dispatched appropriately.
	Command.Handle(m, strings.ToLower(m.Content))
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
