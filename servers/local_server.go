package servers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"alex-j-butler.com/tf2-booking/config"
)

func (s *Server) setupLocalServer() (string, string, error) {
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

func (s *Server) startLocalServer() error {
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

func (s *Server) stopLocalServer() error {
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

func (s *Server) uploadSTVLocalServer() (string, error) {
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

func (s *Server) sendCommandLocalServer(command string) error {
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
