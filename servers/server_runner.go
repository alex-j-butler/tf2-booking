package servers

import "alex-j-butler.com/tf2-booking/models"

type ServerRunner interface {
	Setup(server *Server) (rconPassword string, srvPassword string, err error)
	Start(server *Server) (err error)
	Stop(server *Server) (err error)
	UploadSTV(server *Server) (demos []models.Demo, err error)
	SendCommand(server *Server, command string) (err error)
	Console(server *Server) (consoleLines []string, err error)

	IsAvailable(server *Server) bool
	IsBooked(server *Server) bool
}
