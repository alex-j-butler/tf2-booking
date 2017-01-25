package main

import (
	"math/rand"
	"time"

	"alex-j-butler.com/tf2-booking/config"
)

func ChooseRandomTip() string {
	// Set the random seed (this doesn't need to be secure since we're just using it for a tip message).
	rand.Seed(time.Now().UTC().UnixNano())

	return config.Conf.Tips[rand.Intn(len(config.Conf.Tips))]
}
