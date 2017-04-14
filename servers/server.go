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
	"alex-j-butler.com/tf2-booking/models"
	"github.com/bwmarrin/discordgo"
	"github.com/james4k/rcon"
	"github.com/vattle/sqlboiler/queries/qm"
	"gopkg.in/nullbio/null.v6"
	redis "gopkg.in/redis.v5"
)

// DefaultRunner is the default server runner.
var DefaultRunner = ScriptServerRunner{}

type Server struct {
	// Interface implementation to use to run this server.
	Runner ServerRunner `json:"-"`

	Name        string `json:"name" yaml:"name"`
	Path        string `json:"path" yaml:"path"`
	Address     string `json:"address" yaml:"address"`
	STVAddress  string `json:"stv_address" yaml:"stv_address"`
	SessionName string `json:"session_name" yaml:"session_name"`

	// Whether this server has been sent the unbooking warning.
	SentWarning bool

	// Whether this server has been sent the idle unbooking warning.
	SentIdleWarning bool

	// Whether this server has been sent the TF2Center/TF2Stadium lobby warning message,
	// informing them that they need 2 players on the server to prevent idle unbooking.
	SentLobbyWarning bool

	// Timestamp indicating when the server is to be returned.
	ReturnDate time.Time

	// Last known RCON password.
	// If this RCON password is invalid, the server can send a tmux command to reset it.
	RCONPassword string

	// Average tick rate reported by the server.
	TickRate float32

	// Number of tick rate measurements (used internally for calculating a new average).
	TickRateMeasurements int

	// Time that the next performance warning will occur.
	NextPerformanceWarning time.Time

	// Specifies whether the server is currently booked.
	Booked bool

	// Specifies when the server was booked.
	BookedDate time.Time

	// The ID of the Discord user who booked the server.
	Booker string

	// The mention string of the Discord user who booked the server.
	BookerMention string

	// Booking ID that the server is currently associated with.
	BookingID int

	// IdleMinutes is the number of minutes the server has been idle for.
	IdleMinutes int

	// ErrorMinutes is the number of minutes the server has been in an errored state for.
	ErrorMinutes int
}

// Init sets the default runner if it hasn't already been set.
func (s *Server) Init() {
	if s.Runner == nil {
		s.Runner = DefaultRunner
	}
}

// Update performs an update of the server into the specified Redis client.
func (s *Server) Update(redisClient *redis.Client) error {
	// Serialise the server as JSON.
	serialised, err := json.Marshal(s)
	if err != nil {
		return err
	}

	// Perform a SET command on the Redis client.
	err = redisClient.Set(fmt.Sprintf("server.%s", s.SessionName), serialised, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// Synchronise performs a synchronise of the server, retrieving the server data from the specified Redis client.
func (s *Server) Synchronise(redisClient *redis.Client) error {
	result, err := redisClient.Get(fmt.Sprintf("server.%s", s.SessionName)).Result()
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

// IsAvailable returns whether the server is currently available for booking.
// Currently a server is available for booking if it is not being booked by another user,
// in the future, this could be extended to block servers from being used (for example, if they are down).
func (s *Server) IsAvailable() bool {
	return !s.Booked
}

// BEGIN Deprecated

func (s *Server) GetBookedTime() time.Time {
	return s.BookedDate
}

func (s *Server) GetBooker() string {
	return s.Booker
}

func (s *Server) GetBookerMention() string {
	return s.BookerMention
}

func (s *Server) GetIdleMinutes() int {
	return s.IdleMinutes
}

func (s *Server) AddIdleMinute() {
	s.IdleMinutes++

	s.Update(globals.RedisClient)
}

func (s *Server) ResetIdleMinutes() {
	s.IdleMinutes = 0

	s.Update(globals.RedisClient)
}

// END Deprecated

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
	// Reset the warning notification so that it can be sent again.
	s.SentWarning = false

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

	// Run the stop function from the runner implementation.
	err := s.Runner.Stop(s)

	return err
}

func (s *Server) Book(user *discordgo.User, duration time.Duration) (string, string, error) {
	if s.Booked == true {
		return "", "", errors.New("Server is already booked")
	}

	// TODO: Move the database handling to after the server setup
	// in case an error occurs before then.

	// Tries to select the user by discord id,
	// if no record is found, insert a new record.
	dbUser, err := models.Users(globals.DB, qm.Where("discord_id=?", user.ID)).One()
	if err != nil {
		// Insert new record.
		var newUser models.User
		newUser.DiscordID = null.StringFrom(user.ID)
		newUser.Name = null.StringFrom(user.Username)

		err = newUser.Insert(globals.DB)

		if err != nil {
			log.Println("Database error:", err)
			return "", "", errors.New("User record could not be created")
		}

		dbUser = &newUser
	}

	// Adds a new booking to the database
	// and set the booking id.
	var booking models.Booking
	booking.SetBooker(globals.DB, false, dbUser)
	booking.ServerName = s.Name
	booking.BookedTime = null.TimeFrom(time.Now())
	err = booking.Insert(globals.DB)

	if err != nil {
		log.Println("Database error:", err)
		return "", "", errors.New("Server record could not be created")
	}

	s.BookingID = booking.BookingID

	// Update the server to Redis
	// This is deferred to make sure it happens whether the server is setup or not.
	defer s.Update(globals.RedisClient)

	// Set the server variables.
	s.SetServerVars(duration, user.ID)

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

func (s *Server) SetServerVars(duration time.Duration, userID string) {
	s.ReturnDate = time.Now().Add(duration)
	s.Booked = true
	s.BookedDate = time.Now()
	s.Booker = userID
	s.BookerMention = fmt.Sprintf("<@%s>", userID)
	s.NextPerformanceWarning = time.Now().Add(5 * time.Minute)
	s.SentWarning = false
	s.SentIdleWarning = false
	s.SentLobbyWarning = false
	s.IdleMinutes = 0
	s.ErrorMinutes = 0
}

func (s *Server) ResetServerVars() {
	s.ReturnDate = time.Time{}
	s.Booked = false
	s.BookedDate = time.Time{}
	s.Booker = ""
	s.BookerMention = ""
	s.NextPerformanceWarning = time.Time{}
	s.SentWarning = false
	s.SentIdleWarning = false
	s.SentLobbyWarning = false
	s.IdleMinutes = 0
	s.ErrorMinutes = 0
}

func (s *Server) Unbook() error {
	if s.Booked == false {
		return errors.New("Server is not booked")
	}

	booking, err := models.FindBooking(globals.DB, s.BookingID)
	if err != nil {
		return errors.New("Server record could not be updated")
	}

	booking.UnbookedTime = null.TimeFrom(time.Now())
	booking.Update(globals.DB)

	// Reset server variables.
	s.ResetServerVars()

	// Update the server in Redis.
	s.Update(globals.RedisClient)

	return nil
}

func (s *Server) ExtendBooking(amount time.Duration) {
	// Add duration to the return date.
	s.ReturnDate = s.ReturnDate.Add(amount)

	// Update the server in Redis.
	s.Update(globals.RedisClient)
}

func (s *Server) generateSTVReply(demos []models.Demo) string {
	message := "STV Demo(s) uploaded:"
	for i := 0; i < len(demos); i++ {
		message = fmt.Sprintf("%s\n\t%s", message, demos[i].URL)
	}

	return message
}

func (s *Server) UploadSTV() (string, error) {
	// Run the uploadSTV function from the runner implementation.
	demos, err := s.Runner.UploadSTV(s)

	message := s.generateSTVReply(demos)

	// Grab the current booking.
	booking, err := models.FindBooking(globals.DB, s.BookingID)
	if err != nil {
		log.Println("FindBooking failed")
		return "", errors.New("Server record could not be updated")
	}

	// Add demos to booking.
	for i := 0; i < len(demos); i++ {
		booking.AddDemos(globals.DB, true, &demos[i])
	}

	// Update booking.
	err = booking.Update(globals.DB)
	if err != nil {
		log.Println("Update failed")
		return "", errors.New("Server record could not be updated")
	}

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
