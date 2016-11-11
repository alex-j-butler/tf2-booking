package loghandler

import (
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"

	"alex-j-butler.com/tf2-booking/config"
	"alex-j-butler.com/tf2-booking/servers"
)

// type CommandCallback func(matches []string)
// Server, UserID, Username, SteamID, Team, Message
type CommandCallback func(*servers.Server, string, string, string, string, string)

type LogHandler struct {
	Address  string
	Port     int
	Callback CommandCallback
	conn     *net.UDPConn
}

func Dial(address string, port int) (*LogHandler, error) {
	lh := &LogHandler{
		Address: address,
		Port:    port,
	}

	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", lh.Address, lh.Port))
	if err != nil {
		return nil, err
	}

	lh.conn, err = net.ListenUDP("udp", serverAddr)
	if err != nil {
		return nil, err
	}

	go lh.handle()

	return lh, nil
}

func (lh LogHandler) handle() {
	buf := make([]byte, 1024)

	for {
		n, addr, err := lh.conn.ReadFromUDP(buf)
		log.Println(fmt.Sprintf("Received %d bytes from %s", n, addr))
		log.Println(fmt.Sprintf("String: %s", string(buf[:n])))

		if err != nil {
			log.Println("LogHandler error:", err)
		}

		data := string(buf[:n])

		matches, err := lh.ParseLine(data)
		if err != nil {
			continue
		}
		if lh.Callback != nil {
			// Find a server with the same IP and Port.
			server, err := servers.GetServerByAddress(config.Conf.Servers, addr.String())
			if err != nil {
				// Ignore this log line, we don't recognise the server.
				log.Println("Unrecognised server:", err)
				continue
			}

			// Notify the callback with the appropriate parameters.
			// matches[0] = Username
			// matches[1] = UserID
			// matches[2] = SteamID
			// matches[3] = Team
			// matches[4] = Message
			lh.Callback(server, matches[1], matches[0], matches[2], matches[3], matches[4])
		}
	}
}

func (lh LogHandler) ParseLine(data string) ([]string, error) {
	regex, err := regexp.Compile("\"(.+)<(\\d+)><(.+)><(Blue|Red|Unassigned|Spectator)>\" say \"(.+)\"")

	if err != nil {
		return []string{}, err
	}

	matches := regex.FindStringSubmatch(data)

	if len(matches) > 0 {
		return matches, nil
	}

	return []string{}, errors.New("No match found")
}
