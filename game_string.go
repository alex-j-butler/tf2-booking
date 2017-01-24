package main

import (
	"fmt"

	"alex-j-butler.com/tf2-booking/servers"
)

func GetGameString(num int) string {
	if num == 0 {
		return "No servers available"
	} else if num == 1 {
		return "1 server available"
	} else {
		return fmt.Sprintf("%d servers available", num)
	}
}

func UpdateGameString() error {
	availableServers := len(servers.GetAvailableServers(servers.Servers))

	if availableServers == 0 {
		return Session.UpdateStatus(1, GetGameString(availableServers))
	} else {
		return Session.UpdateStatus(0, GetGameString(availableServers))
	}
}
