// TODO: Rewrite the booking API client & test it.
// This should also be renamed to the server API client.
package booking_api

import (
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

func (c BookingClient) GetServer(uuid string) (Server, error) {
	req, _ := http.NewRequest(
		"GET",
		c.getAPIPath("v1", "servers", "list"),
		nil,
	)

	c.createQuery(req, map[string]string{
		"uuid": uuid,
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
