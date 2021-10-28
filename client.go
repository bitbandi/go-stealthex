package stealthex

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type client struct {
	apiKey      string
	httpClient  *http.Client
	httpTimeout time.Duration
	debug       bool
}

// NewClient return a new StealthEX HTTP client
func NewClient(apiKey string) (c *client) {
	return &client{apiKey, &http.Client{}, 30 * time.Second, false}
}

// NewClientWithCustomHttpConfig returns a new StealthEX HTTP client using the predefined http client
func NewClientWithCustomHttpConfig(apiKey string, httpClient *http.Client) (c *client) {
	timeout := httpClient.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &client{apiKey, httpClient, timeout, false}
}

// NewClient returns a new StealthEX HTTP client with custom timeout
func NewClientWithCustomTimeout(apiKey string, timeout time.Duration) (c *client) {
	return &client{apiKey, &http.Client{}, timeout, false}
}

func (c client) dumpRequest(r *http.Request) {
	if r == nil {
		log.Print("dumpReq ok: <nil>")
		return
	}
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Print("dumpReq err:", err)
	} else {
		log.Print("dumpReq ok:", string(dump))
	}
}

func (c client) dumpResponse(r *http.Response) {
	if r == nil {
		log.Print("dumpResponse ok: <nil>")
		return
	}
	dump, err := httputil.DumpResponse(r, true)
	if err != nil {
		log.Print("dumpResponse err:", err)
	} else {
		log.Print("dumpResponse ok:", string(dump))
	}
}

// doTimeoutRequest do a HTTP request with timeout
func (c *client) doTimeoutRequest(timer *time.Timer, req *http.Request) (*http.Response, error) {
	// Do the request in the background so we can check the timeout
	type result struct {
		resp *http.Response
		err  error
	}
	done := make(chan result, 1)
	go func() {
		if c.debug {
			c.dumpRequest(req)
		}
		resp, err := c.httpClient.Do(req)
		if c.debug {
			c.dumpResponse(resp)
		}
		done <- result{resp, err}
	}()
	// Wait for the read or the timeout
	select {
	case r := <-done:
		return r.resp, r.err
	case <-timer.C:
		return nil, errors.New("timeout on reading data from StealthEX API")
	}
}

// do prepare and process HTTP request to StealthEX API
func (c *client) do(method string, resource string, payload map[string]string, authNeeded bool) (response []byte, err error) {
	connectTimer := time.NewTimer(c.httpTimeout)

	var rawurl string
	if strings.HasPrefix(resource, "http") {
		rawurl = resource
	} else {
		rawurl = fmt.Sprintf("%s/%s", API_BASE, resource)
	}

	var formData string
	var URL *url.URL
	URL, err = url.Parse(rawurl)
	if err != nil {
		return
	}
	q := URL.Query()
	// Auth
	if authNeeded {
		if len(c.apiKey) == 0 {
			err = errors.New("You need to set API Key to call this method")
			return
		}
		q.Set("api_key", c.apiKey)
	}

	if method == "GET" {
		for key, value := range payload {
			q.Set(key, value)
		}
	} else {
		formValues := url.Values{}
		for key, value := range payload {
			formValues.Set(key, value)
		}
		formData = formValues.Encode()
	}
	URL.RawQuery = q.Encode()
	rawurl = URL.String()
	req, err := http.NewRequest(method, rawurl, strings.NewReader(formData))
	if err != nil {
		return
	}

	req.Header.Add("Accept", "application/json")

	resp, err := c.doTimeoutRequest(connectTimer, req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 401 {
		err = errors.New(resp.Status)
	}
	return response, err
}
