package servers

import "net/http"

type BookingAPIServerRunner struct {
	baseURL string
	client  http.Client
}

func (b BookingAPIServerRunner) Setup(server *Server) (rconPassword string, srvPassword string, err error) {
	// Generate an RCON and server password.
	rconPassword = "example"
	srvPassword = "example"

	return "", "", nil
}
