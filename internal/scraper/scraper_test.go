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

type mockUserAgentClient struct {
	Instance *http.Client
}

func (client *mockUserAgentClient) DispatchRequest(req *http.Request) (*http.Response, error) {
	return client.Instance.Do(req)
}

func (client *mockUserAgentClient) NewRequest(url string, method string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)

	req.Header.Set("User-Agent", "ooh-directory-random-bot")

	return req, err
}

type mockAcceptLanguageClient struct {
	Instance *http.Client
}

func (client *mockAcceptLanguageClient) DispatchRequest(req *http.Request) (*http.Response, error) {
	return client.Instance.Do(req)
}

func (client *mockAcceptLanguageClient) NewRequest(url string, method string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)

	req.Header.Set("Accept-Language", "en-us, en-gb, en")

	return req, err
}

type mockAcceptClient struct {
	Instance *http.Client
}

func (client *mockAcceptClient) DispatchRequest(req *http.Request) (*http.Response, error) {
	return client.Instance.Do(req)
}

func (client *mockAcceptClient) NewRequest(url string, method string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)

	req.Header.Set("Accept", "text/html")

	return req, err
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

func TestSuccesfulResponse(t *testing.T) {
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

func TestUserAgentHeader(t *testing.T) {
	expectedUserAgent := "ooh-directory-random-bot"
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userAgentHeader := r.Header.Get("User-Agent")
			if userAgentHeader != expectedUserAgent {
				t.Errorf("Expected user agent header \"%s\", got \"%s\"", expectedUserAgent, userAgentHeader)
			}

			io.WriteString(w, "ok")
		}),
	)

	defer server.Close()

	ScrapeRandom(server.URL, &mockUserAgentClient{
		Instance: server.Client(),
	})
}

func TestAcceptLanguageHeader(t *testing.T) {
	expectedAcceptLanguage := "en-us, en-gb, en"

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			acceptLanguageHeader := r.Header.Get("Accept-Language")
			if acceptLanguageHeader != expectedAcceptLanguage {
				t.Errorf("Expected accept language header \"%s\", got \"%s\"", expectedAcceptLanguage, acceptLanguageHeader)
			}

			io.WriteString(w, "ok")
		}),
	)

	defer server.Close()

	ScrapeRandom(server.URL, &mockAcceptLanguageClient{
		Instance: server.Client(),
	})
}

func TestAcceptHeader(t *testing.T) {
	expectedAccept := "text/html"

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			acceptHeader := r.Header.Get("Accept")
			if acceptHeader != expectedAccept {
				t.Errorf("Expected accept header \"%s\", got \"%s\"", expectedAccept, acceptHeader)
			}

			io.WriteString(w, "ok")
		}),
	)

	defer server.Close()

	ScrapeRandom(server.URL, &mockAcceptClient{
		Instance: server.Client(),
	})
}

func TestNonSuccesfulResponse(t *testing.T) {
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
