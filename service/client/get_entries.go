package client

import (
	"net/http"
	"path"
	"strconv"

	"github.com/pkg/errors"
	"github.com/tchaudhry91/zsh-archaeologist/history"
)

// GetEntries fetches history events from the server
func (c *HistoryClient) GetEntries(req GetEntriesRequest) ([]history.Entry, error) {
	url := c.remoteURL
	url.Path = path.Join(url.Path, getEntriesPath)

	// Encode Values
	url.Query().Add("start", strconv.Itoa(int(req.Start)))
	url.Query().Add("end", strconv.Itoa(int(req.End)))
	url.Query().Add("limit", strconv.Itoa(int(req.Limit)))

	// Build Request
	r, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	// Headers
	c.attachHeaders(r)

	// Decode Response
	response := GetEntriesResponse{}
	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = decodeResponse(resp.Body, &response)
	if err != nil {
		return nil, err
	}

	if response.Err != "" {
		return response.Entries, errors.Errorf(response.Err)
	}

	return response.Entries, nil
}

// GetEntriesRequest contains the request params to query for history entries
type GetEntriesRequest struct {
	Start uint64
	End   uint64
	Limit int64
}

// GetEntriesResponse is the struct to unmarshal JSON output from the server
type GetEntriesResponse struct {
	Entries []history.Entry `json:"entries,omitempty"`
	Err     string          `json:"err,omitempty"`
}
