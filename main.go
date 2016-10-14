package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var BotID string

func main() {
	InitialiseConfiguration()
	SetupServers()

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
		Serv := GetAvailableServer()

		if Serv != nil {
			// Book the server.
			_, _, err := Serv.Book(m.Author)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Something went wrong while trying to book your server, please try again later.", User.GetMention()))
			} else {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Server details have been sent via private message.", User.GetMention()))

				// Send message.
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: No servers are currently available.", User.GetMention()))
		}

		// RCON, Server, err := Serv.Book(m.Author)
		// if err != nil {
		//	log.Println("Booking server failed:", err)
		//}
		// log.Println("Server booked:", RCON, Server)
	case "return", "return a server", "unbook", "unbook a server":
		log.Println("Unbooking a server for", m.Author.Username, "!")
	}
}
