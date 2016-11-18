package servers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/james4k/rcon"
)

type Server struct {
	Name        string `json:"name" yaml:"name"`
	Path        string `json:"path" yaml:"path"`
	Address     string `json:"address" yaml:"address"`
	SessionName string `json:"session_name" yaml:"session_name"`

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
	// Retrieve the RCON password & server password.
	process := exec.Command("sh", "-c", fmt.Sprintf("cd %s; %s/book_server.sh", s.Path, s.Path))
	stdout, _ := process.StdoutPipe()
	stderr, _ := process.StderrPipe()

	var err error
	err = process.Start()

	if err != nil {
		log.Println("Failed to setup server:", err)
		return "", "", errors.New("Your server could not be setup")
	}

	stdoutBytes, _ := ioutil.ReadAll(stdout)
	stderrBytes, _ := ioutil.ReadAll(stderr)

	err = process.Wait()

	if err != nil {
		log.Println("Failed to setup server:", err)
		return "", "", errors.New("Your server could not be setup")
	}

	s.SentWarning = false

	// Trim passwords.
	RCONPassword := strings.TrimSpace(string(stdoutBytes))
	ServerPassword := strings.TrimSpace(string(stderrBytes))

	s.RCONPassword = RCONPassword

	return RCONPassword, ServerPassword, nil
}

// Start the server using a bash script.
// Returns:
//  error - Error of a failed start, or nil if none
func (s *Server) Start() error {
	process := exec.Command("sh", "-c", fmt.Sprintf("cd %s; %s/run r", s.Path, s.Path))

	var err error
	err = process.Start()

	if err != nil {
		log.Println("Process failed to start:", err)
		return errors.New("Your server could not be started")
	}

	err = process.Wait()

	if err != nil {
		log.Println("Process failed to wait:", err)
		return errors.New("Your server could not be started")
	}

	return nil
}

func (s *Server) Stop() error {
	// Stop the STV recording and kick all players cleanly.
	s.SendCommand("tv_stop; kickall \"Server has been unbooked! Thanks for using Qixalite's bookable servers!\"")

	process := exec.Command("sh", "-c", fmt.Sprintf("cd %s; %s/run sp", s.Path, s.Path))

	var err error
	err = process.Start()

	if err != nil {
		log.Println("Process failed to start:", err)
		return errors.New("Your server could not be stopped")
	}

	err = process.Wait()

	if err != nil {
		log.Println("Process failed to wait:", err)
		return errors.New("Your server could not be stopped")
	}

	return nil
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
	// Run upload STV demo script.
	process := exec.Command("sh", "-c", fmt.Sprintf("cd %s; %s/stv.sh", s.Path, s.Path))
	stdout, _ := process.StdoutPipe()

	var err error
	err = process.Start()

	if err != nil {
		log.Println("Failed to upload STV:", err)
		return "", errors.New("Your server failed to upload STV")
	}

	stdoutBytes, _ := ioutil.ReadAll(stdout)

	err = process.Wait()

	if err != nil {
		log.Println("Failed to upload STV:", err)
		return "", errors.New("Your server failed to upload STV")
	}

	Files := strings.Split(string(stdoutBytes), "\n")
	for i := 0; i < len(Files); i++ {
		Files[i] = strings.TrimSpace(Files[i])
	}

	Message := "STV Demo(s) uploaded:"
	for i := 0; i < len(Files); i++ {
		Message = fmt.Sprintf("%s\n\t%s", Message, Files[i])
	}

	return Message, nil
}

func (s *Server) SendCommand(command string) error {
	process := exec.Command("tmux", "send-keys", "-t", s.SessionName, "C-m", command, "C-m")

	var err error
	err = process.Start()

	if err != nil {
		log.Println("Failed to send command:", err)
		return errors.New("Your server failed to respond to commands")
	}

	err = process.Wait()

	if err != nil {
		log.Println("Failed to send command:", err)
		return errors.New("Your server failed to respond to commands")
	}

	return nil
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
