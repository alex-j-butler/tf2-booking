package main

import (
	"fmt"
	"log"
	"time"

	"bytes"

	"alex-j-butler.com/tf2-booking/config"
	"alex-j-butler.com/tf2-booking/globals"
	"alex-j-butler.com/tf2-booking/servers"
	"alex-j-butler.com/tf2-booking/util"
	"alex-j-butler.com/tf2-booking/wait"
	"github.com/bwmarrin/discordgo"
	"github.com/google/go-github/github"
	"github.com/olekukonko/tablewriter"
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
				Name: ChooseRandomTip(),
			},
		},
	)
}

func Version(m *discordgo.MessageCreate, command string, args []string) {
	User := &util.PatchUser{m.Author}
	Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: `tf2-booking` running git revision `%s`", User.GetMention(), version))
}

func SyncServers(m *discordgo.MessageCreate, command string, args []string) {
	User := &util.PatchUser{m.Author}

	for i, server := range servers.Servers {
		// Synchronise the server from Redis, to get information for existing servers.
		err := server.Synchronise(globals.RedisClient)
		if err != nil {
			panic(err)
		}

		// Put the modified server back.
		servers.Servers[i] = server
	}

	Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Synchronised all servers.", User.GetMention()))
}

func Help(m *discordgo.MessageCreate, command string, args []string) {
	User := &util.PatchUser{m.Author}

	helpMessage := `book            - Book a new server
unbook          - Unbook your current server
send password   - Send the updated server details
demos           - Send the link to the uploaded demos (soon!)
help            - Display the help message (you're reading it!)

For help, ping @Alex#4240 in this channel.

Note: Ozfortress booking commands also are accepted by this bot.`

	Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: **Qixalite Bookable Help**:```%s```", User.GetMention(), helpMessage))
}

// DemoLink command handler.
// Called when the user types the 'demo' command into the Discord channel.
// This function should send them the link to the Qixalite demo store.
func DemoLink(m *discordgo.MessageCreate, command string, args []string) {
	User := &util.PatchUser{m.Author}

	Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: https://stv.qixalite.com/", User.GetMention()))
}

// BookServer command handler
// Called when a user types the 'book' command into the Discord channel.
// This function checks whether the user has a server booked, if not,
// it books a new server, preventing it from being used by another user,
// sets up the RCON password & Server Password and finally starts the TF2 server.
func BookServer(m *discordgo.MessageCreate, command string, args []string) {
	User := &util.PatchUser{m.Author}

	bookingInfo, err := GetDefaultValue.Run(globals.RedisClient, []string{fmt.Sprintf("user.%s", m.Author.ID)}, nil).Result()
	if err != nil {
		// Send a message to let the user know an error occurred.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Oops, looked like an error has occurred. Please contact an admin for assistance.", User.GetMention()))
		return
	}

	if len(bookingInfo.(string)) != 0 {
		// Send a message to let the user know they've already booked a server.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: You've already booked a server. Type `unbook` to return the server.", User.GetMention()))
		return
	}

	// Get the next available server.
	Serv := pool.GetAvailableServer()

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
					log.Println("Failed to book server:", err)

					// Reset the user's booked state.
					if err := globals.RedisClient.Set(fmt.Sprintf("user.%s", m.Author.ID), "", 0).Err(); err != nil {
						log.Println("Redis error:", err)
						log.Println("Failed to set user information for user:", m.Author.ID)
						return
					}

					UpdateGameString()

					log.Println(fmt.Sprintf("Failed to start server \"%s\" from \"%s\"", Serv.Name, m.Author.ID))
				}
			}(Serv, m)

			// Send message to public channel, without server details.
			Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Server details have been sent via private message.", User.GetMention()))

			// Create the private DM channel, and then send the server details (and a small tip).
			UserChannel, _ := Session.UserChannelCreate(m.Author.ID)
			sendServerDetails(UserChannel.ID, Serv, ServerPassword, RCONPassword)

			// Add the user's booked state.
			if err := globals.RedisClient.Set(fmt.Sprintf("user.%s", m.Author.ID), Serv.SessionName, 0).Err(); err != nil {
				log.Println("Redis error:", err)
				log.Println("Failed to set user information for user:", m.Author.ID)
				return
			}

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

	bookingInfo, err := GetDefaultValue.Run(globals.RedisClient, []string{fmt.Sprintf("user.%s", m.Author.ID)}, nil).Result()
	if err != nil {
		// Send a message to let the user know an error occurred.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Oops, looked like an error has occurred. Please contact an admin for assistance.", User.GetMention()))
		return
	}
	bookingInfoStr := bookingInfo.(string)

	if len(bookingInfoStr) == 0 {
		// Send a message to let the user know they do not have a server booked.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: You haven't booked a server. Type `book` to book a server.", User.GetMention()))
		return
	}

	Serv, err := pool.GetServerBySessionName(bookingInfoStr)

	if err == nil && Serv != nil {
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
		if err := globals.RedisClient.Set(fmt.Sprintf("user.%s", m.Author.ID), "", 0).Err(); err != nil {
			log.Println("Redis error:", err)
			log.Println("Failed to set user information for user:", m.Author.ID)
			return
		}

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
		if err := globals.RedisClient.Set(fmt.Sprintf("user.%s", m.Author.ID), "", 0).Err(); err != nil {
			log.Println("Redis error:", err)
			log.Println("Failed to set user information for user:", m.Author.ID)
			return
		}

		return
	}
}

// ExtendServer command handler
// Called when a user types the 'extend' command into the Discord channel.
// This function checks whether the user has a server booked out, if so,
// it will extend the booking by adding time onto the servers return time.
func ExtendServer(m *discordgo.MessageCreate, command string, args []string) {
	User := &util.PatchUser{m.Author}

	bookingInfo, err := GetDefaultValue.Run(globals.RedisClient, []string{fmt.Sprintf("user.%s", m.Author.ID)}, nil).Result()
	if err != nil {
		// Send a message to let the user know an error occurred.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Oops, looked like an error has occurred. Please contact an admin for assistance.", User.GetMention()))
		return
	}
	bookingInfoStr := bookingInfo.(string)

	if len(bookingInfoStr) == 0 {
		// Send a message to let the user know they do not have a server booked.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: You haven't booked a server. Type `book` to book a server.", User.GetMention()))
		return
	}

	Serv, err := pool.GetServerBySessionName(bookingInfoStr)

	if err == nil && Serv != nil {
		// Extend the booking.
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
		if err := globals.RedisClient.Set(fmt.Sprintf("user.%s", m.Author.ID), "", 0).Err(); err != nil {
			log.Println("Redis error:", err)
			log.Println("Failed to set user information for user:", m.Author.ID)
			return
		}

		return
	}
}

func SendPassword(m *discordgo.MessageCreate, command string, args []string) {
	User := &util.PatchUser{m.Author}

	bookingInfo, err := GetDefaultValue.Run(globals.RedisClient, []string{fmt.Sprintf("user.%s", m.Author.ID)}, nil).Result()
	if err != nil {
		// Send a message to let the user know an error occurred.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Oops, looked like an error has occurred. Please contact an admin for assistance.", User.GetMention()))
		return
	}
	bookingInfoStr := bookingInfo.(string)

	if len(bookingInfoStr) == 0 {
		// Send a message to let the user know they do not have a server booked.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: You haven't booked a server. Type `book` to book a server.", User.GetMention()))
		return
	}

	Serv, err := pool.GetServerBySessionName(bookingInfoStr)

	if err == nil && Serv != nil {
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
		if err := globals.RedisClient.Set(fmt.Sprintf("user.%s", m.Author.ID), "", 0).Err(); err != nil {
			log.Println("Redis error:", err)
			log.Println("Failed to set user information for user:", m.Author.ID)
			return
		}

		return
	}
}

func getServerStatusString(server *servers.Server) string {
	if server.IsBooked() {
		return "Running"
	}
	return "Stopped"
}

func PrintStats(m *discordgo.MessageCreate, command string, args []string) {
	User := &util.PatchUser{m.Author}

	servs := pool.GetServers()
	// Sync the servers from redis.
	for i, server := range servs {
		// Synchronise the server from Redis, to get information for existing servers.
		err := server.Synchronise(globals.RedisClient)
		if err != nil {
			panic(err)
		}

		// Put the modified server back.
		servs[i] = server
	}

	message := "Server stats:"
	count := 0

	data := make([][]string, 0, len(servs))
	for _, serv := range servs {
		// Retrieve the name of the Discord user who booked the server. Uses empty string if no one has booked the server.
		bookerID := serv.GetBooker()
		bookerUser, err := Session.User(bookerID)

		var username string
		if err != nil {
			username = ""
		} else {
			username = bookerUser.Username
		}

		data = append(data, []string{serv.Name, getServerStatusString(serv), serv.ReturnDate.String(), username, serv.Booker})
	}

	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"Server name", "Status", "Unbook time", "Booker name", "Booker ID"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoFormatHeaders(false)
	table.AppendBulk(data)
	table.Render()

	/*
		for i := 0; i < len(servs); i++ {
			server := servs[i]
			if server != nil {
				bookerID := server.GetBooker()
				bookerUser, err := Session.User(bookerID)

				var username string
				if err != nil {
					username = "Unknown"
				} else {
					username = bookerUser.Username
				}

				message = fmt.Sprintf("%s\n\t%s (Booked by %s)", message, server.Name, username)
				count++
			}
		}
	*/

	message = fmt.Sprintf("%s\n```%s```", message, buf.String())
	message = fmt.Sprintf("%s\n\n%d out of %d servers booked", message, count, len(servers.Servers))

	// This command seems to be taking a long time, so for debugging, we'll see how long this SQL query takes to run.
	dbqueryStartTime := time.Now()

	stmt, err := globals.DB.Prepare("SELECT server_name, sum(age(unbooked_time, booked_time)) FROM bookings WHERE booked_time > (current_date - $1::interval) GROUP BY server_name ORDER BY server_name ASC;")
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

	// Get the db query time.
	dbqueryTimeElapsed := time.Since(dbqueryStartTime)

	message = fmt.Sprintf("%s\n\n%s", message, "7 day history:")

	var serverName string
	var duration string
	for rows.Next() {
		rows.Scan(&serverName, &duration)
		message = fmt.Sprintf("%s\n\t%s: %s", message, serverName, duration)
	}

	Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: %s\nQuery took %s", User.GetMention(), message, dbqueryTimeElapsed))
}

func Update(m *discordgo.MessageCreate, command string, args []string) {
	User := &util.PatchUser{m.Author}

	if len(args) <= 0 {
		// Send usage.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Usage: `update <release tag>`", User.GetMention()))
		return
	}

	// Create a GitHub API client.
	client := github.NewClient(nil)
	// Tag name
	tagName := args[0]

	// Get release by tag.
	release, _, err := client.Repositories.GetReleaseByTag("alex-j-butler", "tf2-booking", tagName)
	if err != nil {
		// Send error message.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Failed to retrieve release.", User.GetMention()))
		return
	}

	asset, err := util.GetReleaseAsset(release.Assets, "tf2-booking-amd64")
	if err != nil {
		// Send error message.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Failed to retrieve release asset.", User.GetMention()))
		return
	}

	//
	Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Starting update to release %s", User.GetMention(), *release.TagName))

	go func(asset github.ReleaseAsset) {
		// Update the executable.
		UpdateExecutable(*asset.BrowserDownloadURL)

		// Send the success notification.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Updated `tf2-booking` & restarting now.", User.GetMention()))

		// Annnnnd, exit.
		wait.Exit()
	}(asset)
}

func Exit(m *discordgo.MessageCreate, command string, args []string) {
	User := &util.PatchUser{m.Author}

	Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Shutting down `tf2-booking`.", User.GetMention()))

	wait.Exit()
}
