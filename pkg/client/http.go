/*
Modify the http package to allow easy calls
*/
package client

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type HTTP struct {
	client *http.Client
}

// timeout the http request to a level and sets a proxy is passed
func New(timeout time.Duration, proxy string) (*HTTP, error) {
	httpConn := &HTTP{
		client: http.DefaultClient,
	}

	if timeout > 0 {
		httpConn.client.Timeout = timeout
	}

	if proxy != "" {
		// parse the url
		proxy, err := url.Parse(proxy)

		if err != nil {
			return nil, err
		}

		t := http.DefaultTransport.(*http.Transport)
		t.Proxy = func(*http.Request) (*url.URL, error) {
			return proxy, nil
		}
		httpConn.client.Transport = t
	}

	return httpConn, nil
}

// make a GET request and return some info about the request
func Get(url string) (int, []byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return -1, nil, err
	}

	// read the response body
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		return resp.StatusCode, nil, err
	}

	return resp.StatusCode, bodyBytes, nil
}
