package servers

type ServerPool interface {
	// Initialise creates the server pool.
	Initialise() error

	GetServers() []*Server
	GetAvailableServer() *Server
	GetAvailableServers() []*Server
	GetBookedServers() []*Server

	GetServerByAddress(address string) (*Server, error)
	GetServerBySessionName(sessionName string) (*Server, error)
}
