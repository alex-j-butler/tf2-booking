package loghandler

import (
	"fmt"
	"log"
	"net"
)

type LogHandler struct {
	Address string
	Port    int
	conn    *net.UDPConn
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
	}
}
