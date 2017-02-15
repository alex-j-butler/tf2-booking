package servers

import (
	"errors"
	"math"
	"strings"
)

var Servers map[string]*Server

func InitialiseServers() {
	Servers = make(map[string]*Server)
}

func GetAvailableServer(serverList map[string]*Server) *Server {
	var bestServer *Server
	var bestDiff float64
	servers := GetAvailableServers(serverList)

	// Higher than the maximum a TF2 tickrate can differ.
	bestDiff = 4096.0
	for _, v := range servers {
		if diff := math.Abs(float64(v.TickRate - 66.6666)); diff < bestDiff {
			bestServer = v
			bestDiff = diff
		}
	}

	return bestServer
}

func GetAvailableServers(serverList map[string]*Server) map[string]*Server {
	servers := make(map[string]*Server)
	for k, v := range serverList {
		if v.IsAvailable() {
			servers[k] = v
		}
	}
	return servers
}

func GetBookedServers(serverList map[string]*Server) map[string]*Server {
	servers := make(map[string]*Server)
	for k, v := range serverList {
		if !v.IsBooked() {
			servers[k] = v
		}
	}
	return servers
}

func GetServerByUUID(serverList map[string]*Server, uuid string) (*Server, error) {
	for k, v := range serverList {
		if strings.EqualFold(v.UUID, uuid) {
			return serverList[k], nil
		}
	}

	return nil, errors.New("Server not found.")
}

func GetServerByName(serverList map[string]*Server, name string) (*Server, error) {
	for k, v := range serverList {
		if strings.EqualFold(v.Name, name) {
			return serverList[k], nil
		}
	}

	return nil, errors.New("Server not found.")
}

func GetServerByAddress(serverList map[string]*Server, address string) (*Server, error) {
	for k, v := range serverList {
		if v.Address == address {
			return serverList[k], nil
		}
	}

	return nil, errors.New("Server not found.")
}

func GetServerBySessionName(serverList map[string]*Server, sessionName string) (*Server, error) {
	for k, v := range serverList {
		if v.SessionName == sessionName {
			return serverList[k], nil
		}
	}

	return nil, errors.New("Server not found.")
}
