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

func RunInitialCommands() {
	// Iterate through servers.
	for i := 0; i < len(Conf.Servers); i++ {
		Serv := Conf.Servers[i]

		if Serv.IsAvailable() {
			return
		}

		// Only run if the initial commands have never been run before.
		if !Serv.InitialCommandsRun {
			commands := []string{
				fmt.Sprintf("rcon_password %s", Serv.RCONPassword),
				fmt.Sprintf("sv_password %s", Serv.Password),
				fmt.Sprintf("logaddress_add %s:%d", config.Conf.LogAddressRemote, config.Conf.LogPort),
			}

			failed := false
			for _, command := range commands {
				_, err := Serv.SendRCONCommand(command)

				if err != nil {
					failed = true
				}
			}

			if !failed {
				// Set
				Serv.InitialCommandsRun = true
			}
		}
	}
}

// Check if any servers are ready to be unbooked by the automatic timeout after 4 hours.
func CheckUnbookServers() {
	// Iterate through servers.
	for i := 0; i < len(servers.Servers); i++ {
		Serv := &servers.Servers[i]

		if Serv.IsAvailable() {
			return
		}

		if !Serv.IsAvailable() && !Serv.SentWarning && (Serv.ReturnDate.Add(config.Conf.Booking.WarningDuration.Duration)).Before(time.Now()) {
			// Only allow this message to be sent once.
			Serv.SentWarning = true

			// Send warning message.
			Serv.SendCommand(
				fmt.Sprintf(
					"say Your booking will expire in %s, type 'extend' into Discord to extend the booking.",
					util.ToHuman(&config.Conf.Booking.WarningDuration.Duration),
				),
			)
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
			Session.ChannelMessageSend(config.Conf.Discord.DefaultChannel, fmt.Sprintf("%s: Your server was automatically unbooked.", UserMention))

			// Send 'stv' message, if it uploaded successfully.
			if err == nil {
				Session.ChannelMessageSend(config.Conf.Discord.DefaultChannel, fmt.Sprintf("%s: %s", UserMention, STVMessage))
			}

			UpdateGameString()

			log.Println(fmt.Sprintf("Automatically unbooked server \"%s\" from \"%s\", Reason: Booking timelimit reached", Serv.Name, UserID))
		}
	}
}

func CheckIdleMinutes() {
	// Iterate through servers.
	for i := 0; i < len(servers.Servers); i++ {
		Serv := &servers.Servers[i]

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

				if info.Players < config.Conf.Booking.MinPlayers {
					s.AddIdleMinute()
				} else {
					s.ResetIdleMinutes()
				}

				if s.GetIdleMinutes() >= config.Conf.Booking.MaxIdleMinutes {
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
					STVMessage, err := s.UploadSTV()

					// Send 'returned' message
					Session.ChannelMessageSend(config.Conf.Discord.DefaultChannel, fmt.Sprintf("%s: Your server was automatically unbooked.", UserMention))

					// Send 'stv' message, if it uploaded successfully.
					if err == nil {
						Session.ChannelMessageSend(config.Conf.Discord.DefaultChannel, fmt.Sprintf("%s: %s", UserMention, STVMessage))
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
	for i := 0; i < len(servers.Servers); i++ {
		Serv := &servers.Servers[i]

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
