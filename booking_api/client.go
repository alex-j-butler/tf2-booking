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

func (c BookingClient) NextServer(tag string) (Server, error) {
	req, _ := http.NewRequest(
		"GET",
		c.getAPIPath("v1", "servers", "next"),
		nil,
	)

	q := req.URL.Query()
	q.Add("tag", tag)
	req.URL.RawQuery = q.Encode()

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
