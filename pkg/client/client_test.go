package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	baseURL := "http://example.com"
	apiKey := "test-api-key"

	client := NewClient(baseURL, apiKey)

	if client.baseURL != baseURL {
		t.Errorf("Expected baseURL to be '%s', got '%s'", baseURL, client.baseURL)
	}

	if client.apiKey != apiKey {
		t.Errorf("Expected apiKey to be '%s', got '%s'", apiKey, client.apiKey)
	}

	if client.httpClient == nil {
		t.Error("Expected httpClient to be initialized, got nil")
	}
}

func TestClientGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected request method to be GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/v3/test" {
			t.Errorf("Expected request path to be '/api/v3/test', got '%s'", r.URL.Path)
		}

		if r.URL.Query().Get("apikey") != "test-api-key" {
			t.Errorf("Expected apikey query parameter to be 'test-api-key', got '%s'", r.URL.Query().Get("apikey"))
		}

		if r.URL.Query().Get("param1") != "value1" {
			t.Errorf("Expected param1 query parameter to be 'value1', got '%s'", r.URL.Query().Get("param1"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-api-key")

	params := map[string]string{"param1": "value1"}
	data, err := client.Get(context.Background(), "test", params)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := `{"result":"success"}`
	if string(data) != expected {
		t.Errorf("Expected response data to be '%s', got '%s'", expected, string(data))
	}
}
