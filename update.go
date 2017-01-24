package main

import (
	"net/http"

	update "github.com/inconshreveable/go-update"
)

func UpdateExecutable(address string) error {
	resp, err := http.Get(address)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		return err
	}

	return nil
}
