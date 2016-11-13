package main

import (
	"fmt"
	"log"
	"time"

	"alex-j-butler.com/tf2-booking/config"
	"alex-j-butler.com/tf2-booking/servers"
	"alex-j-butler.com/tf2-booking/util"

	"github.com/kidoman/go-steam"
)

// Check if any servers are ready to be unbooked by the automatic timeout after 4 hours.
func CheckUnbookServers() {
	// Iterate through servers.
	for i := 0; i < len(config.Conf.Servers); i++ {
		Serv := &config.Conf.Servers[i]

		if Serv.IsAvailable() {
			return
		}

		if !Serv.IsAvailable() && !Serv.SentWarning && (Serv.ReturnDate.Add(config.Conf.BookingWarningDuration.Duration)).Before(time.Now()) {
			// Only allow this message to be sent once.
			Serv.SentWarning = true

			// Send warning message.
			Serv.SendCommand(fmt.Sprintf("say Your booking will expire in %s, type 'extend' into Discord to extend the booking.", config.Conf.BookingWarningDurationText))
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
			booking, err := Serv.GetBooking()

			STVMessage := "STV Demo(s) uploaded:"
			for i := 0; i < len(booking.Demos); i++ {
				demo := booking.Demos[i]
				STVMessage = fmt.Sprintf("%s\n\t%s", STVMessage, demo.URL)
			}

			Session.ChannelMessageSend(config.Conf.DefaultChannel, fmt.Sprintf("%s: Players: %+v", UserMention, booking.Players))

			// Send 'returned' message
			Session.ChannelMessageSend(config.Conf.DefaultChannel, fmt.Sprintf("%s: Your server was automatically unbooked.", UserMention))

			// Send 'stv' message, if it uploaded successfully.
			if err == nil {
				Session.ChannelMessageSend(config.Conf.DefaultChannel, fmt.Sprintf("%s: %s", UserMention, STVMessage))
			}

			UpdateGameString()

			log.Println(fmt.Sprintf("Automatically unbooked server \"%s\" from \"%s\", Reason: Booking timelimit reached", Serv.Name, UserID))
		}
	}
}

func CheckPlayers() {
	// Iterate through servers.
	for i := 0; i < len(config.Conf.Servers); i++ {
		Serv := &config.Conf.Servers[i]

		if !Serv.IsAvailable() {
			go func(s *servers.Server) {
				statusLine, err := s.SendRCONCommand("status")
				if err != nil {
					log.Println(fmt.Sprintf("Failed to send RCON status to server \"%s\"", s.Name), err)
					return
				}

				status, _ := util.ParseStatus(statusLine)

				for _, user := range status.Users {
					if Serv.Players != nil {
						Serv.Players[user.ID.CommunityID] = true
					}
				}
			}(Serv)
		}
	}
}

func CheckIdleMinutes() {
	// Iterate through servers.
	for i := 0; i < len(config.Conf.Servers); i++ {
		Serv := &config.Conf.Servers[i]

		if !Serv.IsAvailable() {
			go func(s *servers.Server) {
				server, err := steam.Connect(s.Address)
				if err != nil {
					log.Println(fmt.Sprintf("Failed to connect to server \"%s\":", s.Name), err)

					HandleQueryError(s, err)

					return
				}

				defer server.Close()

				info, err := server.Info()
				if err != nil {
					log.Println(fmt.Sprintf("Failed to query server \"%s\":", s.Name), err)

					HandleQueryError(s, err)

					return
				}

				if info.Players < config.Conf.MinPlayers {
					s.AddIdleMinute()
				} else {
					s.ResetIdleMinutes()
				}

				if s.GetIdleMinutes() >= config.Conf.MaxIdleMinutes {
					UserID := s.GetBooker()
					UserMention := s.GetBookerMention()

					// Reset the idle minutes.
					s.ResetIdleMinutes()

					// Remove the user's booked state.
					Users[UserID] = false
					UserServers[UserID] = nil

					// Unbook the server.
					s.Unbook()
					s.Stop()

					// Upload STV demos
					booking, err := Serv.GetBooking()

					STVMessage := "STV Demo(s) uploaded:"
					for i := 0; i < len(booking.Demos); i++ {
						demo := booking.Demos[i]
						STVMessage = fmt.Sprintf("%s\n\t%s", STVMessage, demo.URL)
					}

					Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: Players: %+v", User.GetMention(), booking.Players))

					// Send 'returned' message
					Session.ChannelMessageSend(config.Conf.DefaultChannel, fmt.Sprintf("%s: Your server was automatically unbooked.", UserMention))

					// Send 'stv' message, if it uploaded successfully.
					if err == nil {
						Session.ChannelMessageSend(config.Conf.DefaultChannel, fmt.Sprintf("%s: %s", UserMention, STVMessage))
					}

					UpdateGameString()

					log.Println(fmt.Sprintf("Automatically unbooked server \"%s\" from \"%s\", Reason: Idle timeout from too little players", s.Name, UserID))
				}
			}(Serv)
		}
	}
}

func CheckStats() {
	// Iterate through servers.
	for i := 0; i < len(config.Conf.Servers); i++ {
		Serv := &config.Conf.Servers[i]

		if !Serv.IsAvailable() {
			go func(s *servers.Server) {
				stats, err := Serv.SendRCONCommand("stats")

				if err != nil {
					log.Println("Stats query error:", err)
					return
				}

				// log.Println("Stats query: ", stats)
				st, err := util.ParseStats(stats)
				if err != nil {
					log.Println("Stats parse error:", err)
					return
				}

				// Calculate new average.
				if s.TickRateMeasurements == 0 || s.TickRateMeasurements > 20 {
					s.TickRate = st.FPS
					s.TickRateMeasurements = 1
				} else {
					s.TickRate = ((s.TickRate*float32(s.TickRateMeasurements) + st.FPS) / float32(s.TickRateMeasurements+1))
					s.TickRateMeasurements++
				}
			}(Serv)
		}
	}
}
