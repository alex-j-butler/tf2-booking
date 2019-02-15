package servers

import (
	"errors"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

var Servers []*Server

func InitialiseServers() {
	configuration, _ := ioutil.ReadFile("./servers.yml")
	err := yaml.Unmarshal(configuration, &Servers)

	if err != nil {
		log.Println("Failed to initialise server configuration:", err)
	}
}

func GetAvailableServer(serverList []*Server) *Server {
	servers := GetAvailableServers(serverList)
	if len(servers) > 0 {
		return servers[0]
	}

	return nil
}

func GetAvailableServers(serverList []*Server) []*Server {
	servers := make([]*Server, 0, len(serverList))
	for i := 0; i < len(serverList); i++ {
		if !serverList[i].IsBooked() {
			servers = append(servers, serverList[i])
		}
	}
	return servers
}

func GetBookedServers(serverList []*Server) []*Server {
	servers := make([]*Server, 0, len(serverList))
	for i := 0; i < len(serverList); i++ {
		if serverList[i].IsBooked() {
			servers = append(servers, serverList[i])
		}
	}
	return servers
}

func GetServerByAddress(serverList []*Server, address string) (*Server, error) {
	for i := 0; i < len(serverList); i++ {
		if serverList[i].Address == address {
			return serverList[i], nil
		}
	}

	return nil, errors.New("Server not found.")
}

func GetServerBySessionName(serverList []*Server, sessionName string) (*Server, error) {
	for i := 0; i < len(serverList); i++ {
		if serverList[i].UUID == sessionName {
			return serverList[i], nil
		}
	}

	return nil, errors.New("Server not found.")
}
