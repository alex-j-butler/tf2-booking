package config

import (
	"io/ioutil"
	"log"

	"alex-j-butler.com/tf2-booking/util"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Discord struct {
		Token          string `yaml:"token"`
		DefaultChannel string `yaml:"default_channel"`
		Debug          bool   `yaml:"debug"`

		AcceptableChannels []string `yaml:"acceptable_channels"`
		NotificationUsers  []string `yaml:"notification_users"`
	} `yaml:"discord"`

	Booking struct {
		Duration        util.DurationUtil `yaml:"duration"`
		ExtendDuration  util.DurationUtil `yaml:"extend_duration"`
		WarningDuration util.DurationUtil `yaml:"warning_duration"`

		KickMessage      string `yaml:"kick_message"`
		SetupCommand     string `yaml:"setup_command"`
		StartCommand     string `yaml:"start_command"`
		StopCommand      string `yaml:"stop_command"`
		UploadSTVCommand string `yaml:"upload_stv_command"`

		MaxIdleMinutes int `yaml:"max_idle_minutes"`
		MinPlayers     int `yaml:"min_players"`

		ErrorThreshold int `yaml:"error_threshold"`
	} `yaml:"booking"`

	Commands struct {
		ReportDuration util.DurationUtil `yaml:"report_duration"`
	}
}

var Conf Config

func InitialiseConfiguration() {
	configuration, _ := ioutil.ReadFile("./config.yml")
	err := yaml.Unmarshal(configuration, &Conf)

	if err != nil {
		log.Println("Failed to initialise configuration:", err)
	}
}
