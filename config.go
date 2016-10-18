package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"alex-j-butler.com/tf2-booking/util"
)

type Config struct {
	DiscordToken       string   `json:"discord_token"`
	DefaultChannel     string   `json:"default_channel"`
	AcceptableChannels []string `json:"acceptable_channels"`
	MaxIdleMinutes     int      `json:"max_idle_minutes"`
	MinPlayers         int      `json:"min_players"`

	ErrorThreshold    int      `json:"error_threshold"`
	NotificationUsers []string `json:"notification_users"`

	BookingDuration        util.DurationUtil `json:"booking_duration"`
	BookingExtendDuration  util.DurationUtil `json:"booking_extend_duration"`
	BookingWarningDuration util.DurationUtil `json:"booking_warning_duration"`

	BookingDurationText        string `json:"booking_duration_text"`
	BookingExtendDurationText  string `json:"booking_extend_duration_text"`
	BookingWarningDurationText string `json:"booking_warning_duration_text"`

	Servers []Server `json:"servers"`
}

var Conf Config

func InitialiseConfiguration() {
	configuration, _ := ioutil.ReadFile("./config.json")
	err := json.Unmarshal(configuration, &Conf)

	if err != nil {
		log.Println("Failed to initialise configuration:", err)
	}
}
