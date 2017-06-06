package servers

import (
	"math/rand"
	"time"

	"github.com/Qixalite/booking-api/client"

	null "gopkg.in/nullbio/null.v6"

	"alex-j-butler.com/tf2-booking/models"
)

type BookingAPIServerRunner struct {
	APIClient     *client.Client
	cachedServers map[string]*client.RPCServer
}

func NewBookingAPIRunner(apiClient *client.Client) *BookingAPIServerRunner {
	return &BookingAPIServerRunner{
		APIClient:     apiClient,
		cachedServers: make(map[string]*client.RPCServer),
	}
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

func (b BookingAPIServerRunner) getServer(uuid string) (*client.RPCServer, error) {
	if server, ok := b.cachedServers[uuid]; ok {
		return server, nil
	}

	server, err := b.APIClient.GetServer(uuid)
	if err != nil {
		return nil, err
	}
	b.cachedServers[uuid] = &server
	return &server, nil
}

func (b BookingAPIServerRunner) Setup(server *Server) (rconPassword string, srvPassword string, err error) {
	// Generate an RCON and server password.
	rconPassword = b.generatePassword()
	srvPassword = b.generatePassword()

	// Retrieve the API server instance from the API client.
	apiServer, err := b.getServer(server.UUID)
	if err != nil {
		return "", "", err
	}

	// Set the password on the server.
	err = apiServer.SetPassword(b.APIClient, rconPassword, srvPassword)

	return rconPassword, srvPassword, err
}

func (b BookingAPIServerRunner) Start(server *Server) error {
	// Retrieve the API server instance from the API client.
	apiServer, err := b.getServer(server.UUID)
	if err != nil {
		return err
	}

	// Start the server.
	err = apiServer.Start(b.APIClient)

	return err
}

func (b BookingAPIServerRunner) Stop(server *Server) error {
	// Retrieve the API server instance from the API client.
	apiServer, err := b.getServer(server.UUID)
	if err != nil {
		return err
	}

	// Stop the server.
	err = apiServer.Stop(b.APIClient)

	return err
}

func (b BookingAPIServerRunner) UploadSTV(server *Server) ([]models.Demo, error) {
	// Retrieve the API server instance from the API client.
	apiServer, err := b.getServer(server.UUID)
	if err != nil {
		return nil, err
	}

	// Upload demos.
	demoURLs, err := apiServer.UploadDemos(b.APIClient)
	if err != nil {
		return nil, err
	}

	demos := make([]models.Demo, 0, len(demoURLs))
	for _, demoURL := range demoURLs {
		demos = append(demos, models.Demo{UploadedTime: null.TimeFrom(time.Now()), URL: demoURL})
	}

	return demos, nil
}

func (b BookingAPIServerRunner) SendCommand(server *Server, command string) error {
	// Retrieve the API server instance from the API client.
	apiServer, err := b.getServer(server.UUID)
	if err != nil {
		return err
	}

	// Send the command.
	err = apiServer.SendCommand(b.APIClient, command)

	return err
}

func (b BookingAPIServerRunner) Console(server *Server, lines int) ([]string, error) {
	// Retrieve the API server instance from the API client.
	apiServer, err := b.getServer(server.UUID)
	if err != nil {
		return nil, err
	}
	consoleLines, err := apiServer.Console(b.APIClient, lines)

	return consoleLines, err
}

func (b BookingAPIServerRunner) IsAvailable(server *Server) bool {
	// Attempt to request the server information, if it fails, the server is unavailable.
	_, err := b.APIClient.GetServer(server.UUID)
	if err != nil {
		// Unavailable!
		return false
	}

	return true
}

func (b BookingAPIServerRunner) IsBooked(server *Server) bool {
	apiServer, err := b.APIClient.GetServer(server.UUID)
	if err != nil {
		return false
	}

	return apiServer.Running
}
