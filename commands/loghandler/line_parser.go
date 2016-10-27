package loghandler

import (
	"errors"
	"log"
	"regexp"
)

func ParseLine(data string) (string, error) {
	regex, err := regexp.Compile("\"(.+)<(\\d+)><(.+)><(Blue|Red|Unassigned|Spectator)>\" say \"(.+)\"")

	if err != nil {
		return "", err
	}

	matches := regex.FindAllString(data, -1)

	if len(matches) > 0 {
		for _, match := range matches {
			log.Println(match)
		}
		// log.Println(fmt.Sprintf("User %s (%s): %s", matches[0], matches[2], matches[4]))

		return "", nil
	}

	return "", errors.New("No match found")
}
