package main

import (
	"fmt"
	"log"
)

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
		log.Println(fmt.Sprintf("Server %s is available: %d", Servers[i].Name, Servers[i].IsAvailable()))
		if Servers[i].IsAvailable() {
			return &Servers[i]
		}
	}

	return nil
}
