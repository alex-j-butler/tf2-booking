package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type StringServerMap map[string]*Server
type StringStringMap map[string]string

func (m StringServerMap) ToStringMap() StringStringMap {
	newMap := make(StringStringMap)
	for k, v := range m {
		newMap[k] = v.SessionName
	}

	return newMap
}

func (m StringStringMap) ToServerMap(servers []*Server) StringServerMap {
	newMap := make(StringServerMap)
	for k, v := range m {
		var serv *Server
		for _, s := range servers {
			if s.SessionName == v {
				serv = s
			}
		}

		newMap[k] = serv
	}

	return newMap
}

type State struct {
	Servers     []Server
	Users       map[string]bool
	UserStrings StringStringMap
}

func HasState(save string) bool {
	_, err := os.Stat(save)
	return err == nil
}

func SaveState(save string, users map[string]bool, userServers StringServerMap) error {
	state := State{
		Servers:     Conf.Servers,
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

func LoadState(save string) (error, []Server, map[string]bool, map[string]*Server) {
	j, err := ioutil.ReadFile(save)

	if err != nil {
		return err, nil, nil, nil
	}

	state := State{}
	err = json.Unmarshal(j, &state)

	servers := make([]*Server, len(Conf.Servers))
	for i, j := range state.Servers {
		servers[i] = &j
	}
	userServers := state.UserStrings.ToServerMap(servers)

	return err, state.Servers, state.Users, userServers
}
