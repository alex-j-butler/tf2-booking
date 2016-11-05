package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	update "github.com/inconshreveable/go-update"
)

type StringServerMap map[string]*Server
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

func (m StringStringMap) ToServerMap(servers []*Server) StringServerMap {
	newMap := make(StringServerMap)
	for k, v := range m {
		var serv *Server

		log.Println(fmt.Sprintf("Searching for server matching '%s'", v))
		for _, s := range servers {
			log.Println(fmt.Sprintf("Trying server: %v", s))
			if s.SessionName == v {
				serv = s
				log.Println(fmt.Sprintf("Found server matching '%s': %v", v, s))
			}
		}

		log.Println(fmt.Sprintf("Setting map key '%s': %v", k, serv))
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

func DeleteState(save string) error {
	return os.Remove(save)
}

func SaveState(save string, servers []Server, users map[string]bool, userServers StringServerMap) error {
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
	log.Println("Loaded servers:", servers)
	userServers := state.UserStrings.ToServerMap(servers)

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
