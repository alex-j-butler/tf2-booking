package main

import (
	"bytes"
	"encoding/json"
	"time"
)

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
