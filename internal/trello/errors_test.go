package trello_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gerlv/trello-cli/internal/contract"
	"github.com/gerlv/trello-cli/internal/trello"
)

func TestHTTPError401(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	var result map[string]any
	err := client.Get(context.Background(), "/1/members/me", nil, &result)

	var ce *contract.ContractError
	if !errors.As(err, &ce) {
		t.Fatalf("error should be *ContractError, got %T: %v", err, err)
	}
	if ce.Code != contract.AuthInvalid {
		t.Errorf("Code = %q, want %q", ce.Code, contract.AuthInvalid)
	}
}

func TestHTTPError404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	var result map[string]any
	err := client.Get(context.Background(), "/1/boards/nonexistent", nil, &result)

	var ce *contract.ContractError
	if !errors.As(err, &ce) {
		t.Fatalf("error should be *ContractError, got %T", err)
	}
	if ce.Code != contract.NotFound {
		t.Errorf("Code = %q, want %q", ce.Code, contract.NotFound)
	}
}

func TestHTTPError429(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	opts := trello.DefaultClientOptions()
	opts.MaxRetries = 0
	client := trello.NewClient(server.URL, "k", "t", opts)

	var result map[string]any
	err := client.Get(context.Background(), "/1/boards/b1", nil, &result)

	var ce *contract.ContractError
	if !errors.As(err, &ce) {
		t.Fatalf("error should be *ContractError, got %T", err)
	}
	if ce.Code != contract.RateLimited {
		t.Errorf("Code = %q, want %q", ce.Code, contract.RateLimited)
	}
}

func TestHTTPError500(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	var result map[string]any
	err := client.Get(context.Background(), "/1/boards/b1", nil, &result)

	var ce *contract.ContractError
	if !errors.As(err, &ce) {
		t.Fatalf("error should be *ContractError, got %T", err)
	}
	if ce.Code != contract.HTTPError {
		t.Errorf("Code = %q, want %q", ce.Code, contract.HTTPError)
	}
}
