package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var BotID string
var Users map[string]bool

func main() {
	InitialiseConfiguration()
	SetupServers()

	Users = make(map[string]bool)

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

	BotID = u.ID

	// Register a message create handler.
	dg.AddHandler(MessageCreate)

	// Open the Discord websocket.
	err = dg.Open()
	if err != nil {
		log.Println("Failed to open Discord websocket:", err)
		return
	}

	log.Println("Discord bot successfully started.")

	// Keep running until Control-C pressed.
	<-make(chan struct{})
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}

	Message := strings.ToLower(m.Content)

	switch Message {
	case "book a server", "book":
		log.Println("Booking a server for", m.Author.Username, "!")

		User := &PatchUser{m.Author}

		if value, _ := Users[m.Author.ID]; value == true {
			// User has already booked a server.
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
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: No servers are currently available.", User.GetMention()))
		}
	case "return", "return a server", "unbook", "unbook a server":
		log.Println("Unbooking a server for", m.Author.Username, "!")
	}
}
