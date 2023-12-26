package client

import (
	"errors"
	"fmt"
	"io"
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

func addRequiredHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "ooh-directory-random-bot")
	req.Header.Set("Accept-Language", "en-us, en-gb, en")
	req.Header.Set("Accept", "text/html")
}

func (*Client) NewRequest(url string, method string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}

	addRequiredHeaders(req)

	return req, nil
}

func (*Client) NewRequestWithBody(url string, method string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	addRequiredHeaders(req)

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
