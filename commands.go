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

func sendServerDetails(channelID string, serv *servers.Server, serverPassword, rconPassword string) {
	Session.ChannelMessageSendEmbed(
		channelID,
		"**Here are the details for your booked server:**",
		&discordgo.MessageEmbed{
			Color: 12763842,
			Type:  "rich",
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   "Server Address",
					Value:  fmt.Sprintf("`%s`", serv.Address),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Server Password",
					Value:  fmt.Sprintf("`%s`", serverPassword),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "RCON Password",
					Value:  fmt.Sprintf("`%s`", rconPassword),
					Inline: true,
				},
			},
		},
	)
	Session.ChannelMessageSendEmbed(
		channelID,
		"",
		&discordgo.MessageEmbed{
			Color: 321378,
			Type:  "rich",
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   "Connect String",
					Value:  fmt.Sprintf("`connect %s; password %s; rcon_password %s`", serv.Address, serverPassword, rconPassword),
					Inline: false,
				},
				&discordgo.MessageEmbedField{
					Name:   "STV String",
					Value:  fmt.Sprintf("`connect %s`", serv.STVAddress),
					Inline: false,
				},
			},
		},
	)
	Session.ChannelMessageSendEmbed(
		channelID,
		"",
		&discordgo.MessageEmbed{
			Color: 12763842,
			Type:  "rich",
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Did you know you can report a server by typing !report into ingame chat?",
				IconURL: "https://tf2-au.qixalite.com/stv/help_icon.png",
			},
		},
	)
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
	Serv := servers.GetAvailableServer(servers.Servers)

	if Serv != nil {
		// Book the server.
		RCONPassword, ServerPassword, err := Serv.Book(m.Author, config.Conf.Booking.Duration.Duration)
		if err != nil {
			Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Something went wrong while trying to book your server, please try again later.", User.GetMention()))
			log.Println(fmt.Sprintf("Failed to book server \"%s\" from \"%s\":", Serv.Name, m.Author.ID), err)
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
			/*
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
			*/

			// Send the server details (and a small tip).
			sendServerDetails(UserChannel.ID, Serv, ServerPassword, RCONPassword)

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

func SendPassword(m *discordgo.MessageCreate, command string, args []string) {
	User := &util.PatchUser{m.Author}

	// Check if the user has already booked a server out.
	if value, ok := Users[m.Author.ID]; !ok || value == false {
		// Send a message to let the user know they do not have a server booked.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: You haven't booked a server. Type `book` to book a server.", User.GetMention()))

		return
	}

	if Serv, ok := UserServers[m.Author.ID]; ok && Serv != nil {
		serverPassword, err := Serv.GetCurrentPassword()
		if err != nil {
			Session.ChannelMessageSend(
				m.ChannelID,
				fmt.Sprintf(
					"%s: We failed to retrieve your server password.",
					User.GetMention(),
				),
			)

			return
		}

		// Send message to public channel, without server details.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Server password have been sent via private message.", User.GetMention()))

		// Send message to private DM, with server details.
		UserChannel, _ := Session.UserChannelCreate(m.Author.ID)
		Session.ChannelMessageSend(
			UserChannel.ID,
			fmt.Sprintf(
				"Here is your server details:\n\tServer address: %s\n\tPassword: %s\n\tConnect string: `connect %s; password %s`",
				Serv.Address,
				serverPassword,
				Serv.Address,
				serverPassword,
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
			bookerID := server.GetBooker()
			bookerUser, err := Session.User(bookerID)

			var username string
			if err != nil {
				username = "Unknown"
			} else {
				username = bookerUser.Username
			}

			message = fmt.Sprintf("%s\n\t%s (Booked by %s): %f", message, server.Name, username, server.TickRate)
			count++
		}
	}

	message = fmt.Sprintf("%s\n\n%d out of %d servers booked", message, count, len(servers))

	if count == 0 {
		message = "No servers are currently booked."
	}

	stmt, err := database.DB.Prepare("SELECT server_name, sum(age(unbooked_time, booked_time)) FROM bookings WHERE booked_time > (current_date - $1::interval) GROUP BY server_name ORDER BY server_name ASC;")
	defer stmt.Close()
	if err != nil {
		log.Println("Prepare error:", err)
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: %s", User.GetMention(), "Something went wrong retrieving server history!"))
		return
	}

	rows, err := stmt.Query("7 days")
	defer rows.Close()
	if err != nil {
		log.Println("Query error:", err)
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: %s", User.GetMention(), "Something went wrong retrieving server history!"))
		return
	}

	message = fmt.Sprintf("%s\n\n%s", message, "7 day history:")

	var serverName string
	var duration string
	for rows.Next() {
		rows.Scan(&serverName, &duration)
		message = fmt.Sprintf("%s\n\t%s: %s", message, serverName, duration)
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
