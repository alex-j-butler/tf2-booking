package util

import (
	"bytes"
	"encoding/json"
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
