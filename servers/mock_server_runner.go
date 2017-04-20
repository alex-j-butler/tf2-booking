package servers

import (
	"log"
	"time"

	null "gopkg.in/nullbio/null.v6"

	"alex-j-butler.com/tf2-booking/models"
)

type MockServerRunner struct {
}

// Setup the server by running the 'book server' bash script.
func (s MockServerRunner) Setup(server *Server) (rconPassword string, srvPassword string, err error) {
	log.Println("[Mock]", "Server setup")

	return "rcon_password", "server_password", nil
}

func (s MockServerRunner) Start(server *Server) (err error) {
	log.Println("[Mock]", "Server start")

	return nil
}

func (s MockServerRunner) Stop(server *Server) (err error) {
	log.Println("[Mock]", "Server stop")

	return nil
}

func (s MockServerRunner) UploadSTV(server *Server) (demos []models.Demo, err error) {
	log.Println("[Mock]", "Server uploadSTV")

	return []models.Demo{
		{
			URL:          "https://dl.qixalite.com/stv/mock.dem",
			UploadedTime: null.TimeFrom(time.Now()),
		},
	}, nil
}

func (s MockServerRunner) SendCommand(server *Server, command string) (err error) {
	log.Println("[Mock]", "Server sendCommand")

	return nil
}
