// TODO: Rewrite the booking API client & test it.
// This should also be renamed to the server API client.
package booking_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
)

type BookingClient struct {
	baseURL string
	client  *http.Client
}

func New(baseURL string) *BookingClient {
	return &BookingClient{
		baseURL: baseURL,
		client:  http.DefaultClient,
	}
}

func (c BookingClient) getAPIPath(version, subname, endpoint string) string {
	return fmt.Sprintf("%s/%s", c.baseURL, path.Join("api", version, subname, endpoint))
}

func (c BookingClient) createQuery(req *http.Request, queryParams map[string]string) {
	q := req.URL.Query()
	for k, v := range queryParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
}

func (c BookingClient) GetServers() ([]Server, error) {
	req, _ := http.NewRequest(
		"GET",
		c.getAPIPath("v1", "servers", "listall"),
		nil,
	)

	resp, err := c.client.Do(req)
	if err != nil {
		return []Server{}, err
	}

	if resp.StatusCode == http.StatusOK {
		// Decode the response from JSON to a struct.
		var listAllResp ListAllResponse
		j := json.NewDecoder(resp.Body)
		err := j.Decode(&listAllResp)
		if err != nil {
			return []Server{}, err
		}

		return listAllResp.Servers, err
	} else if resp.StatusCode == http.StatusNoContent {
		// Return no servers.
		return []Server{}, nil
	} else {
		var errResp ErrorResponse
		j := json.NewDecoder(resp.Body)
		err := j.Decode(&errResp)
		if err != nil {
			return []Server{}, err
		}
		return []Server{}, errors.New(errResp.Message)
	}
}

func (c BookingClient) StartServer(name string) error {
	var buf bytes.Buffer
	j := json.NewEncoder(&buf)
	err := j.Encode(StartServerReq{Name: name})
	if err != nil {
		return err
	}

	req, _ := http.NewRequest(
		"POST",
		c.getAPIPath("v1", "servers", "start"),
		&buf,
	)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil
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

func (c BookingClient) StopServer(name string) error {
	var buf bytes.Buffer
	j := json.NewEncoder(&buf)
	err := j.Encode(StopServerReq{Name: name})
	if err != nil {
		return err
	}

	req, _ := http.NewRequest(
		"POST",
		c.getAPIPath("v1", "servers", "stop"),
		&buf,
	)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil
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

func (c BookingClient) SetPassword(name string, rconPassword string, srvPassword string) error {
	var buf bytes.Buffer
	j := json.NewEncoder(&buf)
	err := j.Encode(SetPasswordReq{Name: name, RCONPassword: rconPassword, ServerPassword: srvPassword})
	if err != nil {
		return err
	}

	req, _ := http.NewRequest(
		"POST",
		c.getAPIPath("v1", "servers", "setpassword"),
		&buf,
	)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil
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

func (c BookingClient) SendCommand(name string, command string) error {
	var buf bytes.Buffer
	j := json.NewEncoder(&buf)
	err := j.Encode(SendCommandReq{Name: name, Command: command})
	if err != nil {
		return err
	}

	req, _ := http.NewRequest(
		"POST",
		c.getAPIPath("v1", "servers", "sendcommand"),
		&buf,
	)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil
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

func (c BookingClient) GetServersByTag(tag string) ([]Server, error) {
	req, _ := http.NewRequest(
		"GET",
		c.getAPIPath("v1", "servers", "listall"),
		nil,
	)

	c.createQuery(req, map[string]string{
		"tag": tag,
	})

	resp, err := c.client.Do(req)
	if err != nil {
		return []Server{}, err
	}

	if resp.StatusCode == http.StatusOK {
		// Decode the response from JSON to a struct.
		var listAllResp ListAllResponse
		j := json.NewDecoder(resp.Body)
		err := j.Decode(&listAllResp)
		if err != nil {
			return []Server{}, err
		}

		return listAllResp.Servers, err
	} else if resp.StatusCode == http.StatusNoContent {
		// Return no servers.
		return []Server{}, nil
	} else {
		var errResp ErrorResponse
		j := json.NewDecoder(resp.Body)
		err := j.Decode(&errResp)
		if err != nil {
			return []Server{}, err
		}
		return []Server{}, errors.New(errResp.Message)
	}
}

func (c BookingClient) NextServer(tag string) (Server, error) {
	req, _ := http.NewRequest(
		"GET",
		c.getAPIPath("v1", "servers", "next"),
		nil,
	)

	c.createQuery(req, map[string]string{
		"tag": tag,
	})

	resp, err := c.client.Do(req)
	if err != nil {
		return Server{}, err
	}

	// Server was found.
	if resp.StatusCode == 200 {
		// Decode the response from JSON to a struct.
		var i Server
		j := json.NewDecoder(resp.Body)
		err := j.Decode(&i)
		if err != nil {
			return Server{}, err
		}

		return i, nil
	} else if resp.StatusCode == 204 {
		return Server{}, errors.New("No servers available")
	} else {
		var errResp ErrorResponse
		j := json.NewDecoder(resp.Body)
		err := j.Decode(&errResp)
		if err != nil {
			return Server{}, err
		}
		return Server{}, errors.New(errResp.Message)
	}
}

func (c BookingClient) GetServer(name string) (Server, error) {
	req, _ := http.NewRequest(
		"GET",
		c.getAPIPath("v1", "servers", "list"),
		nil,
	)

	c.createQuery(req, map[string]string{
		"name": name,
	})

	resp, err := c.client.Do(req)
	if err != nil {
		return Server{}, err
	}

	// Server was found.
	if resp.StatusCode == 200 {
		// Decode the response from JSON to a struct.
		var i Server
		j := json.NewDecoder(resp.Body)
		err := j.Decode(&i)
		if err != nil {
			return Server{}, err
		}

		return i, nil
	}

	var errResp ErrorResponse
	j := json.NewDecoder(resp.Body)
	err = j.Decode(&errResp)
	if err != nil {
		return Server{}, err
	}
	return Server{}, errors.New(errResp.Message)
}
