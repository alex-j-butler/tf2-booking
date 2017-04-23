package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"alex-j-butler.com/tf2-booking/config"
	"alex-j-butler.com/tf2-booking/globals"
	"alex-j-butler.com/tf2-booking/servers"
	"alex-j-butler.com/tf2-booking/util"

	"strings"

	"github.com/kidoman/go-steam"
)

// Check if any servers are ready to be unbooked by the automatic timeout after 4 hours.
func CheckUnbookServers() {
	// Iterate through servers.
	for _, Serv := range pool.GetBookedServers() {
		// Send the timelimit warning notification, if required.
		if !Serv.SentWarning && (Serv.ReturnDate.Add(config.Conf.Booking.WarningDuration.Duration)).Before(time.Now()) {
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

		// Notify the user that their booking is about to timeout due to idle.
		// This will happen x minutes before the max idle minute time (set to 15), configurable through the IdleWarningDuration option in the configuration file.
		// This will only happen once, unless the idle timeout is reset.

		// TODO: Move this to the configuration file?
		maxIdleMinutes := 15
		if !Serv.SentIdleWarning && (maxIdleMinutes-Serv.IdleMinutes) <= config.Conf.Booking.IdleWarningDuration {
			// Only allow this message to be sent once.
			Serv.SentIdleWarning = true

			// Calculate minutes remaining before the idle timeout.
			minutesRemaining := maxIdleMinutes - Serv.IdleMinutes

			// Send warning message in server.
			Serv.SendCommand(
				fmt.Sprintf(
					"say Your booking will timeout in %d %s, to prevent this, make sure 2 players are on the server.",
					minutesRemaining,
					util.PluralMinutes(minutesRemaining),
				),
			)

			UserMention := Serv.GetBookerMention()
			// Send warning message in Discord.
			Session.ChannelMessageSend(
				config.Conf.Discord.DefaultChannel,
				fmt.Sprintf("%s: Your booking will timeout in %d %s. To prevent this, make sure 2 players are on the server.", UserMention, minutesRemaining, util.PluralMinutes(minutesRemaining)),
			)
		}

		// Check if their server is past the return date.
		if Serv.ReturnDate.Before(time.Now()) {
			UserID := Serv.GetBooker()
			UserMention := Serv.GetBookerMention()

			// Remove the user's booked state.
			if err := globals.RedisClient.Set(fmt.Sprintf("user.%s", UserID), "", 0).Err(); err != nil {
				log.Println("Redis error:", err)
				log.Println("Failed to set user information for user:", UserID)
				return
			}

			// Unbook the server.
			Serv.Unbook()
			Serv.Stop()

			// Upload STV demos
			STVMessage, err := Serv.UploadSTV()

			// Send 'returned' message
			Session.ChannelMessageSend(config.Conf.Discord.DefaultChannel, fmt.Sprintf("%s: Your server was automatically unbooked (timelimit reached).", UserMention))

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
	for _, Serv := range pool.GetBookedServers() {
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
				// Reset the number of idle minutes, and allow the timeout warning message to be sent again.
				s.SentIdleWarning = false
				s.ResetIdleMinutes()
			}

			if s.GetIdleMinutes() >= config.Conf.Booking.MaxIdleMinutes {
				UserID := s.GetBooker()
				UserMention := s.GetBookerMention()

				// Reset the idle minutes.
				s.ResetIdleMinutes()

				// Remove the user's booked state.
				if err := globals.RedisClient.Set(fmt.Sprintf("user.%s", UserID), "", 0).Err(); err != nil {
					log.Println("Redis error:", err)
					log.Println("Failed to set user information for user:", UserID)
					return
				}

				// Unbook the server.
				s.Unbook()
				s.Stop()

				// Upload STV demos
				STVMessage, err := s.UploadSTV()

				// Send 'returned' message
				Session.ChannelMessageSend(config.Conf.Discord.DefaultChannel, fmt.Sprintf("%s: Your server was automatically unbooked (not enough players).", UserMention))

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

func Cron10Seconds() {
	// Iterate through servers.
	for _, Serv := range pool.GetBookedServers() {
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

			tickrate := 1000.0 / 15.0
			if math.Abs(float64(s.TickRate)-tickrate) > 5.0 && Serv.NextPerformanceWarning.Before(time.Now()) {
				// Only allow this message to be sent once.
				Serv.NextPerformanceWarning = time.Now().Add(5 * time.Minute)

				// The 'var' of the server is too high, notify admins.
				message := fmt.Sprintf(
					"The server `%s` may be performing poorly (tickrate %f, var %f). Check the ensure the server is not lagging.",
					s.Name,
					s.TickRate,
					math.Abs(float64(s.TickRate)-tickrate),
				)
				for _, notificationUser := range config.Conf.Discord.NotificationUsers {
					UserChannel, _ := Session.UserChannelCreate(notificationUser)
					Session.ChannelMessageSend(UserChannel.ID, message)
				}
			}

			s.Update(globals.RedisClient)
		}(Serv)

		// TF2Center/TF2Stadium checking.
		// Here, we want to get the tags of the server, and check if they contain the words 'TF2Center' or 'TF2Stadium', and if they do
		// we want to send a message to the default channel letting the user know that they're using a Qixalite server for a lobby,
		// and they should ensure they get 2 people in the server, otherwise the server will unbook in 15 minutes, which cannot be extended using
		// the 'extend' command.
		go func(s *servers.Server) {

			// Query the server for tags.
			server, err := steam.Connect(Serv.Address)
			if err != nil {
				log.Println(fmt.Sprintf("Failed to connect to server \"%s\":", s.Name), err)
				return
			}

			resp, err := server.Info()
			if err != nil {
				log.Println(fmt.Sprintf("Failed to connect to server \"%s\":", s.Name), err)
				return
			}

			// Check for matches of 'tf2center' or 'tf2stadium' in the tags.
			lowercase := strings.ToLower(resp.Keywords)
			if strings.Contains(lowercase, "tf2center") || strings.Contains(lowercase, "tf2stadium") {
				// Send the lobby warning.
				if !Serv.SentLobbyWarning {
					// Don't allow this message again.
					Serv.SentLobbyWarning = true

					// Send a warning message.
					Session.ChannelMessageSend(config.Conf.Discord.DefaultChannel, fmt.Sprintf("%s: We noticed you're running a TF2Center lobby, make sure you have 2 people on the server, otherwise your server will unbook after 15 minutes! If you need to get the password after it's been changed, type `send password` into this channel and we'll send the updated password.", Serv.BookerMention))
				}
			}

		}(Serv)
	}
}
