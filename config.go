package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	DiscordToken       string   `json:"discord_token"`
	DefaultChannel     string   `json:"default_channel"`
	AcceptableChannels []string `json:"acceptable_channels"`
	MaxIdleMinutes     int      `json:"max_idle_minutes"`
	MinPlayers         int      `json:"min_players"`
	Servers            []Server `json:"servers"`
}

var Conf Config

func InitialiseConfiguration() {
	configuration, _ := ioutil.ReadFile("./config.json")
	err := json.Unmarshal(configuration, &Conf)

	if err != nil {
		log.Println("Failed to initialise configuration:", err)
	}
}
