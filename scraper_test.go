package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSuccesfulResponse(t *testing.T) {
	expected := "!DOCTYPE html\n<html><body>test</body></html>"
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, expected)
		}),
	)

	client := *server.Client()

	defer server.Close()

	body, _ := ScrapeRandom(server.URL, client)

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

	client := *server.Client()

	defer server.Close()

	ScrapeRandom(server.URL, client)
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

	client := *server.Client()

	defer server.Close()

	ScrapeRandom(server.URL, client)
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

	client := *server.Client()

	defer server.Close()

	ScrapeRandom(server.URL, client)
}

func TestNonSuccesfulResponse(t *testing.T) {
	expectedLog := fmt.Sprint("Request failed: 503")

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}),
	)

	client := *server.Client()

	defer server.Close()

	_, err := ScrapeRandom(server.URL, client)

	if err.Error() != expectedLog {
		t.Errorf("Expected log \"%s\", got \"%s\"", expectedLog, err.Error())
	}
}
