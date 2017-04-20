package servers

import (
	"errors"
	"fmt"
	"math"
	"path"

	"alex-j-butler.com/tf2-booking/booking_api"
)

// APIServerPool is a server pool that is loaded from the booking API.
type APIServerPool struct {
	Tag           string
	CachedServers map[string]*Server
	APIClient     *booking_api.BookingClient
}

func (asp *APIServerPool) Initialise() error {
	asp.CachedServers = make(map[string]*Server)

	return nil
}

func (asp *APIServerPool) GetServers() []*Server {
	// Create a slice of servers.
	servers := make([]*Server, 0)

	// Update the server cache.
	asp.updateCache()

	for _, server := range asp.CachedServers {
		servers = append(servers, server)
	}

	return servers
}

func (asp *APIServerPool) GetAvailableServer() *Server {
	var bestServer *Server
	var bestDiff float64
	servers := asp.GetAvailableServers()

	// Higher than the maximum a TF2 tickrate can differ.
	bestDiff = 4096.0
	for _, server := range servers {
		if diff := math.Abs(float64(server.TickRate - 66.6666)); diff < bestDiff {
			bestServer = server
			bestDiff = diff
		}
	}

	// Return the best available server, may be nil if no servers are available.
	return bestServer
}

func (asp *APIServerPool) updateCache() error {
	apiServers, err := asp.APIClient.GetServersByTag(asp.Tag)
	if err != nil {
		return err
	}

	// Convert all of the servers returned from the API to
	// a booking server.
	for _, apiServer := range apiServers {
		// Check if we've seen this server before, and get the server it's mapped to.
		if _, ok := asp.CachedServers[apiServer.Name]; !ok {
			server := &Server{
				Name:         apiServer.Name,
				Path:         path.Dir(apiServer.Executable),
				Address:      fmt.Sprintf("%s:%d", apiServer.IPAddress, apiServer.Port),
				STVAddress:   fmt.Sprintf("%s:%d", apiServer.IPAddress, apiServer.STVPort),
				SessionName:  apiServer.Name,
				RCONPassword: apiServer.RCONPassword,
			}
			// server.Init()
			server.Runner = &BookingAPIServerRunner{APIClient: asp.APIClient}
			asp.CachedServers[apiServer.Name] = server
		}
	}

	return nil
}

func (asp *APIServerPool) GetAvailableServers() []*Server {
	// Create a slice of servers.
	servers := make([]*Server, 0)

	// Update the server cache.
	asp.updateCache()

	// Convert all of the servers returned from the API to
	// a booking server.
	for _, server := range asp.CachedServers {
		// Server is unavailable if it's already running.
		if !server.IsAvailable() {
			continue
		}

		servers = append(servers, server)
	}

	return servers
}

func (asp *APIServerPool) GetBookedServers() []*Server {
	// Create a slice of servers.
	servers := make([]*Server, 0)

	// Update the server cache.
	asp.updateCache()

	// Convert all of the servers returned from the API to
	// a booking server.
	for _, server := range asp.CachedServers {
		// Server is unavailable if it's already running.
		if server.IsAvailable() {
			continue
		}

		servers = append(servers, server)
	}

	return servers
}

func (asp *APIServerPool) GetServerByAddress(address string) (*Server, error) {
	// Update server cache.
	asp.updateCache()

	for _, server := range asp.CachedServers {
		if server.Address == address {
			return server, nil
		}
	}

	return nil, errors.New("Server not found")
}

func (asp *APIServerPool) GetServerBySessionName(sessionName string) (*Server, error) {
	// Update server cache.
	asp.updateCache()

	for _, server := range asp.CachedServers {
		if server.SessionName == sessionName {
			return server, nil
		}
	}

	return nil, errors.New("Server not found")
}
