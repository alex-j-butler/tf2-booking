package main

import (
	"fmt"
	"log"
	"time"
)

// Check if any servers are ready to be unbooked by the automatic timeout after 4 hours.
func CheckUnbookServers() {
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
