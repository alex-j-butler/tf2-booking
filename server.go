package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Server struct {
	Name    string
	Path    string
	Address string

	booked     bool
	bookedDate time.Time

	booker        string
	bookerMention string

	host string
	port int
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

	// Trim passwords.
	RCONPassword := strings.TrimSpace(string(stdoutBytes))
	ServerPassword := strings.TrimSpace(string(stderrBytes))
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

func (s *Server) Book(user *discordgo.User) (string, string, error) {
	if s.booked == true {
		return "", "", errors.New("Server is already booked")
	}

	// Set the server variables.
	s.booked = true
	s.bookedDate = time.Now()
	s.booker = user.ID
	s.bookerMention = fmt.Sprintf("<@%s>", user.ID)

	var err error

	// Setup the server.
	RCONPassword, ServerPassword, err := s.Setup()

	if err != nil {
		return "", "", err
	}

	// Start the server.
	err = s.Start()

	return RCONPassword, ServerPassword, err
}

func (s *Server) Unbook() error {
	if s.booked == false {
		return errors.New("Server is not booked")
	}

	// Set the server variables.
	s.booked = false
	s.bookedDate = time.Time{}
	s.booker = ""
	s.bookerMention = ""

	// Stop the server.
	err := s.Stop()

	return err
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
