package main

import (
	"fmt"

	"alex-j-butler.com/tf2-booking/config"
	"alex-j-butler.com/tf2-booking/servers"
)

// HandleQueryError handles incrementing the errorMinute value on the server, and
// notifying an admin via Discord if too many errors occur in a short space of time.
func HandleQueryError(s *servers.Server, err error) {
	s.ErrorMinutes++

	// Too many notifications. Send a message.
	if s.ErrorMinutes >= config.Conf.ErrorThreshold {
		var message string
		bookerName := "Unknown"
		if !s.IsAvailable() {
			u, err := Session.User(s.GetBooker())
			if err == nil {
				bookerName = u.Username
			}

			message = fmt.Sprintf(
				"The server `%s` failed to be contacted after %d retries after being booked by `%s`. Check to ensure the server is correctly working.",
				s.Name,
				s.ErrorMinutes,
				bookerName,
			)
		} else {
			message = fmt.Sprintf(
				"The server `%s` failed to be contacted after %d retries while unbooked. Check to ensure the server is correctly working.",
				s.Name,
				s.ErrorMinutes,
			)
		}

		// Reset the error minutes.
		s.ErrorMinutes = 0

		for _, notificationUser := range config.Conf.NotificationUsers {
			UserChannel, _ := Session.UserChannelCreate(notificationUser)
			Session.ChannelMessageSend(UserChannel.ID, message)
		}
	}
}
