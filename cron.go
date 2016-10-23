package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kidoman/go-steam"
)

// Check if any servers are ready to be unbooked by the automatic timeout after 4 hours.
func CheckUnbookServers() {
	// Iterate through servers.
	for _, Serv := range Conf.Servers {
		if Serv.IsAvailable() {
			return
		}

		if !Serv.IsAvailable() && !Serv.SentWarning && (Serv.ReturnDate.Add(Conf.BookingWarningDuration.Duration)).Before(time.Now()) {
			// Only allow this message to be sent once.
			Serv.SentWarning = true

			// Send warning message.
			Serv.SendCommand(fmt.Sprintf("say Your booking will expire in %s, type 'extend' into Discord to extend the booking.", Conf.BookingWarningDurationText))
		}

		if !Serv.IsAvailable() && Serv.ReturnDate.Before(time.Now()) {
			UserID := Serv.GetBooker()
			UserMention := Serv.GetBookerMention()

			// Remove the user's booked state.
			Users[UserID] = false
			UserServers[UserID] = nil

			// Unbook the server.
			Serv.Unbook()
			Serv.Stop()

			// Upload STV demos
			STVMessage, err := Serv.UploadSTV()

			// Send 'returned' message
			Session.ChannelMessageSend(Conf.DefaultChannel, fmt.Sprintf("%s: Your server was automatically unbooked.", UserMention))

			// Send 'stv' message, if it uploaded successfully.
			if err == nil {
				Session.ChannelMessageSend(Conf.DefaultChannel, fmt.Sprintf("%s: %s", UserMention, STVMessage))
			}

			UpdateGameString()

			log.Println(fmt.Sprintf("Automatically unbooked server \"%s\" from \"%s\", Reason: Booking timelimit reached", Serv.Name, UserID))
		}
	}
}

func CheckIdleMinutes() {
	// Iterate through servers.
	for _, Serv := range Conf.Servers {
		if !Serv.IsAvailable() {
			go func(Serv *Server) {
				server, err := steam.Connect(Serv.Address)
				if err != nil {
					log.Println(fmt.Sprintf("Failed to connect to server \"%s\":", Serv.Name), err)

					HandleQueryError(Serv, err)

					return
				}

				defer server.Close()

				info, err := server.Info()
				if err != nil {
					log.Println(fmt.Sprintf("Failed to query server \"%s\":", Serv.Name), err)

					HandleQueryError(Serv, err)

					return
				}

				if info.Players < Conf.MinPlayers {
					Serv.AddIdleMinute()
				} else {
					Serv.ResetIdleMinutes()
				}

				if Serv.GetIdleMinutes() >= Conf.MaxIdleMinutes {
					UserID := Serv.GetBooker()
					UserMention := Serv.GetBookerMention()

					// Remove the user's booked state.
					Users[UserID] = false
					UserServers[UserID] = nil

					// Unbook the server.
					Serv.Unbook()
					Serv.Stop()

					// Upload STV demos
					STVMessage, err := Serv.UploadSTV()

					// Send 'returned' message
					Session.ChannelMessageSend(Conf.DefaultChannel, fmt.Sprintf("%s: Your server was automatically unbooked.", UserMention))

					// Send 'stv' message, if it uploaded successfully.
					if err == nil {
						Session.ChannelMessageSend(Conf.DefaultChannel, fmt.Sprintf("%s: %s", UserMention, STVMessage))
					}

					UpdateGameString()

					log.Println(fmt.Sprintf("Automatically unbooked server \"%s\" from \"%s\", Reason: Idle timeout from too little players", Serv.Name, UserID))
				}
			}(&Serv)
		}
	}
}
