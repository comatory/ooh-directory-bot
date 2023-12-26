package scraper

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"internal/client"
)

type mockSuccesfulResponseClient struct {
	Instance *http.Client
}

func (mockClient *mockSuccesfulResponseClient) DispatchRequest(req *http.Request) (*http.Response, error) {
	return mockClient.Instance.Do(req)
}

func (*mockSuccesfulResponseClient) NewRequestBuilder(url string) *client.RequestBuilder {
	builder := client.RequestBuilder{}
	builder.New(url)

	return &builder
}

type mockUnsuccesfulResponseClient struct {
	Instance *http.Client
}

func (*mockUnsuccesfulResponseClient) DispatchRequest(req *http.Request) (*http.Response, error) {
	return nil, errors.New("Failed in test")
}

func (*mockUnsuccesfulResponseClient) NewRequestBuilder(url string) *client.RequestBuilder {
	builder := client.RequestBuilder{}
	builder.New(url)

	return &builder
}

func TestSuccesfulScrapeResponse(t *testing.T) {
	expected := "!DOCTYPE html\n<html><body>test</body></html>"
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, expected)
		}),
	)

	defer server.Close()

	body, _ := ScrapeRandom(server.URL, &mockSuccesfulResponseClient{
		Instance: server.Client(),
	})

	if body != expected {
		t.Errorf("Expected response text \"%s\", got \"%s\"", expected, body)
	}
}

func TestNonSuccesfulScrapeResponse(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}),
	)

	defer server.Close()

	_, err := ScrapeRandom(server.URL, &mockUnsuccesfulResponseClient{
		Instance: server.Client(),
	})

	if err == nil {
		t.Errorf("Expected error, got \"%s\"", err)
	}
}
