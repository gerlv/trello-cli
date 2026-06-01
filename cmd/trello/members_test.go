package main

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/gerlv/trello-cli/internal/credentials"
	"github.com/gerlv/trello-cli/internal/trello"
)

func TestMembersListCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		listMembersFn: func(ctx context.Context, boardID string) ([]trello.Member, error) {
			if boardID != "b1" {
				t.Fatalf("boardID = %q, want b1", boardID)
			}
			return []trello.Member{{ID: "m1", Username: "alice"}}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"members", "list", "--board", "b1"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("members list failed: %v", err)
	}
	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}
	data := envelope["data"].([]any)
	if len(data) != 1 {
		t.Fatalf("data = %+v", envelope["data"])
	}
	member := data[0].(map[string]any)
	if member["id"] != "m1" || member["username"] != "alice" {
		t.Fatalf("member = %+v", member)
	}
}

func TestMembersAddCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		addMemberToCardFn: func(ctx context.Context, cardID, memberID string) error {
			if cardID != "c1" || memberID != "m1" {
				t.Fatalf("card/member = %q/%q", cardID, memberID)
			}
			return nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"members", "add", "--card", "c1", "--member", "m1"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("members add failed: %v", err)
	}
	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}
	data := envelope["data"].(map[string]any)
	if data["success"] != true || data["id"] != "m1" {
		t.Fatalf("data = %+v", data)
	}
}

func TestMembersRemoveCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		removeMemberFromCardFn: func(ctx context.Context, cardID, memberID string) error {
			if cardID != "c1" || memberID != "m1" {
				t.Fatalf("card/member = %q/%q", cardID, memberID)
			}
			return nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"members", "remove", "--card", "c1", "--member", "m1"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("members remove failed: %v", err)
	}
	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}
	data := envelope["data"].(map[string]any)
	if data["success"] != true || data["id"] != "m1" {
		t.Fatalf("data = %+v", data)
	}
}
