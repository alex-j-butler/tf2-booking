package servers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"alex-j-butler.com/tf2-booking/config"
	"alex-j-butler.com/tf2-booking/globals"
	"alex-j-butler.com/tf2-booking/util"
	"github.com/bwmarrin/discordgo"
	"github.com/james4k/rcon"
	redis "gopkg.in/redis.v5"
)

type Server struct {
	Runner *ServerRunner

	UUID       string `json:"-"`
	Name       string `json:"-"`
	Path       string `json:"-"`
	Address    string `json:"-"`
	STVAddress string `json:"-"`

	// Whether this server has been sent the idle unbooking warning.
	SentIdleWarning bool

	// Last known RCON password.
	// If this RCON password is invalid, the server can send a tmux command to reset it.
	RCONPassword string

	// Specifies whether the server is currently booked.
	Booked bool

	// Specifies when the server was booked.
	BookedDate time.Time

	// The ID of the Discord user who booked the server.
	Booker string

	// The mention string of the Discord user who booked the server.
	BookerMention string

	// The full name of the Discord user who booked the server.
	BookerFullname string

	// IdleMinutes is the number of minutes the server has been idle for.
	IdleMinutes int

	// ErrorMinutes is the number of minutes the server has been in an errored state for.
	ErrorMinutes int
}

func (s *Server) SetServerVars(userID string, fullname string) {
	s.Booked = true
	s.BookedDate = time.Now()
	s.Booker = userID
	s.BookerMention = fmt.Sprintf("<@%s>", userID)
	s.BookerFullname = fullname
	s.SentIdleWarning = false
	s.IdleMinutes = 0
	s.ErrorMinutes = 0
}

func (s *Server) ResetServerVars() {
	s.Booked = false
	s.BookedDate = time.Time{}
	s.Booker = ""
	s.BookerMention = ""
	s.SentIdleWarning = false
	s.IdleMinutes = 0
	s.ErrorMinutes = 0
}

// Update performs an update of the server into the specified Redis client.
func (s *Server) Update(redisClient *redis.Client) error {
	// Serialise the server as JSON.
	serialised, err := json.Marshal(s)
	if err != nil {
		log.Println("marshal error:", err)
		return err
	}

	// Perform a SET command on the Redis client.
	err = redisClient.Set(fmt.Sprintf("server.%s", s.UUID), serialised, 0).Err()
	if err != nil {
		log.Println("redis error:", err)
		return err
	}

	return nil
}

// Synchronise performs a synchronise of the server, retrieving the server data from the specified Redis client.
func (s *Server) Synchronise(redisClient *redis.Client) error {
	result, err := redisClient.Get(fmt.Sprintf("server.%s", s.UUID)).Result()
	if err != nil {
		return err
	}

	// Deserialise the JSON.
	err = json.Unmarshal([]byte(result), &s)
	if err != nil {
		return err
	}

	return nil
}

// Available returns whether the server is currently bookable,
// or whether it's experiencing an error that would prevent it from being successfully booked.
func (s *Server) Available() bool {
	return s.Runner.IsAvailable(s) && !s.Runner.IsBooked(s)
}

// IsBooked returns whether the server is currently booked
func (s *Server) IsBooked() bool {
	return s.Booked && s.Runner.IsBooked(s)
}

func (s *Server) AddIdleMinute() {
	s.IdleMinutes++

	s.Update(globals.RedisClient)
}

func (s *Server) ResetIdleMinutes() {
	s.IdleMinutes = 0

	s.Update(globals.RedisClient)
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
	// Run the setup function from the runner implementation.
	rconPassword, srvPassword, err := s.Runner.Setup(s)

	// Cache the RCON password, since it can't be changed by the user.
	s.RCONPassword = rconPassword

	// Update the server.
	s.Update(globals.RedisClient)

	return rconPassword, srvPassword, err
}

// Start the server using a bash script.
// Returns:
//  error - Error of a failed start, or nil if none
func (s *Server) Start() error {
	// Run the start function from the runner implementation.
	err := s.Runner.Start(s)

	return err
}

// Stop the server using a bash script.
// Returns:
// 	error - Error of a failed stop, or nil if none
func (s *Server) Stop() error {
	// Stop the STV recording and kick all players cleanly.
	KickCommand := fmt.Sprintf("tv_stop; kickall \"%s\"", config.Conf.Booking.KickMessage)
	s.SendCommand(KickCommand)

	// Wait 1 second to the kick command to properly kick everyone.
	time.Sleep(1 * time.Second)

	// Run the stop function from the runner implementation.
	err := s.Runner.Stop(s)

	return err
}

func (s *Server) Book(user *discordgo.User) (string, string, error) {
	patchUser := &util.PatchUser{user}

	if s.Booked == true {
		return "", "", errors.New("Server is already booked")
	}

	// Update the server to Redis
	// This is deferred to make sure it happens whether the server is setup or not.
	defer s.Update(globals.RedisClient)

	// Set the server variables.
	s.SetServerVars(user.ID, patchUser.GetFullname())

	// Setup the server.
	RCONPassword, ServerPassword, err := s.Setup()

	if err != nil {
		// Reset the server variables so that
		// the booking bot correctly unbooks the server in case of an error.
		s.ResetServerVars()

		return "", "", err
	}

	return RCONPassword, ServerPassword, err
}

func (s *Server) Unbook() error {
	if s.Booked == false {
		return errors.New("Server is not booked")
	}

	// Reset server variables.
	s.ResetServerVars()

	// Update the server in Redis.
	s.Update(globals.RedisClient)

	return nil
}

func (s *Server) ExtendBooking() {
	// TODO: Implement this.
	// Reset the number of idle minutes, and allow the timeout warning message to be sent again.
	s.SentIdleWarning = false
	s.ResetIdleMinutes()

	// Update the server in Redis.
	s.Update(globals.RedisClient)
}

func (s *Server) generateSTVReply(demos []string) string {
	message := "STV Demo(s) uploaded:"
	for i := 0; i < len(demos); i++ {
		message = fmt.Sprintf("%s\n\t%s", message, demos[i])
	}

	return message
}

func (s *Server) UploadSTV() (string, error) {
	// Run the uploadSTV function from the runner implementation.
	demos, err := s.Runner.UploadSTV(s)
	if err != nil {
		return "", err
	}
	if len(demos) == 0 {
		return "", errors.New("No demos")
	}

	message := s.generateSTVReply(demos)

	return message, nil
}

func (s *Server) SendCommand(command string) error {
	// Run the SendCommand function from the runner implementation.
	err := s.Runner.SendCommand(s, command)

	return err
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

// Console queries the server for the latest console lines.
func (s *Server) Console() ([]string, error) {
	return s.Runner.Console(s, 0)
}
