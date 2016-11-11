package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type SteamID struct {
	SteamID     string
	SteamID3    string
	CommunityID string
}

func steamIDToCommunityID(steamID string) string {
	r, _ := regexp.Compile("^STEAM_(\\d+):(\\d+):(\\d+)$")
	matches := r.FindStringSubmatch(steamID)
	y, _ := strconv.ParseInt(matches[2], 10, 32)
	z, _ := strconv.ParseInt(matches[3], 10, 32)

	communityID := z*2 + 0x0110000100000000 + y
	return fmt.Sprintf("%d", communityID)
}

func communityIDToSteamID(communityID string) string {
	communityIDint, _ := strconv.ParseInt(communityID, 10, 64)

	x := 0
	y := communityIDint & 1
	z := communityIDint - y
	z = z - 0x0110000100000000
	z = z / 2

	return fmt.Sprintf("STEAM_%d:%d:%d", x, y, z)
}

func steamIDToSteamID3(steamID string) string {
	r, _ := regexp.Compile("^STEAM_(\\d+):(\\d+):(\\d+)$")
	matches := r.FindStringSubmatch(steamID)
	y, _ := strconv.ParseInt(matches[2], 10, 32)
	z, _ := strconv.ParseInt(matches[3], 10, 32)

	return fmt.Sprintf("[U:1:%d]", z*2+y)
}

func steamID3ToCommunityID(steamID3 string) string {
	args := strings.Split(steamID3, ":")
	accountID, _ := strconv.ParseInt(strings.Trim(args[2], "]"), 10, 64)

	var y int64
	var z int64

	if accountID%2 == 0 {
		y = 0
		z = accountID / 2
	} else {
		y = 1
		z = (accountID - 1) / 2
	}

	return fmt.Sprintf("7656119%d", (z*2)+(7960265728+y))
}

func FromSteamID(steamid string) SteamID {
	return SteamID{
		SteamID:     steamid,
		SteamID3:    steamIDToSteamID3(steamid),
		CommunityID: steamIDToCommunityID(steamid),
	}
}

func FromSteamID3(steamid3 string) SteamID {
	return SteamID{
		SteamID:     communityIDToSteamID(steamID3ToCommunityID(steamid3)),
		SteamID3:    steamid3,
		CommunityID: steamID3ToCommunityID(steamid3),
	}
}

func FromCommunityID(communityid string) SteamID {
	return SteamID{
		SteamID:     communityIDToSteamID(communityid),
		SteamID3:    steamIDToSteamID3(communityIDToSteamID(communityid)),
		CommunityID: communityid,
	}
}
