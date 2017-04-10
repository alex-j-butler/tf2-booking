package servers

import "alex-j-butler.com/tf2-booking/models"

type Feature int

const (
	ServerCommandResponse Feature = 1
)

type ServerRunner interface {
	Supports(feature Feature) bool

	Setup(server *Server) (rconPassword string, srvPassword string, err error)
	Start(server *Server) (err error)
	Stop(server *Server) (err error)
	UploadSTV(server *Server) (demos []models.Demo, err error)
	SendCommand(server *Server, command string) (err error)
}
