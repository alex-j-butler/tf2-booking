package config

import (
	"io/ioutil"
	"log"

	"alex-j-butler.com/tf2-booking/util"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {

	// Settings for the Discord bot
	Discord struct {
		Token          string `yaml:"token"`
		DefaultChannel string `yaml:"default_channel"`
		Debug          bool   `yaml:"debug"`

		AcceptableChannels []string `yaml:"acceptable_channels"`
		NotificationUsers  []string `yaml:"notification_users"`
	} `yaml:"discord"`

	// Settings for the UDP log handling server
	LogServer struct {
		LogAddress       string `yaml:"log_address"`
		LogAddressRemote string `yaml:"log_address_remote"`
		LogPort          int    `yaml:"log_port"`
	} `yaml:"log_server"`

	// Settings for the bookings
	Booking struct {
		IdleWarningDuration int `yaml:"idle_warning_duration"`

		KickMessage string `yaml:"kick_message"`

		// Settings used for Bash server runner.
		SetupCommand     string `yaml:"setup_command"`
		StartCommand     string `yaml:"start_command"`
		StopCommand      string `yaml:"stop_command"`
		UploadSTVCommand string `yaml:"upload_stv_command"`

		// Settings used for the server API server runner.
		Tag        string `yaml:"tag"`
		APIAddress string `yaml:"api_address"`
		APIPort    int    `yaml:"api_port"`
		APIKey     string `yaml:"api_key"`

		MaxIdleMinutes int `yaml:"max_idle_minutes"`
		MinPlayers     int `yaml:"min_players"`

		ErrorThreshold int `yaml:"error_threshold"`
	} `yaml:"booking"`

	Commands struct {
		DemoLink string `yaml:"demo_link"`

		ReportDuration util.DurationUtil `yaml:"report_duration"`
	}

	Database struct {
		DSN string `yaml:"dsn"`
	} `yaml:"database"`

	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	}

	Tips []string `yaml:"tips"`
}

var Conf Config

func InitialiseConfiguration() {
	configuration, _ := ioutil.ReadFile("./config.yml")
	err := yaml.Unmarshal(configuration, &Conf)

	if err != nil {
		log.Println("Failed to initialise configuration:", err)
	}
}
