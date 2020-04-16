package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"path"

	"github.com/pkg/errors"
	"github.com/tchaudhry91/zsh-archaeologist/history"
)

// PutEntries sends the history events to the server
func (c *HistoryClient) PutEntries(req PutEntriesRequest) (updated int64, err error) {
	url := c.remoteURL
	url.Path = path.Join(url.Path, getEntriesPath)

	// Build Request
	data, err := json.Marshal(&req)
	if err != nil {
		return 0, err
	}
	r, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(data))
	if err != nil {
		return 0, err
	}

	// Make the request
	resp, err := c.client.Do(r)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	response := PutEntriesResponse{}
	err = decodeResponse(resp.Body, &response)
	if err != nil {
		return 0, err
	}

	if response.Err != "" {
		return response.Updated, errors.Errorf(response.Err)
	}
	return response.Updated, nil
}

// PutEntriesRequest contains the parameters to send history entries
type PutEntriesRequest struct {
	Entries []history.Entry `json:"entries,omitempty"`
}

// PutEntriesResponse contains the response for PutEntries
type PutEntriesResponse struct {
	Updated int64  `json:"updated,omitempty"`
	Err     string `json:"err,omitempty"`
}
