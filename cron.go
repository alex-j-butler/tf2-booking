package main

import (
	"fmt"
	"log"

	"alex-j-butler.com/tf2-booking/config"
	"alex-j-butler.com/tf2-booking/globals"
	"alex-j-butler.com/tf2-booking/servers"

	"github.com/kidoman/go-steam"
)

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

			if s.IdleMinutes >= config.Conf.Booking.MaxIdleMinutes {
				UserID := s.Booker
				UserMention := s.BookerMention

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

func Cron1Minute() {
	err := UpdateGameString()
	if err != nil {
		log.Println("Failed to update game string:", err)
	}
}
