package main

import "fmt"

// HandleQueryError handles incrementing the errorMinute value on the server, and
// notifying an admin via Discord if too many errors occur in a short space of time.
func HandleQueryError(s *Server, err error) {
	s.errorMinutes++

	// Too many notifications. Send a message.
	if s.errorMinutes >= Conf.ErrorThreshold {
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
				s.errorMinutes,
				bookerName,
			)
		} else {
			message = fmt.Sprintf(
				"The server `%s` failed to be contacted after %d retries while unbooked. Check to ensure the server is correctly working.",
				s.Name,
				s.errorMinutes,
			)
		}

		// Reset the error minutes.
		s.errorMinutes = 0

		for _, notificationUser := range Conf.NotificationUsers {
			UserChannel, _ := Session.UserChannelCreate(notificationUser)
			Session.ChannelMessageSend(UserChannel.ID, message)
		}
	}
}
