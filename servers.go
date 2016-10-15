package main

import ()

/*
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
*/

func GetAvailableServer() *Server {
	for i := 0; i < len(Conf.Servers); i++ {
		if Conf.Servers[i].IsAvailable() {
			return &Conf.Servers[i]
		}
	}

	return nil
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
