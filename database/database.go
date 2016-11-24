package database

import (
	"log"

	"alex-j-butler.com/tf2-booking/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

// Initialise connects to the database.
func Initialise() {
	db, err := gorm.Open(config.Conf.Database.Dialect, config.Conf.Database.URL)
	if err != nil {
		log.Println("Database error:", err)
	}

	db.AutoMigrate(&AuthSecret{})
	DB = db
}
