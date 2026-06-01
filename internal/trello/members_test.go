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

func TestListMembers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/1/boards/b1/members" {
			t.Errorf("path = %s, want /1/boards/b1/members", r.URL.Path)
		}
		if err := json.NewEncoder(w).Encode([]map[string]any{
			{"id": "m1", "username": "alice", "fullName": "Alice Example"},
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	members, err := client.ListMembers(context.Background(), "b1")
	if err != nil {
		t.Fatalf("ListMembers() error: %v", err)
	}
	if len(members) != 1 || members[0].ID != "m1" {
		t.Fatalf("members = %+v", members)
	}
}

func TestAddMemberToCard(t *testing.T) {
	var capturedMethod string
	var capturedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		capturedQuery = r.URL.RawQuery
		if r.URL.Path != "/1/cards/c1/idMembers" {
			t.Errorf("path = %s, want /1/cards/c1/idMembers", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	if err := client.AddMemberToCard(context.Background(), "c1", "m1"); err != nil {
		t.Fatalf("AddMemberToCard() error: %v", err)
	}
	if capturedMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", capturedMethod)
	}
	if !strings.Contains(capturedQuery, "value=m1") {
		t.Errorf("query missing value=m1: %s", capturedQuery)
	}
}

func TestRemoveMemberFromCard(t *testing.T) {
	var capturedMethod string
	var capturedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	if err := client.RemoveMemberFromCard(context.Background(), "c1", "m1"); err != nil {
		t.Fatalf("RemoveMemberFromCard() error: %v", err)
	}
	if capturedMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", capturedMethod)
	}
	if capturedPath != "/1/cards/c1/idMembers/m1" {
		t.Errorf("path = %s, want /1/cards/c1/idMembers/m1", capturedPath)
	}
}

func TestGetMe(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/1/members/me" {
			t.Errorf("path = %s, want /1/members/me", r.URL.Path)
		}
		if err := json.NewEncoder(w).Encode(map[string]any{
			"id": "me1", "username": "alice", "fullName": "Alice Example",
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	member, err := client.GetMe(context.Background())
	if err != nil {
		t.Fatalf("GetMe() error: %v", err)
	}
	if member.ID != "me1" {
		t.Fatalf("member = %+v", member)
	}
}
