package servers

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/james4k/rcon"
)

type Server struct {
	Name        string `json:"name" yaml:"name"`
	Path        string `json:"path" yaml:"path"`
	Address     string `json:"address" yaml:"address"`
	SessionName string `json:"session_name" yaml:"session_name"`
	Type        string `json:"type" yaml:"type"`

	SSHAddress    string `json:"ssh_address" yaml:"ssh_address"`
	SSHUsername   string `json:"ssh_username" yaml:"ssh_username"`
	SSHPassword   string `json:"ssh_password" yaml:"ssh_password"`
	SSHPrivateKey string `json:"ssh_private_key" yaml:"ssh_private_key"`

	SentWarning bool
	ReturnDate  time.Time

	// Last known RCON password.
	// If this RCON password is invalid, the server can send a tmux command to reset it.
	RCONPassword string

	// Average tick rate reported by the server.
	TickRate float32

	// Number of tick rate measurements (used internally for calculating a new average).
	TickRateMeasurements int

	booked     bool
	bookedDate time.Time

	booker        string
	bookerMention string

	host string
	port int

	IdleMinutes  int
	ErrorMinutes int
}

// Returns whether the server is currently available for booking.
func (s *Server) IsAvailable() bool {
	return !s.booked
}

func (s *Server) GetBookedTime() time.Time {
	return s.bookedDate
}

func (s *Server) GetBooker() string {
	return s.booker
}

func (s *Server) GetBookerMention() string {
	return s.bookerMention
}

func (s *Server) GetIdleMinutes() int {
	return s.IdleMinutes
}

func (s *Server) AddIdleMinute() {
	s.IdleMinutes++
}

func (s *Server) ResetIdleMinutes() {
	s.IdleMinutes = 0
}

// Setup the server with a randomised RCON password & server password from a bash script.
// Returns:
//  string - RCON password
//  string - Server password
//  error - Error of a failed setup, or nil if none
func (s *Server) Setup() (string, string, error) {
	switch s.Type {
	case "local":
		return s.setupLocalServer()

	case "ssh":
		return s.setupSSHServer()

	default:
		return "", "", errors.New("Invalid server type.")
	}
}

// Start the server using a bash script.
// Returns:
//  error - Error of a failed start, or nil if none
func (s *Server) Start() error {
	switch s.Type {
	case "local":
		return s.startLocalServer()

	case "ssh":
		return s.startRemoteServer()

	default:
		return errors.New("Invalid server type.")
	}
}

func (s *Server) Stop() error {
	switch s.Type {
	case "local":
		return s.stopLocalServer()

	case "ssh":
		return s.stopRemoteServer()

	default:
		return errors.New("Invalid server type.")
	}
}

func (s *Server) Book(user *discordgo.User, duration time.Duration) (string, string, error) {
	if s.booked == true {
		return "", "", errors.New("Server is already booked")
	}

	// Set the server variables.
	s.ReturnDate = time.Now().Add(duration)
	s.booked = true
	s.bookedDate = time.Now()
	s.booker = user.ID
	s.bookerMention = fmt.Sprintf("<@%s>", user.ID)
	s.SentWarning = false
	s.IdleMinutes = 0
	s.ErrorMinutes = 0

	var err error

	// Setup the server.
	RCONPassword, ServerPassword, err := s.Setup()

	if err != nil {
		return "", "", err
	}

	return RCONPassword, ServerPassword, err
}

func (s *Server) Unbook() error {
	if s.booked == false {
		return errors.New("Server is not booked")
	}

	// Set the server variables.
	s.ReturnDate = time.Time{}
	s.booked = false
	s.bookedDate = time.Time{}
	s.booker = ""
	s.bookerMention = ""
	s.SentWarning = false
	s.IdleMinutes = 0
	s.ErrorMinutes = 0

	return nil
}

func (s *Server) ExtendBooking(amount time.Duration) {
	// Add duration to the return date.
	s.ReturnDate = s.ReturnDate.Add(amount)
}

func (s *Server) UploadSTV() (string, error) {
	switch s.Type {
	case "local":
		return s.uploadSTVLocalServer()

	case "ssh":
		return s.uploadSTVRemoteServer()

	default:
		return "", errors.New("Invalid server type.")
	}
}

func (s *Server) SendCommand(command string) error {
	switch s.Type {
	case "local":
		return s.sendCommandLocalServer(command)

	case "ssh":
		return s.sendCommandRemoteServer(command)

	default:
		return errors.New("Invalid server type.")
	}
}

func (s *Server) SendRCONCommand(command string) (string, error) {
	rc, err := rcon.Dial(s.Address, s.RCONPassword)

	if err == rcon.ErrAuthFailed {
		// Attempt to reset RCON password.
		s.SendCommand(fmt.Sprintf("rcon_password %s", s.RCONPassword))

		rc, err = rcon.Dial(s.Address, s.RCONPassword)
	}

	if err != nil {
		return "", err
	}

	// Run the command.
	_, err = rc.Write(command)

	if err != nil {
		return "", err
	}

	// Grab the output.
	output, _, err := rc.Read()

	if err != nil {
		return "", err
	}

	return output, nil
}
