package servers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"alex-j-butler.com/tf2-booking/config"

	"golang.org/x/crypto/ssh"
)

func publicKeyFile(file string) ssh.AuthMethod {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(data)
	if err != nil {
		return nil
	}

	return ssh.PublicKeys(key)
}

// connectSSH creates a client connection to an SSH server.
func (s *Server) connectSSH() (*ssh.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User: s.SSHUsername,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.SSHPassword),
			publicKeyFile(s.SSHPrivateKey),
		},
	}

	connection, err := ssh.Dial("tcp", s.SSHAddress, sshConfig)
	if err != nil {
		return nil, err
	}

	return connection, nil
}

// sendSSH opens a new SSH session and runs the specified command.
func (s *Server) sendSSH(conn *ssh.Client, command string) (string, string, error) {
	// New SSH session.
	sess, err := conn.NewSession()
	if err != nil {
		return "", "", err
	}

	// Request pseudo-terminal.
	if err = sess.RequestPty("xterm", 80, 40, ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}); err != nil {
		sess.Close()
		return "", "", err
	}

	stdout, _ := sess.StdoutPipe()
	stderr, _ := sess.StderrPipe()

	err = sess.Run(command)
	if err != nil {
		return "", "", err
	}

	stdoutBytes, _ := ioutil.ReadAll(stdout)
	stderrBytes, _ := ioutil.ReadAll(stderr)

	return string(stdoutBytes), string(stderrBytes), nil
}

// Internal function that connects to the SSH server and executes the command
// to setup the TF2 server.
func (s *Server) setupSSHServer() (string, string, error) {
	// Connect to the SSH server.
	conn, err := s.connectSSH()

	if err != nil {
		return "", "", errors.New("Your server could not be setup")
	}

	// Send remote command.
	stdout, _, err := s.sendSSH(conn, fmt.Sprintf(
		"cd %s; %s/%s",
		s.Path,
		s.Path,
		config.Conf.Booking.SetupCommand,
	))

	if err != nil {
		log.Println("Failed to setup server:", err)
		return "", "", errors.New("Your server could not be setup")
	}

	slices := strings.Split(stdout, "\n")

	s.SentWarning = false

	// Trim passwords.
	RCONPassword := strings.TrimSpace(slices[0])
	ServerPassword := strings.TrimSpace(slices[1])

	s.RCONPassword = RCONPassword

	return RCONPassword, ServerPassword, nil
}

func (s *Server) startRemoteServer() error {
	// Connect to the SSH server.
	conn, err := s.connectSSH()

	if err != nil {
		return errors.New("Your server could not be started")
	}

	// Send remote command.
	_, _, err = s.sendSSH(conn,
		fmt.Sprintf(
			"cd %s; %s/%s",
			s.Path,
			s.Path,
			config.Conf.Booking.StartCommand,
		),
	)

	if err != nil {
		return errors.New("Your server could not be started")
	}

	return nil
}

func (s *Server) stopRemoteServer() error {
	// Connect to the SSH server.
	conn, err := s.connectSSH()

	if err != nil {
		return errors.New("Your server could not be stopped")
	}

	// Stop the STV recording and kick all players cleanly.
	KickCommand := fmt.Sprintf("tv_stop; kickall \"%s\"", config.Conf.Booking.KickMessage)
	s.SendCommand(KickCommand)

	// Send remote command.
	_, _, err = s.sendSSH(conn,
		fmt.Sprintf(
			"cd %s; %s/%s",
			s.Path,
			s.Path,
			config.Conf.Booking.StopCommand,
		),
	)

	if err != nil {
		return errors.New("Your server could not be stopped")
	}

	return nil
}

func (s *Server) uploadSTVRemoteServer() (string, error) {
	// Connect to the SSH server.
	conn, err := s.connectSSH()

	if err != nil {
		return "", errors.New("Your server failed to upload STV")
	}

	// Run upload STV demo script.
	// Send remote command.
	stdout, _, err := s.sendSSH(conn,
		fmt.Sprintf(
			"cd %s; %s/%s",
			s.Path,
			s.Path,
			config.Conf.Booking.UploadSTVCommand,
		),
	)

	if err != nil {
		log.Println("Failed to upload STV:", err)
		return "", errors.New("Your server failed to upload STV")
	}

	Files := strings.Split(stdout, "\n")
	for i := 0; i < len(Files); i++ {
		Files[i] = strings.TrimSpace(Files[i])
	}

	Message := "STV Demo(s) uploaded:"
	for i := 0; i < len(Files); i++ {
		Message = fmt.Sprintf("%s\n\t%s", Message, Files[i])
	}

	return Message, nil
}

func (s *Server) sendCommandRemoteServer(command string) error {
	// Connect to the SSH server.
	conn, err := s.connectSSH()

	if err != nil {
		return errors.New("Your server failed to respond to commands")
	}

	// Send remote command.
	_, _, err = s.sendSSH(conn,
		fmt.Sprintf(
			"tmux send-keys -t %s C-m %s C-m",
			s.SessionName,
			command,
		),
	)

	if err != nil {
		log.Println("Failed to send command:", err)
		return errors.New("Your server failed to respond to commands")
	}

	return nil
}
