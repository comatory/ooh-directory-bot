package client

import (
	"errors"
	"fmt"
	"net/http"
)

type HttpClient interface {
	NewRequest(url string, method string) (*http.Request, error)
	DispatchRequest(req *http.Request) (*http.Response, error)
}

type Client struct {
	Instance *http.Client
}

func CreateHttpClient() Client {
	return Client{
		Instance: &http.Client{},
	}
}

func (*Client) NewRequest(url string, method string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "ooh-directory-random-bot")
	req.Header.Set("Accept-Language", "en-us, en-gb, en")
	req.Header.Set("Accept", "text/html")

	return req, nil
}

func (client *Client) DispatchRequest(req *http.Request) (*http.Response, error) {
	res, err := client.Instance.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Request failed: %d", res.StatusCode))
	}

	return res, nil
}
