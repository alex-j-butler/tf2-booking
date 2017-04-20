package servers

import (
	"errors"
	"io/ioutil"
	"math"

	yaml "gopkg.in/yaml.v2"
)

// ConfigServerPool is a server pool that is loaded from a configuration file
// and available servers are managed by the booking bot.
type ConfigServerPool struct {
	Servers []*Server
}

// Initialise creates the server pool
// Loads the 'servers.yml' file and sets the default runner for all the servers.
func (csp *ConfigServerPool) Initialise() error {
	configuration, _ := ioutil.ReadFile("./servers.yml")
	err := yaml.Unmarshal(configuration, &csp.Servers)

	if err != nil {
		return err
	}

	// Use the default server runner for all servers.
	for _, server := range csp.Servers {
		server.Init()
	}

	return nil
}

func (csp *ConfigServerPool) GetServers() []*Server {
	return csp.Servers
}

// GetAvailableServer gets the next available server from the server pool.
func (csp *ConfigServerPool) GetAvailableServer() *Server {
	var bestServer *Server
	var bestDiff float64
	servers := csp.GetAvailableServers()

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

// GetAvailableServers gets a slice of all available servers from the server pool.
func (csp *ConfigServerPool) GetAvailableServers() []*Server {
	servers := make([]*Server, 0, len(csp.Servers))
	for _, server := range csp.Servers {
		if server.IsAvailable() {
			servers = append(servers, server)
		}
	}

	return servers
}

func (csp *ConfigServerPool) GetBookedServers() []*Server {
	servers := make([]*Server, 0, len(csp.Servers))
	for _, server := range csp.Servers {
		if !server.IsAvailable() {
			servers = append(servers, server)
		}
	}

	return servers
}

func (csp *ConfigServerPool) GetServerByAddress(address string) (*Server, error) {
	for _, server := range csp.Servers {
		if server.Address == address {
			return server, nil
		}
	}

	return nil, errors.New("Server not found")
}

func (csp *ConfigServerPool) GetServerBySessionName(sessionName string) (*Server, error) {
	for _, server := range csp.Servers {
		if server.SessionName == sessionName {
			return server, nil
		}
	}

	return nil, errors.New("Server not found")
}
