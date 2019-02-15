package servers

import (
	"errors"
	"fmt"
	"log"
	"path"

	"github.com/Qixalite/booking-api/client"
)

// APIServerPool is a server pool that is loaded from the booking API.
type APIServerPool struct {
	Tag       string
	APIClient *client.Client
}

func (asp *APIServerPool) Initialise() error {
	return nil
}

func (asp *APIServerPool) GetServers() []*Server {
	// Create a slice of servers.
	allServers, err := asp.servers()
	if err != nil {
		return []*Server{}
	}

	return allServers
}

func (asp *APIServerPool) GetAvailableServer() *Server {
	servers := asp.GetAvailableServers()
	if len(servers) > 0 {
		return servers[0]
	}

	return nil
}

func (asp *APIServerPool) servers() ([]*Server, error) {
	apiServers, err := asp.APIClient.GetServersByTag(asp.Tag)
	if err != nil {
		return nil, err
	}

	log.Println("servers lookup:")

	// Convert all of the servers returned from the API to
	// a booking server.
	servers := make([]*Server, 0, len(apiServers))
	for _, apiServer := range apiServers {
		log.Println(fmt.Sprintf("server (name:%s, running:%t, reserved:%t)", apiServer.Name, apiServer.Running, apiServer.Reserved))

		server := &Server{
			UUID:         apiServer.UUID,
			Name:         apiServer.Name,
			Path:         path.Dir(apiServer.Executable),
			Address:      fmt.Sprintf("%s:%d", apiServer.IPAddress, apiServer.Port),
			STVAddress:   fmt.Sprintf("%s:%d", apiServer.IPAddress, apiServer.STVPort),
			RCONPassword: apiServer.RCONPassword,
			Running:      apiServer.Running,
			Reserved:     apiServer.Reserved,
		}
		server.Runner = NewRunner(asp.APIClient)
		servers = append(servers, server)
	}

	return servers, nil
}

func (asp *APIServerPool) GetAvailableServers() []*Server {
	// Create a slice of servers.
	allServers, err := asp.servers()
	if err != nil {
		return []*Server{}
	}

	servers := make([]*Server, 0)

	// Convert all of the servers returned from the API to
	// a booking server.
	for _, server := range allServers {
		if !server.Reserved {
			servers = append(servers, server)
		}
	}

	return servers
}

func (asp *APIServerPool) GetBookedServers() []*Server {
	// Create a slice of servers.
	allServers, err := asp.servers()
	if err != nil {
		return []*Server{}
	}

	servers := make([]*Server, 0)

	// Convert all of the servers returned from the API to
	// a booking server.
	for _, server := range allServers {
		if server.Reserved {
			servers = append(servers, server)
		}
	}

	return servers
}

func (asp *APIServerPool) GetServerByAddress(address string) (*Server, error) {
	allServers, err := asp.servers()
	if err != nil {
		allServers = []*Server{}
	}

	for _, server := range allServers {
		if server.Address == address {
			return server, nil
		}
	}

	return nil, errors.New("Server not found")
}

func (asp *APIServerPool) GetServerByName(name string) (*Server, error) {
	allServers, err := asp.servers()
	if err != nil {
		allServers = []*Server{}
	}

	for _, server := range allServers {
		if server.Name == name {
			return server, nil
		}
	}

	return nil, errors.New("Server not found")
}

func (asp *APIServerPool) GetServerByUUID(uuid string) (*Server, error) {
	allServers, err := asp.servers()
	if err != nil {
		allServers = []*Server{}
	}

	for _, server := range allServers {
		if server.UUID == uuid {
			return server, nil
		}
	}

	return nil, errors.New("Server not found")
}
