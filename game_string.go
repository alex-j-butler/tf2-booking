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

func UpdateGameString() {
	availableServers := len(GetAvailableServers())

	if availableServers == 0 {
		Session.UpdateStatus(1, GetGameString(availableServers))
	} else {
		Session.UpdateStatus(0, GetGameString(availableServers))
	}
}
