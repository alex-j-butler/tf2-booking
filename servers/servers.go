package servers

import (
	"errors"
	"io/ioutil"
	"log"
	"math"

	yaml "gopkg.in/yaml.v2"
)

var Servers []*Server

func InitialiseServers() {
	configuration, _ := ioutil.ReadFile("./servers.yml")
	err := yaml.Unmarshal(configuration, &Servers)

	if err != nil {
		log.Println("Failed to initialise server configuration:", err)
	}

	// Use the default server runner for all servers.
	for _, server := range Servers {
		server.Init()
	}
}

func GetAvailableServer(serverList []*Server) *Server {
	var bestServer *Server
	var bestDiff float64
	servers := GetAvailableServers(serverList)

	// Higher than the maximum a TF2 tickrate can differ.
	bestDiff = 4096.0
	for i := 0; i < len(servers); i++ {
		server := servers[i]

		if diff := math.Abs(float64(server.TickRate - 66.6666)); diff < bestDiff {
			bestServer = server
			bestDiff = diff
		}
	}

	return bestServer
}

func GetAvailableServers(serverList []*Server) []*Server {
	servers := make([]*Server, 0, len(serverList))
	for i := 0; i < len(serverList); i++ {
		if !serverList[i].IsBooked() && serverList[i].Available() {
			servers = append(servers, serverList[i])
		}
	}
	return servers
}

func GetBookedServers(serverList []*Server) []*Server {
	servers := make([]*Server, 0, len(serverList))
	for i := 0; i < len(serverList); i++ {
		if serverList[i].IsBooked() && serverList[i].Available() {
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
		if serverList[i].SessionName == sessionName {
			return serverList[i], nil
		}
	}

	return nil, errors.New("Server not found.")
}
