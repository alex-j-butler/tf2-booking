package util

import (
	"regexp"
	"strconv"
)

type PlayerStatus struct {
	UserID        int
	Username      string
	SteamID       SteamID
	ConnectedTime string
	Ping          int
	Loss          int
	State         string
	IP            string
	Port          int
}

func ParsePlayerStatus(statusLine string) ([]*PlayerStatus, error) {
	r, _ := regexp.Compile(`#\s+(\d+)\s+"(.+)"\s+(\[U:\d:\d+\])\s+(.+)\s+(\d+)\s+(\d+)\s+active\s+(\d+\.\d+\.\d+\.\d+):(\d+)`)
	matches := r.FindAllStringSubmatch(statusLine, -1)

	var statuses []*PlayerStatus
	for _, match := range matches {
		userID, err := strconv.ParseInt(match[1], 10, 32)
		if err != nil {
			return nil, err
		}

		username := match[2]
		steamID := FromSteamID3(match[3])
		connectedTime := match[4]

		ping, err := strconv.ParseInt(match[5], 10, 32)
		if err != nil {
			return nil, err
		}

		loss, err := strconv.ParseInt(match[6], 10, 32)
		if err != nil {
			return nil, err
		}

		state := "active" // match[7]
		ip := match[7]

		port, err := strconv.ParseInt(match[8], 10, 32)
		if err != nil {
			return nil, err
		}

		statuses = append(statuses, &PlayerStatus{
			UserID:        int(userID),
			Username:      username,
			SteamID:       steamID,
			ConnectedTime: connectedTime,
			Ping:          int(ping),
			Loss:          int(loss),
			State:         state,
			IP:            ip,
			Port:          int(port),
		})
	}

	return statuses, nil
}
