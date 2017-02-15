package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"alex-j-butler.com/tf2-booking/globals"
	"alex-j-butler.com/tf2-booking/servers"
	"alex-j-butler.com/tf2-booking/util"
	"github.com/bwmarrin/discordgo"
	uuid "github.com/satori/go.uuid"
)

func SynchroniseServers(message *discordgo.MessageCreate, input string, args []string) bool {
	User := &util.PatchUser{message.Author}

	for i, server := range servers.Servers {
		// Synchronise the server from Redis, to get information for existing servers.
		err := server.Synchronise(globals.RedisClient)
		if err != nil {
			// panic(err)
		}

		// Put the modified server back.
		servers.Servers[i] = server
	}

	Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Synchronised all servers.", User.GetMention()))

	// We've handled everything we need to.
	return true
}

func PrintStats(m *discordgo.MessageCreate, input string, args []string) bool {
	User := &util.PatchUser{m.Author}

	servs := servers.GetBookedServers(servers.Servers)
	message := "Server stats:"
	count := 0

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

			message = fmt.Sprintf("%s\n\t%s (Booked by %s): %f", message, server.Name, username, server.TickRate)
			count++
		}
	}

	message = fmt.Sprintf("%s\n\n%d out of %d servers booked", message, count, len(servers.Servers))

	if count == 0 {
		message = "No servers are currently booked."
	}

	// This command seems to be taking a long time, so for debugging, we'll see how long this SQL query takes
	// to run.
	dbqueryStartTime := time.Now()

	stmt, err := globals.DB.Prepare("SELECT server_name, sum(age(unbooked_time, booked_time)) FROM bookings WHERE booked_time > (current_date - $1::interval) GROUP BY server_name ORDER BY server_name ASC;")
	defer stmt.Close()
	if err != nil {
		log.Println("Prepare error:", err)
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: %s", User.GetMention(), "Something went wrong retrieving server history!"))
		return true
	}

	rows, err := stmt.Query("7 days")
	defer rows.Close()
	if err != nil {
		log.Println("Query error:", err)
		Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: %s", User.GetMention(), "Something went wrong retrieving server history!"))
		return true
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

	// We've handled everything we need to.
	return true
}

// AddLocalServer is the command handler for the '-b add local' command, that creates the
// a new local server and saves it to Redis.
func AddLocalServer(message *discordgo.MessageCreate, input string, args []string) bool {
	User := &util.PatchUser{message.Author}

	if len(args) > 3 {
		// Generate UUID for the server.
		serverUUID := uuid.NewV4()

		// Get server details from command.
		name := args[0]
		path := args[1]
		address := args[2]
		stvAddress := args[3]

		// Create server struct.
		server := servers.Server{
			UUID:       serverUUID.String(),
			Name:       name,
			Type:       "local",
			Path:       path,
			Address:    address,
			STVAddress: stvAddress,
		}

		// Serialise the server as JSON.
		serialised, err := json.Marshal(server)
		if err != nil {
			Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Failed to save server", User.GetMention()))
			return true
		}

		// Save the server in redis.
		err = globals.RedisClient.Set(fmt.Sprintf("server.%s", serverUUID.String()), serialised, 0).Err()
		if err != nil {
			Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Failed to save server", User.GetMention()))
			return true
		}

		// Send success message.
		Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Added local server: %s", User.GetMention(), serverUUID.String()))
	} else {
		// Print usage.
		Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Invalid command, usage: '-b add local <name> <path> <address> <stv address>'", User.GetMention()))
	}

	// We've handled everything we need to.
	return true
}

func AddRemoteServer(message *discordgo.MessageCreate, input string, args []string) bool {
	User := &util.PatchUser{message.Author}
	Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Added remote server", User.GetMention()))

	// We've handled everything we need to.
	return true
}

func ConfirmServer(message *discordgo.MessageCreate, input string, args []string) bool {
	User := &util.PatchUser{message.Author}
	Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Confirmed server creation", User.GetMention()))

	// We've handled everything we need to.
	return true
}

func DeleteServer(message *discordgo.MessageCreate, input string, args []string) bool {
	User := &util.PatchUser{message.Author}
	Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Deleted server", User.GetMention()))

	// We've handled everything we need to.
	return true
}

func ListServers(message *discordgo.MessageCreate, input string, args []string) bool {
	User := &util.PatchUser{message.Author}
	Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Servers:", User.GetMention()))

	// We've handled everything we need to.
	return true
}
