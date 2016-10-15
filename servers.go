package main

import ()

type ServerConfig struct {
	Name    string
	Path    string
	Address string
}

var Servers []Server

func SetupServers() {
	for i := 0; i < len(Conf.Servers); i++ {
		s := Conf.Servers[i]
		Servers = append(Servers, Server{
			Name:    s.Name,
			Path:    s.Path,
			Address: s.Address,
		})
	}
}

func GetAvailableServer() *Server {
	for i := 0; i < len(Servers); i++ {
		if Servers[i].IsAvailable() {
			return &Servers[i]
		}
	}

	return nil
}

func GetAvailableServers() []*Server {
	servers := make([]*Server, 0, len(Servers))
	for i := 0; i < len(Servers); i++ {
		if Servers[i].IsAvailable() {
			servers = append(servers, &Servers[i])
		}
	}
	return servers
}
