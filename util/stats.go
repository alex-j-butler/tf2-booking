package util

import (
	"regexp"
	"strconv"
	"strings"
)

type Stats struct {
	CPU        float32
	InKBs      float32
	OutKBs     float32
	Uptime     int
	MapChanges int
	FPS        float32
	Players    int
	Connects   int
}

func ParseStats(statsLine string) (*Stats, error) {
	line := strings.Split(statsLine, "\n")[1]
	r, _ := regexp.Compile("([0-9\\.]+)")
	matches := r.FindAllString(line, -1)

	cpu, err := strconv.ParseFloat(matches[0], 32)
	if err != nil {
		return nil, err
	}

	inkbs, err := strconv.ParseFloat(matches[1], 32)
	if err != nil {
		return nil, err
	}

	outkbs, err := strconv.ParseFloat(matches[2], 32)
	if err != nil {
		return nil, err
	}

	uptime, err := strconv.ParseInt(matches[3], 10, 32)
	if err != nil {
		return nil, err
	}

	mapchanges, err := strconv.ParseInt(matches[4], 10, 32)
	if err != nil {
		return nil, err
	}

	fps, err := strconv.ParseFloat(matches[5], 32)
	if err != nil {
		return nil, err
	}

	players, err := strconv.ParseInt(matches[6], 10, 32)
	if err != nil {
		return nil, err
	}

	connects, err := strconv.ParseInt(matches[7], 10, 32)
	if err != nil {
		return nil, err
	}

	return &Stats{
		CPU:        float32(cpu),
		InKBs:      float32(inkbs),
		OutKBs:     float32(outkbs),
		Uptime:     int(uptime),
		MapChanges: int(mapchanges),
		FPS:        float32(fps),
		Players:    int(players),
		Connects:   int(connects),
	}, nil
}
