package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"path"

	"github.com/pkg/errors"
)

// Login logs into the service and retrieves the token for subsequent use
func (c *HistoryClient) Login(req LoginRequest) (tokenStr string, err error) {
	url := c.remoteURL
	url.Path = path.Join(url.Path, "login")

	// Build Request
	data, err := json.Marshal(&req)
	if err != nil {
		return "", err
	}
	r, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("Failed with status: %d", resp.StatusCode)
	}

	response := loginResponse{}
	err = decodeResponse(resp.Body, &response)
	if err != nil {
		return "", err
	}

	if response.Err != "" {
		return "", errors.Errorf(response.Err)
	}
	return response.Token, nil
}

// LoginRequest contains user information to register
type LoginRequest struct {
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}

type loginResponse struct {
	Token string `json:"token,omitempty"`
	Err   string `json:"err,omitempty"`
}
