package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"path"

	"github.com/pkg/errors"
)

// Register creates an account on the remote service
func (c *HistoryClient) Register(req RegisterRequest) error {
	url := c.remoteURL
	url.Path = path.Join(url.Path, "register")

	// Build Request
	data, err := json.Marshal(&req)
	if err != nil {
		return err
	}
	r, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("Failed with status: %d", resp.StatusCode)
	}

	response := registerResponse{}
	err = decodeResponse(resp.Body, &response)
	if err != nil {
		return err
	}

	if response.Err != "" {
		return errors.Errorf(response.Err)
	}
	return nil
}

// RegisterRequest contains user information to register
type RegisterRequest struct {
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}

type registerResponse struct {
	Err string `json:"err,omitempty"`
}
