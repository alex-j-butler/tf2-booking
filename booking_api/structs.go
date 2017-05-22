package booking_api

type ErrorResponse struct {
	Status  int    `json:"code"`
	Message string `json:"message"`
}

// NextResponse is the response from the 'Next' endpoint
// indicating the next available server.
type NextResponse struct {
	Server
}

type ListAllResponse struct {
	Servers          []Server `json:"servers"`
	ServersAvailable int      `json:"servers_available"`
}

type ConsoleServerResponse struct {
	ConsoleLines []string `json:"console_lines"`
	Lines        int      `json:"lines"`
}

type StartServerReq struct {
	UUID string `json:"uuid"`
}

type StopServerReq struct {
	UUID string `json:"uuid"`
}

type SetPasswordReq struct {
	UUID           string `json:"uuid"`
	RCONPassword   string `json:"rcon_password"`
	ServerPassword string `json:"server_password"`
}

type SendCommandReq struct {
	UUID    string `json:"uuid"`
	Command string `json:"command"`
}
