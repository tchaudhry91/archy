package client

import (
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/pkg/errors"
	"github.com/tchaudhry91/archy/history"
)

// GetEntries fetches history events from the server
func (c *HistoryClient) GetEntries(req GetEntriesRequest) ([]history.Entry, error) {
	uri := c.remoteURL
	uri.Path = path.Join(uri.Path, getEntriesPath)

	q := uri.Query()
	// Encode Values
	q.Add("start", strconv.Itoa(int(req.Start)))
	q.Add("end", strconv.Itoa(int(req.End)))
	q.Add("limit", strconv.Itoa(int(req.Limit)))
	if req.Command != "" {
		q.Add("command", url.QueryEscape(req.Command))
	}
	if req.Machine != "" {
		q.Add("machine", url.QueryEscape(req.Machine))
	}

	uri.RawQuery = q.Encode()

	// Build Request
	r, err := http.NewRequest("GET", uri.String(), nil)
	if err != nil {
		return nil, err
	}

	// Headers
	c.attachHeaders(r)

	// Decode Response
	response := getEntriesResponse{}
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
	Start   uint64
	End     uint64
	Limit   int64
	Machine string
	Command string
}

// GetEntriesResponse is the struct to unmarshal JSON output from the server
type getEntriesResponse struct {
	Entries []history.Entry `json:"entries,omitempty"`
	Err     string          `json:"err,omitempty"`
}
