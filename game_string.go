package main

import (
	"fmt"
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
	availableServers := len(pool.GetAvailableServers())

	if availableServers == 0 {
		return Session.UpdateStatus(1, GetGameString(availableServers))
	}

	return Session.UpdateStatus(0, GetGameString(availableServers))
}
