package servers

import (
	"math/rand"
	"time"

	"github.com/Qixalite/booking-api/client"
)

type ServerRunner struct {
	APIClient     *client.Client
	cachedServers map[string]*client.ServerResource
}

func NewRunner(apiClient *client.Client) *ServerRunner {
	return &ServerRunner{
		APIClient:     apiClient,
		cachedServers: make(map[string]*client.ServerResource),
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func (sr ServerRunner) generatePassword() string {
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

func (sr ServerRunner) getServer(uuid string) (*client.ServerResource, error) {
	if server, ok := sr.cachedServers[uuid]; ok {
		return server, nil
	}

	server, err := sr.APIClient.GetServer(uuid)
	if err != nil {
		return nil, err
	}
	sr.cachedServers[uuid] = &server
	return &server, nil
}

func (sr ServerRunner) Setup(server *Server) (rconPassword string, srvPassword string, err error) {
	// Generate an RCON and server password.
	rconPassword = sr.generatePassword()
	srvPassword = sr.generatePassword()

	// Retrieve the API server instance from the API client.
	apiServer, err := sr.getServer(server.UUID)
	if err != nil {
		return "", "", err
	}

	// Reserve the server, yo.
	err = apiServer.Reserve(sr.APIClient, true)
	if err != nil {
		return "", "", err
	}

	// Set the password on the server.
	err = apiServer.SetPassword(sr.APIClient, rconPassword, srvPassword)

	return rconPassword, srvPassword, err
}

func (sr ServerRunner) Destroy(server *Server) (err error) {
	// Retrieve the API server instance from the API client.
	apiServer, err := sr.getServer(server.UUID)
	if err != nil {
		return err
	}

	// Unreserve the server, yo.
	err = apiServer.Reserve(sr.APIClient, false)
	return err
}

func (sr ServerRunner) Start(server *Server) error {
	// Retrieve the API server instance from the API client.
	apiServer, err := sr.getServer(server.UUID)
	if err != nil {
		return err
	}

	// Start the server.
	err = apiServer.Start(sr.APIClient)

	return err
}

func (sr ServerRunner) Stop(server *Server) error {
	// Retrieve the API server instance from the API client.
	apiServer, err := sr.getServer(server.UUID)
	if err != nil {
		return err
	}

	// Stop the server.
	err = apiServer.Stop(sr.APIClient)

	return err
}

func (sr ServerRunner) UploadSTV(server *Server, uploaderName string) ([]string, error) {
	// Retrieve the API server instance from the API client.
	apiServer, err := sr.getServer(server.UUID)
	if err != nil {
		return nil, err
	}

	// Upload demos.
	demoURLs, err := apiServer.UploadDemos(sr.APIClient, uploaderName)
	if err != nil {
		return nil, err
	}

	demos := make([]string, 0, len(demoURLs))
	for _, demoURL := range demoURLs {
		demos = append(demos, demoURL)
	}

	return demos, nil
}

func (sr ServerRunner) SendCommand(server *Server, command string) error {
	// Retrieve the API server instance from the API client.
	apiServer, err := sr.getServer(server.UUID)
	if err != nil {
		return err
	}

	// Send the command.
	err = apiServer.SendCommand(sr.APIClient, command)

	return err
}

func (sr ServerRunner) Console(server *Server, lines int) ([]string, error) {
	// Retrieve the API server instance from the API client.
	apiServer, err := sr.getServer(server.UUID)
	if err != nil {
		return nil, err
	}
	consoleLines, err := apiServer.Console(sr.APIClient, lines)

	return consoleLines, err
}

func (sr ServerRunner) IsBooked(server *Server) bool {
	apiServer, err := sr.APIClient.GetServer(server.UUID)
	if err != nil {
		return false
	}

	return apiServer.Reserved
}
