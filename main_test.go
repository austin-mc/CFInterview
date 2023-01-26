package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			healthResponse := Health{
				NumRequests: 1,
				NumErrors:   1,
			}
			json, _ := json.Marshal(healthResponse)
			w.Write(json)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		url := server.URL
		err := healthCheck(url, 5)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("404 error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		url := server.URL
		err := healthCheck(url, 5)
		if err == nil {
			t.Error("Expected invalid request error, got nil")
		}
		if err.Error() != "invalid request" {
			t.Error(err)
		}
	})
}
