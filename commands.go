package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"alex-j-butler.com/tablewriter"
	"alex-j-butler.com/tf2-booking/commands"
	"alex-j-butler.com/tf2-booking/config"
	"alex-j-butler.com/tf2-booking/globals"
	"alex-j-butler.com/tf2-booking/servers"
	"alex-j-butler.com/tf2-booking/util"
	"github.com/bwmarrin/discordgo"
)

func sendServerDetails(channelID string, serv *servers.Server, serverPassword, rconPassword string) {
	Session.ChannelMessageSendEmbedWithContent(
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
	Session.ChannelMessageSendEmbedWithContent(
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
	Session.ChannelMessageSendEmbedWithContent(
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

func Version(m *discordgo.MessageCreate, command string, permissions commands.CommandPermissions, args commands.CommandArgList) {
	User := &util.PatchUser{m.Author}
	Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: `tf2-booking` running git revision `%s`", User.GetMention(), version))
}

func SyncServers(m *discordgo.MessageCreate, command string, permissions commands.CommandPermissions, args commands.CommandArgList) {
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

func Help(m *discordgo.MessageCreate, command string, permissions commands.CommandPermissions, args commands.CommandArgList) {
	User := &util.PatchUser{m.Author}

	helpMessage := `book            - Book a new server
unbook          - Unbook your current server
send password   - Send the updated server details
demos           - Send the link to the uploaded demos
help            - Display the help message (you're reading it!)

Admin commands:
stats    - Shows run status for all servers.
exit     - Exits the booking bot.

For help, ping @Alex_#7324 in this channel.

Note: Ozfortress booking commands also are accepted by this bot.`

	Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: **Qixalite Bookable Help**:```%s```", User.GetMention(), helpMessage))
}

// DemoLink command handler.
// Called when the user types the 'demo' command into the Discord channel.
// This function should send them the link to the Qixalite demo store.
func DemoLink(m *discordgo.MessageCreate, command string, permissions commands.CommandPermissions, args commands.CommandArgList) {
	User := &util.PatchUser{m.Author}
	Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Sorry, demos command is currently unavailable.", User.GetMention()))

	return

	// Fix the god damn demos-host.
	/*
		demosTarget := strings.Join(args.ToSlice()[1:], " ")
		if demosTarget == "" {
			demosTarget = User.GetFullname()
		}

		Session.ChannelMessageSend(m.ChannelID,
			// Sort of dangerous - format string from a config file, but the config file should only be accessible by admins.
			fmt.Sprintf(
				"%s: "+config.Conf.Commands.DemoLink,
				User.GetMention(),
				url.QueryEscape(demosTarget),
			),
		)
	*/
}

func Console(m *discordgo.MessageCreate, command string, permissions commands.CommandPermissions, args commands.CommandArgList) {
	User := &util.PatchUser{m.Author}

	serv, err := pool.GetServerByName(strings.Join(args.ToSlice()[1:], " "))
	serv.Synchronise(globals.RedisClient)

	if err != nil {
		Session.ChannelMessageSend(m.ChannelID,
			fmt.Sprintf(
				"%s: Server not found\n%s",
				User.GetMention(),
				strings.Join(args.ToSlice()[1:], " "),
			),
		)

		return
	}

	logLines, err := serv.ConsoleMax(10)
	if err != nil {
		Session.ChannelMessageSend(m.ChannelID,
			fmt.Sprintf(
				"%s: Console view for \"%s\" failed: %s",
				User.GetMention(),
				serv.Name,
				err,
			),
		)
		return
	}

	if len(logLines) > 0 {
		Session.ChannelMessageSend(m.ChannelID,
			fmt.Sprintf(
				"```%s```",
				strings.Join(logLines, "\n"),
			),
		)
		return
	}

	Session.ChannelMessageSend(m.ChannelID,
		fmt.Sprintf(
			"%s: No log lines for \"%s\"",
			User.GetMention(),
			serv.Name,
		),
	)
}

// BookServer command handler
// Called when a user types the 'book' command into the Discord channel.
// This function checks whether the user has a server booked, if not,
// it books a new server, preventing it from being used by another user,
// sets up the RCON password & Server Password and finally starts the TF2 server.
func BookServer(m *discordgo.MessageCreate, command string, permissions commands.CommandPermissions, args commands.CommandArgList) {
	User := &util.PatchUser{m.Author}
	target := m.Author

	// handle admins unbooking servers for other users.
	// `unbook @Alex` -> Unbooks a server from @Alex
	if permissions.Test(discordgo.PermissionManageServer) && args.Num() > 1 {
		// No point trying to decode if there are no valid mentions.
		if len(m.Mentions) > 0 {
			discordMention, _ := args.GetArg(1)
			var dirtyDiscordID string
			n, err := fmt.Sscanf(discordMention, "<@%s>", &dirtyDiscordID)
			discordID := dirtyDiscordID[:len(dirtyDiscordID)-1]
			if n == 1 || err == nil { // Ensure no error occurred, and that we only scanned 1 ID
				user, err := Session.User(discordID)
				if err == nil {
					Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Using Discord user `%s`", User.GetMention(), user.Username))
					target = user
				}
			}
		}
	}

	bookingInfo, err := GetDefaultValue.Run(globals.RedisClient, []string{fmt.Sprintf("user.%s", target.ID)}, nil).Result()
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
		RCONPassword, ServerPassword, err := Serv.Book(target)
		if err != nil {
			Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Something went wrong while trying to book your server, please try again later.", User.GetMention()))
			log.Println(fmt.Sprintf("Failed to book server \"%s\" from \"%s\":", Serv.Name, target.ID), err)
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
					if err := globals.RedisClient.Set(fmt.Sprintf("user.%s", target.ID), "", 0).Err(); err != nil {
						log.Println("Redis error:", err)
						log.Println("Failed to set user information for user:", target.ID)
						return
					}

					UpdateGameString()

					log.Println(fmt.Sprintf("Failed to start server \"%s\" from \"%s\"", Serv.Name, target.ID))
				}
			}(Serv, m)

			// Send message to public channel, without server details.
			Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Server details have been sent via private message.", User.GetMention()))

			// Create the private DM channel, and then send the server details (and a small tip).
			UserChannel, _ := Session.UserChannelCreate(target.ID)
			sendServerDetails(UserChannel.ID, Serv, ServerPassword, RCONPassword)

			// Add the user's booked state.
			if err := globals.RedisClient.Set(fmt.Sprintf("user.%s", target.ID), Serv.UUID, 0).Err(); err != nil {
				log.Println("Redis error:", err)
				log.Println("Failed to set user information for user:", target.ID)
				return
			}

			UpdateGameString()

			log.Println(fmt.Sprintf("Booked server \"%s\" from \"%s\"", Serv.Name, target.ID))
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
func UnbookServer(m *discordgo.MessageCreate, command string, permissions commands.CommandPermissions, args commands.CommandArgList) {
	User := &util.PatchUser{m.Author}
	target := m.Author

	// handle admins unbooking servers for other users.
	// `unbook @Alex` -> Unbooks a server from @Alex
	if permissions.Test(discordgo.PermissionManageServer) && args.Num() > 1 {
		// No point trying to decode if there are no valid mentions.
		if len(m.Mentions) > 0 {
			discordMention, _ := args.GetArg(1)
			var dirtyDiscordID string
			n, err := fmt.Sscanf(discordMention, "<@%s>", &dirtyDiscordID)
			discordID := dirtyDiscordID[:len(dirtyDiscordID)-1]
			if n == 1 || err == nil { // Ensure no error occurred, and that we only scanned 1 ID
				user, err := Session.User(discordID)
				if err == nil {
					Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Using Discord user `%s`", User.GetMention(), user.Username))
					target = user
				}
			}
		}
	}

	bookingInfo, err := GetDefaultValue.Run(globals.RedisClient, []string{fmt.Sprintf("user.%s", target.ID)}, nil).Result()
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

	Serv, err := pool.GetServerByUUID(bookingInfoStr)
	Serv.Synchronise(globals.RedisClient)

	if err == nil && Serv != nil {
		// Get the booker name, before the server gets unbooked
		bookerName := Serv.BookerFullname

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
		if err := globals.RedisClient.Set(fmt.Sprintf("user.%s", target.ID), "", 0).Err(); err != nil {
			log.Println("Redis error:", err)
			log.Println("Failed to set user information for user:", target.ID)
			return
		}

		// Unbook the server.
		err = Serv.Unbook()
		if err != nil {
			Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Unbook failed. Please contact an admin for assistance. (Failed to unreserve the server)", User.GetMention()))
			return
		}

		// Upload STV demos
		STVMessage, err := Serv.UploadSTV(bookerName)

		// Send 'returned' message.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Server unbooked.", User.GetMention()))

		// Send 'stv' message, if it uploaded successfully.
		if err == nil {
			Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: %s", User.GetMention(), STVMessage))
		}

		UpdateGameString()

		log.Println(fmt.Sprintf("Unbooked server \"%s\" from \"%s\"", Serv.Name, target.ID))
	} else {
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: You haven't booked a server. Type `book` to book a server.", User.GetMention()))

		// We're in an invalid state, reset back to normal.
		if err := globals.RedisClient.Set(fmt.Sprintf("user.%s", target.ID), "", 0).Err(); err != nil {
			log.Println("Redis error:", err)
			log.Println("Failed to set user information for user:", target.ID)
			return
		}

		return
	}
}

// ExtendServer command handler
// Called when a user types the 'extend' command into the Discord channel.
// This function checks whether the user has a server booked out, if so,
// it will extend the booking by adding time onto the servers return time.
func ExtendServer(m *discordgo.MessageCreate, command string, permissions commands.CommandPermissions, args commands.CommandArgList) {
	User := &util.PatchUser{m.Author}
	target := m.Author

	// handle admins unbooking servers for other users.
	// `unbook @Alex` -> Unbooks a server from @Alex
	if permissions.Test(discordgo.PermissionManageServer) && args.Num() > 1 {
		// No point trying to decode if there are no valid mentions.
		if len(m.Mentions) > 0 {
			discordMention, _ := args.GetArg(1)
			var dirtyDiscordID string
			n, err := fmt.Sscanf(discordMention, "<@%s>", &dirtyDiscordID)
			discordID := dirtyDiscordID[:len(dirtyDiscordID)-1]
			if n == 1 || err == nil { // Ensure no error occurred, and that we only scanned 1 ID
				user, err := Session.User(discordID)
				if err == nil {
					Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Using Discord user `%s`", User.GetMention(), user.Username))
					target = user
				}
			}
		}
	}

	bookingInfo, err := GetDefaultValue.Run(globals.RedisClient, []string{fmt.Sprintf("user.%s", target.ID)}, nil).Result()
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

	Serv, err := pool.GetServerByUUID(bookingInfoStr)
	Serv.Synchronise(globals.RedisClient)

	if err == nil && Serv != nil {
		Serv.ExtendBooking()

		// Notify server of successful operation.
		Serv.SendCommand(
			fmt.Sprintf(
				"say @%s: Your booking has been extended by %s.",
				target.Username,
				fmt.Sprintf("%d minutes", config.Conf.Booking.MaxIdleMinutes),
			),
		)

		// Notify Discord channel of successful operation.
		Session.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf(
				"%s: Your booking has been extended by %s.",
				User.GetMention(),
				fmt.Sprintf("%d minutes", config.Conf.Booking.MaxIdleMinutes),
			),
		)
	} else {
		// Notify Discord channel of failed operation.
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: You haven't booked a server. Type `book` to book a server.", User.GetMention()))

		// If the program execution reaches here, the state of the users & user-servers map
		// is invalid and should be reset to the 'unbooked' state.
		if err := globals.RedisClient.Set(fmt.Sprintf("user.%s", target.ID), "", 0).Err(); err != nil {
			log.Println("Redis error:", err)
			log.Println("Failed to set user information for user:", target.ID)
			return
		}

		return
	}
}

func SendPassword(m *discordgo.MessageCreate, command string, permissions commands.CommandPermissions, args commands.CommandArgList) {
	User := &util.PatchUser{m.Author}
	target := m.Author

	// handle admins unbooking servers for other users.
	// `unbook @Alex` -> Unbooks a server from @Alex
	if permissions.Test(discordgo.PermissionManageServer) && args.Num() > 1 {
		// No point trying to decode if there are no valid mentions.
		if len(m.Mentions) > 0 {
			discordMention, _ := args.GetArg(1)
			var dirtyDiscordID string
			n, err := fmt.Sscanf(discordMention, "<@%s>", &dirtyDiscordID)
			discordID := dirtyDiscordID[:len(dirtyDiscordID)-1]
			if n == 1 || err == nil { // Ensure no error occurred, and that we only scanned 1 ID
				user, err := Session.User(discordID)
				if err == nil {
					Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Using Discord user `%s`", User.GetMention(), user.Username))
					target = user
				}
			}
		}
	}

	bookingInfo, err := GetDefaultValue.Run(globals.RedisClient, []string{fmt.Sprintf("user.%s", target.ID)}, nil).Result()
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

	Serv, err := pool.GetServerByUUID(bookingInfoStr)
	Serv.Synchronise(globals.RedisClient)

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
		UserChannel, _ := Session.UserChannelCreate(target.ID)
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
		if err := globals.RedisClient.Set(fmt.Sprintf("user.%s", target.ID), "", 0).Err(); err != nil {
			log.Println("Redis error:", err)
			log.Println("Failed to set user information for user:", target.ID)
			return
		}

		return
	}
}

func getServerStatusString(server *servers.Server) string {
	if server.Reserved && server.Running {
		return "Booked"
	}
	if server.Reserved && !server.Running {
		return "Errored"
	}
	return "Stopped"
}

func PrintStats(m *discordgo.MessageCreate, command string, permissions commands.CommandPermissions, args commands.CommandArgList) {
	User := &util.PatchUser{m.Author}

	now := time.Now()
	to := now.Add(-4 * time.Hour)

	rpcLatencyGraph, err := rpcLatency(now, to)
	if err != nil {
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Graph error: %s", User.GetMention(), err))
	} else {
		Session.ChannelFileSend(m.ChannelID, "rpc_latency.png", &rpcLatencyGraph)
	}

	rpcRequestsGraph, err := rpcRequests(now, to)
	if err != nil {
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Graph error: %s", User.GetMention(), err))
	} else {
		Session.ChannelFileSend(m.ChannelID, "rpc_requests.png", &rpcRequestsGraph)
	}

	servs := pool.GetServers()
	// Sync the servers from redis.
	for i, server := range servs {
		// Synchronise the server from Redis, to get information for existing servers.
		err := server.Synchronise(globals.RedisClient)
		if err != nil {
			log.Println("Server", server.Name, "failed to Redis sync")
			continue
		}

		// Put the modified server back.
		servs[i] = server
	}

	message := "Server stats:"

	data := make([][]string, 0, len(servs))
	for _, serv := range servs {
		log.Println(fmt.Sprintf("Server: %s, Reserved: %t, Running: %t", serv.Name, serv.Reserved, serv.Running))

		// Retrieve the name of the Discord user who booked the server. Uses empty string if no one has booked the server.
		bookerID := serv.Booker
		bookerUser, err := Session.User(bookerID)

		var username string
		if err != nil {
			username = "Unknown"
		} else {
			username = bookerUser.Username
		}

		bookedTime := serv.BookedDate.String()
		if serv.BookedDate.Equal(time.Time{}) {
			bookedTime = ""
		}

		data = append(data, []string{serv.Name, getServerStatusString(serv), bookedTime, username})
	}

	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"Server name", "Status", "Book time", "Booker name"})
	table.SetHeaderLine(false)
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetAutoFormatHeaders(false)
	table.SetRowSeparator("|")
	table.AppendBulk(data)
	table.Render()

	message = fmt.Sprintf("%s\n```%s```", message, buf.String())
	Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: %s", User.GetMention(), message))
}

func Graph(m *discordgo.MessageCreate, command string, permissions commands.CommandPermissions, args commands.CommandArgList) {
	User := &util.PatchUser{m.Author}
	numArgs := args.Num()

	if numArgs == 3 {
		now := time.Now()
		to := now.Add(-24 * time.Hour)

		graphClient := NewGraphClient("dd5930a2b34093f052aea1eeb290f11b", "a138ebaed3072d4c04e063b8ee66f686980aa794")
		graph, err := graphClient.Graph(args.ToSlice()[1], args.ToSlice()[2], now, to)
		if err != nil {
			Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Graph err, %s", User.GetMention(), err))
			return
		}

		Session.ChannelFileSend(m.ChannelID, "graph.png", &graph)
	} else {
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: %s", User.GetMention(), "graph <title> <dd query>"))
	}
}
