package loghandler

import (
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"
)

type CommandCallback func(matches []string)

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
			log.Println("ParseLine error:", err)
			return
		}
		if lh.Callback != nil {
			lh.Callback(matches)
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
		for _, match := range matches {
			log.Println(match)
		}
		// log.Println(fmt.Sprintf("User %s (%s): %s", matches[0], matches[2], matches[4]))

		return matches, nil
	}

	return []string{}, errors.New("No match found")
}
