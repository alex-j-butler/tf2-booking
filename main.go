package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
)

var c *cron.Cron
var Session *discordgo.Session
var BotID string
var Users map[string]bool
var UserServers map[string]*Server

func main() {
	InitialiseConfiguration()
	SetupCron()

	Users = make(map[string]bool)
	UserServers = make(map[string]*Server)

	dg, err := discordgo.New(fmt.Sprintf("Bot %s", Conf.DiscordToken))
	if err != nil {
		log.Println("Failed to create Discord session:", err)
		return
	}

	// Get bot information.
	u, err := dg.User("@me")
	if err != nil {
		log.Println("Failed to obtain Discord bot information:", err)
		return
	}

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

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}

	Message := strings.ToLower(m.Content)

	switch Message {
	case "book a server", "book":
		User := &PatchUser{m.Author}

		// Check if the user has already booked a server out.
		if value, _ := Users[m.Author.ID]; value == true {
			// Send a message to let the user know they've already booked a server.
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: You've already booked a server. Type `unbook` to return the server.", User.GetMention()))
			return
		}

		Serv := GetAvailableServer()

		if Serv != nil {
			// Book the server.
			RCONPassword, ServerPassword, err := Serv.Book(m.Author)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Something went wrong while trying to book your server, please try again later.", User.GetMention()))
			} else {
				// Send message to public channel, without server details.
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Server details have been sent via private message.", User.GetMention()))

				// Send message to private DM, with server details.
				UserChannel, _ := s.UserChannelCreate(m.Author.ID)
				s.ChannelMessageSend(
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
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: No servers are currently available.", User.GetMention()))
		}
	case "return", "return a server", "unbook", "unbook a server":
		User := &PatchUser{m.Author}

		// Check if the user has already booked a server out.
		if value, ok := Users[m.Author.ID]; !ok || value == false {
			// Send a message to let the user know they do not have a server booked.
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: You haven't booked a server. Type `book` to book a server.", User.GetMention()))

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
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Server returned.", User.GetMention()))

			// Send 'stv' message, if it uploaded successfully.
			if err == nil {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: %s", User.GetMention(), STVMessage))
			}

			UpdateGameString()

			log.Println(fmt.Sprintf("Unbooked server \"%s\" from \"%s\"", Serv.Name, m.Author.ID))

			return
		} else {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: You haven't booked a server. Type `book` to book a server.", User.GetMention()))

			// We're in an invalid state, reset back to normal.
			Users[m.Author.ID] = false
			UserServers[m.Author.ID] = nil

			return
		}

		log.Println("Unbooking a server for", m.Author.Username, "!")
	}
}

func SetupCron() {
	c = cron.New()
	c.AddFunc("*/1 * * * *", CheckServers)
	c.Start()
}

func CheckServers() {
	// Iterate through servers.
	for i := 0; i < len(Conf.Servers); i++ {
		since := time.Since(Conf.Servers[i].GetBookedTime())
		if !Conf.Servers[i].IsAvailable() && since > (4*time.Hour) {
			UserID := Conf.Servers[i].GetBooker()
			UserMention := Conf.Servers[i].GetBookerMention()

			// Remove the user's booked state.
			Users[UserID] = false
			UserServers[UserID] = nil

			// Unbook the server.
			Conf.Servers[i].Unbook()

			// Upload STV demos
			STVMessage, err := Conf.Servers[i].UploadSTV()

			// Send 'returned' message
			Session.ChannelMessageSend(Conf.DefaultChannel, fmt.Sprintf("%s: Your server was automatically unbooked.", UserMention))

			// Send 'stv' message, if it uploaded successfully.
			if err == nil {
				Session.ChannelMessageSend(Conf.DefaultChannel, fmt.Sprintf("%s: %s", UserMention, STVMessage))
			}

			UpdateGameString()

			log.Println(fmt.Sprintf("Automatically unbooked server \"%s\" from \"%s\"", Conf.Servers[i].Name, UserID))
		}
	}
}
