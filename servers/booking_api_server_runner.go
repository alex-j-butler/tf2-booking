package servers

import (
	"math/rand"
	"time"

	"alex-j-butler.com/tf2-booking/booking_api"
	"alex-j-butler.com/tf2-booking/models"
)

type BookingAPIServerRunner struct {
	APIClient *booking_api.BookingClient
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func (br BookingAPIServerRunner) generatePassword() string {
	n := 10
	b := make([]byte, n)

	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func (b BookingAPIServerRunner) Setup(server *Server) (rconPassword string, srvPassword string, err error) {
	// Generate an RCON and server password.
	rconPassword = b.generatePassword()
	srvPassword = b.generatePassword()

	// Set the password on the server.
	b.APIClient.SetPassword(server.Name, rconPassword, srvPassword)

	return rconPassword, srvPassword, nil
}

func (b BookingAPIServerRunner) Start(server *Server) error {
	err := b.APIClient.StartServer(server.Name)

	return err
}

func (b BookingAPIServerRunner) Stop(server *Server) error {
	err := b.APIClient.StopServer(server.Name)

	return err
}

func (b BookingAPIServerRunner) UploadSTV(server *Server) ([]models.Demo, error) {
	return []models.Demo{
		models.Demo{
			URL: "WARNING! Implement UploadSTV in BookingAPIServerRunner!",
		},
	}, nil
}

func (b BookingAPIServerRunner) SendCommand(server *Server, command string) error {
	// Send the command.
	b.APIClient.SendCommand(server.Name, command)
	return nil
}

func (b BookingAPIServerRunner) IsAvailable(server *Server) bool {
	// Attempt to request the server information, if it fails, the server is unavailable.
	_, err := b.APIClient.GetServer(server.Name)
	if err != nil {
		// Unavailable!
		return false
	}

	return true
}

func (b BookingAPIServerRunner) IsBooked(server *Server) bool {
	apiServer, err := b.APIClient.GetServer(server.Name)
	if err != nil {
		return false
	}

	return apiServer.Running
}
