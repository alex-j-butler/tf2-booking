package servers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"

	null "gopkg.in/nullbio/null.v6"

	"alex-j-butler.com/tf2-booking/config"
	"alex-j-butler.com/tf2-booking/models"
)

// ScriptServerRunner is an implementation of ServerRunner that
// runs servers using a modified LinuxGSM script file.
type ScriptServerRunner struct {
}

// Setup the server by running the 'book server' bash script.
func (s ScriptServerRunner) Setup(server *Server) (rconPassword string, srvPassword string, err error) {
	// Retrieve the RCON password & server password.
	process := exec.Command(
		"sh",
		"-c",
		fmt.Sprintf(
			"cd %s; %s/%s",
			server.Path,
			server.Path,
			config.Conf.Booking.SetupCommand,
		),
	)

	// Create pipes for stdout & stderr.
	stdout, _ := process.StdoutPipe()
	stderr, _ := process.StderrPipe()

	// Start the process.
	err = process.Start()

	if err != nil {
		log.Println("Failed to setup server:", err)
		return "", "", errors.New("Your server could not be setup")
	}

	// Read stdout & stderr.
	stdoutBytes, _ := ioutil.ReadAll(stdout)
	stderrBytes, _ := ioutil.ReadAll(stderr)

	// Wait for the process to complete.
	err = process.Wait()

	if err != nil {
		log.Println("Failed to setup server:", err)
		return "", "", errors.New("Your server could not be setup")
	}

	// Trim passwords.
	rconPassword = strings.TrimSpace(string(stdoutBytes))
	srvPassword = strings.TrimSpace(string(stderrBytes))

	return
}

func (s ScriptServerRunner) Start(server *Server) (err error) {
	process := exec.Command(
		"sh",
		"-c",
		fmt.Sprintf(
			"cd %s; %s/%s",
			server.Path,
			server.Path,
			config.Conf.Booking.StartCommand,
		),
	)

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

	return
}

func (s ScriptServerRunner) Stop(server *Server) (err error) {
	process := exec.Command(
		"sh",
		"-c",
		fmt.Sprintf(
			"cd %s; %s/%s",
			server.Path,
			server.Path,
			config.Conf.Booking.StopCommand,
		),
	)

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

	return
}

func (s ScriptServerRunner) UploadSTV(server *Server) (demos []models.Demo, err error) {
	// Run upload STV demo script.
	process := exec.Command(
		"sh",
		"-c",
		fmt.Sprintf(
			"cd %s; %s/%s %s",
			server.Path,
			server.Path,
			config.Conf.Booking.UploadSTVCommand,
			server.Booker,
		),
	)
	stdout, _ := process.StdoutPipe()

	err = process.Start()

	if err != nil {
		log.Println("Failed to upload STV:", err)
		return []models.Demo{}, errors.New("Your server failed to upload STV")
	}

	stdoutBytes, _ := ioutil.ReadAll(stdout)

	err = process.Wait()

	if err != nil {
		log.Println("Failed to upload STV:", err)
		return []models.Demo{}, errors.New("Your server failed to upload STV")
	}

	Files := strings.Split(strings.TrimSpace(string(stdoutBytes)), "\n")
	for i := 0; i < len(Files); i++ {
		var demo models.Demo
		demo.UploadedTime = null.TimeFrom(time.Now())
		demo.URL = Files[i]
		demos = append(demos, demo)

		Files[i] = strings.TrimSpace(Files[i])
	}

	return
}

func (s ScriptServerRunner) SendCommand(server *Server, command string) (err error) {
	process := exec.Command("tmux", "send-keys", "-t", server.GetRedisName(), "C-m", command, "C-m")

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

	return
}

func (s ScriptServerRunner) Console(server *Server, lines int) (consoleLines []string, err error) {
	return []string{"WARNING! Implement Console in ScriptServerRunner!"}, nil
}

func (s ScriptServerRunner) IsAvailable(server *Server) bool {
	return true
}

func (s ScriptServerRunner) IsBooked(server *Server) bool {
	return true
}
