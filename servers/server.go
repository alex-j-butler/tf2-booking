package servers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"alex-j-butler.com/tf2-booking/config"
	"alex-j-butler.com/tf2-booking/database"
	"alex-j-butler.com/tf2-booking/models"
	"github.com/bwmarrin/discordgo"
	"github.com/james4k/rcon"
	"github.com/vattle/sqlboiler/queries/qm"
	"gopkg.in/nullbio/null.v6"
)

type Server struct {
	Name        string `json:"name" yaml:"name"`
	Path        string `json:"path" yaml:"path"`
	Address     string `json:"address" yaml:"address"`
	STVAddress  string `json:"stv_address" yaml:"stv_address"`
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

	// Booking ID that the server is currently associated with.
	bookingID int

	IdleMinutes  int
	ErrorMinutes int
}

// IsAvailable returns whether the server is currently available for booking.
// Currently a server is available for booking if it is not being booked by another user,
// in the future, this could be extended to block servers from being used (for example, if they are down).
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

// GetCurrentPassword retrieves the current server password from the server.
func (s *Server) GetCurrentPassword() (string, error) {
	svPasswordResp, err := s.SendRCONCommand("sv_password")
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile("\"sv_password\" = \"(.+)\" \\( def")
	matches := re.FindStringSubmatch(svPasswordResp)

	if len(matches) == 2 {
		return matches[1], nil
	}

	return "", errors.New("Invalid sv_password response")
}

// Setup the server with a randomised RCON password & server password from a bash script.
// Returns:
//  string - RCON password
//  string - Server password
//  error - Error of a failed setup, or nil if none
func (s *Server) Setup() (string, string, error) {
	// Retrieve the RCON password & server password.
	process := exec.Command(
		"sh",
		"-c",
		fmt.Sprintf(
			"cd %s; %s/%s",
			s.Path,
			s.Path,
			config.Conf.Booking.SetupCommand,
		),
	)
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
	process := exec.Command(
		"sh",
		"-c",
		fmt.Sprintf(
			"cd %s; %s/%s",
			s.Path,
			s.Path,
			config.Conf.Booking.StartCommand,
		),
	)

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
	KickCommand := fmt.Sprintf("tv_stop; kickall \"%s\"", config.Conf.Booking.KickMessage)
	s.SendCommand(KickCommand)

	process := exec.Command(
		"sh",
		"-c",
		fmt.Sprintf(
			"cd %s; %s/%s",
			s.Path,
			s.Path,
			config.Conf.Booking.StopCommand,
		),
	)

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

	// Tries to select the user by discord id,
	// if no record is found, insert a new record.
	dbUser, err := models.Users(database.DB, qm.Where("discord_id=?", user.ID)).One()
	if err != nil {
		// Insert new record.
		var newUser models.User
		newUser.DiscordID = null.StringFrom(user.ID)
		newUser.Name = null.StringFrom(user.Username)

		err = newUser.Insert(database.DB)

		if err != nil {
			log.Println("Database error:", err)
			return "", "", errors.New("User record could not be created")
		}

		dbUser = &newUser
	}

	// Adds a new booking to the database
	// and set the booking id.
	var booking models.Booking
	booking.SetBooker(database.DB, false, dbUser)
	booking.ServerName = s.Name
	booking.BookedTime = null.TimeFrom(time.Now())
	err = booking.Insert(database.DB)

	if err != nil {
		log.Println("Database error:", err)
		return "", "", errors.New("Server record could not be created")
	}

	s.bookingID = booking.BookingID

	// Set the server variables.
	s.ReturnDate = time.Now().Add(duration)
	s.booked = true
	s.bookedDate = time.Now()
	s.booker = user.ID
	s.bookerMention = fmt.Sprintf("<@%s>", user.ID)
	s.SentWarning = false
	s.IdleMinutes = 0
	s.ErrorMinutes = 0

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

	booking, err := models.FindBooking(database.DB, s.bookingID)
	if err != nil {
		return errors.New("Server record could not be updated")
	}

	booking.UnbookedTime = null.TimeFrom(time.Now())
	booking.Update(database.DB)

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
	process := exec.Command(
		"sh",
		"-c",
		fmt.Sprintf(
			"cd %s; %s/%s",
			s.Path,
			s.Path,
			config.Conf.Booking.UploadSTVCommand,
		),
	)
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

	Files := strings.Split(strings.TrimSpace(string(stdoutBytes)), "\n")
	for i := 0; i < len(Files); i++ {
		Files[i] = strings.TrimSpace(Files[i])
	}

	var demos []models.Demo

	Message := "STV Demo(s) uploaded:"
	for i := 0; i < len(Files); i++ {
		Message = fmt.Sprintf("%s\n\t%s", Message, Files[i])

		// Create the demo model.
		var demo models.Demo
		demo.UploadedTime = null.TimeFrom(time.Now())
		demo.URL = Files[i]

		demos = append(demos, demo)
	}

	// Grab the current booking.
	booking, err := models.FindBooking(database.DB, s.bookingID)
	if err != nil {
		log.Println("FindBooking failed")
		return "", errors.New("Server record could not be updated")
	}

	// Add demos to booking.
	for i := 0; i < len(demos); i++ {
		booking.AddDemos(database.DB, true, &demos[i])
	}

	// Update booking.
	err = booking.Update(database.DB)
	if err != nil {
		log.Println("Update failed")
		return "", errors.New("Server record could not be updated")
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
