package loghandler

import (
	"errors"
	"fmt"
	"log"
	"regexp"
)

func ParseLine(data string) error {
	regex, err := regexp.Compile("\"(.+)<(\\d+)><(.+)><(Blue|Red|Unassigned|Spectator)>\" say \"(.+)\"")

	if err != nil {
		return err
	}

	matches := regex.FindAllString(data, -1)

	if len(matches) > 0 {
		log.Println(fmt.Sprintf("User %s (%s): %s", matches[0], matches[2], matches[4]))

		return nil
	}

	return errors.New("No match found")
}
