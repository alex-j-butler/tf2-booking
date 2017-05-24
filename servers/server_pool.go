package servers

type ServerPool interface {
	// Initialise creates the server pool.
	Initialise() error

	GetServers() []*Server
	GetAvailableServer() *Server
	GetAvailableServers() []*Server
	GetBookedServers() []*Server

	GetServerByAddress(address string) (*Server, error)
	GetServerByName(name string) (*Server, error)
	GetServerByRedisName(redisName string) (*Server, error)
}
