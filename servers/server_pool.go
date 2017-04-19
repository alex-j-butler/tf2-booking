package servers

type ServerPool interface {
	// Initialise creates the server pool.
	Initialise() error

	GetAvailableServer() *Server
	GetAvailableServers() []*Server
	GetBookedServers() []*Server

	GetServerByAddress(address string) *Server
	GetServerBySessionName(sessionName string) *Server
}
