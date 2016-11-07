package main

import "math"

func GetAvailableServer() *Server {
	var bestServer *Server
	var bestDiff float64
	servers := GetAvailableServers()

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

func GetAvailableServers() []*Server {
	servers := make([]*Server, 0, len(Conf.Servers))
	for i := 0; i < len(Conf.Servers); i++ {
		if Conf.Servers[i].IsAvailable() {
			servers = append(servers, &Conf.Servers[i])
		}
	}
	return servers
}

func GetBookedServers() []*Server {
	servers := make([]*Server, 0, len(Conf.Servers))
	for i := 0; i < len(Conf.Servers); i++ {
		if !Conf.Servers[i].IsAvailable() {
			servers = append(servers, &Conf.Servers[i])
		}
	}
	return servers
}
