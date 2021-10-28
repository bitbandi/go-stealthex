package stealthex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	API_BASE = "https://api.stealthex.io/api/v2" // StealthEX API endpoint
)

// New returns an instantiated StealthEX struct
func New(apiKey string) *StealthEX {
	client := NewClient(apiKey)
	return &StealthEX{client}
}

// NewWithCustomHttpClient returns an instantiated StealthEX struct with custom http client
func NewWithCustomHttpClient(apiKey string, httpClient *http.Client) *StealthEX {
	client := NewClientWithCustomHttpConfig(apiKey, httpClient)
	return &StealthEX{client}
}

// NewWithCustomTimeout returns an instantiated StealthEX struct with custom timeout
func NewWithCustomTimeout(apiKey string, timeout time.Duration) *StealthEX {
	client := NewClientWithCustomTimeout(apiKey, timeout)
	return &StealthEX{client}
}

// StealthEX represent a StealthEX client
type StealthEX struct {
	client *client
}

// SetDebug sets enable/disable http request/response dump
func (b *StealthEX) SetDebug(enable bool) {
	b.client.debug = enable
}

// GetTrade is used to get a trade at StealthEX along with other meta data.
func (b *StealthEX) GetTrade(id string) (trade Trade, err error) {
	r, err := b.client.do("GET", fmt.Sprintf("https://stealthex.io/api/exchange/%s", id), nil, false)
	if err != nil {
		return
	}
	var response interface{}
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	err = json.Unmarshal(r, &trade)
	return
}
