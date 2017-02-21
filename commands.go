package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"bytes"

	"alex-j-butler.com/tf2-booking/globals"
	"alex-j-butler.com/tf2-booking/servers"
	"alex-j-butler.com/tf2-booking/util"
	"alex-j-butler.com/tf2-booking/wait"
	"github.com/bwmarrin/discordgo"
	"github.com/google/go-github/github"
	"github.com/olekukonko/tablewriter"
	uuid "github.com/satori/go.uuid"
)

// SynchroniseServers is the command handler function for the command
// that synchronises all the locally cached servers from the Redis database.
func SynchroniseServers(message *discordgo.MessageCreate, input string, args []string) bool {
	User := &util.PatchUser{message.Author}

	// Properly reloads all available servers from the Redis store, and attempts to update the status text.
	ReloadServers()
	UpdateGameString()

	Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Synchronised all servers.", User.GetMention()))

	// We've handled everything we need to.
	return true
}

// PrintStats is the command handler function that prints the stats of all the currently booked servers
// as well as providing an overview of the number of hours each server has been booked in the last 7 days.
func PrintStats(m *discordgo.MessageCreate, input string, args []string) bool {
	User := &util.PatchUser{m.Author}

	servs := servers.GetBookedServers(servers.Servers)
	message := "Server stats:"
	count := 0

	for _, server := range servs {
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

	// Do not even ask what in the hell is going on here.
	// This is to allow spaces in the 'name' argument.
	firstArg := args[0]
	if len(args) > 5 {
		if strings.HasPrefix(firstArg, "\"") {
			firstArg = firstArg[1:]
			deleted := 0

			// Iterate over the arguments.
			for k, arg := range args {
				if k == 0 {
					continue
				}

				// Delete from the slice.
				j := k - deleted
				args = args[:j+copy(args[j:], args[j+1:])]
				deleted++

				if strings.HasSuffix(arg, "\"") {
					firstArg = fmt.Sprintf("%s %s", firstArg, arg[:len(arg)-1])
					break
				} else {
					firstArg = fmt.Sprintf("%s %s", firstArg, arg)
				}
			}
		}
	}

	if len(args) >= 5 {
		// Generate UUID for the server.
		serverUUID := uuid.NewV4()

		// Get server details from command.
		name := firstArg
		path := args[1]
		sessionName := args[2]
		address := args[3]
		stvAddress := args[4]

		// Create server struct.
		server := servers.Server{
			UUID:        serverUUID.String(),
			Name:        name,
			Type:        "local",
			Path:        path,
			SessionName: sessionName,
			Address:     address,
			STVAddress:  stvAddress,

			// Servers are inactive until they are confirmed using the '-b confirm <uuid|name>' command.
			Active: false,
		}

		// Add the server to our cached list of servers.
		servers.Servers[serverUUID.String()] = &server

		// Save the server to Redis.
		server.Update(globals.RedisClient)

		// Send success message.
		Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Added local server: %s", User.GetMention(), serverUUID.String()))
	} else {
		// Print usage.
		Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Invalid command, usage: '-b add local <name> <path> <session name> <address> <stv address>'", User.GetMention()))
	}

	// We've handled everything we need to.
	return true
}

func AddRemoteServer(message *discordgo.MessageCreate, input string, args []string) bool {
	User := &util.PatchUser{message.Author}
	Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Not yet implemented.", User.GetMention()))

	// We've handled everything we need to.
	return true
}

func ConfirmServer(message *discordgo.MessageCreate, input string, args []string) bool {
	User := &util.PatchUser{message.Author}

	if len(args) == 1 {
		// Is the argument passed in a UUID or a server name?
		isUUID := util.IsUUID4(args[0])

		var server *servers.Server
		var err error

		if isUUID {
			// Nice and simple - UUID's are unique to a server.
			server, err = servers.GetServerByUUID(servers.Servers, args[0])
		} else {
			// This should later be renamed to GetServersByName, and will return a slice of servers that match
			// since server names do not need to be unique.
			// If more than one server is found in this situation, we should reply and tell them that the name is ambiguous (and provide the appropriate commands for every possible server).
			server, err = servers.GetServerByName(servers.Servers, args[0])
		}

		if err != nil {
			// Yo, what. We can't find that server.
			Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Server was not found.", User.GetMention()))
			return true
		}

		// Activate that server & save it in Redis.
		server.Active = true
		err = server.Update(globals.RedisClient)
		if err != nil {
			// Redis error, oh no!
			Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Server could not be saved.", User.GetMention()))
			return true
		}

		// Send a message.
		Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Server '%s' was successfully activated!", User.GetMention(), server.Name))
	} else {
		// Print usage.
		Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Invalid command, usage: '-b confirm <uuid|name>'", User.GetMention()))
	}

	// We've handled everything we need to.
	return true
}

func DeleteServer(message *discordgo.MessageCreate, input string, args []string) bool {
	User := &util.PatchUser{message.Author}

	// Is the argument passed in a UUID or a server name?
	isUUID := util.IsUUID4(args[0])

	var server *servers.Server
	var err error

	if isUUID {
		// Nice and simple - UUID's are unique to a server.
		server, err = servers.GetServerByUUID(servers.Servers, args[0])
	} else {
		// This should later be renamed to GetServersByName, and will return a slice of servers that match
		// since server names do not need to be unique.
		// If more than one server is found in this situation, we should reply and tell them that the name is ambiguous (and provide the appropriate commands for every possible server).
		server, err = servers.GetServerByName(servers.Servers, args[0])
	}

	if err != nil {
		// Yo, what. We can't find that server.
		Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Server was not found.", User.GetMention()))
		return true
	}

	// Delete the server in Redis.
	err = globals.RedisClient.Del(fmt.Sprintf("server.%s", server.UUID)).Err()
	if err != nil {
		// Error deleting the server in Redis.
		Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Server could not be deleted.", User.GetMention()))
		return true
	}

	// Delete the server.
	delete(servers.Servers, server.UUID)

	// Send a message.
	Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Server '%s' was successfully deleted!", User.GetMention(), server.Name))

	// We've handled everything we need to.
	return true
}

func ListServers(message *discordgo.MessageCreate, input string, args []string) bool {
	User := &util.PatchUser{message.Author}

	serverSlice := servers.ServersToSlice(servers.Servers)
	sort.Sort(serverSlice)

	var data [][]string
	var buf bytes.Buffer

	for _, server := range serverSlice {
		data = append(data, []string{
			server.UUID,
			server.Name,
			server.Type,
			server.Address,
			server.STVAddress,
			fmt.Sprintf("%t", server.Active),
		})
	}

	table := tablewriter.NewWriter(&buf)
	table.SetHeader(
		[]string{
			"UUID",
			"Name",
			"Type",
			"Address",
			"STV Address",
			"Active",
		},
	)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(data)
	table.Render()

	Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Servers:\n```%s```", User.GetMention(), buf.String()))

	// We've handled everything we need to.
	return true
}

func Update(message *discordgo.MessageCreate, input string, args []string) bool {
	User := &util.PatchUser{message.Author}

	if len(args) == 1 {
		// Create a GitHub API client.
		client := github.NewClient(nil)
		// Tag name
		tagName := args[0]

		// Get release by tag.
		release, _, err := client.Repositories.GetReleaseByTag("alex-j-butler", "tf2-booking", tagName)
		if err != nil {
			// Send error message.
			Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Failed to retrieve release.", User.GetMention()))
			return true
		}

		asset, err := util.GetReleaseAsset(release.Assets, "tf2-booking-amd64")
		if err != nil {
			// Send error message.
			Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Failed to retrieve release asset.", User.GetMention()))
			return true
		}

		Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Starting update to release %s", User.GetMention(), *release.TagName))

		go func(asset github.ReleaseAsset, message *discordgo.MessageCreate) {
			// Update the executable.
			UpdateExecutable(*asset.BrowserDownloadURL)

			// Send the success notification.
			Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Updated `tf2-booking` & restarting now.", User.GetMention()))

			// Annnnnd, exit.
			wait.Exit()
		}(asset, message)
	}

	return true
}

func Exit(message *discordgo.MessageCreate, input string, args []string) bool {
	User := &util.PatchUser{message.Author}

	Session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s: Shutting down `tf2-booking`.", User.GetMention()))

	wait.Exit()

	return true
}
