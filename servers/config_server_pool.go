package servers

import (
	"io/ioutil"
	"math"

	yaml "gopkg.in/yaml.v2"
)

// ConfigServerPool is a server pool that is loaded from a configuration file
// and available servers are managed by the booking bot.
type ConfigServerPool struct {
	Servers []*Server
}

func (csp ConfigServerPool) Initialise() error {
	configuration, _ := ioutil.ReadFile("./servers.yml")
	err := yaml.Unmarshal(configuration, &csp.Servers)

	if err != nil {
		return err
	}

	// Use the default server runner for all servers.
	for _, server := range Servers {
		server.Init()
	}

	return nil
}

func (csp ConfigServerPool) GetAvailableServer() *Server {
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

func (csp ConfigServerPool) GetAvailableServers() []*Server {
	servers := make([]*Server, 0, len(csp.Servers))
	for _, server := range csp.Servers {
		if server.IsAvailable() {
			servers = append(servers, server)
		}
	}

	return servers
}

func (csp ConfigServerPool) GetBookedServers() []*Server {
	servers := make([]*Server, 0, len(csp.Servers))
	for _, server := range csp.Servers {
		if !server.IsAvailable() {
			servers = append(servers, server)
		}
	}

	return servers
}

func (csp ConfigServerPool) GetServerByAddress(address string) *Server {
	for _, server := range csp.Servers {
		if server.Address == address {
			return server
		}
	}

	return nil
}

func (csp ConfigServerPool) GetServerBySessionName(sessionName string) *Server {
	for _, server := range csp.Servers {
		if server.SessionName == sessionName {
			return server
		}
	}

	return nil
}
