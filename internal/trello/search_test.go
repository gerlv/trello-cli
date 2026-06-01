package trello_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gerlv/trello-cli/internal/trello"
)

func TestSearchCards(t *testing.T) {
	var capturedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/1/search" {
			t.Errorf("path = %s, want /1/search", r.URL.Path)
		}
		capturedQuery = r.URL.RawQuery
		if err := json.NewEncoder(w).Encode(map[string]any{
			"cards": []map[string]any{
				{"id": "c1", "name": "Alpha card", "idBoard": "b1", "idList": "l1"},
			},
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	result, err := client.SearchCards(context.Background(), "alpha")
	if err != nil {
		t.Fatalf("SearchCards() error: %v", err)
	}
	for _, want := range []string{"query=alpha", "modelTypes=cards"} {
		if !strings.Contains(capturedQuery, want) {
			t.Errorf("query missing %s: %s", want, capturedQuery)
		}
	}
	if result.Query != "alpha" || len(result.Cards) != 1 || result.Cards[0].ID != "c1" {
		t.Fatalf("result = %+v", result)
	}
}

func TestSearchCardsPreservesCommaInQuery(t *testing.T) {
	var capturedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		if err := json.NewEncoder(w).Encode(map[string]any{
			"cards": []map[string]any{},
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	result, err := client.SearchCards(context.Background(), "alpha,beta")
	if err != nil {
		t.Fatalf("SearchCards() error: %v", err)
	}
	if !strings.Contains(capturedQuery, "query=alpha%2Cbeta") {
		t.Fatalf("query should preserve comma: %s", capturedQuery)
	}
	if result.Query != "alpha,beta" {
		t.Fatalf("result = %+v", result)
	}
}

func TestSearchBoards(t *testing.T) {
	var capturedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/1/search" {
			t.Errorf("path = %s, want /1/search", r.URL.Path)
		}
		capturedQuery = r.URL.RawQuery
		if err := json.NewEncoder(w).Encode(map[string]any{
			"boards": []map[string]any{
				{"id": "b1", "name": "Alpha board"},
			},
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	result, err := client.SearchBoards(context.Background(), "alpha")
	if err != nil {
		t.Fatalf("SearchBoards() error: %v", err)
	}
	for _, want := range []string{"query=alpha", "modelTypes=boards"} {
		if !strings.Contains(capturedQuery, want) {
			t.Errorf("query missing %s: %s", want, capturedQuery)
		}
	}
	if result.Query != "alpha" || len(result.Boards) != 1 || result.Boards[0].ID != "b1" {
		t.Fatalf("result = %+v", result)
	}
}
