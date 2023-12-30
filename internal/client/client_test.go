package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createMockClient(instance *http.Client) Client {
	return Client{
		Instance: instance,
	}
}

func TestSuccesfulResponse(t *testing.T) {
	expected := "!DOCTYPE html\n<html><body>test</body></html>"
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, expected)
		}),
	)

	defer server.Close()

	client := createMockClient(server.Client())
	reqConfig := client.NewRequestBuilder(server.URL)
	req, _ := reqConfig.Build()

	res, _ := client.DispatchRequest(req)

	bytes, _ := io.ReadAll(res.Body)
	body := string(bytes)

	if body != expected {
		t.Errorf("Expected response text \"%s\", got \"%s\"", expected, body)
	}
}

func TestNonSuccesfulResponse(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}),
	)

	defer server.Close()

	client := createMockClient(server.Client())
	reqConfig := client.NewRequestBuilder(server.URL)
	req, _ := reqConfig.Build()

	_, err := client.DispatchRequest(req)

	if err == nil {
		t.Errorf("Expected error, got \"%s\"", err)
	}
}

func TestUserAgentHeader(t *testing.T) {
	expectedUserAgent := "ooh-directory-random-bot"
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "<!DOCTYPE html><html><body>test</body></html>")
		}),
	)

	defer server.Close()

	client := createMockClient(server.Client())
	reqConfig := client.NewRequestBuilder(server.URL).Header("User-Agent", expectedUserAgent)
	req, _ := reqConfig.Build()

	if req.Header.Get("User-Agent") != expectedUserAgent {
		t.Errorf("Expected user agent header \"%s\", got \"%s\"", expectedUserAgent, req.Header.Get("User-Agent"))
	}
}

func TestAcceptLanguageHeader(t *testing.T) {
	expectedAcceptLanguage := "en-us, en-gb, en"

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "<!DOCTYPE html><html><body>test</body></html>")
		}),
	)

	defer server.Close()

	client := createMockClient(server.Client())
	reqConfig := client.NewRequestBuilder(server.URL).Header("Accept-Language", expectedAcceptLanguage)
	req, _ := reqConfig.Build()

	if req.Header.Get("Accept-Language") != expectedAcceptLanguage {
		t.Errorf("Expected accept language header \"%s\", got \"%s\"", expectedAcceptLanguage, req.Header.Get("Accept-Language"))
	}
}

func TestAcceptHeader(t *testing.T) {
	expectedAccept := "text/html"

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "<!DOCTYPE html><html><body>test</body></html>")
		}),
	)

	defer server.Close()

	client := createMockClient(server.Client())
	reqConfig := client.NewRequestBuilder(server.URL).Header("Accept", expectedAccept)
	req, _ := reqConfig.Build()

	if req.Header.Get("Accept") != expectedAccept {
		t.Errorf("Expected accept header \"%s\", got \"%s\"", expectedAccept, req.Header.Get("Accept"))
	}
}
