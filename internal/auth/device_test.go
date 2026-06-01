package auth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gerlv/trello-cli/internal/auth"
)

func TestDeviceClient_RequestCode(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/device/code" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"device_code":      "abc123",
			"user_code":        "WDJBMJHT",
			"verification_uri": "https://example.com/help",
			"expires_in":       900,
			"interval":         5,
		})
	}))
	defer srv.Close()

	client := auth.NewDeviceClient(srv.URL)
	resp, err := client.RequestCode()
	if err != nil {
		t.Fatalf("RequestCode: %v", err)
	}
	if resp.UserCode != "WDJBMJHT" {
		t.Errorf("expected WDJBMJHT, got %s", resp.UserCode)
	}
	if resp.DeviceCode != "abc123" {
		t.Errorf("expected abc123, got %s", resp.DeviceCode)
	}
}

func TestDeviceClient_PollToken_Pending(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "authorization_pending"})
	}))
	defer srv.Close()

	client := auth.NewDeviceClient(srv.URL)
	_, err := client.PollToken("device123")
	if err == nil || err.Error() != "authorization_pending" {
		t.Errorf("expected authorization_pending error, got %v", err)
	}
}

func TestDeviceClient_PollToken_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "trello-token-xyz",
			"api_key":      "api-key-123",
		})
	}))
	defer srv.Close()

	client := auth.NewDeviceClient(srv.URL)
	result, err := client.PollToken("device123")
	if err != nil {
		t.Fatalf("PollToken: %v", err)
	}
	if result.AccessToken != "trello-token-xyz" {
		t.Errorf("expected trello-token-xyz, got %s", result.AccessToken)
	}
	if result.APIKey != "api-key-123" {
		t.Errorf("expected api-key-123, got %s", result.APIKey)
	}
}

func TestDeviceClient_PollToken_Expired(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "expired_token"})
	}))
	defer srv.Close()

	client := auth.NewDeviceClient(srv.URL)
	_, err := client.PollToken("device123")
	if err == nil || err.Error() != "expired_token" {
		t.Errorf("expected expired_token error, got %v", err)
	}
}
