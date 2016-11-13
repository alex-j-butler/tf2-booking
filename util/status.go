package util

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	STATUS_LINES   = "(?m:^(.+?)\\s*: (.+)\\s*)"
	STATUS_PLAYERS = "(?m:^#\\s+(\\d+)(?:\\s+)*\\s+\"(.+?)\"\\s+(\\[U:1:\\d+\\])\\s+(.+?)\\s+(\\d+)\\s+(\\d+)\\s+active\\s+([0-9,.]+):(\\d+))"
)

type Status struct {
	Users []User
	Lines map[string]string
}

type User struct {
	UserID int
	ID     SteamID
	Name   string
	Ping   int
	Loss   int
	IP     string
	Port   int
}

func ParseStatus(statusLine string) (*Status, error) {
	r, _ := regexp.Compile(STATUS_PLAYERS)
	matches := r.FindAllStringSubmatch(statusLine, -1)

	users := make([]User, 0, len(matches))
	for _, match := range matches {
		userID, _ := strconv.ParseInt(match[1], 10, 32)
		ping, _ := strconv.ParseInt(match[5], 10, 32)
		loss, _ := strconv.ParseInt(match[6], 10, 32)
		port, _ := strconv.ParseInt(match[8], 10, 32)

		users = append(users, User{
			UserID: int(userID),
			ID:     FromSteamID3(match[3]),
			Name:   match[2],
			Ping:   int(ping),
			Loss:   int(loss),
			IP:     match[7],
			Port:   int(port),
		})
	}

	r, _ = regexp.Compile(STATUS_LINES)
	matches = r.FindAllStringSubmatch(statusLine, -1)

	lines := make(map[string]string)
	for _, match := range matches {
		lines[strings.TrimSpace(match[1])] = strings.TrimSpace(match[2])
	}

	return &Status{
		Users: users,
		Lines: lines,
	}, nil
}
