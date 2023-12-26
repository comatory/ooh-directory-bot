package client

import (
	"errors"
	"fmt"
	"log"
	"io"
	"net/http"
)

type HttpClient interface {
	DispatchRequest(req *http.Request) (*http.Response, error)
	NewRequestBuilder(url string) *RequestBuilder
}

type Client struct {
	Instance *http.Client
}

func CreateHttpClient() Client {
	return Client{
		Instance: &http.Client{},
	}
}

func (client *Client) DispatchRequest(req *http.Request) (*http.Response, error) {
	log.Println(fmt.Sprintf("Request %s %s", req.Method, req.URL.String()))
	res, err := client.Instance.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Request failed: %d", res.StatusCode))
	}

	return res, nil
}

func (*Client) NewRequestBuilder(url string) *RequestBuilder {
	builder := RequestBuilder{}

	return builder.New(url)
}

type RequestBuilder struct {
	url     string
	method  string
	body    io.Reader
	headers map[string]string
}

func (builder *RequestBuilder) New(url string) *RequestBuilder {
	builder.url = url
	builder.method = http.MethodGet
	builder.body = nil
	builder.headers = make(map[string]string)
	builder.headers["User-Agent"] = "ooh-directory-random-bot"

	return builder
}

func (builder *RequestBuilder) Method(method string) *RequestBuilder {
	builder.method = method

	return builder
}

func (builder *RequestBuilder) Header(key string, value string) *RequestBuilder {
	builder.headers[key] = value

	return builder
}

func (builder *RequestBuilder) Body(body io.Reader) *RequestBuilder {
	builder.body = body

	return builder
}

func (builder *RequestBuilder) Build() (*http.Request, error) {
	req, err := http.NewRequest(builder.method, builder.url, builder.body)

	if err != nil {
		return nil, err
	}

	for key, value := range builder.headers {
		req.Header.Set(key, value)
	}

	return req, err
}
