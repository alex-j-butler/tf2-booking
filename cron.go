package main

import (
	"fmt"
	"log"
	"time"

	"alex-j-butler.com/tf2-booking/util"

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

			log.Println(fmt.Sprintf("Querying server %s.", Serv.Name))

			go func() {
				log.Println(fmt.Sprintf("Querying server %s in goroutine.", Serv.Name))

				server, err := steam.Connect(Serv.Address)
				if err != nil {
					log.Println(fmt.Sprintf("Failed to connect to server \"%s\":", Serv.Name), err)

					HandleQueryError(&Serv, err)

					return
				}

				defer server.Close()

				info, err := server.Info()
				if err != nil {
					log.Println(fmt.Sprintf("Failed to query server \"%s\":", Serv.Name), err)

					HandleQueryError(&Serv, err)

					return
				}

				if info.Players < Conf.MinPlayers {
					Serv.AddIdleMinute()
					log.Println(fmt.Sprintf("Added idle minute for server %s", Serv.Name))
				} else {
					Serv.ResetIdleMinutes()
					log.Println(fmt.Sprintf("Reset idle minutes for server %s", Serv.Name))
				}

				log.Println(fmt.Sprintf("Current idle minutes for server %s: %d out of %d", Serv.Name, Serv.GetIdleMinutes(), Conf.MaxIdleMinutes))

				if Serv.GetIdleMinutes() >= Conf.MaxIdleMinutes {
					log.Println(fmt.Sprintf("Idle unbooked for server %s: %d out of %d", Serv.Name, Serv.GetIdleMinutes(), Conf.MaxIdleMinutes))

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
			}()
		}
	}
}

func CheckStats() {
	// Iterate through servers.
	for i := 0; i < len(Conf.Servers); i++ {
		if !Conf.Servers[i].IsAvailable() {
			go func(Serv *Server) {
				stats, err := Serv.SendRCONCommand("stats")

				if err != nil {
					log.Println("Stats query error:", err)
					return
				}

				// log.Println("Stats query: ", stats)
				s, err := util.ParseStats(stats)
				if err != nil {
					log.Println("Stats parse error:", err)
					return
				}

				// Calculate new average.
				if Serv.TickRate == 0.0 {
					Serv.TickRate = s.FPS
					Serv.TickRateMeasurements = 1
				} else {
					Serv.TickRate = ((Serv.TickRate*float32(Serv.TickRateMeasurements) + s.FPS) / float32(Serv.TickRateMeasurements+1))
					Serv.TickRateMeasurements++
				}
			}(&Conf.Servers[i])
		}
	}
}
