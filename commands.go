package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"alex-j-butler.com/tf2-booking/config"
	"alex-j-butler.com/tf2-booking/database"
	"alex-j-butler.com/tf2-booking/servers"
	"alex-j-butler.com/tf2-booking/util"
	"alex-j-butler.com/tf2-booking/wait"
	"github.com/bwmarrin/discordgo"
)

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
	Serv := servers.GetAvailableServer(servers.Servers)

	if Serv != nil {
		// Book the server.
		RCONPassword, ServerPassword, err := Serv.Book(m.Author, config.Conf.Booking.Duration.Duration)
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
		// Serv.ExtendBooking(config.Conf.BookingExtendDuration.Duration)
		Serv.ExtendBooking(config.Conf.Booking.ExtendDuration.Duration)

		// Notify server of successful operation.
		Serv.SendCommand(
			fmt.Sprintf(
				"say @%s: Your booking has been extended by %s.",
				m.Author.Username,
				util.ToHuman(&config.Conf.Booking.ExtendDuration.Duration),
			),
		)

		// Notify Discord channel of successful operation.
		Session.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf(
				"%s: Your booking has been extended by %s.",
				User.GetMention(),
				util.ToHuman(&config.Conf.Booking.ExtendDuration.Duration),
			),
		)
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

	servers := servers.GetBookedServers(servers.Servers)
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
		SaveState(".state.json", servers.Servers, Users, UserServers)
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

	SaveState(".state.json", servers.Servers, Users, UserServers)
	wait.Exit()
}

func Link(m *discordgo.MessageCreate, command string, args []string) {
	User := &util.PatchUser{m.Author}

	secret := util.RandStringBytes(24)
	database.DB.Create(&database.AuthSecret{Secret: secret, DiscordID: User.ID})
	Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Visit %s/auth/%s to link your Steam account.", User.GetMention(), config.Conf.SteamAuthServer.RootURL, secret))
}
