package scraper

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockSuccesfulResponseClient struct {
	Instance *http.Client
}

func (client *mockSuccesfulResponseClient) DispatchRequest(req *http.Request) (*http.Response, error) {
	return client.Instance.Do(req)
}

func (client *mockSuccesfulResponseClient) NewRequest(url string, method string) (*http.Request, error) {
	return http.NewRequest(method, url, nil)
}

type mockUnsuccesfulResponseClient struct {
	Instance *http.Client
}

func (client *mockUnsuccesfulResponseClient) DispatchRequest(req *http.Request) (*http.Response, error) {
	return nil, errors.New("Failed in test")
}

func (client *mockUnsuccesfulResponseClient) NewRequest(url string, method string) (*http.Request, error) {
	return http.NewRequest(method, url, nil)
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
