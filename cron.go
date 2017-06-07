package main

import (
	"fmt"
	"log"

	"alex-j-butler.com/tf2-booking/config"
	"alex-j-butler.com/tf2-booking/globals"
	"alex-j-butler.com/tf2-booking/servers"

	"strings"

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

func Cron10Seconds() {
	// Iterate through servers.
	for _, Serv := range pool.GetBookedServers() {
		if Serv.IsBooked() {
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
}

func Cron1Minute() {
	err := UpdateGameString()
	if err != nil {
		log.Println("Failed to update game string:", err)
	}
}
