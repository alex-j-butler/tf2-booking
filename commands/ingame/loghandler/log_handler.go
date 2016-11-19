package loghandler

import (
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"regexp"
	"sync"

	"alex-j-butler.com/tf2-booking/servers"
)

type LogHandler struct {
	Address    string
	Port       int
	conn       *net.UDPConn
	handlersMu sync.RWMutex
	handlers   map[interface{}][]reflect.Value
}

type UserEvent struct {
	UserID   string
	Username string
	SteamID  string
	Team     string
}

type SayEvent struct {
	UserEvent
	Message string
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

	go lh.handleConn()

	return lh, nil
}

func (lh *LogHandler) handleConn() {
	buf := make([]byte, 1024)

	for {
		n, addr, err := lh.conn.ReadFromUDP(buf)

		if err != nil {
			log.Println("LogHandler error:", err)
		}

		data := string(buf[:n])

		matches, err := lh.ParseLine(data)
		if err != nil {
			continue
		}

		// Find a server with the same IP and Port.
		server, err := servers.GetServerByAddress(servers.Servers, addr.String())
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
		lh.handle(server, &SayEvent{
			UserEvent: UserEvent{
				Username: matches[1],
				UserID:   matches[2],
				SteamID:  matches[3],
				Team:     matches[4],
			},
			Message: matches[5],
		})
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

func (lh *LogHandler) AddHandler(handler interface{}) func() {
	lh.initialise()

	eventType := lh.validateHandler(handler)

	lh.handlersMu.Lock()
	defer lh.handlersMu.Unlock()

	h := reflect.ValueOf(handler)

	lh.handlers[eventType] = append(lh.handlers[eventType], h)

	return func() {
		lh.handlersMu.Lock()
		defer lh.handlersMu.Unlock()

		handlers := lh.handlers[eventType]
		for i, v := range handlers {
			if h == v {
				lh.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
				return
			}
		}
	}
}

func (lh *LogHandler) initialise() {
	lh.handlersMu.Lock()
	if lh.handlers != nil {
		lh.handlersMu.Unlock()
		return
	}

	lh.handlers = make(map[interface{}][]reflect.Value)
	lh.handlersMu.Unlock()
}

func (lh *LogHandler) validateHandler(handler interface{}) reflect.Type {
	handlerType := reflect.TypeOf(handler)

	if handlerType.NumIn() != 3 {
		panic("Unable to add event handler, handler must be of type func(*loghandler.LogHandler, *servers.Server, *loghandler.EventType)")
	}

	if handlerType.In(0) != reflect.TypeOf(lh) {
		panic("Unable to add event handler, first argument must be of type *loghandler.LogHandler")
	}

	eventType := handlerType.In(2)

	if eventType.Kind() == reflect.Interface {
		eventType = nil
	}

	return eventType
}

func (lh *LogHandler) handle(server *servers.Server, event interface{}) {
	lh.handlersMu.RLock()
	defer lh.handlersMu.RUnlock()

	if lh.handlers == nil {
		return
	}

	handlerParameters := []reflect.Value{reflect.ValueOf(lh), reflect.ValueOf(server), reflect.ValueOf(event)}

	if handlers, ok := lh.handlers[nil]; ok {
		for _, handler := range handlers {
			go handler.Call(handlerParameters)
		}
	}

	if handlers, ok := lh.handlers[reflect.TypeOf(event)]; ok {
		for _, handler := range handlers {
			go handler.Call(handlerParameters)
		}
	}
}
