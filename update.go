package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"alex-j-butler.com/tf2-booking/servers"

	update "github.com/inconshreveable/go-update"
)

type StringServerMap map[string]*servers.Server
type StringStringMap map[string]string

func (m StringServerMap) ToStringMap() StringStringMap {
	newMap := make(StringStringMap)
	for k, v := range m {
		if v != nil {
			newMap[k] = v.SessionName
		}
	}

	return newMap
}

func (m StringStringMap) ToServerMap(serverMap []*servers.Server) StringServerMap {
	newMap := make(StringServerMap)
	for k, v := range m {
		var serv *servers.Server
		for _, s := range serverMap {
			if s.SessionName == v {
				serv = s
			}
		}

		newMap[k] = serv
	}

	return newMap
}

type State struct {
	Servers     []servers.Server
	Users       map[string]bool
	UserStrings StringStringMap
}

func HasState(save string) bool {
	_, err := os.Stat(save)
	return err == nil
}

func DeleteState(save string) error {
	return os.Remove(save)
}

func SaveState(save string, servers []servers.Server, users map[string]bool, userServers StringServerMap) error {
	state := State{
		Servers:     servers,
		Users:       users,
		UserStrings: userServers.ToStringMap(),
	}

	j, err := json.Marshal(state)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(save, j, 0644)

	return err
}

func LoadState(save string) (error, []servers.Server, map[string]bool, map[string]*servers.Server) {
	j, err := ioutil.ReadFile(save)

	if err != nil {
		return err, nil, nil, nil
	}

	state := State{}
	err = json.Unmarshal(j, &state)

	servers_ := make([]*servers.Server, len(servers.Servers))
	for i := 0; i < len(servers.Servers); i++ {
		serv := servers.Servers[i]
		servers_[i] = &serv
	}
	userServers := state.UserStrings.ToServerMap(servers_)

	return err, state.Servers, state.Users, userServers
}

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
