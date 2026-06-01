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

func TestListComments(t *testing.T) {
	var capturedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/1/cards/c1/actions" {
			t.Errorf("path = %s, want /1/cards/c1/actions", r.URL.Path)
		}
		capturedQuery = r.URL.RawQuery
		if err := json.NewEncoder(w).Encode([]map[string]any{
			{"id": "a1", "type": "commentCard", "date": "2026-03-13T12:00:00Z", "memberCreator": map[string]any{"id": "m1"}, "data": map[string]any{"text": "hello"}},
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	comments, err := client.ListComments(context.Background(), "c1")
	if err != nil {
		t.Fatalf("ListComments() error: %v", err)
	}
	if !strings.Contains(capturedQuery, "filter=commentCard") {
		t.Errorf("query missing filter=commentCard: %s", capturedQuery)
	}
	if len(comments) != 1 || comments[0].ID != "a1" {
		t.Fatalf("comments = %+v", comments)
	}
}

func TestAddComment(t *testing.T) {
	var capturedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/1/cards/c1/actions/comments" {
			t.Errorf("path = %s, want /1/cards/c1/actions/comments", r.URL.Path)
		}
		capturedQuery = r.URL.RawQuery
		if err := json.NewEncoder(w).Encode(map[string]any{
			"id": "a1", "type": "commentCard", "date": "2026-03-13T12:00:00Z", "memberCreator": map[string]any{"id": "m1"}, "data": map[string]any{"text": "hello"},
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	comment, err := client.AddComment(context.Background(), "c1", "hello")
	if err != nil {
		t.Fatalf("AddComment() error: %v", err)
	}
	if !strings.Contains(capturedQuery, "text=hello") {
		t.Errorf("query missing text=hello: %s", capturedQuery)
	}
	if comment.ID != "a1" {
		t.Errorf("ID = %q, want a1", comment.ID)
	}
}

func TestUpdateComment(t *testing.T) {
	var capturedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/1/actions/a1/text" {
			t.Errorf("path = %s, want /1/actions/a1/text", r.URL.Path)
		}
		capturedQuery = r.URL.RawQuery
		if err := json.NewEncoder(w).Encode(map[string]any{
			"id": "a1", "type": "commentCard", "date": "2026-03-13T12:00:00Z", "memberCreator": map[string]any{"id": "m1"}, "data": map[string]any{"text": "updated"},
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	comment, err := client.UpdateComment(context.Background(), "a1", "updated")
	if err != nil {
		t.Fatalf("UpdateComment() error: %v", err)
	}
	if !strings.Contains(capturedQuery, "value=updated") {
		t.Errorf("query missing value=updated: %s", capturedQuery)
	}
	if comment.Data.Text != "updated" {
		t.Errorf("Text = %q, want updated", comment.Data.Text)
	}
}

func TestDeleteComment(t *testing.T) {
	var capturedMethod string
	var capturedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	if err := client.DeleteComment(context.Background(), "a1"); err != nil {
		t.Fatalf("DeleteComment() error: %v", err)
	}
	if capturedMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", capturedMethod)
	}
	if capturedPath != "/1/actions/a1" {
		t.Errorf("path = %s, want /1/actions/a1", capturedPath)
	}
}
