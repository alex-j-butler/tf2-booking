package database

import "github.com/jinzhu/gorm"

type AuthSecret struct {
	gorm.Model
	Secret    string
	DiscordID string
}
