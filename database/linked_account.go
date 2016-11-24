package database

import "github.com/jinzhu/gorm"

type LinkedAccount struct {
	gorm.Model
	SteamID   string `gorm:"unique"`
	DiscordID string `gorm:"unique"`
}
