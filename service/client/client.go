package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

const getEntriesPath = "/entries"

// HistoryClient provides the client SDK to talk to the HistoryServer
type HistoryClient struct {
	remoteURL url.URL
	userToken string
	client    *http.Client
}

// NewHistoryClient initializes a new HistoryClient
func NewHistoryClient(remoteURL string, userToken string, timeout int) (*HistoryClient, error) {
	u, err := url.Parse(remoteURL)
	if err != nil {
		return nil, err
	}
	return &HistoryClient{
		remoteURL: *u,
		userToken: userToken,
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}, nil
}

func (c *HistoryClient) attachHeaders(req *http.Request) {
	req.Header.Add("token", c.userToken)
}

func decodeResponse(reader io.Reader, v interface{}) error {
	return json.NewDecoder(reader).Decode(v)
}
