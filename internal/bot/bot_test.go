package bot

import (
	"testing"
	"io"
	"errors"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"internal/parser"
	"internal/client"
)

type mockClient struct {
	Instance *http.Client
}

func (mockClient *mockClient) DispatchRequest(req *http.Request) (*http.Response, error) {
	return mockClient.Instance.Do(req)
}

func (*mockClient) NewRequestBuilder(url string) *client.RequestBuilder {
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

func TestSuccesfulPosting(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "{\"ok\":true}")
		}),
	)

	defer server.Close()

	result := parser.Result{
		Url: "https://ooh.directory/random/1",
		Title: "Random 1",
	}
	config := Config{
		AccessToken: "123",
		BotServerUrl: server.URL,
	}

	err := PostResult(&result, &config, &mockClient{
		Instance: server.Client(),
	})

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}

func TestUnsuccesfulPosting(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
			io.WriteString(w, "{\"ok\":false}")
		}),
	)

	defer server.Close()

	result := parser.Result{
		Url: "https://ooh.directory/random/1",
		Title: "Random 1",
	}
	config := Config{
		AccessToken: "123",
		BotServerUrl: server.URL,
	}

	err := PostResult(&result, &config, &mockUnsuccesfulResponseClient{
		Instance: server.Client(),
	})

	if err == nil {
		t.Errorf("Expected error, got %v", err)
	}
}

func TestAuthorizationHeader(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "Bearer u=7q" {
				t.Errorf("Expected Authorization header to be \"Bearer u=7q\", got \"%s\"", authHeader)
			}
			io.WriteString(w, "{\"ok\":true}")
		}),
	)

	defer server.Close()

	result := parser.Result{
		Url: "https://ooh.directory/random/1",
		Title: "Random 1",
	}
	config := Config{
		AccessToken: "u=7q",
		BotServerUrl: server.URL,
	}

	err := PostResult(&result, &config, &mockClient{
		Instance: server.Client(),
	})

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}

func TestBaseStatusPayload(t *testing.T) {
	expectedStatus := "https://ooh.directory/random/1 Random 1"
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			payload := Payload{}
			err := json.NewDecoder(r.Body).Decode(&payload)

			if err != nil {
				t.Errorf("Expected body to be non-empty, got \"%v\"", err)
			}

			if payload.Status != expectedStatus {
			  t.Errorf("Expected \"%s\", got \"%s\"", expectedStatus, payload.Status)
			}
			
			io.WriteString(w, "{\"ok\":true}")
		}),
	)

	defer server.Close()

	result := parser.Result{
		Url: "https://ooh.directory/random/1",
		Title: "Random 1",
	}
	config := Config{
		AccessToken: "u=7q",
		BotServerUrl: server.URL,
	}

	err := PostResult(&result, &config, &mockClient{
		Instance: server.Client(),
	})

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}

func TestStatusPayloadWithAuthor(t *testing.T) {
	expectedStatus := "https://ooh.directory/random/1 Random 1 (by John Doe)"
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			payload := Payload{}
			err := json.NewDecoder(r.Body).Decode(&payload)

			if err != nil {
				t.Errorf("Expected body to be non-empty, got \"%v\"", err)
			}

			if payload.Status != expectedStatus {
			  t.Errorf("Expected \"%s\", got \"%s\"", expectedStatus, payload.Status)
			}
			
			io.WriteString(w, "{\"ok\":true}")
		}),
	)

	defer server.Close()

	result := parser.Result{
		Url: "https://ooh.directory/random/1",
		Title: "Random 1",
		AuthorName: "John Doe",
	}
	config := Config{
		AccessToken: "u=7q",
		BotServerUrl: server.URL,
	}

	err := PostResult(&result, &config, &mockClient{
		Instance: server.Client(),
	})

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}

func TestStatusPayloadWithAuthorAndTimestamp(t *testing.T) {
	expectedStatus := "https://ooh.directory/random/1 Random 1 (by John Doe, Feb 02 2021)"
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			payload := Payload{}
			err := json.NewDecoder(r.Body).Decode(&payload)

			if err != nil {
				t.Errorf("Expected body to be non-empty, got \"%v\"", err)
			}

			if payload.Status != expectedStatus {
			  t.Errorf("Expected \"%s\", got \"%s\"", expectedStatus, payload.Status)
			}
			
			io.WriteString(w, "{\"ok\":true}")
		}),
	)

	defer server.Close()

	result := parser.Result{
		Url: "https://ooh.directory/random/1",
		Title: "Random 1",
		AuthorName: "John Doe",
		UpdatedAt: 1612345678,
	}
	config := Config{
		AccessToken: "u=7q",
		BotServerUrl: server.URL,
	}

	err := PostResult(&result, &config, &mockClient{
		Instance: server.Client(),
	})

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}
