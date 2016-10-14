package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	DiscordToken string         `json:"discord_token"`
	Servers      []ServerConfig `json:"servers"`
}

var Conf Config

func InitialiseConfiguration() {
	configuration, _ := ioutil.ReadFile("./config.json")
	err := json.Unmarshal(configuration, &Conf)

	if err != nil {
		log.Println("Failed to initialise configuration:", err)
	}
}
