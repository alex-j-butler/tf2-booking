package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func HasState(save string) bool {
	_, err := os.Stat(save)
	return err == nil
}

func SaveState(save string) error {
	j, err := json.Marshal(Conf.Servers)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(save, j, 0644)

	return err
}

func LoadState(save string) error {
	j, err := ioutil.ReadFile(save)

	if err != nil {
		return err
	}

	Conf.Servers = []Server{}
	err = json.Unmarshal(j, &Conf.Servers)

	return err
}
