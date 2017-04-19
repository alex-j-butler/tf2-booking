package booking_api

type Server struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`

	IPAddress string `json:"ip_address"`
	Port      int    `json:"port"`
	STVPort   int    `json:"stv_port"`

	ServerPassword string `json:"server_password"`
	RCONPassword   string `json:"rcon_password"`

	Executable string   `json:"executable"`
	Options    []string `json:"options"`

	Running bool `json:"running"`
}

type ErrorResponse struct {
	Status  int    `json:"code"`
	Message string `json:"message"`
}

// NextResponse is the response from the 'Next' endpoint
// indicating the next available server.
type NextResponse struct {
	Server
}
