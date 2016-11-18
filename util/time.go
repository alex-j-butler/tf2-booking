package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"
)

// DurationUtil is a wrapper object around time.Duration
// allowing it to be unmarshaled from JSON.
type DurationUtil struct {
	time.Duration
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (du *DurationUtil) UnmarshalJSON(buf []byte) error {
	var str string
	b := new(bytes.Buffer)
	b.Write(buf)
	json.NewDecoder(b).Decode(&str)

	duration, err := time.ParseDuration(str)
	du.Duration = duration

	return err
}

func (du *DurationUtil) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	unmarshal(&str)

	duration, err := time.ParseDuration(str)
	du.Duration = duration

	return err
}

// ToHuman converts a duration into a human-readable string with support for only hours and minutes.
// Prints in the format '1 hour 45 minutes'
func ToHuman(duration *time.Duration) string {
	var str string

	hours := int(math.Floor(duration.Hours()))
	minutes := int(math.Floor(duration.Minutes()) - float64(hours*60))

	hourString := "%d hours"
	minuteString := "%d minutes"

	if hours == 1 {
		hourString = "%d hour"
	}
	if minutes == 1 {
		minuteString = "%d minute"
	}

	if hours > 0 {
		str = fmt.Sprintf("%s%s ", str, fmt.Sprintf(hourString, hours))
	}
	if minutes > 0 {
		str = fmt.Sprintf("%s%s ", str, fmt.Sprintf(minuteString, minutes))
	}

	return strings.TrimSpace(str)
}
