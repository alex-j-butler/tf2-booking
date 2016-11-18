package main

import (
	"fmt"
	"strings"
	"time"

	"alex-j-butler.com/tf2-booking/commands/ingame"
	"alex-j-butler.com/tf2-booking/config"
	"alex-j-butler.com/tf2-booking/util"
)

func ReportServer(commandInfo ingame.CommandInfo, command string, args []string) {
	if len(args) < 1 {
		// No reason given, send error message to server.
		commandInfo.Server.SendCommand(fmt.Sprintf("say Please give a reason: !report <reason>"))
		return
	}

	if timeout, ok := UserReportTimeouts[commandInfo.SteamID]; ok && time.Since(timeout).Nanoseconds() < 0 {
		// User can't report right now.
		commandInfo.Server.SendCommand(fmt.Sprintf("say You can't report that quickly! Try again in a few minutes."))
		return
	}

	reason := strings.Join(args, " ")

	// Convert the steam id provided by the log handler.
	steamID := util.FromSteamID3(commandInfo.SteamID)

	// Construct the message.
	message := fmt.Sprintf(
		"Server '%s' (%s) has been reported by '%s' (%s) with reason: '%s'",
		commandInfo.Server.Name,
		commandInfo.Server.SessionName,
		commandInfo.Username,
		steamID.GetCommunityURL(),
		reason,
	)

	// Send the message to the notification users.
	for _, notificationUser := range config.Conf.Discord.NotificationUsers {
		UserChannel, _ := Session.UserChannelCreate(notificationUser)
		Session.ChannelMessageSend(UserChannel.ID, message)
	}

	// Set the report timeout for this user.
	UserReportTimeouts[commandInfo.SteamID] = time.Now().Add(config.Conf.Commands.ReportDuration.Duration)

	// Reply to the command.
	commandInfo.Server.SendCommand(fmt.Sprintf("say Server reported! Thank you for your input."))
}

func TimeLeft(commandInfo ingame.CommandInfo, command string, args []string) {
	duration := -time.Since(commandInfo.Server.ReturnDate)
	commandInfo.Server.SendCommand(fmt.Sprintf("say %s remaining in booking.", util.ToHuman(&duration)))
}
