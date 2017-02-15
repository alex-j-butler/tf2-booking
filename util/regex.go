package util

import (
	"log"
	"regexp"
)

func IsUUID4(uuid string) bool {
	success, err := regexp.Match("(?i)[a-f0-9]{8}-?[a-f0-9]{4}-?4[a-f0-9]{3}-?[89ab][a-f0-9]{3}-?[a-f0-9]{12}", []byte(uuid))
	if err != nil {
		log.Println("Regex error:", err)
		return false
	}

	return success
}
