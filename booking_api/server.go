package booking_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type Server struct {
	UUID string   `json:"uuid"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`

	IPAddress string `json:"ip_address"`
	Port      int    `json:"port"`
	STVPort   int    `json:"stv_port"`

	ServerPassword string `json:"server_password"`
	RCONPassword   string `json:"rcon_password"`

	Executable string   `json:"executable"`
	Options    []string `json:"options"`

	Running bool `json:"running"`
}

// Start sends a request to the booking API to start the server.
func (s *Server) Start(client *BookingClient) error {
	var buf bytes.Buffer
	j := json.NewEncoder(&buf)
	err := j.Encode(StartServerReq{UUID: s.UUID})
	if err != nil {
		return err
	}

	req, _ := http.NewRequest(
		"POST",
		client.getAPIPath("v1", "servers", "start"),
		&buf,
	)

	resp, err := client.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	var errResp ErrorResponse
	jsonDecoder := json.NewDecoder(resp.Body)
	err = jsonDecoder.Decode(&errResp)
	if err != nil {
		return err
	}
	return errors.New(errResp.Message)
}

// Stop sends a request to the booking API to stop the server.
func (s *Server) Stop(client *BookingClient) error {
	var buf bytes.Buffer
	j := json.NewEncoder(&buf)
	err := j.Encode(StopServerReq{UUID: s.UUID})
	if err != nil {
		return err
	}

	req, _ := http.NewRequest(
		"POST",
		client.getAPIPath("v1", "servers", "stop"),
		&buf,
	)

	resp, err := client.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	var errResp ErrorResponse
	jsonDecoder := json.NewDecoder(resp.Body)
	err = jsonDecoder.Decode(&errResp)
	if err != nil {
		return err
	}
	return errors.New(errResp.Message)
}

func (s *Server) SetPassword(client *BookingClient, rconPassword, srvPassword string) error {
	var buf bytes.Buffer
	j := json.NewEncoder(&buf)
	err := j.Encode(SetPasswordReq{UUID: s.UUID, RCONPassword: rconPassword, ServerPassword: srvPassword})
	if err != nil {
		return err
	}

	req, _ := http.NewRequest(
		"POST",
		client.getAPIPath("v1", "servers", "setpassword"),
		&buf,
	)

	resp, err := client.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	var errResp ErrorResponse
	jsonDecoder := json.NewDecoder(resp.Body)
	err = jsonDecoder.Decode(&errResp)
	if err != nil {
		return err
	}
	return errors.New(errResp.Message)
}

func (s *Server) SendCommand(client *BookingClient, command string) error {
	var buf bytes.Buffer
	j := json.NewEncoder(&buf)
	err := j.Encode(SendCommandReq{UUID: s.UUID, Command: command})
	if err != nil {
		return err
	}

	req, _ := http.NewRequest(
		"POST",
		client.getAPIPath("v1", "servers", "sendcommand"),
		&buf,
	)

	resp, err := client.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	var errResp ErrorResponse
	jsonDecoder := json.NewDecoder(resp.Body)
	err = jsonDecoder.Decode(&errResp)
	if err != nil {
		return err
	}
	return errors.New(errResp.Message)
}

func (s *Server) Update(client *BookingClient) (bool, error) {
	var buf bytes.Buffer
	j := json.NewEncoder(&buf)
	err := j.Encode(UpdateReq{UUID: s.UUID})
	if err != nil {
		return false, err
	}

	req, _ := http.NewRequest(
		"POST",
		client.getAPIPath("v1", "servers", "update"),
		&buf,
	)

	resp, err := client.client.Do(req)
	if err != nil {
		return false, err
	}

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	if resp.StatusCode == http.StatusAccepted {
		return false, nil
	}

	var errResp ErrorResponse
	jsonDecoder := json.NewDecoder(resp.Body)
	err = jsonDecoder.Decode(&errResp)
	if err != nil {
		return false, err
	}
	return false, errors.New(errResp.Message)
}

func (s *Server) Console(client *BookingClient) ([]string, error) {
	req, _ := http.NewRequest(
		"GET",
		client.getAPIPath("v1", "servers", "console"),
		nil,
	)

	client.createQuery(req, map[string]string{
		"uuid": s.UUID,
	})

	resp, err := client.client.Do(req)
	if err != nil {
		return []string{}, err
	}

	if resp.StatusCode == http.StatusOK {
		// Decode the response from JSON to a struct.
		var consoleResp ConsoleServerResponse
		j := json.NewDecoder(resp.Body)
		err := j.Decode(&consoleResp)
		if err != nil {
			return []string{}, err
		}

		return consoleResp.ConsoleLines, err
	} else if resp.StatusCode == http.StatusNoContent {
		// Return no console lines.
		return []string{}, nil
	} else {
		var errResp ErrorResponse
		j := json.NewDecoder(resp.Body)
		err := j.Decode(&errResp)
		if err != nil {
			return []string{}, err
		}
		return []string{}, errors.New(errResp.Message)
	}
}
